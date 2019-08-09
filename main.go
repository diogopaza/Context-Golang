package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	http.HandleFunc("/", EchoHandler)

	errs := make(chan error, 2)

	go func() {
		fmt.Println("escutando em http://localhost:8081")
		errs <- http.ListenAndServe(":8081", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println("Finalizado ", <-errs)

}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	echo := r.URL.Query().Get("key")
	w.WriteHeader(200)
	w.Write([]byte(echo))
}
