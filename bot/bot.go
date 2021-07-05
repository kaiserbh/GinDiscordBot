package bot

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/config"
	"github.com/kaiserbh/gin-bot-go/database"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
)

var db = database.Connect()
var (
	red            = 0xff0000
	green          = 0x11ff00
	previousAuthor string
)
var Uptime = time.Now()

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// Logging Method Name
	//log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Couldn't initiate bot:  ", err)
		return
	}

	_, err = goBot.User("@me")
	if err != nil {
		log.Fatal("Couldn't get botID:  ", err)
	}

	// intent or what to store for bot?
	goBot.Identify.Intents = discordgo.IntentsAll

	// Register handlers here.
	goBot.AddHandler(guildJoinInit)
	goBot.AddHandler(pingMessageHandler)
	goBot.AddHandler(setPrefixHandler)
	goBot.AddHandler(setBotChannelHandler)
	go goBot.AddHandler(helpMessageHandler)
	goBot.AddHandler(setNicknameDuration)
	goBot.AddHandler(stats)
	goBot.AddHandler(setNick)

	// Start bot with chan.
	err = goBot.Open()
	if err != nil {
		log.Fatal("Couldn't Connect bot:  ", err)
		return
	}

	log.Info("Bot is running")
}

func guildJoinInit(s *discordgo.Session, g *discordgo.GuildCreate) {
	guild, err := s.Guild(g.ID)
	if err != nil {
		log.Error("Getting guild information from Session: ", err)
		return
	}

	guildChannels := g.Channels
	var guildIDs []string
	for _, guild := range guildChannels {
		guildIDs = append(guildIDs, guild.ID)
	}
	_, err = db.FindGuildByID(guild.ID)
	if err != nil {
		log.Error("Guild not found in DB creating one with default values... ", err)
		guildSetting := model.GuildSettings{
			GuildID:               guild.ID,
			GuildName:             guild.Name,
			GuildPrefix:           config.BotPrefix,
			GuildBotChannelsID:    guildIDs,
			GuildNicknameDuration: "30",
			TimeStamp:             time.Now().UTC(),
		}
		err := db.InsertOrUpdateGuild(&guildSetting)
		if err != nil {
			log.Error("Error inserting default values into DB", err)
			return
		}
	}
	log.Info("init successful")
}

// help menu
func helpMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
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
								AddField("Invite", "https://www.google.com").
								AddField("Support Server", "https://www.google.com").
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
								SetDescription("My default prefix is `!`. Use `!help <command>` to get more information on a command.").
								AddField("prefix", "Change the prefix or view the current prefix.").
								AddField("botchannel", "sets the current channel as bot channel or set multiple channel as bot channel.").
								AddField("nickname", "set duration for nickname changes in days").
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
								SetDescription("My default prefix is `!`. Use `!help <command>` to get more information on a command.").
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
								SetTitle("Miscellaneous").
								SetThumbnail(botImage).
								SetDescription("My default prefix is `!`. Use `!help <command>` to get more information on a command.").
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
						case 5:
							previousAuthor = m.Author.ID
							// get the time to check if it's idle or not
							currentTime := time.Now()
							// start embed
							embed := NewEmbed().
								SetTitle("Anilist").
								SetThumbnail(botImage).
								SetDescription("My default prefix is `!`. Use `!help <command>` to get more information on a command.").
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
								AddField("Invite", "https://www.google.com").
								AddField("Support Server", "https://www.google.com").
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

// ping
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

// resetMessageHandler used to reset nickname to default value, does not change duration..
func resetMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
				}
			}
		}
	}
}

