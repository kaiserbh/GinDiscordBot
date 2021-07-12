package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func ban(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Command name
	// !<command>
	command := "ban"

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

	}

	allowedChannels := checkAllowedChannel(m.ChannelID, guild)
	messageContent := strings.ToLower(m.Content)
	parameter := getArguments(messageContent)
	if allowedChannels {
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
			// If command matches syntax do blah

			if parameter[0] == guild.GuildPrefix+command {
				switch {
				case len(parameter) == 2:
					err := s.GuildBanCreate(m.GuildID, strings.Trim(parameter[1], "<@!>"), 0)
					if err != nil {
						log.Error("Error banning member %s", strings.Trim(parameter[1], "<@!>"))
					}
					embed := NewEmbed().
						SetDescription(fmt.Sprintf("Banned member %s", strings.Trim(parameter[1], "<@!>"))).
						SetColor(red).
						MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

					if err != nil {
						log.Error("Error sending embed: ", err)
						return
					}

				case len(parameter) > 2:
					if strings.ContainsAny(parameter[2], "1234567890") && strings.ContainsAny(strings.TrimLeft(parameter[3], "1234567890"), "mhdy") {

						reason := strings.Join(parameter[3:], " ")

						err := s.GuildBanCreateWithReason(m.GuildID, parameter[1], reason, 0)

						if err != nil {
							log.Error("Error banning member %s with reason %s", strings.Trim(parameter[1], "<@!>"), reason)
							return
						}

						embed := NewEmbed().
							SetDescription(fmt.Sprintf("Banned member `%s` for `%s` with reason `%s`", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2], reason)).
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Error sending embed: ", err)
							return
						}

					} else {
						reason := strings.Join(parameter[2:], " ")

						err := s.GuildBanCreateWithReason(m.GuildID, parameter[1], reason, 0)

						if err != nil {
							log.Error("Error banning member %s with reason %s", strings.Trim(parameter[1], "<@!>"), reason)
							return
						}

						embed := NewEmbed().
							SetDescription(fmt.Sprintf("Banned member `%s` with reason `%s`", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), reason)).
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Error sending embed: ", err)
							return
						}
					}

				default:
					embed := NewEmbed().
						SetTitle(fmt.Sprintf("Command: %sban", guild.GuildPrefix)).
						SetDescription("Ban a member, optional time limit").
						AddField("Syntax", fmt.Sprintf("%sban [user] (ban length) (reason)", guild.GuildPrefix)).
						SetColor(0x000000).
						MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Error sending embed: ", err)
						return
					}
				}
			}
		}
	}
}
func pardon(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Command name
	// !<command>
	command := "pardon"

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

	}

	allowedChannels := checkAllowedChannel(m.ChannelID, guild)
	messageContent := strings.ToLower(m.Content)
	parameter := getArguments(messageContent)

	if allowedChannels {
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
			// If command matches syntax do blah

			if parameter[0] == guild.GuildPrefix+command {
				err := s.GuildBanDelete(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"))
				if err != nil {
					log.WithFields(log.Fields{
						"User": strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"),
					}).Error("Error Pardoning member: ", err)
					return
				}
				log.WithFields(log.Fields{
					"User": strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"),
				}).Info("Pardoned member")
				embed := NewEmbed().
					SetDescription(fmt.Sprintf("Pardoned member %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"))).
					SetColor(green).
					MessageEmbed

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

				if err != nil {
					log.Error("Error sending embed: ", err)
					return
				}
			}
		}
	}
}

// Gaki command disabled

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

// func chatFilter(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	Filter()
// }
