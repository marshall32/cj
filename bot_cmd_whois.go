package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func commandWhois(cm CommandManager, args string, message discordgo.Message, contextual bool) (bool, bool, error) {
	var (
		verified bool
		err      error
		count    = 0
		username string
		link     string
		result   string
	)

	if len(message.Mentions) == 0 {
		var userID string
		userID, err = cm.App.GetDiscordUserFromForumName(args)
		if err != nil {
			return false, false, err
		}

		result += fmt.Sprintf("**%s** is here as <@%s>", args, userID)
	} else {
		for _, user := range message.Mentions {
			if count == 5 {
				break
			}
			count++

			if user.ID == cm.App.config.BotID {
				result += "I am Carl Johnson, co-leader of Grove Street Families. "
				continue
			}

			verified, err = cm.App.IsUserVerified(user.ID)
			if err != nil {
				result += err.Error()
				continue
			}

			if !verified {
				result += fmt.Sprintf("The user <@%s> is not verified. ", user.ID)
			} else {
				username, err = cm.App.GetForumNameFromDiscordUser(user.ID)
				if err != nil {
					return false, false, err
				}

				link, err = cm.App.GetForumUserFromDiscordUser(user.ID)
				if err != nil {
					return false, false, err
				}

				result += fmt.Sprintf("<@%s> is **%s** (%s) on SA-MP forums. ", user.ID, username, link)
			}
		}
	}

	_, err = cm.App.discordClient.ChannelMessageSend(message.ChannelID, result)
	if err != nil {
		return false, false, err
	}

	return true, false, nil
}
