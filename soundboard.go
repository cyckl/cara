//
//	Soundboard module
//

package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/dgvoice"
)

func soundboard (s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command and see if the feature is even enabled
	if i.Data.Name != "soundboard" || config.Modules.Diegoism != true || config.AudioFiles == "" {
		return
	}
	
	// Get directory files
	files, err := ioutil.ReadDir(config.AudioFiles)
	if err != nil {
		log.Println("[Error][Soundboard] Failed to read directory:", err)
		return
	}
	
	// Randomize list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	
	// Prepend path with filename
	clip := config.AudioFiles + files[0].Name()
	
	// Interaction respond
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:     		3454719,
					Author: &discordgo.MessageEmbedAuthor{
						Name:   	"Now playing",
						IconURL:	"https://cdn.discordapp.com/avatars/712495459070902274/36f9d5bcf2586633791b60dc2f34dc64.png?size=512",
					},
					Description:	files[0].Name(),
				},
			},
		},
	})
	
	// Find the guild for that interaction
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		log.Println("[Error][Soundboard] Failed to get guild state:", err)
		return
	}
	
	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			// Join main voice channel for server (config)
			v, err := s.ChannelVoiceJoin(config.ID.Guild, vs.ChannelID, false, true)
				if err != nil {
				log.Println("[Error][Soundboard] Failed to connect to channel:", err)
				return
			}
		
			// Respond to command event with sound
			dgvoice.PlayAudioFile(v, clip, make(chan bool))

			v.Disconnect()
		}
	}
}
