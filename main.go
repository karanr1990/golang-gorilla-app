package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.ibm.com/Quest-CIO/go-micro-app/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS",false,":4000","Bind Address For the Server")

func main()  {

	l := log.New(os.Stdout,"product-api",log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()
	getRouter :=sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/",ph.GetProducts)

	putRouter :=sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}",ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter :=sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/",ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	s := &http.Server{
		Addr: ":4000",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}

	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan,os.Interrupt)
	signal.Notify(sigchan,os.Kill)

	sig := <- sigchan

	l.Println("recieved terminate...",sig)

	tc, _ := context.WithTimeout(context.Background(),30*time.Second)

	s.Shutdown(tc)
}

