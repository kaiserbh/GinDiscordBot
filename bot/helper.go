package bot

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kaiserbh/anilistgo"

	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
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

// check if it's guild owner
func checkGuildOwner(session *discordgo.Session, msgEvent *discordgo.MessageCreate) (bool, error) {
	guild, err := session.Guild(msgEvent.GuildID)
	if err != nil {
		return false, err

	}
	authorID := msgEvent.Author.ID
	ownerID := guild.OwnerID

	if authorID == ownerID {
		return true, nil
	}

	return false, nil
}

// check latency or how long a function takes to execute
// func measureTime(funcName string) func() {
// 	start := time.Now()
// 	return func() {
// 		log.WithFields(log.Fields{
// 			"funcName": funcName,
// 			"time":     time.Since(start),
// 		}).Info("Time taken by function is completed")
// 	}
// }

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
		log.Error("Failed to get messages from channel:", err)
		return "", err
	}

	if len(channelMessages) <= 0 {
		return "", err
	}

	botMessageID := channelMessages[0].ID
	return botMessageID, nil
}

// getAllBotMessagesID returns the last 100 messages from the bot.
func getAllBotMessagesID(session *discordgo.Session, msgEvent *discordgo.MessageCreate) ([]string, error) {

	var botMessagesID []string

	// check the message if it's from the bot if it is ignore.
	channelMessages, err := session.ChannelMessages(msgEvent.ChannelID, 100, "", "", "")
	if err != nil {
		log.Error("Failed to get messages from channel")
		return []string{}, err
	}

	for _, message := range channelMessages {
		if message.Author.ID == session.State.User.ID {
			botMessagesID = append(botMessagesID, message.ID)
		}
	}

	return botMessagesID, nil
}

// checkMessageReactionAuthor check reactions for help menu
func checkMessageReactionAuthor(session *discordgo.Session, channelID, botMessageID, emojiID, authorID string, limit int) (bool, error) {
	checkReaction, err := session.MessageReactions(channelID, botMessageID, emojiID, limit, botMessageID, "")
	if err != nil {
		log.Error("Failed to get message reactions: ", err)
		return false, err
	}
	for _, reactions := range checkReaction {
		if reactions.ID == authorID {
			return true, nil
		}
	}
	return false, nil
}

// checkHelpMenuReactions check reactions for help menu
func checkHelpMenuReactions(session *discordgo.Session, msgEvent *discordgo.MessageCreate, botMessageID string) (map[string]bool, error) {

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
				log.Error("Failed to remove embeds", err)
				return 0, err
			}
			previousAuthor = ""
			return errorVal, err
		}
		// check if the reaction matches the author ID aka sender
		checkReaction, err := checkHelpMenuReactions(session, msgEvent, botMessageID)
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
			if page == 6 {
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
			page = 6
			return page, nil
		}
	}
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Error("Failed to read file /proc/stat: ", err)
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					log.WithFields(log.Fields{
						"index":  i,
						"fields": fields[i],
					}).Error("Error parsing: ", err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func getCpuUsage() (string, error) {
	idle0, total0 := getCPUSample()
	time.Sleep(1 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

	convertToString := strconv.FormatFloat(cpuUsage, 'f', 2, 64)

	return convertToString + "%", nil
}

//func getMemInfo() (string, error) {
//	memInfo, err := linux.ReadMemInfo("/proc/meminfo")
//	if err != nil {
//		log.Error("Failed to Read memory info: ", err)
//		return "", err
//	}
//
//	// it's in kB (Kilo Byte) convert to Gigabyte
//	memoryAvailable := memInfo.MemAvailable
//	memoryFree := memInfo.MemFree
//
//	// memory used and convert to MiB
//	memoryUsed := float64(memoryAvailable-memoryFree) / float64(memoryAvailable) * 100
//
//	convertToString := strconv.FormatFloat(memoryUsed, 'f', 2, 64)
//
//	return convertToString + "%", nil
//}

func getTimeLeftForNick(s *discordgo.Session, authorID, guildID, channelID string, message string) error {
	// get guild info from DB
	guild, err := db.FindGuildByID(guildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return err
	}

	// guild member used to retrieve username
	guildMember, err := s.GuildMember(guildID, authorID)
	if err != nil {
		log.Error("Failed to get member details: ", err)
		return err
	}

	userDB, err := db.FindUserByID(guildID, authorID)
	if err != nil {
		log.Error("Failed to get user: ", err)
		// let them know when they can reset their nickname.
		embed := NewEmbed().
			SetDescription(message + "Can change nickname").
			SetColor(green).MessageEmbed
		_, err = s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Error("On sending parameter error message to channel: ", err)
			return err
		}
		return err
	}

	// calculate how long left and reset the duration if it's up.
	// seconds for the time since last update changes every month or whatever the owner or admin set the nick change days
	userLastNickUpdate := time.Since(userDB.Date).Seconds()

	// convertStringToInt
	guildDurationToFloat, err := strconv.ParseFloat(userDB.Guild.GuildNicknameDuration, 64)
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
	// and let them know they can change their nick.
	// updates the allowedNickChange to True if it's full filled
	if userLastNickUpdate >= guildNickDaysDurationToSeconds {
		updateUserDB := model.User{
			UserID:            authorID,
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
	}

	// if user can't change nickname check if they can change nick.
	if !userDB.AllowedNickChange {
		// let them know when they can reset their nickname.
		embed := NewEmbed().
			SetDescription(message +
				fmt.Sprintf(" Can change nickname in `%d%s %d%s %d%s %d%s`.",
					days, "d",
					hours, "h",
					minutes, "m",
					seconds, "s")).
			SetColor(green).MessageEmbed
		_, err = s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Error("On sending parameter error message to channel: ", err)
			return err
		}
		return nil
	} else {
		// let them know when they can reset their nickname.
		embed := NewEmbed().
			SetDescription(message + "Can change nickname").
			SetColor(green).MessageEmbed
		_, err = s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Error("On sending parameter error message to channel: ", err)
			return err
		}
	}
	return nil
}

