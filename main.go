package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infoln("launching example app")

	rootMux := setupRoutes()

	// spew.Dump(rootMux)

	rootMux.Walk(walk)

	err := http.ListenAndServe("localhost:8080", rootMux)
	if err != nil {
		logrus.Infoln(err)
	}
	logrus.Infoln("all done")
}

func walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	tmpl, _ := route.GetPathTemplate()
	fmt.Printf("route: %s, handler: %v\n", tmpl, route.GetHandler())
	return nil
}

// gmux
// |
// |\
// |  > s1
// |    |
// |    |\
// |    |  > s1a
// |    |
// |     \
// |       > s1b
// |
// |\
// |  > s2
// |    |
// |     \
// |       > s2a
// |
//  \
//    > s3

// things to test
// which middleware is run when?
// what if different submuxes have the same route
// what if different submuxes have a prefix and a full path that match

func setupRoutes() *mux.Router {
	gmux := mux.NewRouter()
	s1 := gmux.PathPrefix("/s1").Subrouter()
	s2 := gmux.PathPrefix("/s2").Subrouter()
	s3 := gmux.PathPrefix("/s3").Subrouter()

	indMux := mux.NewRouter()
	indMux.Path("/echo").HandlerFunc(e3)
	indMux.PathPrefix("/").HandlerFunc(echo)

	gmux.Handle("/ind", indMux)

	gmux.PathPrefix("/").HandlerFunc(e1)
	gmux.Path("/")

	// s2 := gmux.NewRoute().Subrouter()

	// s1a := s1.NewRoute().Subrouter()

	// s1b := s1.NewRoute().Subrouter()

	// s2a := s2.NewRoute().Subrouter()

	s1.Use(m1)
	s2.Use(m2)
	s3.Use(m3)
	gmux.Use(mroot)

	// first match wins. s2/actuallys1 is routable, s2/actuallys3 is not. It is
	// shadowed by the s2 path prefix, which always matches. Note that the order
	// of these three lines doesn't matter.  First match happens _within a
	// router_ and submuxes are matched in the order they are added to the root
	// mux. s1 will *always* come befroe s2, regardless of where the lines
	// appear.  But within s2, insertion order of routes matters.
	s1.Path("/s2/actuallys1").HandlerFunc(e1) // curl s2/actuallys1
	// s2.NewRoute().HandlerFunc(e2)               // curl /s2/*
	s2.Path("/s2/reallys2").HandlerFunc(e2real) // curl /s2/reallys2
	s3.Path("/s2/wontwork").HandlerFunc(e3)     // curl s2/wontwork won't match
	s3.Path("/s3/easy").HandlerFunc(e3)         // curl s3/easy
	s3.Path("/{myvar}").HandlerFunc(e3var)
	s3.Path("/foo").HandlerFunc(e3fixed)

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

func e2real(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e2real")
	w.Write([]byte("found e2real\n"))
	logrus.Infoln("ending e2real")
}

func e3(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e3")
	w.Write([]byte("found e3\n"))
	logrus.Infoln("ending e3")
}
func e3var(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e3var")
	word := mux.Vars(r)["myvar"]
	w.Write([]byte("found e3var: " + word + "\n"))
	logrus.Infoln("ending e3var")
}
func e3fixed(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e3fixed")
	w.Write([]byte("found e3fixed\n"))
	logrus.Infoln("ending e3fixed")
}

func echo(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting echo")
	w.Write([]byte("found echo\n"))
	route := mux.CurrentRoute(r)
	spew.Dump(route)
	logrus.Infoln("ending echo")
}
