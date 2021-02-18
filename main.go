package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var (
	token string
	version string
	serverTarget string
)

func init() {
	serverTarget = "638862415408398351"
	version = "2.1.0"
	flag.StringVar(&token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	// Create new Discord session using bot token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating discord session: ", err)
		return
	}

	// Register the functions as the callbacks for events
	dg.AddHandler(cmd)
	dg.AddHandler(roleDel)
	dg.AddHandler(roleFix)
	dg.AddHandler(nickFix)
	dg.AddHandler(vomitComet)

	// Receiving message events, guild events, guild member events, DM events
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsDirectMessages)

	// Open websocket connection to Discord and listen
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	
	dg.UpdateStatus(0, version)
	
	// Wait here until EOF
	fmt.Println(fmt.Sprintf("cara %s", version))

	// Close Discord session cleanly
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

// General commands
func cmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages, good practice
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	// Almond hating
	if strings.Contains(strings.ToLower(m.Content), "almond") {
		s.ChannelMessageSend(m.ChannelID, "fuck you")
	}
	
	// Self-service nick
	if strings.HasPrefix(m.Content, "c.nick") {
		// Restrict to me
		if m.Author.ID == "249052912762880000" {
			s.GuildMemberNickname(m.GuildID, "@me", strings.TrimPrefix(m.Content, "c.nick "))
			s.ChannelMessageSend(m.ChannelID, "Nickname set! I'll get used to it.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Unauthorized!")
		}
	}
	
	// Self-service status
	if strings.HasPrefix(m.Content, "c.status") {
		// Restrict to me
		if m.Author.ID == "249052912762880000" {
			s.UpdateStatus(0, strings.TrimPrefix(m.Content, "c.status "))
			s.ChannelMessageSend(m.ChannelID, "Status set!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Unauthorized!")
		}
	}
	
	// Add word to blacklist db
	// NOTE: Please restrict to DMs in future
	if strings.HasPrefix(m.Content, "c.dbadd") {
		dictAddWord(s, strings.TrimPrefix(m.Content, "c.dbadd "))
		s.ChannelMessageSend(m.ChannelID, "Added \"" + strings.TrimPrefix(m.Content, "c.dbadd ") + "\" to the database.")
		// Log to testing server
		s.ChannelMessageSend("798605461669543967", "Added \"" + strings.TrimPrefix(m.Content, "c.dbadd ") + "\" to the database.")
	}
	
	// Remove words from blacklist db
	// NOTE: Please restrict to DMs in future
	if strings.HasPrefix(m.Content, "c.dbrm") {
		dictRemoveWord(strings.TrimPrefix(m.Content, "c.dbrm "))
		s.ChannelMessageSend(m.ChannelID, "Removed \"" + strings.TrimPrefix(m.Content, "c.dbrm ") + "\" from the database.")
	}
	
	// Get words from blacklist db
	// NOTE: Please restrict to DMs in future
	//if strings.HasPrefix(m.Content, "c.dbget") {
	//	dictGetWords()
	//	s.ChannelMessageSend(m.ChannelID, "Added \"" + strings.TrimPrefix(m.Content, "c.dbget ") + "\" to the database.")
	//}
	
	// Speak though bot
	if strings.HasPrefix(m.Content, "c.say") {
		if m.Author.ID == "249052912762880000" {
			s.ChannelMessageSend(serverTarget, strings.TrimPrefix(m.Content, "c.say "))
		} else {
			s.ChannelMessageSend(m.ChannelID, "Unauthorized!")
		}
	}
	
	// Privilege escalation in case of server kick/ban
	if m.Content == "c.issue" {
		if m.Author.ID == "249052912762880000" {
			s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, "638877378159837185") // foot fetish
			s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, "698763659186536468") // the boys
			s.ChannelMessageSend(m.ChannelID, "Re-issue complete.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Unauthorized!")
		}
	}

	// Help dialog
	if m.Content == "c.help" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Commands",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Help page",
			},
			Color:       16724804,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Are you purposely dense?",
			},

			Fields: []*discordgo.MessageEmbedField{
				{Name: "c.about", Value: "Daniel's ego boost", Inline: false},
				{Name: "c.help", Value: "Both self-explanatory and redundant", Inline: false},
			},
		})
	}
	
	// About dialog
	if m.Content == "c.about" || m.Content == "c.version" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Project Cara",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "About me",
			},
			Color:       16724804,
			Description: "I am a general purpose bot for this stupid fucking hellhole of a server",
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Your benevolent omnipotent ruler!",
			},

			Fields: []*discordgo.MessageEmbedField{
				{Name: "Bot version", Value: version, Inline: true},
			},
		})
	}
}

// Automatic role deletion
func roleDel(s *discordgo.Session, event *discordgo.GuildRoleUpdate) {
	dictDelRole(s, event.GuildID)
}

// Vomit image on command
// NOTE: Please limit to DMs in future
func vomitComet(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "c.vomit" {
		// Whitelisted users
		ids := []string{
			"249052912762880000",
		}
		
		// Image read
		img, err := os.Open("./resources/vomit.jpg")
		if err != nil {
			fmt.Println("Error opening image: ", err)
		}
		
		// Iterate through IDs to authenticate
		for _, id := range ids {
			if m.Author.ID == id {
				s.ChannelFileSend(serverTarget, "./resources/vomit.jpg", img)
			}
		}
	}
}

// Adjusts roles to preference
func roleFix(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	// Remove "baby" from user
	s.GuildMemberRoleRemove(event.GuildID, "321301712600039424", "771114327780491296")
	// Add "perfection"
	s.GuildMemberRoleAdd(event.GuildID, "321301712600039424", "787560189574119445") // A
	s.GuildMemberRoleAdd(event.GuildID, "634893718637248512", "787560189574119445") // S
}

// Restricts user from setting nick to "baby"
func nickFix(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	// ID coded to user
	member, _ := s.GuildMember(event.GuildID, "321301712600039424")
	if strings.EqualFold(member.Nick, "baby") {
		s.GuildMemberNickname(event.GuildID, "321301712600039424", "not baby")
	}
}
