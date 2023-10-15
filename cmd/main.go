package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	// flag to check if we want to panic to try to lose connection in void
	mustPanic := flag.Bool("must-panic", false, "if argument is true then application will crash after 15s")
	flag.Parse()

	fmt.Println("debug ", *mustPanic)

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("port is not a number")
	}

	dsn := fmt.Sprintf("user=%s password='%s' host=%s port=%d dbname=%s search_path=%s sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		port,
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_SCHEMA"),
	)
	isConnected := TestDBConnection(dsn)

	if isConnected {
		log.Printf("connection is alive")
	}

	time.Sleep(time.Second * 15)

	if *mustPanic == true {
		panic("crash after 15s")
	}
}

func TestDBConnection(connectionString string) bool {
	// Open a database connection
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to the database: %v", err)
		return false
	}

	// Ping the database to test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping the database: %v", err)
		return false
	}

	// Connection is successful
	return true
}
