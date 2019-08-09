package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ContextKey string

const (
	AuthContextKey = ContextKey("username")
)

func main() {

	http.Handle("/", WithAuth(echoHandler))

	errs := make(chan error, 1)

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

var echoHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value(ContextKey(AuthContextKey)).(string)
	fmt.Println(username)
	echo := r.URL.Query().Get("text")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprint(username, " ", echo)))

})

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.Header.Get("authorization")
		ctx := context.WithValue(r.Context(), AuthContextKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
