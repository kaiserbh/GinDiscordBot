package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

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
				//guildOwner, err := checkGuildOwner(s, m)
				//if err != nil {
				//	log.Error("Failed to check guild owner: ", err)
				//	return
				//}
				if permission {
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
				//guildOwner, err := checkGuildOwner(s, m)
				//if err != nil {
				//	log.Error("Failed to check guild owner: ", err)
				//	return
				//}
				// check if they have permission
				if permission {
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

// setNicknameCooldown setting the nickname days
func setNicknameCooldown(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			messageContent := strings.ToLower(m.Content)
			parameter := getArguments(messageContent)
			if parameter[0] == guild.GuildPrefix+"cooldown" {
				if strings.Contains(messageContent, guild.GuildPrefix+"cooldown") {
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
					//guildOwner, err := checkGuildOwner(s, m)
					//if err != nil {
					//	log.Error("Failed to check guild owner: ", err)
					//	return
					//}
					if permission {
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
								guildData.GuildNicknameDuration + " days`").
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
}
