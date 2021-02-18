package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	reset := flag.NewFlagSet("reset", flag.ExitOnError)
	print := flag.NewFlagSet("print", flag.ExitOnError)
	switch os.Args[1] {
	case "reset":
		reset.Parse(os.Args[2:])
		resetDB()
	case "print":
		print.Parse(os.Args[2:])
		printDB()
	default:
		fmt.Println("expected reset or print command")
		os.Exit(1)
	}
}

func printDB() {
	// Open db
	db, _ := sql.Open("sqlite3", "../resources/cara.db")
	
	rows, _ := db.Query("SELECT id, word FROM words")
    var id int
    var word string
    for rows.Next() {
        rows.Scan(&id, &word)
        fmt.Println(strconv.Itoa(id) + ": " + word)
    }
}

func resetDB() {
	// Open db
	db, _ := sql.Open("sqlite3", "../resources/cara.db")
	
	// Wipe old words table
	statement, _ := db.Prepare("DROP TABLE words")
    statement.Exec()
    // Create new words table
	statement, _ = db.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, word TEXT)")
    statement.Exec()
	
	// Set word-list
	targets := []string{
		// Suicide
		"test1",
		"test2",
	}
	
	for _, word := range targets {
		// For every word, import it to db
		statement, _ := db.Prepare("INSERT INTO words (word) VALUES (?)")
		statement.Exec(word)
 	}
}
