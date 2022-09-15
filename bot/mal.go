package bot

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nstratos/go-myanimelist/mal"
	log "github.com/sirupsen/logrus"
)

type clientIDTransport struct {
	Transport http.RoundTripper
	ClientID  string
}

func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.ClientID)
	return c.Transport.RoundTrip(req)
}

var publicInfoClient = &http.Client{
	Transport: &clientIDTransport{ClientID: "8153b32dcfa0299d57e6a1c5299e69d2"},
}

// anime query media from MAL by id or name.
func malAnime(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			if animeArgs == guild.GuildPrefix+"anime_mal" || animeArgs == guild.GuildPrefix+"am" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				animeQuery := args[1]

				// set mal
				c := mal.NewClient(publicInfoClient)
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
				defer cancel()

				// check if it's a number or a strings.
				anime_id, err := strconv.Atoi(animeQuery)
				if err != nil {
					// query media by title
					// join multiple search query
					animeQuery = strings.Join(args[1:], " ")
					anime_list, _, err := c.Anime.List(ctx, animeQuery,
						mal.Fields{"rank", "popularity", "synopsis", "mean", "background", "status", "source", "media_type", "num_episodes", "average_episode_duration", "start_date", "end_date", "alternative_titles", "start_season", "genres", "studios"},
						mal.Limit(5),
					)
					if err != nil {
						log.Error("Failed to filter media by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Anime not found!\n maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}
					// check if there is more than one result
					if len(anime_list) > 1 {
						// append the title then send it as embed
						var animeList []string
						for index, anime := range anime_list {
							index += 1
							item := strconv.Itoa(index)
							animeList = append(animeList, item+". "+anime.Title)
						}
						// start embed
						embed := NewEmbed().
							SetTitle("More than one result found (ANIME)").
							SetDescription(strings.Join(animeList, "\n")).SetColor(green).MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						reactions := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

						botMessageID, err := getBotMessageID(s, m)
						if err != nil {
							log.Error("Failed to get bot message ID: ", err)
							return
						}

						// add reaction to the bot message based on the length of anime_list
						for _, emoji := range reactions[:len(anime_list)] {
							err = s.MessageReactionAdd(m.ChannelID, botMessageID, emoji)
							if err != nil {
								log.Error("Failed to add reaction: ", err)
								return
							}
						}

						var chosen_anime int

						startTimer := time.Now()

						for {
							passedTimer := time.Since(startTimer).Seconds()
							checkAuthorReactionOne, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "1️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionOne {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_anime = 0
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionTwo, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "2️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionTwo {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_anime = 1
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionThree, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "3️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionThree {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_anime = 2
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionFour, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "4️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionFour {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_anime = 3
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionFive, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "5️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionFive {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_anime = 4
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							// if no reactions is added then just remove reactions from the message.
							if passedTimer >= 10 {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								return
							}
						}

						genres_list := anime_list[chosen_anime].Genres
						var genres []string
						for _, genre := range genres_list {
							genres = append(genres, genre.Name)
						}

						genres_join := strings.Join(genres, ",")
						if genres_join == "" {
							genres_join = "\u200b"
						}

						studios_list := anime_list[chosen_anime].Studios
						var studios []string
						for _, studio := range studios_list {
							studios = append(studios, studio.Name)
						}

						studios_join := strings.Join(studios, ",")
						if studios_join == "" {
							studios_join = "\u200b"
						}

						if anime_list[chosen_anime].EndDate == "" {
							anime_list[chosen_anime].EndDate = "Ongoing"
						}

						// start embed
						embed = NewEmbed().
							SetTitle(anime_list[chosen_anime].Title).
							SetURL("https://myanimelist.net/anime/"+strconv.Itoa(anime_list[chosen_anime].ID)).
							SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
							SetThumbnail(anime_list[chosen_anime].MainPicture.Large).
							SetDescription(cutDescription(anime_list[chosen_anime].Synopsis)).
							AddField("Format", anime_list[chosen_anime].MediaType).
							AddField("Episodes", strconv.Itoa(anime_list[chosen_anime].NumEpisodes)).
							AddField("Episode Duration", strconv.Itoa(anime_list[chosen_anime].AverageEpisodeDuration/60)+" min").
							AddField("Status", anime_list[chosen_anime].Status).
							AddField("Start Date", anime_list[chosen_anime].StartDate).
							AddField("End Date", anime_list[chosen_anime].EndDate).
							//AddField("Season", anime_list[chosen_anime].StartSeason.Season).
							AddField("Mean Score", strconv.FormatFloat(anime_list[chosen_anime].Mean, 'f', 3, 64)).
							//AddField("Popularity", strconv.Itoa(anime_list[chosen_anime].Popularity)).
							//AddField("Rank", strconv.Itoa(anime_list[chosen_anime].Rank)).
							//AddField("Source", anime_list[chosen_anime].Source).
							AddField("Genres", genres_join).
							AddField("Studio", studios_join).
							SetFooter(anime_list[chosen_anime].AlternativeTitles.Ja, anime_list[chosen_anime].MainPicture.Medium).
							InlineAllFields().
							SetColor(green).MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						// botMessageID, err = getBotMessageID(s, m)
						// if err != nil {
						// 	log.Error("Failed to get bot message ID: ", err)
						// 	return
						// }
						// // check timer and reaction
						// checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)

					}
				} else {
					anime_id, _, err := c.Anime.Details(ctx, anime_id,
						mal.Fields{"rank", "popularity", "synopsis", "mean", "background", "status", "source", "media_type", "num_episodes", "average_episode_duration", "start_date", "end_date", "alternative_titles", "start_season", "genres", "studios"},
					)
					if err != nil {
						log.Error("Failed to filter media by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Anime not found!\n maybe try using Title?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					genres_list := anime_id.Genres
					var genres []string
					for _, genre := range genres_list {
						genres = append(genres, genre.Name)
					}

					genres_join := strings.Join(genres, ",")
					if genres_join == "" {
						genres_join = "\u200b"
					}

					studios_list := anime_id.Studios
					var studios []string
					for _, studio := range studios_list {
						studios = append(studios, studio.Name)
					}

					studios_join := strings.Join(studios, ",")
					if studios_join == "" {
						studios_join = "\u200b"
					}

					if anime_id.EndDate == "" {
						anime_id.EndDate = "Ongoing"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(anime_id.Title).
						SetURL("https://myanimelist.net/anime/"+strconv.Itoa(anime_id.ID)).
						SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
						SetThumbnail(anime_id.MainPicture.Large).
						SetDescription(cutDescription(anime_id.Synopsis)).
						AddField("Format", anime_id.MediaType).
						AddField("Episodes", strconv.Itoa(anime_id.NumEpisodes)).
						AddField("Episode Duration", strconv.Itoa(anime_id.AverageEpisodeDuration/60)+" min").
						AddField("Status", anime_id.Status).
						AddField("Start Date", anime_id.StartDate).
						AddField("End Date", anime_id.EndDate).
						//AddField("Season", anime_id.StartSeason.Season).
						AddField("Mean Score", strconv.FormatFloat(anime_id.Mean, 'f', 3, 64)).
						//AddField("Popularity", strconv.Itoa(anime_id.Popularity)).
						//AddField("Rank", strconv.Itoa(anime_id.Rank)).
						//AddField("Source", anime_id.Source).
						AddField("Genres", genres_join).
						AddField("Studio", studios_join).
						SetFooter(anime_id.AlternativeTitles.Ja, anime_id.MainPicture.Medium).
						InlineAllFields().
						SetColor(green).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

				}
			}
		}
	}
}

// manga query media from MAL by id or name.
func malManga(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	args := getArguments(messageContent)
	mangaArgs := args[0]

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if mangaArgs == guild.GuildPrefix+"manga_mal" || mangaArgs == guild.GuildPrefix+"mm" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				animeQuery := args[1]

				// set mal
				c := mal.NewClient(publicInfoClient)
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
				defer cancel()

				// check if it's a number or a strings.
				manga_id, err := strconv.Atoi(animeQuery)
				if err != nil {
					// query media by title
					// join multiple search query
					animeQuery = strings.Join(args[1:], " ")
					manga_list, _, err := c.Manga.List(ctx, animeQuery,
						mal.Fields{"id", "title", "main_picture", "alternative_titles", "start_date", "end_date", "synopsis", "mean", "rank", "popularity", "num_list_users", "num_scoring_users", "nsfw", "genres", "media_type", "status", "num_volumes", "num_chapters", "authors"},
						mal.Limit(5),
					)
					if err != nil {
						log.Error("Failed to filter manga by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Manga not found!\n maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}
					// check if there is more than one result
					if len(manga_list) > 1 {
						// append the title then send it as embed
						var mangaList []string
						for index, manga := range manga_list {
							index += 1
							item := strconv.Itoa(index)
							mangaList = append(mangaList, item+". "+manga.Title)
						}
						// start embed
						embed := NewEmbed().
							SetTitle("More than one result found (MANGA)").
							SetDescription(strings.Join(mangaList, "\n")).SetColor(green).MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						reactions := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}

						botMessageID, err := getBotMessageID(s, m)
						if err != nil {
							log.Error("Failed to get bot message ID: ", err)
							return
						}

						// add reaction to the bot message based on the length of anime_list
						for _, emoji := range reactions[:len(manga_list)] {
							err = s.MessageReactionAdd(m.ChannelID, botMessageID, emoji)
							if err != nil {
								log.Error("Failed to add reaction: ", err)
								return
							}
						}

						var chosen_manga int

						startTimer := time.Now()

						for {
							passedTimer := time.Since(startTimer).Seconds()
							checkAuthorReactionOne, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "1️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionOne {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_manga = 0
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionTwo, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "2️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionTwo {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_manga = 1
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionThree, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "3️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionThree {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_manga = 2
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionFour, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "4️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionFour {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_manga = 3
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							checkAuthorReactionFive, err := checkMessageReactionAuthor(s, m.ChannelID, botMessageID, "5️⃣", m.Author.ID, 10)
							if err != nil {
								log.Error(err)
								return
							}
							if checkAuthorReactionFive {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								chosen_manga = 4
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								break
							}

							// if no reactions is added then just remove reactions from the message.
							if passedTimer >= 10 {
								err = s.MessageReactionsRemoveAll(m.ChannelID, botMessageID)
								if err != nil {
									log.Error("Failed to remove reactions from bot message: ", err)
									return
								}
								s.ChannelMessageDelete(m.ChannelID, botMessageID)
								return
							}
						}

						genres_list := manga_list[chosen_manga].Genres
						var genres []string
						for _, genre := range genres_list {
							genres = append(genres, genre.Name)
						}

						genres_join := strings.Join(genres, ",")
						if genres_join == "" {
							genres_join = "\u200b"
						}
						var nsfw_string string
						if manga_list[chosen_manga].Nsfw == "white" {
							nsfw_string = "white (This work is safe for work)"
						} else if manga_list[chosen_manga].Nsfw == "" {
							nsfw_string = "Unknown"
						} else if manga_list[chosen_manga].Nsfw == "gray" {
							nsfw_string = "grey (This work may be not safe for work)"
						} else {
							nsfw_string = "black (This work is not safe for work)"
						}

						// start embed
						embed = NewEmbed().
							SetTitle(manga_list[chosen_manga].Title).
							SetURL("https://myanimelist.net/manga/"+strconv.Itoa(manga_list[chosen_manga].ID)).
							SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
							SetThumbnail(manga_list[chosen_manga].MainPicture.Large).
							SetDescription(cutDescription(manga_list[chosen_manga].Synopsis[:])).
							AddField("Type", manga_list[chosen_manga].MediaType).
							AddField("Volumes", strconv.Itoa(manga_list[chosen_manga].NumVolumes)).
							AddField("Chapters", strconv.Itoa(manga_list[chosen_manga].NumChapters)).
							AddField("Status", manga_list[chosen_manga].Status).
							AddField("Published", manga_list[chosen_manga].StartDate).
							AddField("Mean Score", strconv.FormatFloat(manga_list[chosen_manga].Mean, 'f', 3, 64)).
							//AddField("Popularity", strconv.Itoa(manga_list[chosen_manga].Popularity)).
							//AddField("Rank", strconv.Itoa(manga_list[chosen_manga].Rank)).
							AddField("Genres", genres_join).
							AddField("NSFW", nsfw_string).
							SetFooter(manga_list[chosen_manga].AlternativeTitles.Ja, manga_list[chosen_manga].MainPicture.Medium).
							InlineAllFields().
							SetColor(green).MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

					}
				} else {
					manga, _, err := c.Manga.Details(ctx, manga_id,
						mal.Fields{"id", "title", "main_picture", "alternative_titles", "start_date", "end_date", "synopsis", "mean", "rank", "popularity", "num_list_users", "num_scoring_users", "nsfw", "genres", "media_type", "status", "num_volumes", "num_chapters", "authors"},
					)
					if err != nil {
						log.Error("Failed to filter media by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Anime not found!\n maybe try using Title?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					genres_list := manga.Genres
					var genres []string
					for _, genre := range genres_list {
						genres = append(genres, genre.Name)
					}

					genres_join := strings.Join(genres, ",")
					if genres_join == "" {
						genres_join = "\u200b"
					}

					var nsfw_string string
					if manga.Nsfw == "white" {
						nsfw_string = "white (This work is safe for work)"
					} else if manga.Nsfw == "" {
						nsfw_string = "Unknown"
					} else if manga.Nsfw == "gray" {
						nsfw_string = "grey (This work may be not safe for work)"
					} else {
						nsfw_string = "black (This work is not safe for work)"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(manga.Title).
						SetURL("https://myanimelist.net/manga/"+strconv.Itoa(manga.ID)).
						SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
						SetThumbnail(manga.MainPicture.Large).
						SetDescription(cutDescription(manga.Synopsis[:])).
						AddField("Type", manga.MediaType).
						AddField("Volumes", strconv.Itoa(manga.NumVolumes)).
						AddField("Chapters", strconv.Itoa(manga.NumChapters)).
						AddField("Status", manga.Status).
						AddField("Published", manga.StartDate).
						AddField("Mean Score", strconv.FormatFloat(manga.Mean, 'f', 3, 64)).
						//AddField("Popularity", strconv.Itoa(manga.Popularity)).
						//AddField("Rank", strconv.Itoa(manga.Rank)).
						AddField("Genres", genres_join).
						AddField("NSFW", nsfw_string).
						SetFooter(manga.AlternativeTitles.Ja, manga.MainPicture.Medium).
						InlineAllFields().
						SetColor(green).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}
				}
			}
		}
	}
}
