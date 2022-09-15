package bot

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
)

// helpMessageHandler help menu
func helpMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		reactions := []string{"⏮️", "◀️", "⏹️", "▶️", "⏭️"}

		page := 1

		// check the message if it's from the bot if it is ignore.
		if m.Author.ID == s.State.User.ID {
			return
		}
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if strings.HasPrefix(m.Content, guild.GuildPrefix) {
				if messageContent == guild.GuildPrefix+"help" {
					// check if the previous instance is still running.
					if m.Author.ID == previousAuthor {
						embed := NewEmbed().
							SetDescription("hmm, make sure you end the last instance of help menu before executing another one MADAO...").
							SetColor(red).MessageEmbed

						_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					// bot messageID
					var botMessageID string
					var botImage = s.State.User.AvatarURL("")
					//var ok bool
					for {
						if page == 10 {
							break
						}
						switch page {
						// page one About page
						case 1:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Gin Help Menu").
								SetThumbnail(botImage).
								SetDescription("Gin is a feature rich Discord bot designed to bring FUN into your server or one would hope so...").
								AddField("Invite", fmt.Sprintf("[Invite %s](https://discord.com/oauth2/authorize?"+
									"client_id=854839186287493151&permissions=4228906231&scope=bot)", s.State.User.Username)).
								AddField("Support Server", "[Gin Support](https://discord.gg/SD2D6Y8RaC)").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err = getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}
							// add reaction to the bot message with for loop?
							for _, emoji := range reactions {
								err = s.MessageReactionAdd(m.ChannelID, botMessageID, emoji)
								if err != nil {
									log.Error("Failed to add reaction: ", err)
									return
								}
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						case 2:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Configuration").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("prefix", "Change the prefix or view the current prefix.").
								AddField("botchannel", "sets the current channel as bot channel or set multiple channel as bot channel.").
								AddField("cooldown", "set duration for nickname changes in days").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						case 3:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("General").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("help", "Display help menu").
								AddField("ping", "Pong! Get my latency.").
								AddField("stats", "See some super cool statistics about me.").
								AddField("nick", "Change nickname").
								AddField("nickdur", "check other @users nick duration").
								AddField("reset", "resets nickname (doesn't reset duration)").
								AddField("invite", "Get a link to invite me.").
								AddField("support", "Get a link to my support server.").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						case 4:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Anilist").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("anime | a", "Query anime from Anilist").
								AddField("manga | m", "Query manga from Anilist").
								AddField("character | c", "Query character from Anilist").
								AddField("staff | s", "Query person/staff from Anilist").
								AddField("user | u", "Query user from Anilist").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						case 5:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("MyAnimeList").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("anime_mal | am", "Query anime from MAL").
								AddField("manga_mal | mm", "Query manga from Anilist").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						case 6:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Miscellaneous").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("permissions", "Show your permissions or the member specified.").
								AddField("userinfo", "Show some information about yourself or the member specified.").
								AddField("serverinfo", "Get some information about this server.").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed
							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}

							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						default:
							// reset page
							page = 1

							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Gin Help Menu").
								SetThumbnail(botImage).
								SetDescription("Gin is a feature rich Discord bot designed to bring FUN into your server or one would hope so...").
								AddField("Invite", fmt.Sprintf("[Invite %s](https://discord.com/oauth2/authorize?"+
									"client_id=854839186287493151&permissions=4228906231&scope=bot)", s.State.User.Username)).
								AddField("Support Server", "[Gin Support](https://discord.gg/nkGvkUUqHZ)").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/6)").
								SetColor(green).MessageEmbed

							// add reaction to the message author
							_, err = s.ChannelMessageEditEmbed(m.ChannelID, botMessageID, embed)
							if err != nil {
								log.Error("Failed to send embed to the channel: ", err)
								return
							}
							// gets bot Message ID
							botMessageID, err := getBotMessageID(s, m)
							if err != nil {
								log.Error("Failed to get botID")
								return
							}
							// execute checkUserReactionSelect basically while loop that checks or waits for user reaction
							page, err = checkUserReactionSelect(page, currentTime, botMessageID, s, m)
							if err != nil {
								log.Error("Failed to check user select Reaction: ", err)
								return
							}
						}
					}
				}
			}
		}
	}
}

// pingLatency pings the bot to get latency to discord server.
func pingLatency(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if messageContent == guild.GuildPrefix+"ping" {
				// start embed
				embed := NewEmbed().
					SetDescription("pong!").
					SetColor(green).MessageEmbed

				// add reaction to the message author
				lastMessage := m.Message.ID
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					log.Error("Failed to send embed to the channel: ", err)
				}
				err = s.MessageReactionAdd(m.ChannelID, lastMessage, "🏓")
				if err != nil {
					log.Error("Failed to add reaction: ", err)
				}
			}
		}
	}
}

