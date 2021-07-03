package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"strconv"
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
		log.WithFields(log.Fields{
			"funcName": funcName,
			"time":     time.Since(start),
		}).Info("Time taken by function is completed")
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

// slows down the page changing by a lot :shrug:
//func removeMultipleReaction(botMessageID string, session *discordgo.Session, msgEvent *discordgo.MessageCreate)(error){
//	reactions := []string{"⏮️", "◀️", "▶️", "⏭️"}
//	for _, reaction := range reactions{
//		// remove reactions message
//		err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, reaction, msgEvent.Author.ID)
//		if err != nil {
//			log.Error("Failed to remove user reaction from bot message:", err)
//			return err
//		}
//	}
//	return nil
//}

// checks user reaction selection.
func checkUserReactionSelect(page int, currentTime time.Time, botMessageID string, session *discordgo.Session, msgEvent *discordgo.MessageCreate) (int, error) {
	errorVal := 10
	for {
		timePassed := time.Since(currentTime)
		if timePassed.Seconds() >= 30 {
			log.WithFields(log.Fields{
				"Time passed": timePassed,
			}).Info("Removing reactions time has been passed.")
			err := session.MessageReactionsRemoveAll(msgEvent.ChannelID, botMessageID)
			if err != nil {
			}
			previousAuthor = ""
			return errorVal, err
		}
		// check if the reaction matches the author ID aka sender
		checkReaction, err := checkMessageReaction(session, msgEvent, botMessageID)
		if err != nil {
			log.Error("Failed to check emoji from bot message:", err)
			return errorVal, err
		}
		// if true then remove it uwu.
		if checkReaction["Stop"] {
			err = session.MessageReactionsRemoveAll(msgEvent.ChannelID, botMessageID)
			if err != nil {
				log.Error("Failed to remove all reaction from bot message stop reaction.", err)
				return 0, err
			}
			previousAuthor = ""
			return errorVal, err
		} else if checkReaction["FastBack"] {
			page = 0
			// remove user reactions before going to next page.
			err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "⏮️", msgEvent.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message:", err)
				return errorVal, err
			}
			return page, nil
		} else if checkReaction["Back"] {
			if page == 1 {
				// remove user reactions before going to next page.
				err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "◀️", msgEvent.Author.ID)
				if err != nil {
					log.Error("Failed to remove user reaction from bot message:", err)
					return errorVal, err
				}
				return 0, nil
			} else if page == 2 {
				// remove user reactions before going to next page.
				err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "◀️", msgEvent.Author.ID)
				if err != nil {
					log.Error("Failed to remove user reaction from bot message:", err)
					return errorVal, err
				}
				return 0, nil
			}
			// remove user reactions before going to next page.
			err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "◀️", msgEvent.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message:", err)
				return errorVal, err
			}
			page--
			return page, nil
		} else if checkReaction["Forward"] {
			if page == 5 {
				// remove user reactions before going to next page.
				err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "▶️", msgEvent.Author.ID)
				if err != nil {
					log.Error("Failed to remove user reaction from bot message:", err)
					return errorVal, err
				}
				return 0, nil
			}
			// remove user reactions before going to next page.
			err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "▶️", msgEvent.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message:", err)
				return errorVal, err
			}
			page++
			return page, nil
		} else if checkReaction["FastForward"] {
			// remove user reactions before going to next page.
			err := session.MessageReactionRemove(msgEvent.ChannelID, botMessageID, "⏭️", msgEvent.Author.ID)
			if err != nil {
				log.Error("Failed to remove user reaction from bot message:", err)
				return errorVal, err
			}
			page = 5
			return page, nil
		}
	}
}

func getCpuUsage() (string, error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Error("Failed to read stat possibly due not finding /proc/stat: ", err)
		return "", err
	}

	cpuStatSystemFresh := stat.CPUStatAll.System
	cpuStatSystemNotFresh := stat.CPUStatAll.System
	difference := cpuStatSystemFresh - cpuStatSystemNotFresh
	percentage := difference / (uint64(1*time.Second) * 100)

	fmt.Println("cpuUsage:", cpuStatSystemFresh)
	convertToString := strconv.FormatUint(percentage, 10)

	return convertToString + "%", nil
}

func getMemInfo() (string, error) {
	memInfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Error("Failed to Read memory info: ", err)
		return "", err
	}
	memoryFree := memInfo.MemFree
	memoryUsed := memInfo.Active

	memUsagePercentage := (memoryUsed / memoryFree) * 100
	fmt.Println("memUsage:", memoryUsed)

	convertToString := strconv.FormatUint(memUsagePercentage, 10)

	return convertToString + "%", nil
}

func getTimeLeftForNick(s *discordgo.Session, m *discordgo.MessageCreate, message string) error {
	// get guild info from DB
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return err
	}

	// guild member used to retrieve username
	guildMember, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Error("Failed to get member details: ", err)
		return err
	}

	userDB, err := db.FindUserByID(m.GuildID, m.Author.ID)
	if err != nil {
		log.Error("Failed to get user: ", err)
		return err
	}

	// calculate how long left and reset the duration if it's up.
	// seconds for the time since last update changes every month or whatever the owner or admin set the nick change days
	userLastNickUpdate := time.Since(userDB.Date).Seconds()

	// convertStringToInt
	guildDurationToFloat, err := strconv.ParseFloat(userDB.Guild.GuildNicknameDuration, 10)
	if err != nil {
		log.Error("Failed to convert GuildNickname duration to int: ", err)
		return err
	}
	// dynamic guild duration
	guildNickDaysDurationToSeconds := guildDurationToFloat * 86400

	// get the difference.
	remainingSeconds := guildNickDaysDurationToSeconds - userLastNickUpdate
	// convert seconds to clock times
	secondsToDays := remainingSeconds / 86400
	secondsToHours := remainingSeconds / 3600
	secondsToMinutes := remainingSeconds / 60

	// change to readable format time reminders.
	days := int(secondsToDays)
	hours := int(secondsToHours) % 24
	minutes := int(secondsToMinutes) % 60
	seconds := int(remainingSeconds) % 60

	// if the seconds is greater than the duration seconds set by the guild then return
	//and let them know they can change their nick.
	// updates the allowedNickChange to True if it's full filled
	if userLastNickUpdate >= guildNickDaysDurationToSeconds {
		updateUserDB := model.User{
			UserID:            m.Author.ID,
			Guild:             guild,
			NickName:          guildMember.Nick,
			Date:              userDB.Date,
			OldNickNames:      userDB.OldNickNames,
			AllowedNickChange: true,
			TimeStamp:         time.Now(),
		}
		err := db.InsertOrUpdateUser(guild, &updateUserDB)
		if err != nil {
			log.Error("Failed to Update user: ", err)
			return err
		}
		return err
	}

	// let them know when they can reset their nickname.
	embed := NewEmbed().
		SetDescription(message + m.Author.Username +
			fmt.Sprintf(" you can change your nickname in `%d%s %d%s %d%s %d%s`.",
				days, "d",
				hours, "h",
				minutes, "m",
				seconds, "s")).
		SetColor(green).MessageEmbed
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Error("On sending parameter error message to channel: ", err)
		return err

	}
	return nil
}
