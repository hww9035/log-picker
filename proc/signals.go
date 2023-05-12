package proc

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// var done = make(chan struct{})

func init() {
	go func() {
		var profiler Stopper
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
				if profiler == nil {
					profiler = StartProfile()
				} else {
					profiler.Stop()
					profiler = nil
				}
			case syscall.SIGTERM:
			default:
				log.Println("Got unregistered signal:", v)
			}
		}
	}()
}
