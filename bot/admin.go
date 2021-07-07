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

	if allowedChannels {
		messageContent := strings.ToLower(m.Content)
		parameter := getArguments(messageContent)

		// If command matches syntax do blah

		untrimmedTag := parameter[1]

		member := strings.TrimSuffix(strings.TrimPrefix(untrimmedTag, "<@!"), ">")

		if parameter[0] == guild.GuildPrefix+command {

			// Ban member
			err := s.GuildBanCreate(m.GuildID, member, 0)

			if err != nil {
				log.Error(fmt.Sprintf("Error banning member %s: ", member), err)
				return
			}
			log.Info(fmt.Sprintf("Banned member %s", member))

			// Initialize embed
			embed := NewEmbed().
				SetDescription(fmt.Sprintf("Banned member %s", member)).
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

		untrimmedTag := parameter[1]

		member := strings.TrimSuffix(strings.TrimPrefix(untrimmedTag, "<@!"), ">")

		if parameter[0] == guild.GuildPrefix+command {
			err := s.GuildBanDelete(m.GuildID, member)
			if err != nil {
				log.Error(fmt.Sprintf("Error pardoning member %s: ", member), err)
				return
			}
			log.Info(fmt.Sprintf("Pardoned member %s", member))

			embed := NewEmbed().
				SetDescription(fmt.Sprintf("Pardoned member %s", member)).
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