// stats get bot statics.
func stats(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}

	guildsInDB, err := db.GetAllGuild()
	if err != nil {
		log.Error("Failed to get all Guild from DB: ", err)
		return
	}

	if strings.HasPrefix(m.Content, guild.GuildPrefix) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			messageContent := m.Content
			if strings.Contains(messageContent, guild.GuildPrefix+"stats") {
				uptime := time.Since(Uptime).Seconds()
				var botImage = s.State.User.AvatarURL("")
				goVer := runtime.Version()
				numberOfGuildIn := len(guildsInDB)

				// uptime divide to readable format.
				uptimeSeconds := int(uptime) % 60
				uptimeMinutes := int(uptime) / 60
				uptimeHours := uptimeMinutes / 60
				uptimeDays := uptimeHours / 24

				uptimeDaysReminder := uptimeDays % 24
				uptimeHoursReminder := uptimeHours % 24
				uptimeMinutesReminder := uptimeMinutes % 60

				cpuUsage, err := getCpuUsage()
				if err != nil {
					log.Error("Failed to get CPU Usage: ", err)
				}

				//memUsage, err := getMemInfo()
				//if err != nil {
				//	return
				//}

				// start embed
				embed := NewEmbed().
					SetTitle("Gin Statistic").
					SetThumbnail(botImage).
					AddField("Owner", "Kaiser#0101 \n").
					AddField("Bot Version", "0.1").
					AddField("Uptime", fmt.Sprintf("%d%s %d%s %d%s %d%s",
						uptimeDaysReminder, "d",
						uptimeHoursReminder, "h",
						uptimeMinutesReminder, "m",
						uptimeSeconds, "s")).
					AddField("Servers", fmt.Sprintf("%d", numberOfGuildIn)).
					AddField("CPU Usage", cpuUsage).
					//AddField("RAM Usage", memUsage).
					AddField("Go Version", fmt.Sprintf("%v", goVer)).
					InlineAllFields().
					SetColor(green).MessageEmbed

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					log.Warn("Failed to send embed to the channel: ", err)
					return
				}
				return
			}
		}
	}
}

