package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type Blacklist struct {
	Roles	[]string	`json:"roles"`
}

func selfEsteem(s *discordgo.Session, e *discordgo.GuildRoleUpdate) {
	// Check if enabled / file exists
	if config.Modules.SelfEsteem != true || config.BlacklistData == "" {
		return
	}
	
	// Read from JSON list file
	data, err := ioutil.ReadFile(config.BlacklistData)
	if err != nil {
    	log.Println("[Error][Self Esteem] Failed to open file:", err)
    	return
    }
    
    // Convert to go native type
    var blacklist Blacklist
    err = json.Unmarshal(data, &blacklist)
    if err != nil {
    	log.Println("[Error][Self Esteem] Could not unmarshal JSON data:", err)
    	return
    }
    
    // Iterate through entries in blacklist. If the role in event matches, delete.
	for _, role := range blacklist.Roles {
		if strings.Contains(strings.ToLower(e.Role.Name), role) {
			err := s.GuildRoleDelete(e.GuildID, e.Role.ID)
 			if err != nil {
				log.Println("[Error][Self Esteem] Could not delete role:", err)
			}
		}
	}
}
