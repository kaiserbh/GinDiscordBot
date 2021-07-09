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

//helpMessageHandler help menu
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

						// add reaction to the message author
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
								AddField("Invite", fmt.Sprintf("[Invite %s](https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8)", s.State.User.Username, s.State.User.ID)).
								AddField("Support Server", "[Gin Support](https://discord.gg/nkGvkUUqHZ)").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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
								AddField("reset", "resets nickname (doesn't reset duration)").
								AddField("invite", "Get a link to invite me.").
								AddField("support", "Get a link to my support server.").
								AddField("source", "Get the link to Gin's GitHub repository.").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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
								AddField("anime", "Query anime from Anilist").
								AddField("manga", "Query manga from Anilist").
								AddField("character", "Query character from Anilist").
								AddField("staff", "Query person/staff from Anilist").
								AddField("studio", "Query studio from Anilist").
								AddField("user", "Query user from Anilist").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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
								SetTitle("Miscellaneous").
								SetThumbnail(botImage).
								SetDescription(fmt.Sprintf("My default prefix is `%[1]s`. Use `%[1]shelp <command>` to get more information on a command.", guild.GuildPrefix)).
								AddField("permissions", "Show your permissions or the member specified.").
								AddField("userinfo", "Show some information about yourself or the member specified.").
								AddField("serverinfo", "Get some information about this server.").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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
								AddField("Invite", fmt.Sprintf("[Invite %s](https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8)", s.State.User.Username, s.State.User.ID)).
								AddField("Support Server", "[Gin Support](https://discord.gg/nkGvkUUqHZ)").
								SetFooter("Use reactions to flip pages (Page " + strconv.Itoa(page) + "/5)").
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

// pingMessageHandler pings the bot
func pingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

				memUsage, err := getMemInfo()
				if err != nil {
					return
				}

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
					AddField("RAM Usage", memUsage).
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
		if parameter[0] == guild.GuildPrefix+"nick" {
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

					user, err = db.FindUserByID(m.GuildID, m.Author.ID)
					if err != nil {
						log.Error("Failed to get user: ", err)
						return
					}

					// If user can change their nickname.
					allowedNickChange := user.AllowedNickChange
					if allowedNickChange {
						// do nothing if the user didn't provide arguments for nickname
						if len(parameter) < 2 {
							log.Info("Doing nothing since no arguments was provided")
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
						err = getTimeLeftForNick(s, m, "Successfully changed nickname for this server. \n")
						if err != nil {
							log.Error("Failed to get time left for nick change")
							return
						}

						// gets bot Message ID
						botMessageID, err := getBotMessageID(s, m)
						if err != nil {
							log.Error("Failed to get botID")
							return
						}

						// for loop to check time passed before deleting user message and bot message.
						for {
							since := time.Since(timerToRemoveBotMessageAndUser).Seconds()

							if since >= 5 {
								err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								err = s.ChannelMessageDelete(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								return
							}
						}
					} else {
						err = getTimeLeftForNick(s, m, "")
						if err != nil {
							log.Error("Failed to get time left for nick change")
							return
						}
						// gets bot Message ID
						botMessageID, err := getBotMessageID(s, m)
						if err != nil {
							log.Error("Failed to get botID")
							return
						}
						// for loop to check time passed before deleting user message and bot message.
						for {
							since := time.Since(timerToRemoveBotMessageAndUser).Seconds()
							if since >= 5 {
								err = s.ChannelMessageDelete(m.ChannelID, lastMessage)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								err = s.ChannelMessageDelete(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to delete user message: ", err)
									return
								}
								return
							}
						}
					}
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
				err = s.MessageReactionAdd(m.ChannelID, lastMessage, "✅")
				if err != nil {
					log.Error("Failed to add reaction: ", err)
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

// invite fetches bot invite
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
			if m.Author.ID == s.State.User.ID {
				return
			}
			if messageContent == guild.GuildPrefix+"invite" {
				// start embed
				embed := NewEmbed().
					SetTitle(fmt.Sprintf("Invite %s", s.State.User.Username)).
					SetDescription(fmt.Sprintf("[link](https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot)", s.State.User.ID)).
					MessageEmbed

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

// Gaki command disabled
/*
func gaki(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Checks if the message has prefix from the database file.

	messageContent := strings.ToLower(m.Content)

	// check if the channel is bot channel or allowed channel.

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(messageContent, "gaki") {
		// start embed
		embed := NewEmbed().
			SetDescription("Should go back to his dungeon").
			SetImage("https://media.discordapp.net/attachments/703063888241098832/745009187053895831/8CFoRfV0IXScYAAAAASUVORK5CYII.png?width=400&height=226").
			MessageEmbed
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)

		if err != nil {
			log.Error("Failed to send embed to the channel: ", err)
			return
		}

		// add reaction to the message author
		err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, ":smirk~1:862978313655156766")
		if err != nil {
			log.Error("Failed to add reaction: ", err)
			return
		}
	}
}
*/
/*
func test(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Checks if the message has prefix from the database file.

	messageContent := strings.ToLower(m.Content)

	// check if the channel is bot channel or allowed channel.

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(messageContent, "test") {
		// start embed
		_, err := s.ChannelMessageSend(m.ChannelID, "test")

		if err != nil {
			log.Error("Failed to send embed to the channel: ", err)
			return
		}

		// add reaction to the message author
		err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, ":smirk~1:862978313655156766")
		if err != nil {
			log.Error("Failed to add reaction: ", err)
			return
		}
	}
}
*/
