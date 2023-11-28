package app

import (
	"log"
	"net"
	"net/http"
	"sync"
)

func (a *App) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		log.Fatal(a.runGRPC())
	}()

	go func() {
		defer wg.Done()

		log.Fatal(a.runHTTP())
	}()

	go func() {
		defer wg.Done()

		log.Fatal(a.runDocumentation())
	}()

	wg.Wait()

	return nil
}

func (a *App) runGRPC() error {
	listener, err := net.Listen("tcp", a.appConfig.App.PortGRPC)
	if err != nil {
		return err
	}

	log.Println("Admin GRPC server running on port:", a.appConfig.App.PortGRPC)

	return a.grpcAdminServer.Serve(listener)
}

func (a *App) runHTTP() error {
	log.Println("Admin HTTP server running on port:", a.appConfig.App.PortHTTP)

	return http.ListenAndServe(a.appConfig.App.PortHTTP, a.muxAdmin)
}

func (a *App) runDocumentation() error {
	log.Println("Admin Documentation server running on port:", a.appConfig.App.PortDocs)

	return http.ListenAndServe(a.appConfig.App.PortDocs, a.reDoc.Handler())
}
