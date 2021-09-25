//
//	Misc commands module
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"strconv"
	"time"
	"github.com/bwmarrin/discordgo"
)

func guildInfo(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "guild" || config.Modules.Info != true {
		return
	}
	
	// Get guild info
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		log.Println("[Error][Guild] Failed to get guild data", err)
		return
	}
	
	// Respond to command event with embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			// One entry in the array of embeds... Why did they make this a fucking array?
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:			g.Name,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:		fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s", g.ID, g.Icon),
					},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Members", Value: strconv.Itoa(g.MemberCount), Inline: false},
						{Name: "Owner", Value: fmt.Sprintf("<@%s>", g.OwnerID), Inline: false},
						{Name: "VC Region", Value: g.Region, Inline: true},
						{Name: "ID", Value: g.ID, Inline: true},
						{Name: "Created", Value: snowflakeDate(g.ID), Inline: false},
					},
				},
			},
		},
	})
}

func about(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "cara" || config.Modules.Info != true {
		return
	}
	
	// Respond to command event with embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL:	"https://cdn.discordapp.com/avatars/712495459070902274/36f9d5bcf2586633791b60dc2f34dc64.png?size=512",
					},
					Title: "Cara by cyckl",
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "About",
						IconURL: "https://cdn.discordapp.com/avatars/712495459070902274/36f9d5bcf2586633791b60dc2f34dc64.png?size=512",
					},
					Color:       16725166,
					Description: "General purpose bot for this server",
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Why do I even maintain this?",
					},
					Fields: []*discordgo.MessageEmbedField{
						{Name:	"Version",		Value: version,				Inline: true},
						{Name:	"Build date",	Value: buildDate,			Inline: true},
						{Name:	"Commit",		Value: commit,				Inline: true},
						{Name:	"Website",		Value: config.Website,		Inline: false},
					},
				},
			},
		},
	})
}

func help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "help" || config.Modules.Info != true {
		return
	}
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name:    "Help website",
						IconURL: "https://cdn.discordapp.com/avatars/712495459070902274/36f9d5bcf2586633791b60dc2f34dc64.png?size=512",
					},
					Color:       3454719,
					Description: config.Website,
				},
			},
		},
	})
}

func sticks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "sticks" || config.Modules.Sticks != true {
		return
	}
	
	type ServerMember struct {
		ID			string		`json:"id"`
		Name		string		`json:"name"`
		DOB			string		`json:"birthday"`
	}
	
	// Read from JSON list file
	data, err := ioutil.ReadFile(config.MemberData)
	if err != nil {
		log.Println(fmt.Sprintf("[Error][Sticks] Failed to open file: %s", err))
		return
	}
	
	// Link JSON data slice to memberList
	var memberDB []ServerMember
	err = json.Unmarshal(data, &memberDB)
	if err != nil {
		log.Println(fmt.Sprintf("[Error][Sticks] Could not unmarshal JSON data: %s", err))
		return
	}
	
	// Randomize list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(memberDB), func(i, j int) { memberDB[i], memberDB[j] = memberDB[j], memberDB[i] })
	
	// Handle count argument
	var list strings.Builder
	if len(i.Data.Options) == 1 {
		// Handle if argument is larger than slice length
		if i.Data.Options[0].IntValue() > int64(len(memberDB)) {
			// Join slice of members into one string
			for r := 0; r < len(memberDB); r++ {
				list.WriteString(fmt.Sprintf("%s\n", memberDB[r].Name))
			}
		} else {
			for r := 0; int64(r) < i.Data.Options[0].IntValue(); r++ {
				list.WriteString(fmt.Sprintf("%s\n", memberDB[r].Name))
			}
		}
	} else {
		// Print only first of members slice without newline
		list.WriteString(memberDB[0].Name)
	}
	
	// Respond to command event with user(s)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Content: list.String(),
		},
	})
}

func avatar(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command and see if the feature is even enabled
	if i.Data.Name != "avatar" || config.Modules.Avatar != true {
		return
	}
	
	// Get avatar URL
	var url string
	if i.Data.Options[0].UserValue(s).Avatar == "" {
		// Default profile picture handling
		discrim, err := strconv.Atoi(i.Data.Options[0].UserValue(s).Discriminator)
		if err != nil {
			log.Println("[Error][Avatar] Could not convert string to int: %v", err)
		}
		url = fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%s.png", strconv.Itoa(discrim % 5))
	} else {
		url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", i.Data.Options[0].UserValue(s).ID, i.Data.Options[0].UserValue(s).Avatar)
	}
	
	// Handle size arg
	if len(i.Data.Options) == 2 {
		url = fmt.Sprintf("%s?size=%d", url, i.Data.Options[1].IntValue())
	}
	
	// Respond to command event with picture
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			Content: url,
		},
	})
}

// Snowflake timestamp decoder
func snowflakeDate(id string) string {
	// Discord "epoch" since 2015, then get msec in sec
	epoch := 1420070400000 / 1000
	
	// Convert ID to snowflake int
	snowflake, err := strconv.Atoi(id)
	if err != nil {
		log.Println("[Error][Snowflake] Failed to convert ID string to integer:", err)
		return "string conv error"
	}
	// Convert msec to sec
	snowflake = snowflake / 1000
	
	// 1. Right-shift 22 bits to preserve only time data
	// 2. Add Discord epoch as a UNIX timestamp in sec
	// 3. Convert new creation timestamp to int64
	// 4. Return date string as ISO 8601
	date := time.Unix(int64((snowflake >> 22) + epoch), 0)
	return date.Format("2006-01-02 15:04:05 UTC")
}
