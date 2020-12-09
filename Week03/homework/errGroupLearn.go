package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

type httpHandler struct {
}

func (h *httpHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {

}

func main() {
	ctx := context.Background()
	ins, cancelCtx := errgroup.WithContext(ctx)
	ins.Go(func() error {
		return handleSignal(cancelCtx)
	})
	ins.Go(func() error {
		return startServer(cancelCtx, ":8000", &httpHandler{})
	})
	if err := ins.Wait(); err != nil {
		fmt.Println("return error:", err.Error())
	}

	fmt.Println("main routine is done !")
}

func handleSignal(ctx context.Context) error {
	can := make(chan os.Signal)
	signal.Notify(can)
	fmt.Println("handle signal starting...")
	for {
		select {
		case s := <-can:
			return fmt.Errorf("get %v signal !", s)
		case <-ctx.Done():
			return fmt.Errorf("handle signal end !")
		}
	}
}

func startServer(ctx context.Context, address string, handler http.Handler) error {
	ins := http.Server{
		Addr:    address,
		Handler: handler,
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		fmt.Println("http server end !")

		ins.Shutdown(context.Background())
	}(ctx)
	fmt.Println("http server starting...!")

	return ins.ListenAndServe()
}