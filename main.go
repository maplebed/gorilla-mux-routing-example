package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("launching example app")

	brokenMux := setupBrokenMiddleware()
	workingMux := setupWorkingMiddleware()

	go http.ListenAndServe("localhost:8080", brokenMux)
	http.ListenAndServe("localhost:8081", workingMux)

	logrus.Infoln("all done")
}

func setupBrokenMiddleware() *mux.Router {

	// create a mux and some submuxes
	gmux := mux.NewRouter()
	s1 := gmux.NewRoute().Subrouter()
	s2 := gmux.NewRoute().Subrouter()
	s3 := gmux.NewRoute().Subrouter()

	// add some middleware
	gmux.Use(mroot)
	s1.Use(m1)
	s2.Use(m2)
	s3.Use(m3)

	// add some routes
	s1.Path("/s1/hello").HandlerFunc(e1)
	s2.Path("/s2/hello").HandlerFunc(e2)
	s3.Path("/s3/hello").HandlerFunc(e3)

	return gmux
}

func setupWorkingMiddleware() *mux.Router {

	// create a mux and some submuxes, this time with PathPrefixes
	gmux := mux.NewRouter()
	s1 := gmux.PathPrefix("/s1").Subrouter()
	s2 := gmux.PathPrefix("/s2").Subrouter()
	s3 := gmux.PathPrefix("/s3").Subrouter()

	// add some middleware
	gmux.Use(mroot)
	s1.Use(m1)
	s2.Use(m2)
	s3.Use(m3)

	// add some routes (remember the path prefixes)
	s1.Path("/hello").HandlerFunc(e1)
	s2.Path("/hello").HandlerFunc(e2)
	s3.Path("/hello").HandlerFunc(e3)

	return gmux
}

func mroot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mroot")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mroot")
	})
}

func m1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting m1")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving m1")
	})
}

func m2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting m2")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving m2")
	})
}

func m3(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting m3")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving m3")
	})
}

func e1(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e1")
	w.Write([]byte("found e1\n"))
	logrus.Infoln("ending e1")
}

func e2(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e2")
	w.Write([]byte("found e2\n"))
	logrus.Infoln("ending e2")
}

func e3(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e3")
	w.Write([]byte("found e3\n"))
	logrus.Infoln("ending e3")
}
