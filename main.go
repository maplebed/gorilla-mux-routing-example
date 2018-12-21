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

	err := http.ListenAndServe("localhost:8080", rootMux)
	if err != nil {
		logrus.Infoln(err)
	}
	logrus.Infoln("all done")
}

func setupRoutes() *mux.Router {
	gmux := mux.NewRouter()

	// our service is reachable via either a public or a private host
	sPub := gmux.Host("public").Subrouter()
	sPriv := gmux.Host("private").Subrouter()

	mainSub := mux.NewRouter()
	mainSub.HandleFunc("/public/e1", e1)
	mainSub.HandleFunc("/private/e2", e2)

	// the public host can only reach endpionts that begin /public
	sPub.PathPrefix("/public").Handler(mainSub)
	// the private host can reach all endpoints
	sPriv.PathPrefix("/").Handler(mainSub)

	// the global mux always gets some basic middleware
	gmux.Use(mRoot)
	// to help identify which subrouter got traversed, let's add middleware
	sPub.Use(mPub)
	// in reality the private middleware would do authentication
	sPriv.Use(mPriv)
	// and finally, the mainSub router gets its own middleware
	mainSub.Use(mMain)

	return gmux
}

func mRoot(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mRoot")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mRoot")
		fmt.Println("")
	})
}

func mPub(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mPub")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mPub")
	})
}

func mPriv(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mPriv")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mPriv")
	})
}

func mMain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infoln("starting mMain")
		next.ServeHTTP(w, r)
		logrus.Infoln("leaving mMain")
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
