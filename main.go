package main

import (
	"fmt"
	"net/http"

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
	gmux.Use(mroot)

	gmux.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// first subrouter handles /s1/* and has a custom NotFoundHandler
	s1 := gmux.PathPrefix("/s1").Subrouter()
	s1.Use(m1)
	s1.NotFoundHandler = http.HandlerFunc(notFoundSubHandler)
	// but the root mux *also* has a path defined for a specific /s1/ endpoint
	gmux.HandleFunc("/s1/root", e1root)

	// second subrouter handles /s2/* and does _not_ have a custom NotFoundHandler
	s2 := gmux.PathPrefix("/s2").Subrouter()
	s2.Use(m2)
	// don't override the default NotFoundHandler for s2
	// the root mux *also* has a path for a specific /s2/ endpoint
	gmux.HandleFunc("/s2/root", e2root)

	return gmux
}

func mroot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mroot")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mroot\n")
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

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting notFound")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("responding notFound\n"))
	logrus.Infoln("ending notFound")
}

func notFoundSubHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting notFoundSub")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("responding notFoundSub\n"))
	logrus.Infoln("ending notFoundSub")
}

func e1(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e1")
	w.Write([]byte("found e1\n"))
	logrus.Infoln("ending e1")
}

func e1root(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e1root")
	w.Write([]byte("found e1root\n"))
	logrus.Infoln("ending e1root")
}

func e2(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e2")
	w.Write([]byte("found e2\n"))
	logrus.Infoln("ending e2")
}

func e2root(w http.ResponseWriter, r *http.Request) {
	logrus.Infoln("starting e2root")
	w.Write([]byte("found e2root\n"))
	logrus.Infoln("ending e2root")
}
