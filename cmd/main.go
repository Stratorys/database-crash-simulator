package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stratorys/database-crash-simulator/pkg/config"
	"log"
	"time"
)

func main() {

	// flag to check if we want to panic to try to lose connection in void
	mustPanic := flag.Bool("must-panic", false, "if argument is true then application will crash after 15s")
	flag.Parse()

	cfg := config.Load()

	dsn := fmt.Sprintf("user=%s password='%s' host=%s port=%d dbname=%s search_path=%s sslmode=disable",
		cfg.Config.PGUser,
		cfg.Config.PGPassword,
		cfg.Config.PGHost,
		cfg.Config.PGPort,
		cfg.Config.PGDatabase,
		cfg.Config.PGSchema,
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