func convertStringHexColorToInt(data string) (int, error) {
	animeColor := strings.Replace(data, "#", "", -1)
	animeColorHex, err := strconv.ParseInt(animeColor, 16, 64)
	if err != nil {
		log.Error("Failed to convert string to int: ", err)
		return 0, err
	}

	return int(animeColorHex), nil
}

func convMonthIntToStr(year string) string {
	switch year {
	case "1":
		year = "Jan"

	case "2":
		year = "Feb"

	case "3":
		year = "Mar"

	case "4":
		year = "Apr"

	case "5":
		year = "May"

	case "6":
		year = "Jun"

	case "7":
		year = "Jul"

	case "8":
		year = "Aug"

	case "9":
		year = "Sep"

	case "10":
		year = "Oct"

	case "11":
		year = "Nov"

	case "12":
		year = "Dec"
	}
	return year
}

func checkAnilistTimer(session *discordgo.Session, channelID, botMessageID, authorID string) {
	// add reaction to bot message.
	err := session.MessageReactionAdd(channelID, botMessageID, "✅")
	if err != nil {
		log.Error("Failed to add reaction: ", err)
		return
	}
	err = session.MessageReactionAdd(channelID, botMessageID, "❌")
	if err != nil {
		log.Error("Failed to add reaction: ", err)
		return
	}

	startTimer := time.Now()

	for {
		passedTimer := time.Since(startTimer).Seconds()
		checkAuthorReactionOk, err := checkMessageReactionAuthor(session, channelID, botMessageID, "✅", authorID, 10)
		if err != nil {
			log.Error(err)
			return
		}
		if checkAuthorReactionOk {
			err = session.MessageReactionsRemoveAll(channelID, botMessageID)
			if err != nil {
				log.Error("Failed to remove reactions from bot message: ", err)
				return
			}
			return
		}

		// check the delete reaction
		checkAuthorReactionDelete, err := checkMessageReactionAuthor(session, channelID, botMessageID, "❌", authorID, 10)
		if err != nil {
			log.Error("Failed to check author author reaction: ", err)
			return
		}

		if checkAuthorReactionDelete {
			err := session.ChannelMessageDelete(channelID, botMessageID)
			if err != nil {
				log.Error("Failed to delete botMessage: ", err)
				return
			}
			return
		}
		// if no reactions is added then just remove reactions from the message.
		if passedTimer >= 30 {
			err = session.MessageReactionsRemoveAll(channelID, botMessageID)
			if err != nil {
				log.Error("Failed to remove reactions from bot message: ", err)
				return
			}
			return
		}
	}
}

func anilistAnimeData(media *anilistgo.Media) (string, string, string) {
	descriptionCut := cutDescription(media.Description)

	// start date
	animeStartMonth := strconv.Itoa(media.StartDate.Month)
	animeStartDay := strconv.Itoa(media.StartDate.Day) + ","
	animeStartYear := strconv.Itoa(media.StartDate.Year)
	animeStartMonthString := convMonthIntToStr(animeStartMonth) + " "
	startDate := animeStartMonthString + animeStartDay + animeStartYear

	// end date
	animeEndMonth := strconv.Itoa(media.EndDate.Month)
	animeEndDay := strconv.Itoa(media.EndDate.Day) + ","
	animeEndYear := strconv.Itoa(media.EndDate.Year)
	animeEndMonthString := convMonthIntToStr(animeEndMonth) + " "
	endDate := animeEndMonthString + animeEndDay + animeEndYear

	return descriptionCut, startDate, endDate
}

func cutDescription(description string) string {
	if len(description) > 200 {
		description = description[:200]
		description = description + "..."
	}

	// replace with spoiler tag.
	var re = regexp.MustCompile(`(?m)[!~]`)
	description = re.ReplaceAllString(description, "|")

	// replace <br> with new line
	description = strings.Replace(description, "<br>", "\n", -1)

	return description
}
