package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// getArguments split the "parameter from a command
func getArguments(message string) []string {
	trimMessage := strings.TrimSpace(message)
	splitMessage := strings.Split(trimMessage, " ")
	return splitMessage
}

func checkPrefix(message string) bool {
	// check if the prefix or string is greater than 1 character
	if len(message) > 10 {
		log.Warn("Checking prefix; prefix must be less than 10 character")
		return false
	}
	return true
}

// MemberHasPermission checks if a member has the given permission
// for example, If you would like to check if user has the administrator
// permission you would use
// --- MemberHasPermission(s, guildID, userID, discordgo.PermissionAdministrator)
// If you want to check for multiple permissions you would use the bitwise OR
// operator to pack more bits in. (e.g): PermissionAdministrator|PermissionAddReactions
// =================================================================================
//     s          :  discordgo session
//     guildID    :  guildID of the member you wish to check the roles of
//     userID     :  userID of the member you wish to retrieve
//     permission :  the permission you wish to check for
func memberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}
	return false, nil
}

// check latency or how long a function takes to execute
func measureTime(funcName string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Time taken by %s function is %v \n", funcName, time.Since(start))
	}
}

// check allowed channel to use bot command.
func checkAllowedChannel(id string, settings *model.GuildSettings) bool {
	channelIDs := settings.GuildBotChannelsID
	for _, ID := range channelIDs {
		if id == ID {
			return true
		}
	}
	return false
}

// gets bot last message ID from the channel.
func getBotMessageID(session *discordgo.Session, msgEvent *discordgo.MessageCreate) (string, error) {
	// add bot last message to the array after the author ID.
	channelMessages, err := session.ChannelMessages(msgEvent.ChannelID, 1, "", msgEvent.Message.ID, "")
	if err != nil {
		log.Error("Failed to get messages from channel")
		return "", err
	}
	botMessageID := channelMessages[0].ID

	return botMessageID, nil
}

// check skip backward reaction
func checkMessageReaction(session *discordgo.Session, msgEvent *discordgo.MessageCreate, botMessageID string) (map[string]bool, error) {
	checkReaction, err := session.MessageReactions(msgEvent.ChannelID, botMessageID, "⏮️", 10, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message Reactions")
		return map[string]bool{
			"FastBack":    false,
			"Back":        false,
			"Stop":        false,
			"Forward":     false,
			"FastForward": false,
		}, err
	}

	// Check skip back reaction
	for _, reactions := range checkReaction {
		if reactions.ID == msgEvent.Author.ID {
			return map[string]bool{
				"FastBack":    true,
				"Back":        false,
				"Stop":        false,
				"Forward":     false,
				"FastForward": false,
			}, nil
		}
	}
	// Check back reaction
	checkReaction, err = session.MessageReactions(msgEvent.ChannelID, botMessageID, "◀️", 10, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message Reactions")
		return map[string]bool{
			"FastBack":    false,
			"Back":        false,
			"Stop":        false,
			"Forward":     false,
			"FastForward": false,
		}, err
	}

	// Check back reaction
	for _, reactions := range checkReaction {
		if reactions.ID == msgEvent.Author.ID {
			return map[string]bool{
				"FastBack":    false,
				"Back":        true,
				"Stop":        false,
				"Forward":     false,
				"FastForward": false,
			}, nil
		}
	}

	// Check stop reaction
	checkReaction, err = session.MessageReactions(msgEvent.ChannelID, botMessageID, "⏹️", 10, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message Reactions")
		return map[string]bool{
			"FastBack":    false,
			"Back":        false,
			"Stop":        false,
			"Forward":     false,
			"FastForward": false,
		}, err
	}

	// Check stop reaction l
	for _, reactions := range checkReaction {
		if reactions.ID == msgEvent.Author.ID {
			return map[string]bool{
				"FastBack":    false,
				"Back":        false,
				"Stop":        true,
				"Forward":     false,
				"FastForward": false,
			}, nil
		}
	}
	// Check forward reaction
	checkReaction, err = session.MessageReactions(msgEvent.ChannelID, botMessageID, "▶️", 10, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message Reactions")
		return map[string]bool{
			"FastBack":    false,
			"Back":        false,
			"Stop":        false,
			"Forward":     false,
			"FastForward": false,
		}, err
	}

	// Check forward reaction
	for _, reactions := range checkReaction {
		if reactions.ID == msgEvent.Author.ID {
			return map[string]bool{
				"FastBack":    false,
				"Back":        false,
				"Stop":        false,
				"Forward":     true,
				"FastForward": false,
			}, nil
		}
	}

	checkReaction, err = session.MessageReactions(msgEvent.ChannelID, botMessageID, "⏭️", 10, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message Reactions")
		return map[string]bool{
			"FastBack":    false,
			"Back":        false,
			"Stop":        false,
			"Forward":     false,
			"FastForward": false,
		}, err
	}

	// Check back reaction
	for _, reactions := range checkReaction {
		if reactions.ID == msgEvent.Author.ID {
			return map[string]bool{
				"FastBack":    false,
				"Back":        false,
				"Stop":        false,
				"Forward":     false,
				"FastForward": true,
			}, nil
		}
	}
	return nil, nil
}

// checks user reaction selection.
func checkUserReactionSelect(page int, currentTime time.Time, botMessageID string, s *discordgo.Session, m *discordgo.MessageCreate) (int, error) {
	errorVal := 10
	for {
		timePassed := time.Since(currentTime)
		if timePassed.Seconds() >= 30 {
			log.WithFields(log.Fields{
				"Time passed": timePassed,
			}).Info("Removing reactions time has been passed.")
			err := s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
			if err != nil {
				log.Error("Failed to remove all reaction from botMessage: ", err)
				return errorVal, err
			}
			previousAuthor = ""
			return errorVal, err
		}
		// check if the reaction matches the author ID aka sender
		checkReaction, err := checkMessageReaction(s, m, botMessageID)
		if err != nil {
			log.Error("Failed to check emoji from bot message:", err)
			return errorVal, err
		}
		// if true then remove it uwu.
		if checkReaction["Stop"] {
			err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
			if err != nil {
				log.Error("Failed to remove all reaction from bot message stop reaction.", err)
				return 0, err
			}
			previousAuthor = ""
			return errorVal, err
		} else if checkReaction["FastBack"] {
			page = 0
			// remove user reaction before going to next page.
			err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏮️", m.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message:", err)
				return errorVal, err
			}
			return page, nil
		} else if checkReaction["Back"] {
			if page == 1 {
				log.Info("Already on page one not doing anything")
				// remove user reaction before going to next page.
				err = s.MessageReactionRemove(m.ChannelID, botMessageID, "◀️", m.Author.ID)
				if err != nil {
					log.Error("Failed to remove user reaction from bot message: ", err)
					return errorVal, err
				}
				return 0, nil
			}
			page -= 1
			return page, nil
		} else if checkReaction["Forward"] {
			if page == 5 {
				log.Info("Last Page already not doing anything")
				page = 5
			}
			err = s.MessageReactionRemove(m.ChannelID, botMessageID, "▶️", m.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message: ", err)
				return errorVal, err
			}
			page++
			return page, nil
		} else if checkReaction["FastForward"] {
			page = 5
			// remove user reaction before going to next page.
			err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏭️", m.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message: ", err)
				return errorVal, err
			}
			return page, nil
		}
	}
}
