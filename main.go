package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"github.com/gorilla/mux"
)

type task struct {
	
	ID int `json:"ID"`
	Name string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

var tasks = allTasks {
	{
		ID: 1,
		Name: "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request){
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)
	
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func getTask(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {

		if task.ID == taskID {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(task)
		}

	}
}

func deleteTask(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for index, task := range tasks {
		
		if task.ID == taskID {
			tasks = append(tasks[:index], tasks[index + 1:]...)
			fmt.Fprintf(w, "The task with ID %v has been remove succesfully", taskID)
		}

	}
}


func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Please enter valid date")
		return
	}

	json.Unmarshal(reqBody, &updatedTask)

	for index, task := range tasks {
		
		if task.ID == taskID {
			tasks = append(tasks[:index], tasks[index + 1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "The thask with ID %v has been updated succesfully", taskID)
		}

	}

}


func indexRoute(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to my API REST")
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}",getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}