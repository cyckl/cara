package main

import (
	"database/sql"
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

func dictAddWord(s *discordgo.Session, word string) {
	// Open database
	db, err := sql.Open("sqlite3", "./resources/cara.db")
	if err != nil {
		fmt.Println("Error opening database: ", err)
	}
	// Add query
	statement, err := db.Prepare("INSERT INTO words (word) VALUES (?)")
	if err != nil {
		fmt.Println("Error adding entry to database: ", err)
	}
	statement.Exec(word)
	// Run role delete function
	dictDelRole(s, serverTarget)
}

func dictRemoveWord(word string) {
	// Open database
	db, err := sql.Open("sqlite3", "./resources/cara.db")
	if err != nil {
		fmt.Println("Error opening database: ", err)
	}
	// Delete query
	statement, err := db.Prepare("DELETE FROM words WHERE word=?")
	if err != nil {
		fmt.Println("Error deleting entry from database: ", err)
	}
	statement.Exec(word)
}

func dictDelRole(s *discordgo.Session, serverID string) {
	// Open db and make sure blacklist exists; create if missing
	db, err := sql.Open("sqlite3", "./resources/cara.db")
	if err != nil {
		fmt.Println("Error opening database: ", err)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, word TEXT)")
	if err != nil {
		fmt.Println("Error creating table in database: ", err)
	}
    statement.Exec()
	
	// Get role list
	roles, err := s.GuildRoles(serverID)
	if err != nil {
		fmt.Println("Error grabbing roles: ", err)
	}
    
	// Index roles and iterate through them, iterating through all the targets every time; deletes role on match
	for _, role := range roles {
		// Get word list, define word var
		targets, err := db.Query("SELECT word FROM words")
		if err != nil {
			fmt.Println("Error grabbing blacklist word: ", err)
		}
    	var word string
		for targets.Next() {
			// Assign db value to word var
    	    targets.Scan(&word)
    	    if strings.Contains(strings.ToLower(role.Name), word) {
 				err := s.GuildRoleDelete(serverID, role.ID)
 				if err != nil {
					fmt.Println("Error deleting role: ", err)
				}
 			}
    	}
  	}
}
