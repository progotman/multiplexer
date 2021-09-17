package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/progotman/multiplexer/config/provider"
	"github.com/progotman/multiplexer/processor"
	httprequester "github.com/progotman/multiplexer/requester/http"
	"github.com/progotman/multiplexer/router"
	"github.com/progotman/multiplexer/router/handler"
)

func main() {
	cancelChannel := make(chan os.Signal, 1)
	signal.Notify(cancelChannel, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	cnf := provider.FallbackProvider{
		Provider:         provider.EnvironmentProvider{},
		FallbackProvider: provider.InMemoryProvider{},
	}.GetConfig()

	fmt.Printf("Configuration: %+v\n", cnf)

	rtr := router.Router{
		MaxNumberOfIncomingRequests: cnf.MaxNumberOfIncomingRequests,
		Routes: []router.Route{
			{
				Url:    "/process-urls",
				Method: http.MethodPost,
				Handler: &handler.ProcessUrlsHandler{
					Processor: processor.UrlsProcessor{
						Requester: httprequester.Requester{
							Client: http.Client{
								Timeout: time.Second * time.Duration(cnf.RequestTimeout),
							},
						},
						MaxNumberOfUrls:             cnf.MaxNumberOfUrls,
						MaxNumberOfOutgoingRequests: cnf.MaxNumberOfOutgoingRequests,
					},
				},
			},
		},
	}

	srv := http.Server{
		Addr:    cnf.ServerAddress,
		Handler: &rtr,
	}

	go func() {
		<-cancelChannel
		log.Println("Initial graceful shutdown")

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Shutdown error: %v\n", err)
		} else {
			log.Printf("Gracefully stopped\n")
		}

		cancel()
	}()

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()
}
