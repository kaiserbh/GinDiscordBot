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

// anime query media from anilist by id or name.
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
				_, err := strconv.Atoi(animeQuery)
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

					//TODO:Check if there's more than one result.
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
							SetTitle("More than one result found").
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
								break
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
							SetAuthor("MAL", "https://cdn.myanimelist.net/s/common/uploaded_files/1455540405-b2e2bf20e11b68631e8439b44d9a51c7.png").
							SetImage(anime_list[chosen_anime].Background).
							SetThumbnail(anime_list[chosen_anime].MainPicture.Large).
							SetDescription(anime_list[chosen_anime].Synopsis).
							AddField("Format", anime_list[chosen_anime].MediaType).
							AddField("Episodes", strconv.Itoa(anime_list[chosen_anime].NumEpisodes)).
							AddField("Episode Duration", strconv.Itoa(anime_list[chosen_anime].AverageEpisodeDuration)+" seconds").
							AddField("Status", anime_list[chosen_anime].Status).
							AddField("Start Date", anime_list[chosen_anime].StartDate).
							AddField("End Date", anime_list[chosen_anime].EndDate).
							AddField("Season", anime_list[chosen_anime].StartSeason.Season).
							AddField("Mean Score", strconv.FormatFloat(anime_list[chosen_anime].Mean, 'f', 3, 64)).
							AddField("Popularity", strconv.Itoa(anime_list[chosen_anime].Popularity)).
							AddField("Rank", strconv.Itoa(anime_list[chosen_anime].Rank)).
							AddField("Source", anime_list[chosen_anime].Source).
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

						botMessageID, err = getBotMessageID(s, m)
						if err != nil {
							log.Error("Failed to get bot message ID: ", err)
							return
						}
						// check timer and reaction
						checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)

					}

					//TODO:Check if there's no result.
					//TODO:If there is more than one result, send a list of results.
					//TODO:Get bot message ID
					//TODO:Send reaction if there is more than one result.
					//TODO:Check if the reaction is used to select the result.
					//TODO:Remove reaction within set time.
				}

				// 	// making sure color hex is not empty
				// 	var animeColorHex int
				// 	if anime.CoverImage.Color == "" {
				// 		animeColorHex = green
				// 	} else {
				// 		animeColorHex, err = convertStringHexColorToInt(anime.CoverImage.Color)
				// 		if err != nil {
				// 			log.Error("Failed to get media Color hex: ", err)
				// 			return
				// 		}
				// 	}

				// 	averageScore := strconv.Itoa(anime.AverageScore) + "%"
				// 	meanScore := strconv.Itoa(anime.MeanScore) + "%"
				// 	popularity := strconv.Itoa(anime.Popularity)

				// 	genres := strings.Join(anime.Genres, ",")
				// 	if genres == "" {
				// 		genres = "\u200b"
				// 	}

				// 	animeStudios := anime.Studios.Edges
				// 	var mainStudio string
				// 	for _, studio := range animeStudios {
				// 		if studio.IsMain {
				// 			mainStudio = studio.Node.Name
				// 			break
				// 		}
				// 	}

				// 	if mainStudio == "" {
				// 		mainStudio = "\u200b"
				// 	}

				// 	if anime.Season == "" {
				// 		anime.Season = "\u200b"
				// 	}

				// 	if anime.Source == "" {
				// 		anime.Source = "\u200b"
				// 	}

				// 	// making sure the title is not empty.
				// 	var animeTitle string
				// 	if anime.Title.English != "" {
				// 		animeTitle = anime.Title.English
				// 	} else if anime.Title.Romaji != "" {
				// 		animeTitle = anime.Title.Romaji
				// 	} else if anime.Title.Native != "" {
				// 		animeTitle = anime.Title.Native
				// 	} else {
				// 		animeTitle = anime.Title.UserPreferred
				// 	}

				// 	description, startDate, endDate := anilistAnimeData(anime)

				// 	// start embed
				// 	embed := NewEmbed().
				// 		SetTitle(animeTitle).
				// 		SetURL(anime.SiteURL).
				// 		SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
				// 		SetImage(anime.BannerImage).
				// 		SetThumbnail(anime.CoverImage.ExtraLarge).
				// 		SetDescription(description).
				// 		AddField("Format", anime.MediaFormat).
				// 		AddField("Episodes", strconv.Itoa(anime.Episodes)).
				// 		AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
				// 		AddField("Status", anime.Status).
				// 		AddField("Start Date", startDate).
				// 		AddField("End Date", endDate).
				// 		AddField("Season", anime.Season).
				// 		AddField("Average Score", averageScore).
				// 		AddField("Mean Score", meanScore).
				// 		AddField("Popularity", popularity).
				// 		AddField("Favourites", strconv.Itoa(anime.Favourites)).
				// 		AddField("Source", anime.Source).
				// 		AddField("Genres", genres).
				// 		AddField("Studio", mainStudio).
				// 		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
				// 		InlineAllFields().
				// 		SetColor(animeColorHex).MessageEmbed

				// 	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 	if err != nil {
				// 		log.Error("Failed to send embed to the channel: ", err)
				// 		return
				// 	}

				// 	botMessageID, err := getBotMessageID(s, m)
				// 	if err != nil {
				// 		log.Error("Failed to get bot message ID: ", err)
				// 		return
				// 	}
				// 	// check timer and reaction
				// 	checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				// } else {
				// 	// query media by title
				// 	anime := anilistgo.NewMediaQuery()
				// 	_, err := anime.FilterAnimeByID(animeID)
				// 	if err != nil {
				// 		log.Error("Failed to filter media by ID: ", err)
				// 		// start embed
				// 		embed := NewEmbed().
				// 			SetDescription("Anime not found!\n maybe try using title?").
				// 			SetColor(red).
				// 			MessageEmbed

				// 		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 		if err != nil {
				// 			log.Error("Failed to send embed to the channel: ", err)
				// 			return
				// 		}

				// 		return
				// 	}
				// 	// making sure color hex is not empty
				// 	var animeColorHex int
				// 	if anime.CoverImage.Color == "" {
				// 		animeColorHex = green
				// 	} else {
				// 		animeColorHex, err	// making sure color hex is not empty
				// 	var animeColorHex int
				// 	if anime.CoverImage.Color == "" {
				// 		animeColorHex = green
				// 	} else {
				// 		animeColorHex, err = convertStringHexColorToInt(anime.CoverImage.Color)
				// 		if err != nil {
				// 			log.Error("Failed to get media Color hex: ", err)
				// 			return
				// 		}
				// 	}

				// 	averageScore := strconv.Itoa(anime.AverageScore) + "%"
				// 	meanScore := strconv.Itoa(anime.MeanScore) + "%"
				// 	popularity := strconv.Itoa(anime.Popularity)

				// 	genres := strings.Join(anime.Genres, ",")
				// 	if genres == "" {
				// 		genres = "\u200b"
				// 	}

				// 	animeStudios := anime.Studios.Edges
				// 	var mainStudio string
				// 	for _, studio := range animeStudios {
				// 		if studio.IsMain {
				// 			mainStudio = studio.Node.Name
				// 			break
				// 		}
				// 	}

				// 	if mainStudio == "" {
				// 		mainStudio = "\u200b"
				// 	}

				// 	if anime.Season == "" {
				// 		anime.Season = "\u200b"
				// 	}

				// 	if anime.Source == "" {
				// 		anime.Source = "\u200b"
				// 	}

				// 	// making sure the title is not empty.
				// 	var animeTitle string
				// 	if anime.Title.English != "" {
				// 		animeTitle = anime.Title.English
				// 	} else if anime.Title.Romaji != "" {
				// 		animeTitle = anime.Title.Romaji
				// 	} else if anime.Title.Native != "" {
				// 		animeTitle = anime.Title.Native
				// 	} else {
				// 		animeTitle = anime.Title.UserPreferred
				// 	}

				// 	description, startDate, endDate := anilistAnimeData(anime)

				// 	// start embed
				// 	embed := NewEmbed().
				// 		SetTitle(animeTitle).
				// 		SetURL(anime.SiteURL).
				// 		SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
				// 		SetImage(anime.BannerImage).
				// 		SetThumbnail(anime.CoverImage.ExtraLarge).
				// 		SetDescription(description).
				// 		AddField("Format", anime.MediaFormat).
				// 		AddField("Episodes", strconv.Itoa(anime.Episodes)).
				// 		AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
				// 		AddField("Status", anime.Status).
				// 		AddField("Start Date", startDate).
				// 		AddField("End Date", endDate).
				// 		AddField("Season", anime.Season).
				// 		AddField("Average Score", averageScore).
				// 		AddField("Mean Score", meanScore).
				// 		AddField("Popularity", popularity).
				// 		AddField("Favourites", strconv.Itoa(anime.Favourites)).
				// 		AddField("Source", anime.Source).
				// 		AddField("Genres", genres).
				// 		AddField("Studio", mainStudio).
				// 		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
				// 		InlineAllFields().
				// 		SetColor(animeColorHex).MessageEmbed

				// 	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 	if err != nil {
				// 		log.Error("Failed to send embed to the channel: ", err)
				// 		return
				// 	}

				// 	botMessageID, err := getBotMessageID(s, m)
				// 	if err != nil {
				// 		log.Error("Failed to get bot message ID: ", err)
				// 		return
				// 	}
				// 	// check timer and reaction
				// 	checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				// } else {
				// 	// query media by title
				// 	anime := anilistgo.NewMediaQuery()
				// 	_, err := anime.FilterAnimeByID(animeID)
				// 	if err != nil {
				// 		log.Error("Failed to filter media by ID: ", err)
				// 		// start embed
				// 		embed := NewEmbed().
				// 			SetDescription("Anime not found!\n maybe try using title?").
				// 			SetColor(red).
				// 			MessageEmbed

				// 		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 		if err != nil {
				// 			log.Error("Failed to send embed to the channel: ", err)
				// 			return
				// 		}

				// 		return
				// 	}
				// 	// making sure color hex is not empty
				// 	var animeColorHex int
				// 	if anime.CoverImage.Color == "" {
				// 		animeColorHex = green
				// 	} else {
				// 		animeColorHex, err = convertStringHexColorToInt(anime.CoverImage.Color)
				// 		if err != nil {
				// 			log.Error("Failed to get media Color hex: ", err)
				// 			return
				// 		}
				// 	}

				// 	averageScore := strconv.Itoa(anime.AverageScore) + "%"
				// 	meanScore := strconv.Itoa(anime.MeanScore) + "%"
				// 	popularity := strconv.Itoa(anime.Popularity)

				// 	genres := strings.Join(anime.Genres, ",")
				// 	if genres == "" {
				// 		genres = "\u200b"
				// 	}

				// 	animeStudios := anime.Studios.Edges
				// 	var mainStudio string
				// 	for _, studio := range animeStudios {
				// 		if studio.IsMain {
				// 			mainStudio = studio.Node.Name
				// 			break
				// 		}
				// 	}

				// 	if mainStudio == "" {
				// 		mainStudio = "\u200b"
				// 	}

				// 	if anime.Season == "" {
				// 		anime.Season = "\u200b"
				// 	}

				// 	if anime.Source == "" {
				// 		anime.Source = "\u200b"
				// 	}

				// 	description, startDate, endDate := anilistAnimeData(anime)

				// 	// start embed
				// 	embed := NewEmbed().
				// 		SetTitle(anime.Title.English).
				// 		SetURL(anime.SiteURL).
				// 		SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
				// 		SetImage(anime.BannerImage).
				// 		SetThumbnail(anime.CoverImage.ExtraLarge).
				// 		SetDescription(description).
				// 		AddField("Format", anime.MediaFormat).
				// 		AddField("Episodes", strconv.Itoa(anime.Episodes)).
				// 		AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
				// 		AddField("Status", anime.Status).
				// 		AddField("Start Date", startDate).
				// 		AddField("End Date", endDate).
				// 		AddField("Season", anime.Season).
				// 		AddField("Average Score", averageScore).
				// 		AddField("Mean Score", meanScore).
				// 		AddField("Popularity", popularity).
				// 		AddField("Favourites", strconv.Itoa(anime.Favourites)).
				// 		AddField("Source", anime.Source).
				// 		AddField("Genres", genres).
				// 		AddField("Studio", mainStudio).
				// 		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
				// 		InlineAllFields().
				// 		SetColor(animeColorHex).MessageEmbed

				// 	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 	if err != nil {
				// 		log.Error("Failed to send embed to the channel: ", err)
				// 		return
				// 	}

				// 	botMessageID, err := getBotMessageID(s, m)
				// 	if err != nil {
				// 		log.Error("Failed to get bot message ID: ", err)
				// 		return
				// 	}
				// 	// check timer and reaction
				// 	checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				// } = convertStringHexColorToInt(anime.CoverImage.Color)
				// 		if err != nil {
				// 			log.Error("Failed to get media Color hex: ", err)
				// 			return
				// 		}
				// 	}

				// 	averageScore := strconv.Itoa(anime.AverageScore) + "%"
				// 	meanScore := strconv.Itoa(anime.MeanScore) + "%"
				// 	popularity := strconv.Itoa(anime.Popularity)

				// 	genres := strings.Join(anime.Genres, ",")
				// 	if genres == "" {
				// 		genres = "\u200b"
				// 	}

				// 	animeStudios := anime.Studios.Edges
				// 	var mainStudio string
				// 	for _, studio := range animeStudios {
				// 		if studio.IsMain {
				// 			mainStudio = studio.Node.Name
				// 			break
				// 		}
				// 	}

				// 	if mainStudio == "" {
				// 		mainStudio = "\u200b"
				// 	}

				// 	if anime.Season == "" {
				// 		anime.Season = "\u200b"
				// 	}

				// 	if anime.Source == "" {
				// 		anime.Source = "\u200b"
				// 	}

				// 	description, startDate, endDate := anilistAnimeData(anime)

				// 	// start embed
				// 	embed := NewEmbed().
				// 		SetTitle(anime.Title.English).
				// 		SetURL(anime.SiteURL).
				// 		SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
				// 		SetImage(anime.BannerImage).
				// 		SetThumbnail(anime.CoverImage.ExtraLarge).
				// 		SetDescription(description).
				// 		AddField("Format", anime.MediaFormat).
				// 		AddField("Episodes", strconv.Itoa(anime.Episodes)).
				// 		AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
				// 		AddField("Status", anime.Status).
				// 		AddField("Start Date", startDate).
				// 		AddField("End Date", endDate).
				// 		AddField("Season", anime.Season).
				// 		AddField("Average Score", averageScore).
				// 		AddField("Mean Score", meanScore).
				// 		AddField("Popularity", popularity).
				// 		AddField("Favourites", strconv.Itoa(anime.Favourites)).
				// 		AddField("Source", anime.Source).
				// 		AddField("Genres", genres).
				// 		AddField("Studio", mainStudio).
				// 		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
				// 		InlineAllFields().
				// 		SetColor(animeColorHex).MessageEmbed

				// 	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
				// 	if err != nil {
				// 		log.Error("Failed to send embed to the channel: ", err)
				// 		return
				// 	}

				// 	botMessageID, err := getBotMessageID(s, m)
				// 	if err != nil {
				// 		log.Error("Failed to get bot message ID: ", err)
				// 		return
				// 	}
				// 	// check timer and reaction
				// 	checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				// }
			}
		}
	}
}
