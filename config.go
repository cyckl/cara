//
//	Config variables module
//

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const version = "2.0.0"
const configData = "./data/config.json"

// Set build date variable to be set with compilation flag
// 		go build -ldflags "-X main.BuildDate=$(date -I)"
var buildDate string
var commit string

type Config struct {
	Website			string		`json:"website"`
	AudioFiles		string		`json:"audioFiles"`
	RespondData		string		`json:"respondData"`
	MemberData		string		`json:"memberData"`
	BlacklistData	string		`json:"blacklistData"`
	ID				ID			`json:"ID"`
	Modules			Modules		`json:"modules"`
}

type ID struct {
	// Firmcoded IDs
	Owner			string		`json:"owner"`
	Guild			string		`json:"guild"`
}

type Modules struct {
	// Interactive modules
	Respond			bool	`json:"respond"`
	Info			bool	`json:"info"`
	Sticks			bool	`json:"sticks"`
	Avatar			bool	`json:"avatar"`
	Tiktok			bool	`json:"tiktok"`
	Cry				bool	`json:"cry"`
	Lights			bool	`json:"lights"`
	Soundboard		bool	`json:"soundboard"`
	Tweet			bool	`json:"tweet"`
	
	// Background modules
	SelfEsteem		bool	`json:"selfEsteem"`
	Animals			bool	`json:"animals"`
	Birthday		bool	`json:"birthday"`
	
	Abuse			bool	`json:"abuse"`
	Guts			bool	`json:"guts"`
}

var config Config

func configParse() {
	data, err := ioutil.ReadFile(configData)
	if err != nil {
    	log.Fatalf("[Fatal][Config] Failed to open file: %v", err)
    }
    
    // Link JSON data slice to phrases
    err = json.Unmarshal(data, &config)
    if err != nil {
    	log.Fatalf("[Fatal][Config] Could not unmarshal JSON data: %v", err)
    }
    
    // Check if commit and build are empty
    if buildDate == "" {
    	buildDate = "Unknown"
    }
    
    if commit == "" {
    	commit = "Unknown"
    }
}
