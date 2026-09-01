package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


var DB *sql.DB

func Dbinit() {

	if err:=godotenv.Load(); err!=nil{
		log.Println("no .env file found, looking for system environment variables")
	}

	connStr:=os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	var err error
	DB,err=sql.Open("postgres",connStr)
	if err!=nil{
		log.Fatal("Failed to open database:",err)
	}

	if err:=DB.Ping(); err!=nil{
		log.Fatal("Database unreachable:",err)
	}

	log.Println("Database connection successful")


}