package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

const (
	DefaultHTTPAddr = ":8080"
	DefaultHTTPURL  = "http://localhost:8080"
)

var wg sync.WaitGroup

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello"))
	})

	server := &http.Server{Addr: DefaultHTTPAddr}
	wg.Add(1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("%s", err)
		}
		wg.Done()
	}()

	res, err := http.Get(DefaultHTTPURL)
	if err != nil {
		log.Fatalf("failed to get %s", err)
	} else {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	_ = server.Shutdown(context.Background())
	wg.Wait()
}
