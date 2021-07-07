package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/anilistgo"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// anime query anime from anilist by id or name.
func anime(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	args := getArguments(messageContent)
	animeArgs := args[0]

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if animeArgs == guild.GuildPrefix+"anime" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't anime args.")
					return
				}

				animeQuery := args[1]

				// check if it's a number or a strings.
				animeID, err := strconv.Atoi(animeQuery)
				if err != nil {
					// query anime by title
					anime := anilistgo.NewMediaQuery()
					_, err := anime.FilterByTitle(animeQuery)
					if err != nil {
						log.Error("Failed to filter anime by title: ", err)
						return
					}

					// start embed
					embed := NewEmbed().
						SetTitle("Anime").
						SetImage(anime.BannerImage).
						SetThumbnail(anime.CoverImage.Large).
						SetDescription(anime.Description).
						SetColor(green).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
					}

					// make sure the bot messageID is taken after it's being sent.
					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}

					err = s.MessageReactionAdd(m.ChannelID, botMessageID, "✅")
					if err != nil {
						log.Error("Failed to add reaction: ", err)
					}
					err = s.MessageReactionAdd(m.ChannelID, botMessageID, "❌")
					if err != nil {
						log.Error("Failed to add reaction: ", err)
					}
					startTimer := time.Now()

					for {
						passedTimer := time.Since(startTimer).Seconds()
						checkAuthorReactionOk, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "✅", m.Author.ID, 10)
						if err != nil {
							log.Error(err)
							return
						}
						if checkAuthorReactionOk {
							err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to remove reactions from bot message: ", err)
								return
							}
						}

						// check the delete reaction
						checkAuthorReactionDelete, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "❌", m.Author.ID, 10)
						if err != nil {
							log.Error(err)
							return
						}

						if checkAuthorReactionDelete {
							err := s.ChannelMessageDelete(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to delete botMessage: ", err)
								return
							}
						}
						// if no reactions is added then just remove reactions from the message.
						if passedTimer >= 10 {
							err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to remove reactions from bot message: ", err)
								return
							}
						}
						return
					}
				} else {
					// query anime by title
					anime := anilistgo.NewMediaQuery()
					_, err := anime.FilterAnimeByID(animeID)
					if err != nil {
						log.Error("Failed to filter anime by ID: ", err)
						return
					}

					// start embed
					embed := NewEmbed().
						SetTitle("Anime").
						SetImage(anime.BannerImage).
						SetThumbnail(anime.CoverImage.Large).
						SetDescription(anime.Description).
						SetColor(green).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}

					err = s.MessageReactionAdd(m.ChannelID, botMessageID, "✅")
					if err != nil {
						log.Error("Failed to add reaction: ", err)
					}
					err = s.MessageReactionAdd(m.ChannelID, botMessageID, "❌")
					if err != nil {
						log.Error("Failed to add reaction: ", err)
					}
					startTimer := time.Now()

					for {
						passedTimer := time.Since(startTimer).Seconds()
						checkAuthorReactionOk, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "✅", m.Author.ID, 10)
						if err != nil {
							log.Error(err)
							return
						}
						if checkAuthorReactionOk {
							err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to remove reactions from bot message: ", err)
								return
							}
						}

						// check the delete reaction
						checkAuthorReactionDelete, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "❌", m.Author.ID, 10)
						if err != nil {
							log.Error(err)
							return
						}

						if checkAuthorReactionDelete {
							err := s.ChannelMessageDelete(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to delete botMessage: ", err)
								return
							}
						}
						// if no reactions is added then just remove reactions from the message.
						if passedTimer >= 10 {
							err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
							if err != nil {
								log.Error("Failed to remove reactions from bot message: ", err)
								return
							}
						}
						return
					}
				}
			}
		}
	}
}