// setNick changes nickname and gives back duration on how long left until they can change it again.
func setNick(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		// reset message content to normal characters and not lower it
		messageContent = m.Content

		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		parameter := getArguments(messageContent)

		// user message ID
		lastMessage := m.Message.ID

		// check if it's nick or nickname since contains func will return both nickname and nick function.
		if strings.ToLower(parameter[0]) == guild.GuildPrefix+"nick" {

			// check if allowed channel.
			if allowedChannels {
				if strings.Contains(messageContent, guild.GuildPrefix+"nick") {
					timerToRemoveBotMessageAndUser := time.Now()

					// If user can change their nickname (basically if the user is not in DB)
					user, err := db.FindUserByID(m.GuildID, m.Author.ID)
					if err != nil {
						log.Error("Failed to get user: ", err)
						// guild member used to retrieve username
						guildMember, err := s.GuildMember(m.GuildID, m.Author.ID)
						if err != nil {
							log.Error("Failed to get member details: ", err)
							return
						}

						updateUserDB := model.User{
							UserID:            m.Author.ID,
							Guild:             guild,
							NickName:          guildMember.Nick,
							OldNickNames:      []string{guildMember.Nick},
							Date:              time.Now(),
							AllowedNickChange: true,
							TimeStamp:         time.Now(),
						}
						err = db.InsertOrUpdateUser(guild, &updateUserDB)
						if err != nil {
							log.Error("Failed to Update user: ", err)
							return
						}
					}

					// check if the date has already passed the duration set
					oldDate := user.Date
					timePassed := time.Since(oldDate).Seconds()
					guildNickDuration, err := strconv.Atoi(guild.GuildNicknameDuration)
					if err != nil {
						log.Error("Failed to convert guild Nick duration to int: ", err)
						return
					}
					guildNickDuration = guildNickDuration * 86400

					if timePassed >= float64(guildNickDuration) {
						// guild member used to retrieve username
						guildMember, err := s.GuildMember(m.GuildID, m.Author.ID)
						if err != nil {
							log.Error("Failed to get member details: ", err)
							return
						}

						updateUserDB := model.User{
							UserID:            m.Author.ID,
							Guild:             guild,
							NickName:          guildMember.Nick,
							Date:              time.Now(),
							AllowedNickChange: true,
							TimeStamp:         time.Now(),
						}
						err = db.InsertOrUpdateUser(guild, &updateUserDB)
						if err != nil {
							log.Error("Failed to Update user: ", err)
							return
						}
					}

					// retreive the updated change.
					user, _ = db.FindUserByID(m.GuildID, m.Author.ID)

					// If user can change their nickname.
					allowedNickChange := user.AllowedNickChange
					if allowedNickChange {
						// do nothing if the user didn't provide arguments for nickname
						if len(parameter) < 2 {
							timer := time.Now()
							if user.AllowedNickChange {
								embed := NewEmbed().
									SetDescription("You can change your nickname. use `" + guild.GuildPrefix + "nick <desired nickname>` to change it.").
									SetColor(green).MessageEmbed
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
								if err != nil {
									log.Warn("Failed to send embed to the channel: ", err)
									return
								}

								for {
									if time.Since(timer).Seconds() > 5 {
										// bot messageID
										botMessageID, err := getBotMessageID(s, m)
										if err != nil {
											log.Error("Failed to get bot message ID: ", err)
											return
										}
										// delete user message and bot messages.
										err = s.ChannelMessageDelete(m.ChannelID, botMessageID)
										if err != nil {
											log.Error("Failed to remove bot message: ", err)
											return
										}
										err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
										if err != nil {
											log.Error("Failed to remove user message: ", err)
											return
										}
										return
									}
								}
							}
							return
						}
						// get the spaces as well.
						nickname := strings.Join(parameter[1:], " ")
						// discord doesn't allow nickname more than 32 char
						if len(nickname) > 32 {
							embed := NewEmbed().
								SetDescription("Discord does not allow more than 32 characters.").
								SetColor(red).MessageEmbed
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
							if err != nil {
								log.Warn("Failed to send embed to the channel: ", err)
								return
							}
							return
						}

						// add escape characters for \
						if strings.Contains(nickname, "\\") {
							nickname = strings.ReplaceAll(nickname, "\\", "\\")
						}

						// update member nickname
						err = s.GuildMemberNickname(m.GuildID, m.Author.ID, nickname)
						if err != nil {
							log.Error("Failed to change user nickname: ", err)
							embed := NewEmbed().
								SetDescription("Sorry it seem like I do not have the permission to do that.").
								SetColor(red).MessageEmbed
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
							if err != nil {
								log.Warn("Failed to send embed to the channel: ", err)
								return
							}
							return
						}

						// Get old nickname from DB and append it.
						oldNick := user.OldNickNames
						oldNick = append(oldNick, nickname)

						// update DB
						newNickUserDB := model.User{
							UserID:            m.Author.ID,
							Guild:             guild,
							NickName:          nickname,
							OldNickNames:      oldNick,
							Date:              time.Now(),
							AllowedNickChange: false,
							TimeStamp:         time.Now(),
						}

						err = db.InsertOrUpdateUser(guild, &newNickUserDB)
						if err != nil {
							log.Error("Failed to update user in DB: ", err)
							return
						}

						// get how long time left
						err = getTimeLeftForNick(s, m.Author.ID, m.GuildID, m.ChannelID, "Successfully changed nickname for this server. \n"+m.Author.Username)
						if err != nil {
							log.Error("Failed to get time left for nick change: ", err)
							return
						}

						// for loop to check time passed before deleting user message and bot message.
						for {
							since := time.Since(timerToRemoveBotMessageAndUser).Seconds()

							// bot messageIDs
							botMessageIDs, err := getAllBotMessagesID(s, m)
							if err != nil {
								log.Error("Failed to get bot message IDs")
								return
							}

							if since >= 10 {
								log.WithFields(log.Fields{
									"time": since,
								}).Info("Removing bot message and user since time has been passed")
								err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								for _, val := range botMessageIDs {
									err = s.ChannelMessageDelete(m.ChannelID, val)
									if err != nil {
										log.Error("Failed to delete user message: ", err)
										return
									}
								}
								return
							}
						}
					} else {

						// Check if the old date it's passed if it's then change the nickname.

						err = getTimeLeftForNick(s, m.Author.ID, m.GuildID, m.ChannelID, ""+m.Author.Username)
						if err != nil {
							log.Error("Failed to get time left for nick change: ", err)
							return
						}

						// for loop to check time passed before deleting user message and bot message.
						for {
							since := time.Since(timerToRemoveBotMessageAndUser).Seconds()
							// bot messageIDs
							botMessageIDs, err := getAllBotMessagesID(s, m)
							if err != nil {
								log.Error("Failed to get bot message IDs: ", err)
								return
							}

							if since >= 10 {
								log.WithFields(log.Fields{
									"time": since,
								}).Info("Removing bot message and user since time has been passed")
								err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								for _, val := range botMessageIDs {
									err = s.ChannelMessageDelete(m.ChannelID, val)
									if err != nil {
										log.Error("Failed to delete user message: ", err)
										return
									}
								}
								return
							}
						}
					}
				}
			}
			// check other user nickname duration.
		} else if strings.ToLower(parameter[0]) == guild.GuildPrefix+"nickdur" {
			timerToRemoveBotMessageAndUser := time.Now()
			// check param if it's 2 or not.
			if len(parameter) != 2 {
				log.Info("Doing nothing since no params or user have been given or probably more params than needed")
				return
			}
			user := parameter[1]
			// clean param and get the user ID
			cleanUserID := strings.TrimPrefix(strings.TrimSuffix(user, ">"), "<@!")
			err = getTimeLeftForNick(s, cleanUserID, m.GuildID, m.ChannelID, "")
			if err != nil {
				log.Error("Failed to get time left for nick change: ", err)
				return
			}

			// for loop to check time passed before deleting user message and bot message.
			for {
				since := time.Since(timerToRemoveBotMessageAndUser).Seconds()
				// bot messageIDs
				botMessageIDs, err := getAllBotMessagesID(s, m)
				if err != nil {
					log.Error("Failed to get bot message IDs: ", err)
					return
				}

				if since >= 15 {
					log.WithFields(log.Fields{
						"time": since,
					}).Info("Removing bot message and user since time has been passed")
					err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
					if err != nil {
						log.Error("Failed to delete user message: ", err)
						return
					}
					for _, val := range botMessageIDs {
						err = s.ChannelMessageDelete(m.ChannelID, val)
						if err != nil {
							log.Error("Failed to delete user message: ", err)
							return
						}
					}
					return
				}
			}
		}
	}
}

