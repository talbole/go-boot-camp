package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gp, ctx := errgroup.WithContext(context.Background())
	s := &http.Server{
		Addr:           ":9090",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	gp.Go(func() error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello word")
		})
		return s.ListenAndServe()
	})

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)

	gp.Go(func() error {
		for {
			select {
			case <-sg:
				return s.Close()
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	if err := gp.Wait(); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println("game over !")

}
