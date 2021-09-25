//
//	Tweeter module
//

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

// Pueudorandom numbers
func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func tweet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if it's the right command or enabled
	if i.Data.Name != "tweet" || config.Modules.Tweet != true {
		return
	}
	
	// Set default endpoints for data (no selected user and no nickname)
	author := fmt.Sprintf("@%s", i.Member.User.Username)
	avatar := discordgo.EndpointUserAvatar(i.Member.User.ID, i.Member.User.Avatar)
	
	// Handle user selection and nicknames
	if len(i.Data.Options) == 2 {
		// If there is a selected user, set the author and avatar to the selected user instead of default
		author = fmt.Sprintf("@%s", i.Data.Options[1].UserValue(s).Username)
		avatar = discordgo.EndpointUserAvatar(i.Data.Options[1].UserValue(s).ID, i.Data.Options[1].UserValue(s).Avatar)
	} else if i.Member.Nick != "" {
		// If there is no selected user, at least check if the command author has a nickname
		author = fmt.Sprintf("%s (@%s)", i.Member.Nick, i.Member.User.Username)
	}

	// Respond to command event with embed
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type:	discordgo.InteractionResponseChannelMessageWithSource,
		Data:	&discordgo.InteractionApplicationCommandResponseData{
			// One entry in the array of embeds... Why did they make this a fucking array?
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name:		author,
						IconURL:	avatar,
					},
					Color:			1942002,
					Description:	i.Data.Options[0].StringValue(),
					Footer:	&discordgo.MessageEmbedFooter{
						Text:		"Twitter",
						IconURL:	"https://abs.twimg.com/icons/apple-touch-icon-192x192.png",
					},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Retweets", Value: strconv.Itoa(randInt(5000, 50000)), Inline: true},
						{Name: "Likes", Value: strconv.Itoa(randInt(25000, 150000)), Inline: true},
					},
				},
			},
		},
	})
}
