package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/config"
	"github.com/kaiserbh/gin-bot-go/database"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var db = database.Connect()
var (
	red            = 0xff0000
	green          = 0x11ff00
	previousAuthor string
)

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
			GuildID:            guild.ID,
			GuildName:          guild.Name,
			GuildPrefix:        config.BotPrefix,
			GuildBotChannelsID: guildIDs,
			TimeStamp:          time.Now().UTC(),
		}
		err := db.InsertOrUpdateGuild(&guildSetting)
		if err != nil {
			log.Error("Error inserting default values into DB", err)
			return
		}
	}
	log.Info("init successful")
}

func helpMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}

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
			if m.Content == guild.GuildPrefix+"help" {
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

				for {
					switch page {
					// page one About page
					case 1:
						previousAuthor = m.Author.ID
						// get the time to check if it's idle or not
						currentTime := time.Now()
						// start embed
						embed := NewEmbed().
							SetDescription("Help Menu!").
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

						// loop through the time?
						for {
							timePassed := time.Since(currentTime)
							if timePassed.Seconds() >= 30 {
								log.WithFields(log.Fields{
									"Time passed": timePassed,
								}).Info("Removing reactions time has been passed.")
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							}
							// check if the reaction matches the author ID aka sender
							checkReaction, err := checkMessageReaction(s, m, botMessageID)
							if err != nil {
								log.Error("Failed to check emoji from bot message.")
								return
							}
							// if true then remove it uwu.
							if checkReaction["Stop"] {
								fmt.Println("executed stop")
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							} else if checkReaction["FastBack"] {
								page = 1
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏮️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							} else if checkReaction["Back"] {
								if page == 1 {
									log.Info("Already on page one not doing anything")
									// remove user reaction before going to next page.
									err = s.MessageReactionRemove(m.ChannelID, botMessageID, "◀️", m.Author.ID)
									if err != nil {
										log.Error("Failed to remove user reaction from bot message")
										return
									}
								}
								page -= 1
								break
							} else if checkReaction["Forward"] {
								if page == 5 {
									log.Info("Last Page already not doing anything")
									// remove user reaction before going to next page.
								}
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "▶️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								page += 1
								break
							} else if checkReaction["FastForward"] {
								page = 5
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏭️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							}
						}
					case 2:
						previousAuthor = m.Author.ID
						// get the time to check if it's idle or not
						currentTime := time.Now()
						// start embed
						embed := NewEmbed().
							SetDescription("Help Menu! Page 2").
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
						// add reaction to the bot message with for loop?
						for _, emoji := range reactions {
							err = s.MessageReactionAdd(m.ChannelID, botMessageID, emoji)
							if err != nil {
								log.Error("Failed to add reaction: ", err)
								return
							}
						}
						// loop through the time?
						for {
							timePassed := time.Since(currentTime)
							if timePassed.Seconds() >= 30 {
								log.WithFields(log.Fields{
									"Time passed": timePassed,
								}).Info("Removing reactions time has been passed.")
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							}
							// check if the reaction matches the author ID aka sender
							checkReaction, err := checkMessageReaction(s, m, botMessageID)
							if err != nil {
								log.Error("Failed to check emoji from bot message.")
								return
							}
							// if true then remove it uwu.
							if checkReaction["Stop"] {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							} else if checkReaction["FastBack"] {
								page = 1
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏮️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							} else if checkReaction["Back"] {
								if page == 1 {
									log.Info("Already on page one not doing anything")
									// remove user reaction before going to next page.
									err = s.MessageReactionRemove(m.ChannelID, botMessageID, "◀️", m.Author.ID)
									if err != nil {
										log.Error("Failed to remove user reaction from bot message")
										return
									}
								}
								page -= 1
								break
							} else if checkReaction["Forward"] {
								if page == 5 {
									log.Info("Last Page already not doing anything")
									// remove user reaction before going to next page.
								}
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "▶️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								page += 1
								break
							} else if checkReaction["FastForward"] {
								page = 5
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏭️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							}
						}
					default:
						// reset page
						page = 1

						previousAuthor = m.Author.ID
						// get the time to check if it's idle or not
						currentTime := time.Now()
						// start embed
						embed := NewEmbed().
							SetDescription("Help Menu!").
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
						// add reaction to the bot message with for loop?
						for _, emoji := range reactions {
							err = s.MessageReactionAdd(m.ChannelID, botMessageID, emoji)
							if err != nil {
								log.Error("Failed to add reaction: ", err)
								return
							}
						}
						// loop through the time?
						for {
							timePassed := time.Since(currentTime)
							if timePassed.Seconds() >= 30 {
								log.WithFields(log.Fields{
									"Time passed": timePassed,
								}).Info("Removing reactions time has been passed.")
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							}
							// check if the reaction matches the author ID aka sender
							checkReaction, err := checkMessageReaction(s, m, botMessageID)
							if err != nil {
								log.Error("Failed to check emoji from bot message.")
								return
							}
							// if true then remove it uwu.
							if checkReaction["Stop"] {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								previousAuthor = ""
								return
							} else if checkReaction["FastBack"] {
								page = 1
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏮️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							} else if checkReaction["Back"] {
								if page == 1 {
									log.Info("Already on page one not doing anything")
									// remove user reaction before going to next page.
									err = s.MessageReactionRemove(m.ChannelID, botMessageID, "◀️", m.Author.ID)
									if err != nil {
										log.Error("Failed to remove user reaction from bot message")
										return
									}
								}
								page -= 1
								break
							} else if checkReaction["Forward"] {
								if page == 5 {
									log.Info("Last Page already not doing anything")
									// remove user reaction before going to next page.
								}
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "▶️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								page += 1
								break
							} else if checkReaction["FastForward"] {
								page = 5
								// remove user reaction before going to next page.
								err = s.MessageReactionRemove(m.ChannelID, botMessageID, "⏭️", m.Author.ID)
								if err != nil {
									log.Error("Failed to remove user reaction from bot message")
									return
								}
								break
							}
						}
					}
				}
			}
		}
	}
}

func pingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	// check if the channel is bot channel or allowed channel.
	allowedChannels := checkAllowedChannel(m.ChannelID, guild)
	if allowedChannels {
		if strings.HasPrefix(m.Content, guild.GuildPrefix) {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if m.Content == guild.GuildPrefix+"ping" {
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

func setPrefixHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
				if permission {
					prefix := parameter[1]
					newPrefix := checkPrefix(prefix)
					if newPrefix {
						// change prefix code
						// get Guild information
						guild, err := s.Guild(m.GuildID)
						if err != nil {
							log.Error("Failed to get Guild: ", err)
							return
						}

						currentTime := time.Now().UTC()
						guildSettings := &model.GuildSettings{
							GuildID:     m.GuildID,
							GuildName:   guild.Name,
							GuildPrefix: prefix,
							TimeStamp:   currentTime,
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
				}
			}
		}
	}
}

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
			if strings.Contains(messageContent, guild.GuildPrefix+"setbotchannel") {
				// check if the user is admin before executing admin privileged commands.
				permission, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAdministrator)
				if err != nil {
					log.Error("Getting user permission: ", err)
					return
				}
				if permission {
					// if setting one channel only
					if len(parameter) == 1 {
						// add current channel as bot channel
						guildChannels := []string{m.ChannelID}

						currentTime := time.Now().UTC()
						guildSettings := &model.GuildSettings{
							GuildID:            m.GuildID,
							GuildName:          guild.GuildName,
							GuildPrefix:        guild.GuildPrefix,
							GuildBotChannelsID: guildChannels,
							TimeStamp:          currentTime,
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
						GuildID:            m.GuildID,
						GuildName:          guild.GuildName,
						GuildPrefix:        guild.GuildPrefix,
						GuildBotChannelsID: guildChannels,
						TimeStamp:          currentTime,
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
				}
			}
		}
	}
}

//TODO: CREATE NICKNAME CHANGING CAPABILITY WITH 1MONTH TIME LIMIT
