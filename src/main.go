package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// init with logrus
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// main - mux init + router setup
func main() {
	log.Info("** Starting API server **")
	router := mux.NewRouter()
	router.HandleFunc("/apihealth", APIHealth).Methods("GET")
	http.ListenAndServe(":8000", router)
}

// return {"alive":true}
func APIHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health : Success")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
