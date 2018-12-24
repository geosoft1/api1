package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	port = flag.String("port", "8080", "port")
)

/*
curl -I 'localhost:8080/api/'
curl -I 'localhost:8080/api/v1/'
curl -I 'localhost:8080/api/v1/status'
curl -I 'localhost:8080/api/v2/'
curl -I 'localhost:8080/api/v2/status'
curl -I 'localhost:8080/api/v1/status' -H "x-auth-token: admin"
curl -I 'localhost:8080/api/v1/status' -H "x-auth-token: notadmin"
*/
func main() {
	flag.Parse()
	var router = mux.NewRouter()
	var api = router.PathPrefix("/api").Subrouter()
	//	var api = router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
	//		return r.Header.Get("x-auth-token") == "admin"
	//	}).PathPrefix("/api").Subrouter()
	api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("x-auth-token") != "admin" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
	var api1 = api.PathPrefix("/v1").Subrouter()
	api1.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	api1.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	var api2 = api.PathPrefix("/v2").Subrouter()
	api2.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	api2.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	http.ListenAndServe(":"+*port, router)
}