// setPrefixHandler changes the prefix for the current server or show current prefix for the server.
func setPrefixHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)

		if allowedChannels {
			if strings.Contains(messageContent, guild.GuildPrefix+"prefix") {
				parameter := getArguments(messageContent)

				// if parameter is !prefix only
				if len(parameter) == 1 {
					// embed start
					embed := NewEmbed().
						SetDescription("The prefix for this server is `" + guild.GuildPrefix + "`.").
						SetColor(green).MessageEmbed
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("On sending parameter error message to channel: ", err)
					}
					return
				}

				// check if the user is admin before executing admin privileged commands.
				permission, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
				if err != nil {
					log.Error("Getting user permission: ", err)
					return
				}
				guildOwner, err := checkGuildOwner(s, m)
				if err != nil {
					log.Error("Failed to check guild owner: ", err)
					return
				}
				if permission || guildOwner {
					prefix := parameter[1]
					newPrefix := checkPrefix(prefix)
					if newPrefix {
						// change prefix code
						// get Guild information
						updateGuild, err := s.Guild(m.GuildID)
						if err != nil {
							log.Error("Failed to get Guild: ", err)
							return
						}

						currentTime := time.Now().UTC()
						guildSettings := &model.GuildSettings{
							GuildID:               m.GuildID,
							GuildName:             updateGuild.Name,
							GuildPrefix:           prefix,
							GuildBotChannelsID:    guild.GuildBotChannelsID,
							GuildNicknameDuration: guild.GuildNicknameDuration,
							TimeStamp:             currentTime,
						}

						// insert new prefix to database
						err = db.InsertOrUpdateGuild(guildSettings)
						if err != nil {
							log.Warn("Inserting or Updating guild prefix: ", err)
							return
						}
						guildData, err := db.FindGuildByID(m.GuildID)
						if err != nil {
							log.Warn("Couldn't find guild: ", err)
							return
						}

						// start Embed
						embed := NewEmbed().
							SetDescription("Updated successfully prefix now set to `" +
								guildData.GuildPrefix + "`").
							SetColor(green).MessageEmbed
						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Warn("Failed to send embed to the channel: ", err)
							return
						}
					} else {
						// start Embed
						embed := NewEmbed().
							SetDescription("The chosen prefix is too long.").
							SetColor(red).MessageEmbed
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Warn("Failed to send embed to the channel: ", err)
							return
						}
					}
				} else {
					// start Embed
					embed := NewEmbed().
						SetDescription("Sorry you do not have permission to execute that command.").
						SetColor(red).MessageEmbed
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
}

