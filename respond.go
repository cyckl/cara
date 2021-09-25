//
//	Hated *things* module
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"github.com/bwmarrin/discordgo"
)

type phrase struct {
	ID			int			`json:"id"`
	Phrase		string		`json:"phrase"`
	Response	string		`json:"response"`
	Reaction	string		`json:"reaction"`
}

func respond(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages, good practice
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	// Check if option enabled / config file exists
	if config.Modules.Respond != true || config.RespondData == "" {
		return
	}
	
	// Read from JSON list file
	data, err := ioutil.ReadFile(config.RespondData)
	if err != nil {
    	log.Println(fmt.Sprintf("[Error][Respond] Failed to open file: %s", err))
    	return
    }
    
    // Link JSON data slice to phrases
    var phrases []phrase
    err = json.Unmarshal(data, &phrases)
    if err != nil {
    	log.Println(fmt.Sprintf("[Error][Respond] Could not unmarshal JSON data: %s", err))
    	return
    }
	
	// For every slice entry run following:
	for _, entry := range phrases {
		if strings.Contains(strings.ToLower(m.Content), entry.Phrase) {
			if entry.Response != "" {
				s.ChannelMessageSend(m.ChannelID, entry.Response)
			} else if entry.Reaction != "" {
				s.MessageReactionAdd(m.ChannelID, m.ID, entry.Reaction)
			} else {
				s.MessageReactionAdd(m.ChannelID, m.ID, "🤮")
			}
		}
	}
}
