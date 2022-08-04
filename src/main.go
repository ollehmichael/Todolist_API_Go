package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

// db setup
var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

// Task Struct
type TaskStruct struct {
	Id          int `gorm: "primary_key,omitempty"`
	Description string
	Completed   bool
}

// init with logrus
func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

// main - mux init + router setup
func main() {
	defer db.Close()

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
