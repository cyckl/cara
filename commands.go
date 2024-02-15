//
//	Commands list / helpers module
//

package main

import (
	"log"
	"github.com/bwmarrin/discordgo"
)

// Establish available commands
var (
	commands = []*discordgo.ApplicationCommand{
		// Misc
		{
			Name:			"cara",
			Description:	"About Cara",
		},
		{
			Name:			"guild",
			Description:	"About current guild",
		},
		{
			Name:			"help",
			Description:	"Command help",
		},
		{
			Name:			"sticks",
			Description:	"Pick random members",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionInteger,
					Name:			"count",
					Description:	"Amount of members to pick at once",
					Required:		false,
				},
			},
		},
		{
			Name:			"avatar",
			Description:	"Get a member's avatar",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionUser,
					Name:			"member",
					Description:	"Member to steal from",
					Required:		true,
				},
				{
					Type:			discordgo.ApplicationCommandOptionInteger,
					Name:			"size",
					Description:	"Size in px",
					Required:		false,
				},
			},
		},
		
		// Admin
		{
			Name:			"restore",
			Description:	"Restore user role (Admin)",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionRole,
					Name:			"package",
					Description:	"Package to restore",
					Required:		true,
				},
				{
					Type:			discordgo.ApplicationCommandOptionUser,
					Name:			"target",
					Description:	"Select restore target",
					Required:		false,
				},
			},
		},
		{
			Name:			"kick",
			Description:	"Kick user (Admin)",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionRole,
					Name:			"user",
					Description:	"User to kick",
					Required:		true,
				},
			},
		},
		{
			Name:			"status",
			Description:	"Change bot status",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionString,
					Name:			"status",
					Description:	"Write a status",
					Required:		true,
				},
				{
					Type:			discordgo.ApplicationCommandOptionInteger,
					Name:			"type",
					Description:	"type",
					Required:		true,
				},
			},
		},
		
		// Tweet
		{
			Name:			"tweet",
			Description:	"Generate a fake tweet",
			Options:		[]*discordgo.ApplicationCommandOption{
				{
					Type:			discordgo.ApplicationCommandOptionString,
					Name:			"content",
					Description:	"Tweet content",
					Required:		true,
				},
				{
					Type:			discordgo.ApplicationCommandOptionUser,
					Name:			"user",
					Description:	"Select tweet author",
					Required:		false,
				},
			},
		},
	}
)

func registerCmd(s *discordgo.Session) {
	// Register commands
	for _, v := range commands {
		log.Println("[Info] Registering:", v.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, config.ID.Guild, v)
		if err != nil {
			log.Println("[Error] Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func unregisterCmd(s *discordgo.Session) {
	// Get command list
	pubCommands, err := s.ApplicationCommands(s.State.User.ID, config.ID.Guild)
	if err != nil {
		log.Println("[Error] Cannot get application commands: %v", err)
		return
	}
	
	// Unregister commands by ID
	for _, v := range pubCommands {
		log.Println("[Info] Unregistering:", v.Name)
		err := s.ApplicationCommandDelete(s.State.User.ID, config.ID.Guild, v.ID)
		if err != nil {
			log.Println("[Error] Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	log.Println("[Info] Commands unregistered")
}
