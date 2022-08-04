package main

import (
	"encoding/json"
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

	db.Debug().DropTableIfExists(&TaskStruct{})
	db.Debug().AutoMigrate(&TaskStruct{})

	log.Info("** Starting API server **")

	// init router
	router := mux.NewRouter()
	// GET
	router.HandleFunc("/apihealth", APIHealth).Methods("GET")

	// POST
	router.HandleFunc("/createtask", APIHealth).Methods("POST")

	http.ListenAndServe(":8000", router)
}

// return {"alive":true}
func APIHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health : Success")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

// CRUD Functions
// Create (POST)
func CreateTask(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("Description")
	log.WithFields(log.Fields{"Description": description}).Info("Add new Task. Saving to db.")
	task := &TaskStruct{Description: description, Completed: false}
	db.Create(&task)
	result := db.Last(&task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

// Read (GET) - All Tasks
func GetTasks(completed bool) interface{} {
	var tasks []TaskStruct
	TasksbyId := db.Where("completed = ?", completed).Find(&tasks).Value
	return TasksbyId
}

// Read (GET) - Unique Tasks
func GetTaskById(Id int) bool {
	task := &TaskStruct{}
	result := db.First(&task, Id)
	if result.Error != nil {
		log.Warn("Task does not exist")
		return false
	}

	return true
}