// setBotChannelHandler sets bot channel for the given channel or channels by providing channel IDs
func setBotChannelHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
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
			parameter := getArguments(messageContent)
			if strings.Contains(messageContent, guild.GuildPrefix+"botchannel") {
				// check if the user is admin before executing admin privileged commands.
				permission, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
				if err != nil {
					log.Error("Getting user permission: ", err)
					return
				}
				// check if it's guild owner
				guildOwner, err := checkGuildOwner(s, m)
				if err != nil {
					log.Error("Failed to check guild owner: ", err)
					return
				}
				// check if they have permission
				if permission || guildOwner {
					// if setting one channel only
					if len(parameter) == 1 {
						// add current channel as bot channel
						guildChannels := []string{m.ChannelID}

						currentTime := time.Now().UTC()
						guildSettings := &model.GuildSettings{
							GuildID:               m.GuildID,
							GuildName:             guild.GuildName,
							GuildPrefix:           guild.GuildPrefix,
							GuildBotChannelsID:    guildChannels,
							GuildNicknameDuration: guild.GuildNicknameDuration,
							TimeStamp:             currentTime,
						}

						// insert new prefix to database
						err = db.InsertOrUpdateGuild(guildSettings)
						if err != nil {
							log.Warn("Inserting or Updating guild prefix: ", err)
							return
						}
						guildData, err := db.FindGuildByID(m.GuildID)
						if err != nil {
							log.Warn("Couldn't find guild: ", err)
							return
						}

						// start Embed
						embed := NewEmbed().
							SetDescription("Updated successfully this " +
								"channel is set to take bot commands; channel ID: `" +
								guildData.GuildBotChannelsID[0] + "`").
							SetColor(green).MessageEmbed
						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Warn("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					// add multiple channels as bot channel
					var guildChannels []string
					parameterOnly := parameter[1:]
					for _, ids := range parameterOnly {
						if len(ids) < 18 {
							// start Embed
							embed := NewEmbed().
								SetDescription("Make sure the channel ID is correct potential " +
									"issue with: `" + ids + "`" + " Aborting bot channel update").
								SetColor(red).MessageEmbed
							log.Warn("Potential issue with channel ID not equal or greater than 18")
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
							if err != nil {
								log.Warn("Failed to send embed to the channel: ", err)
								return
							}
							return
						}
					}
					guildChannels = append(guildChannels, parameterOnly...)
					currentTime := time.Now().UTC()
					guildSettings := &model.GuildSettings{
						GuildID:               m.GuildID,
						GuildName:             guild.GuildName,
						GuildPrefix:           guild.GuildPrefix,
						GuildBotChannelsID:    guildChannels,
						GuildNicknameDuration: guild.GuildNicknameDuration,
						TimeStamp:             currentTime,
					}
					// insert new prefix to database
					err = db.InsertOrUpdateGuild(guildSettings)
					if err != nil {
						log.Warn("Inserting or Updating guild prefix: ", err)
						return
					}
					guildData, err := db.FindGuildByID(m.GuildID)
					if err != nil {
						log.Warn("Couldn't find guild: ", err)
						return
					}

					guildChannels = guildData.GuildBotChannelsID

					joinedChannelID := strings.Join(guildChannels, ",")
					// start Embed
					embed := NewEmbed().
						SetDescription("Updated successfully the channel IDs: \n`" + joinedChannelID +
							"` \n" +
							"now take bot commands.").
						SetColor(0x11ff00).MessageEmbed
					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Warn("Failed to send embed to the channel: ", err)
						return
					}
				} else {
					// start Embed
					embed := NewEmbed().
						SetDescription("Sorry you do not have permission to execute that command.").
						SetColor(red).MessageEmbed
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
}

// setNicknameDuration setting the nickname days
func setNicknameDuration(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
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
			parameter := getArguments(messageContent)
			if strings.Contains(messageContent, guild.GuildPrefix+"nickname") {
				// if parameter is is none bring out current days set.
				if len(parameter) == 1 {
					// embed start
					embed := NewEmbed().
						SetDescription("The nickname duration for this server is `" + guild.GuildNicknameDuration + " days`.").
						SetColor(green).MessageEmbed
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("On sending parameter error message to channel: ", err)
					}
					return
				}

				// check if the user is admin before executing admin privileged commands.
				permission, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
				if err != nil {
					log.Error("Getting user permission: ", err)
					return
				}
				guildOwner, err := checkGuildOwner(s, m)
				if err != nil {
					log.Error("Failed to check guild owner: ", err)
					return
				}
				if permission || guildOwner {
					enteredDays := parameter[1]
					// check if the argument provided is integer or number only.
					_, err := strconv.ParseInt(enteredDays, 10, 32)
					if err != nil {
						log.Warn(" User error failed to convert string to int: ", err)
						// start Embed
						embed := NewEmbed().
							SetDescription(
								"Are you fucking serious? even a two year old knows what a number is. ").
							SetColor(red).MessageEmbed
						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Warn("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					currentTime := time.Now().UTC()
					guildSettings := &model.GuildSettings{
						GuildID:               m.GuildID,
						GuildName:             guild.GuildName,
						GuildPrefix:           guild.GuildPrefix,
						GuildBotChannelsID:    guild.GuildBotChannelsID,
						GuildNicknameDuration: enteredDays,
						TimeStamp:             currentTime,
					}

					// insert new prefix to database
					err = db.InsertOrUpdateGuild(guildSettings)
					if err != nil {
						log.Warn("Inserting or Updating guild nickname duration: ", err)
						return
					}
					guildData, err := db.FindGuildByID(m.GuildID)
					if err != nil {
						log.Warn("Couldn't find guild: ", err)
						return
					}

					// start Embed
					embed := NewEmbed().
						SetDescription("Updated successfully " +
							"Nickname duration is set to: `" +
							guildData.GuildNicknameDuration + "days`").
						SetColor(green).MessageEmbed
					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Warn("Failed to send embed to the channel: ", err)
						return
					}
					return
				} else {
					// start Embed
					embed := NewEmbed().
						SetDescription("Sorry you do not have permission to execute that command.").
						SetColor(red).MessageEmbed
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

				//TODO:THIS WILL NOT WORK UNLESS IT"S LINUX SYSTEM SO COMMENT UNTIL MIGRATED
				// get cpu and memory usage.
				//cpuUsage, err := getCpuUsage()
				//if err != nil {
				//	log.Error("Failed to get CPU Usage: ", err)
				//}

				//TODO:THIS WILL NOT WORK UNLESS IT"S LINUX SYSTEM SO COMMENT UNTIL MIGRATED
				// get memory usage
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
					AddField("CPU Usage", "need to be worked on").
					AddField("RAM Usage", "need to be worked on").
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

// nick
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

		// check if it's nick or nickname since contains func will return both nickname and nick function.
		if parameter[0] != guild.GuildPrefix+"nick" {
			return
		}

		if allowedChannels {
			if strings.Contains(messageContent, guild.GuildPrefix+"nick") {

				// If user can change their nickname
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
				}
			} else {
				err = getTimeLeftForNick(s, m, "")
				if err != nil {
					log.Error("Failed to get time left for nick change")
					return
				}
			}
		}
	}
}
