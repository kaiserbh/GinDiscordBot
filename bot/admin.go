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

	if allowedChannels {

		messageContent := strings.ToLower(m.Content)

		parameter := getArguments(messageContent)

		// If command matches syntax do blah

		if parameter[0] == guild.GuildPrefix+command {
			err := s.GuildBanDelete(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"))
			if err != nil {
				log.Error(fmt.Sprintf("Error pardoning member %s: ", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")), err)
				return
			}
			log.Info(fmt.Sprintf("Pardoned member %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")))

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