//
//	Admin abuse module
//

package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/bwmarrin/discordgo"
)

var gutEnabled bool

func status(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "status" || config.Modules.Abuse != true {
		return
	}
	
	// Block unauthorized users
	if i.Member.User.ID != config.ID.Owner {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type:	discordgo.InteractionResponseChannelMessageWithSource,
			Data:	&discordgo.InteractionApplicationCommandResponseData{
				Content: "Unauthorized",
			},
		})
		log.Println(fmt.Sprintf("[Info][Status] %s tried to run /status!", i.Member.Nick))
		return
	}
	
	log.Println("[Error][Status] Debug:", i.Data.Options[0].StringValue())
	
	// Send update
	err := s.UpdateStatusComplex(*newUpdateStatusData(0, discordgo.ActivityType(i.Data.Options[1].IntValue()), i.Data.Options[0].StringValue(), ""))
	if err != nil {
		log.Println("[Error][Status] Failed to set status:", err)
		return
	}
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Content: "Done",
		},
	})
}

func gutEnable(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages, good practice
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if option enabled
	if config.Modules.Guts != true {
		return
	}

	if m.Author.ID == "248985105169645568" && strings.Contains(strings.ToLower(m.Content), "it takes real guts to be an organ donor") {
		gutEnabled = !gutEnabled
	}
}

func gutSpam(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages, good practice
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if option enabled
	if config.Modules.Guts != true || gutEnabled != true {
		return
	}
	
	if strings.Contains(m.Content, "u") {
		s.ChannelMessageSend(m.ChannelID, "It takes real guts to be an organ donor")
	}
}

func restore(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "restore" || config.Modules.Abuse != true {
		return
	}
	
	// Block unauthorized users
	if i.Member.User.ID != config.ID.Owner {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type:	discordgo.InteractionResponseChannelMessageWithSource,
			Data:	&discordgo.InteractionApplicationCommandResponseData{
				Content: "Unauthorized",
			},
		})
		log.Println(fmt.Sprintf("[Info][Restore] %s tried to run /restore!", i.Member.Nick))
		return
	}
	// Restore selected role to user
	if len(i.Data.Options) == 2 {
		s.GuildMemberRoleAdd(config.ID.Guild, i.Data.Options[1].UserValue(s).ID, i.Data.Options[0].RoleValue(s, config.ID.Guild).ID)
	} else {
		s.GuildMemberRoleAdd(config.ID.Guild, config.ID.Owner, i.Data.Options[0].RoleValue(s, config.ID.Guild).ID)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Content: "Authorized",
		},
	})
}

func kick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "kick" || config.Modules.Abuse != true {
		return
	}
	
	// Block unauthorized users
	if i.Member.User.ID != config.ID.Owner {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type:	discordgo.InteractionResponseChannelMessageWithSource,
			Data:	&discordgo.InteractionApplicationCommandResponseData{
				Content: "Unauthorized",
			},
		})
		log.Println("[Info][Kick] %s tried to run /kick!", i.Member.Nick)
		return
	}

	// Kick user!
	err := s.GuildMemberDelete(i.GuildID, i.Data.Options[0].UserValue(s).ID)
	if err != nil {
		log.Println("[Info][Kick] Failed to kick user:", err)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Content: "Done",
		},
	})
}

// Helper to set status
func newUpdateStatusData(idle int, activityType discordgo.ActivityType, name, url string) *discordgo.UpdateStatusData {
	usd := &discordgo.UpdateStatusData{
		Status: "online",
	}

	if idle > 0 {
		usd.IdleSince = &idle
	}

	if name != "" {
		usd.Activities = []*discordgo.Activity{{
			Name: name,
			Type: activityType,
			URL:  url,
		}}
	}

	return usd
}
