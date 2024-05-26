package gonweb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Hook func(c context.Context) error

func WaitShutdown(stopwaittime time.Duration,
	hookwaittime time.Duration,
	hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutDownSignals...)
	fmt.Println("listen stop signal......")
	select {
	case sig := <-signals:
		fmt.Printf("recv stop signal:%s,begin stop program\n", sig)
		time.AfterFunc(stopwaittime, func() {
			fmt.Println("Enforcement of closure tasks")
			os.Exit(0)
		})
		for _, hook := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), hookwaittime)
			err := hook(ctx)
			if err != nil {
				fmt.Printf("failure to execute the hook task :%s\n", ctx.Err())
			}
			cancel()
		}
		fmt.Println("all hook task exec finshed")
	}
	os.Exit(0)
}

func TestHookBuilder(s ...Server) Hook {
	return func(c context.Context) error {
		wg := sync.WaitGroup{}
		doneCh := make(chan struct{}, 1)
		wg.Add(len(s))
		for _, svr := range s {
			go func(server Server) {
				err := server.ShutDown()
				if err != nil {
					fmt.Printf("stop server error :%s\n", err)
				}
				wg.Done()
			}(svr)
		}

		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()

		select {
		case <-c.Done():
			return c.Err()
		case <-doneCh:
			return nil
		}
	}
}

type GracefulShutdown struct {
	Cnt     int64
	closing int32
	Done    chan struct{}
}

func (g *GracefulShutdown) RejectFilterBuilder(next GonHandlerFunc) GonHandlerFunc {
	return func(ctx *GonContext) {
		cl := atomic.LoadInt32(&g.closing)
		if cl > 0 {
			ctx.W.WriteHeader(http.StatusInternalServerError)
			ctx.W.Write([]byte("server error"))
			return
		}
		atomic.AddInt64(&g.Cnt, 1)
		next(ctx)
		cnt := atomic.AddInt64(&g.Cnt, -1)
		cl = atomic.LoadInt32(&g.closing)
		if cl > 0 && cnt == 0 {
			g.Done <- struct{}{}
		}
	}
}

func (g *GracefulShutdown) RejectHookBuilder() Hook {
	return func(c context.Context) error {
		atomic.AddInt32(&g.closing, 1)
		if atomic.LoadInt64(&g.Cnt) == 0 {
			return nil
		}
		select {
		case <-c.Done():
			return c.Err()
		case <-g.Done:
			return nil
		}
	}
}

func NewGracefulshutdown() *GracefulShutdown {
	return &GracefulShutdown{
		Done: make(chan struct{}),
	}
}
