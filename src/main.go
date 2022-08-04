package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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
	router.HandleFunc("/tasks-completed", GetCompletedTasks).Methods("GET")
	router.HandleFunc("/tasks-incomplete", GetIncompleteTasks).Methods("GET")

	// POST
	router.HandleFunc("/createtask", APIHealth).Methods("POST")
	router.HandleFunc("/task/{id}", UpdateTask).Methods("POST")

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

// Read (GET) - Completed Tasks
func GetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Completed Tasks")
	completedTasks := GetTasks(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTasks)
}

// Read (GET) - Incomplete Tasks
func GetIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Incomplete Tasks")
	incompleteTasks := GetTasks(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleteTasks)
}

// Update (POST)
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Setup
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := GetTaskById(id)
	if err == false {
		// Task does not exist
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Task does not exist"}`)
	} else {
		// If Task exists, Update Task Status to Completed
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating Task")
		task := &TaskStruct{}
		db.First(&task, id)
		task.Completed = completed
		db.Save(&task)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}
