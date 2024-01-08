package main

import (
	"fmt"
	"log"
	"mysql/controllers"
	"mysql/db"
	"net/http"
)

func main() {
	db.CreateConnection()
	defer db.DBcon.Close()
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/checkmarks", controllers.CheckMarks)
	http.HandleFunc("/studentprofile", controllers.StudentProfile)
	http.HandleFunc("/addmarks", controllers.AddMarks)
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
