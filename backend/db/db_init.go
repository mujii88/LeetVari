package db

import (
	"database/sql"
	"errors"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


var DB *sql.DB

type DB_conn struct {
	Error error

}

func Dbinit() DB_conn {

	if err:=godotenv.Load(); err!=nil{
		return DB_conn {
			Error: err,
		}
	}

	connStr:=os.Getenv("DATABASE_URL")
	if connStr == "" {
		return DB_conn {
			Error: errors.New("DATABASE_URL is not set in .env"),
		}
	}

	var err error
	DB,err=sql.Open("postgres",connStr)
	if err!=nil{
		return DB_conn {
			Error: err,
		}
	}

	if err:=DB.Ping(); err!=nil{
		return DB_conn {
			Error: err,
		}
	}


	return DB_conn{
		Error: nil,
	}

}