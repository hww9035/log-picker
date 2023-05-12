package proc

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync/atomic"
	"syscall"
	"time"
)

const defaultMemProfileRate = 4096

const timeFormat = "20160102150405"

var started uint32

type Profile struct {
	closers []func()

	stopped uint32
}

func (p *Profile) close() {
	for _, closer := range p.closers {
		closer()
	}
}

func (p *Profile) startBlockProfile() {
	fn := createDumpFile("block")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create block profile %q: %v", fn, err)
		return
	}

	runtime.SetBlockProfileRate(1)
	log.Printf("profile: block profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		_ = pprof.Lookup("block").WriteTo(f, 0)
		_ = f.Close()
		runtime.SetBlockProfileRate(0)
		log.Printf("profile: block profiling disabled, %s", fn)
	})
}

func (p *Profile) startCpuProfile() {
	fn := createDumpFile("cpu")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create cpu profile %q: %v", fn, err)
		return
	}

	log.Printf("profile: cpu profiling enabled, %s", fn)
	_ = pprof.StartCPUProfile(f)
	p.closers = append(p.closers, func() {
		pprof.StopCPUProfile()
		_ = f.Close()
		log.Printf("profile: cpu profiling disabled, %s", fn)
	})
}

func (p *Profile) startMemProfile() {
	fn := createDumpFile("mem")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create memory profile %q: %v", fn, err)
		return
	}

	old := runtime.MemProfileRate
	runtime.MemProfileRate = defaultMemProfileRate
	log.Printf("profile: memory profiling enabled (rate %d), %s", runtime.MemProfileRate, fn)
	p.closers = append(p.closers, func() {
		_ = pprof.Lookup("heap").WriteTo(f, 0)
		_ = f.Close()
		runtime.MemProfileRate = old
		log.Printf("profile: memory profiling disabled, %s", fn)
	})
}

func (p *Profile) startMutexProfile() {
	fn := createDumpFile("mutex")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create mutex profile %q: %v", fn, err)
		return
	}

	runtime.SetMutexProfileFraction(1)
	log.Printf("profile: mutex profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		if mp := pprof.Lookup("mutex"); mp != nil {
			_ = mp.WriteTo(f, 0)
		}
		_ = f.Close()
		runtime.SetMutexProfileFraction(0)
		log.Printf("profile: mutex profiling disabled, %s", fn)
	})
}

func (p *Profile) startThreadCreateProfile() {
	fn := createDumpFile("threadcreate")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create threadcreate profile %q: %v", fn, err)
		return
	}

	log.Printf("profile: threadcreate profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		if mp := pprof.Lookup("threadcreate"); mp != nil {
			_ = mp.WriteTo(f, 0)
		}
		_ = f.Close()
		log.Printf("profile: threadcreate profiling disabled, %s", fn)
	})
}

func (p *Profile) startTraceProfile() {
	fn := createDumpFile("trace")
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("profile: could not create trace output file %q: %v", fn, err)
		return
	}

	if err := trace.Start(f); err != nil {
		log.Printf("profile: could not start trace: %v", err)
		return
	}

	log.Printf("profile: trace enabled, %s", fn)
	p.closers = append(p.closers, func() {
		trace.Stop()
		log.Printf("profile: trace disabled, %s", fn)
	})
}

func (p *Profile) Stop() {
	if !atomic.CompareAndSwapUint32(&p.stopped, 0, 1) {
		return
	}
	p.close()
	atomic.StoreUint32(&started, 0)
}

func StartProfile() Stopper {
	if !atomic.CompareAndSwapUint32(&started, 0, 1) {
		log.Printf("profile: Start() already called")
		return noopStopper
	}

	var prof Profile
	prof.startCpuProfile()
	prof.startMemProfile()
	prof.startMutexProfile()
	prof.startBlockProfile()
	prof.startTraceProfile()
	prof.startThreadCreateProfile()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c

		log.Printf("profile: caught interrupt, stopping profiles")
		prof.Stop()

		signal.Reset()
		_ = syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	return &prof
}

func createDumpFile(kind string) string {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	return path.Join(os.TempDir(), fmt.Sprintf("%s-%d-%s-%s.pprof",
		command, pid, kind, time.Now().Format(timeFormat)))
}