// resetNickHandler used to reset nickname to default value, does not change duration..
func resetNickHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if messageContent == guild.GuildPrefix+"reset" {
				err := s.GuildMemberNickname(m.GuildID, m.Author.ID, " ")
				if err != nil {
					log.Error("Failed to reset nickname: ", err)
					return
				}
				// add reaction to the message author
				lastMessage := m.Message.ID
				err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
				if err != nil {
					log.Error("Failed to remove user message: ", err)
					return
				}
			}
		}
	}
}

// botPing handles message that include pings to the bot
func botPing(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)

	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}

	messageContent := strings.ToLower(m.Content)

	// check if the channel is bot channel or allowed channel.
	allowedChannels := checkAllowedChannel(m.ChannelID, guild)

	if allowedChannels {

		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.Contains(messageContent, "<@!"+s.State.User.ID+">") {
			// start embed

			embed := NewEmbed().
				SetDescription("The prefix for this server is " + "`" + guild.GuildPrefix + "`.").
				SetColor(green).MessageEmbed

			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Error("Failed to send embed to the channel: ", err)
				return
			}

			// add reaction to the message author
			err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, "👋")
			if err != nil {
				log.Error("Failed to add reaction: ", err)
				return
			}
		}
	}
}

// invite sends invite link for the bot
func invite(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			// if the message is from the bot
			if m.Author.ID == s.State.User.ID {
				return
			}
			if messageContent == guild.GuildPrefix+"invite" {
				// start embed
				embed := NewEmbed().
					SetTitle("Gin invite link").
					SetURL("https://discord.com/api/oauth2/authorize?client_id=" + s.State.User.ID +
						"&permissions=4228906231&scope=bot").
					SetColor(green).MessageEmbed

				// add reaction to the message author

				_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					log.Error("Failed to send embed to the channel: ", err)
					return
				}

			}
		}
	}
}

// // authorizeAnilist used to authorize user accounts to anilist.
// func authorizeAnilist(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	// Checks if the message has prefix from the database file.
// 	guild, err := db.FindGuildByID(m.GuildID)
// 	if err != nil {
// 		log.Error("Finding Guild: ", err)
// 		return
// 	}
// 	messageContent := strings.ToLower(m.Content)
// 	anilist_authorization_url := "https://anilist.co/api/v2/oauth/authorize?client_id=6823&redirect_uri=https://anilist.co/api/v2/oauth/pin&response_type=code"
// 	//TODO: Need to be done not yet finished
// 	//TODO: Check DM if it's enabled
// 	//TODO: Send message to user DM
// 	//TODO: paste the code inside DM
// 	//TODO: Update DB to store authorizeation code
// 	//TODO: get anilist username and add it to DB
// 	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
// 		// check if the channel is bot channel or allowed channel.
// 		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
// 		if allowedChannels {
// 			if m.Author.ID == s.State.User.ID {
// 				return
// 			}
// 			if messageContent == guild.GuildPrefix+"authorizeAnilist" {

// 			}
// 		}
// 	}
// }
