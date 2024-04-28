package utils

import (
    "fmt"
    "log"
    "os"
    "path"
    "runtime"
    "runtime/pprof"
    "runtime/trace"
    "sync/atomic"
    "time"
)

const defaultMemProfileRate = 4096
const timeFormat = "20160102150405"

var started uint32
var Pf profile

type Stopper interface {
    Stop() bool
}

type profile struct {
    LogFile string
    closers []func()
    stopped uint32
}

func (p *profile) close() {
    for _, closer := range p.closers {
        closer()
    }
    p.closers = make([]func(), 0)
}

// cpu使用分析
func (p *profile) startCpuProfile() {
    fn := createDumpFile("cpu", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create cpu profile %q: %v", fn, err)
        return
    }

    //logCtl.Logger.Info("profile: cpu profiling enabled", zap.Any("cpu_file", fn))
    _ = pprof.StartCPUProfile(f)
    p.closers = append(p.closers, func() {
        pprof.StopCPUProfile()
        _ = f.Close()
        //logCtl.Logger.Info("profile: cpu profiling disabled", zap.Any("cpu_file", fn))
    })
}

// goroutine分析
func (p *profile) startGoroutineProfile() {
    fn := createDumpFile("goroutine", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create goroutine profile %q: %v", fn, err)
        return
    }

    //logCtl.Logger.Info("profile: goroutine profiling enabled", zap.Any("goroutine_file", fn))
    p.closers = append(p.closers, func() {
        _ = pprof.Lookup("goroutine").WriteTo(f, 0)
        _ = f.Close()
        //logCtl.Logger.Info("profile: goroutine profiling disabled", zap.Any("goroutine_file", fn))
    })
}

// 阻塞的调用栈踪迹分析
func (p *profile) startBlockProfile() {
    fn := createDumpFile("block", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create block profile %q: %v", fn, err)
        return
    }

    runtime.SetBlockProfileRate(1)
    //logCtl.Logger.Info("profile: block profiling enabled", zap.Any("block_file", fn))
    p.closers = append(p.closers, func() {
        _ = pprof.Lookup("block").WriteTo(f, 0)
        _ = f.Close()
        runtime.SetBlockProfileRate(0)
        //logCtl.Logger.Info("profile: block profiling disabled", zap.Any("block_file", fn))
    })
}

// 堆内存分配分析
func (p *profile) startMemProfile() {
    fn := createDumpFile("heap", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create memory profile %q: %v", fn, err)
        return
    }

    old := runtime.MemProfileRate
    runtime.MemProfileRate = defaultMemProfileRate
    //logCtl.Logger.Info("profile: memory profiling enabled", zap.Any("rate", runtime.MemProfileRate), zap.Any("mem_file", fn))
    p.closers = append(p.closers, func() {
        _ = pprof.Lookup("heap").WriteTo(f, 0)
        _ = f.Close()
        runtime.MemProfileRate = old
        //logCtl.Logger.Info("profile: memory profiling disabled", zap.Any("mem_file", fn))
    })
}

// 锁分析
func (p *profile) startMutexProfile() {
    fn := createDumpFile("mutex", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create mutex profile %q: %v", fn, err)
        return
    }

    runtime.SetMutexProfileFraction(1)
    //logCtl.Logger.Info("profile: mutex profiling enabled", zap.Any("mutex_file", fn))
    p.closers = append(p.closers, func() {
        if mp := pprof.Lookup("mutex"); mp != nil {
            _ = mp.WriteTo(f, 0)
        }
        _ = f.Close()
        runtime.SetMutexProfileFraction(0)
        //logCtl.Logger.Info("profile: mutex profiling disabled", zap.Any("mutex_file", fn))
    })
}

// 线程创建分析
func (p *profile) startThreadCreateProfile() {
    fn := createDumpFile("threadcreate", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create threadcreate profile %q: %v", fn, err)
        return
    }

    //logCtl.Logger.Info("profile: thread-create profiling enabled", zap.Any("file", fn))
    p.closers = append(p.closers, func() {
        if mp := pprof.Lookup("threadcreate"); mp != nil {
            _ = mp.WriteTo(f, 0)
        }
        _ = f.Close()
        //logCtl.Logger.Info("profile: thread_create profiling disabled", zap.Any("file", fn))
    })
}

// trace分析
func (p *profile) startTraceProfile() {
    fn := createDumpFile("trace", p)
    f, err := os.Create(fn)
    if err != nil {
        log.Printf("profile: could not create trace profile %q: %v", fn, err)
        return
    }

    if err := trace.Start(f); err != nil {
        log.Printf("profile: could not start trace profile %q: %v", fn, err)
        return
    }

    //logCtl.Logger.Info("profile: trace enabled", zap.Any("trace_file", fn))
    p.closers = append(p.closers, func() {
        trace.Stop()
        //logCtl.Logger.Info("profile: trace disabled", zap.Any("trace_file", fn))
    })
}

func (p *profile) Stop() bool {
    if !atomic.CompareAndSwapUint32(&p.stopped, 0, 1) {
        return false
    }
    p.close()
    atomic.StoreUint32(&started, 0)
    return true
}

func Start(logFile string) Stopper {
    if !atomic.CompareAndSwapUint32(&started, 0, 1) {
        log.Printf("profile: Start already called")
        return &Pf
    }

    if logFile == "" {
        Pf.LogFile = os.TempDir()
    } else {
        Pf.LogFile = logFile
    }
    Pf.startCpuProfile()
    Pf.startMemProfile()
    Pf.startMutexProfile()
    Pf.startBlockProfile()
    Pf.startTraceProfile()
    Pf.startThreadCreateProfile()

    //go func() {
    //    c := make(chan os.Signal, 1)
    //    signal.Notify(c, syscall.SIGINT)
    //    <-c
    //
    //    logCtl.Logger.Info("profile: caught interrupt, stopping profiles")
    //    prof.Stop()
    //
    //    signal.Reset()
    //    _ = syscall.Kill(os.Getpid(), syscall.SIGINT)
    //}()

    return &Pf
}

func createDumpFile(kind string, profile *profile) string {
    //command := path.Base(os.Args[0])
    //pid := syscall.Getpid()
    return path.Join(profile.LogFile, fmt.Sprintf("%s-%s.pprof", kind, time.Now().Format(timeFormat)))
}
