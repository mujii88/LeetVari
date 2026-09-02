package main

import (
	"backend_clone/db"
	"fmt"
	"log"
)

func main() {
	conn:=db.Dbinit()
	if conn.Error!=nil{
		log.Fatal("Failed to connect to database:",conn.Error)
	}
	fmt.Println("Database connection successful")

}