package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func ban(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			if parameter[0] == guild.GuildPrefix+"ban" {
				if strings.Contains(messageContent, guild.GuildPrefix+"ban") {
					// if no length for ban, bring out ban syntax
					if len(parameter) == 1 {
						// embed start
						embed := NewEmbed().
							SetTitle(fmt.Sprintf("Command: %sban", guild.GuildPrefix)).
							SetDescription("**Description:** Ban a member, optional time limit").
							AddField("Usage:", fmt.Sprintf("%sban [user] (time) (reason)", guild.GuildPrefix)).
							SetColor(red).MessageEmbed
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

					switch permission || guildOwner {
					case len(parameter) == 2:
						Member := parameter[1]

						retard, err := s.GuildBan(m.GuildID, Member)

						if err != nil {
							log.Error("fucking dumbass bot broken:", err)
							return
						}
						if retard != nil {
							log.Error("fuck this bot", retard)
							return
						}

						return

					default:
						// embed start
						embed := NewEmbed().
							SetTitle(fmt.Sprintf("Command: %sban", guild.GuildPrefix)).
							SetDescription("**Description:** Ban a member, optional time limit").
							AddField("Usage:", fmt.Sprintf("%sban [user] (time) (reason)", guild.GuildPrefix)).
							SetColor(red).MessageEmbed
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("On sending parameter error message to channel: ", err)
						}
						return

					}
					/*
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
						}*/
				}
			}
		}
	}
}
