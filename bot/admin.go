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
				// Ban member
				err := s.GuildBanCreate(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), 0)

				if err != nil {
					log.Error(fmt.Sprintf("Error banning member %s: ", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")), err)
					return
				}
				log.Info(fmt.Sprintf("Banned member %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")))

				// Initialize embed
				embed := NewEmbed().
					SetDescription(fmt.Sprintf("Banned member %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"))).
					SetColor(red).
					MessageEmbed

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

				if err != nil {
					log.Error("Error sending embed: ", err)
					return
				}
				return

			case len(parameter) == 3:

				err := s.GuildBanCreateWithReason(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2], 0)

				if err != nil {
					log.Error(fmt.Sprintf("Error banning member %s with reason", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")))
					return
				}

				log.Info("Banned member %s with reason %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2])

				embed := NewEmbed().
					SetDescription(fmt.Sprintf("Banned member %s with reason %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2])).
					SetColor(red).
					MessageEmbed

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

				if err != nil {
					log.Error("Error sending embed: ", err)
					return
				}

			case len(parameter) == 4:
				err := s.GuildBanCreate(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), 0)
				if err != nil {
					log.Error(fmt.Sprintf("Error banning member %s ", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">")), err)
					return
				}

				embed := NewEmbed().
					SetDescription(fmt.Sprintf("Banned member %s for %s %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2], parameter[3])).
					SetColor(red).
					MessageEmbed

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

				if err != nil {
					log.Error("Error sending embed: ", err)
					return
				}

			case len(parameter) == 5:
				err := s.GuildBanCreateWithReason(m.GuildID, strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[4], 0)

				if err != nil {
					log.Error("Error banning member %s with reason %s", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[4])
					return
				}

				embed := NewEmbed().
					SetDescription(fmt.Sprintf("Banned member **%s** for **%s %s** with reason **%s**", strings.TrimSuffix(strings.TrimPrefix(parameter[1], "<@!"), ">"), parameter[2], parameter[3], parameter[4])).
					AddField("Reason:", parameter[4]).
					SetColor(red).
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
				SetColor(red).
				MessageEmbed

			_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)

			if err != nil {
				log.Error("Error sending embed: ", err)
				return
			}
		}
	}
}
