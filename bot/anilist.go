package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/anilistgo"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// anime query media from anilist by id or name.
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
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				animeQuery := args[1]

				// check if it's a number or a strings.
				animeID, err := strconv.Atoi(animeQuery)
				if err != nil {
					// query media by title
					anime := anilistgo.NewMediaQuery()
					// join multiple search query
					animeQuery = strings.Join(args[1:], " ")
					_, err := anime.FilterAnimeByTitle(animeQuery)
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

					animeColorHex, err := convertStringHexColorToInt(anime.CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					averageScore := strconv.Itoa(anime.AverageScore) + "%"
					meanScore := strconv.Itoa(anime.MeanScore) + "%"
					popularity := strconv.Itoa(anime.Popularity)

					genres := strings.Join(anime.Genres, ",")

					animeStudios := anime.Studios.Edges
					var mainStudio string
					for _, studio := range animeStudios {
						if studio.IsMain {
							mainStudio = studio.Node.Name
							break
						}
					}

					description, startDate, endDate := anilistAnimeData(anime)

					// start embed
					embed := NewEmbed().
						SetTitle(anime.Title.English).
						SetURL(anime.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(anime.BannerImage).
						SetThumbnail(anime.CoverImage.ExtraLarge).
						SetDescription(description).
						AddField("Format", anime.MediaFormat).
						AddField("Episodes", strconv.Itoa(anime.Episodes)).
						AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
						AddField("Status", anime.Status).
						AddField("Start Date", startDate).
						AddField("End Date", endDate).
						AddField("Season", anime.Season).
						AddField("Average Score", averageScore).
						AddField("Mean Score", meanScore).
						AddField("Popularity", popularity).
						AddField("Favourites", strconv.Itoa(anime.Favourites)).
						AddField("Source", anime.Source).
						AddField("Genres", genres).
						AddField("Studio", mainStudio).
						SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
						InlineAllFields().
						SetColor(animeColorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				} else {
					// query media by title
					anime := anilistgo.NewMediaQuery()
					_, err := anime.FilterAnimeByID(animeID)
					if err != nil {
						log.Error("Failed to filter media by ID: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Anime not found!\n maybe try using title?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						return
					}
					animeColorHex, err := convertStringHexColorToInt(anime.CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					averageScore := strconv.Itoa(anime.AverageScore) + "%"
					meanScore := strconv.Itoa(anime.MeanScore) + "%"
					popularity := strconv.Itoa(anime.Popularity)

					genres := strings.Join(anime.Genres, ",")

					animeStudios := anime.Studios.Edges
					var mainStudio string
					for _, studio := range animeStudios {
						if studio.IsMain {
							mainStudio = studio.Node.Name
							break
						}
					}

					description, startDate, endDate := anilistAnimeData(anime)

					// start embed
					embed := NewEmbed().
						SetTitle(anime.Title.English).
						SetURL(anime.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(anime.BannerImage).
						SetThumbnail(anime.CoverImage.ExtraLarge).
						SetDescription(description).
						AddField("Format", anime.MediaFormat).
						AddField("Episodes", strconv.Itoa(anime.Episodes)).
						AddField("Episode Duration", strconv.Itoa(anime.Duration)+" mins").
						AddField("Status", anime.Status).
						AddField("Start Date", startDate).
						AddField("End Date", endDate).
						AddField("Season", anime.Season).
						AddField("Average Score", averageScore).
						AddField("Mean Score", meanScore).
						AddField("Popularity", popularity).
						AddField("Favourites", strconv.Itoa(anime.Favourites)).
						AddField("Source", anime.Source).
						AddField("Genres", genres).
						AddField("Studio", mainStudio).
						SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
						InlineAllFields().
						SetColor(animeColorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				}
			}
		}
	}
}

// media query media from anilist by id or name.
func manga(s *discordgo.Session, m *discordgo.MessageCreate) {
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
			if mangaArgs == guild.GuildPrefix+"manga" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				mangaQuery := args[1]

				// check if it's a number or a strings.
				mangaID, err := strconv.Atoi(mangaQuery)
				if err != nil {
					// query media by title
					manga := anilistgo.NewMediaQuery()
					// join multiple search query
					mangaQuery = strings.Join(args[1:], " ")
					_, err := manga.FilterMangaByTitle(mangaQuery)
					if err != nil {
						log.Error("Failed to filter media by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Media not found!\n Maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						return
					}

					colorHex, err := convertStringHexColorToInt(manga.CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					genres := strings.Join(manga.Genres, ",")
					description, startDate, endDate := anilistAnimeData(manga)

					// start embed
					embed := NewEmbed().
						SetTitle(manga.Title.English).
						SetURL(manga.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(manga.BannerImage).
						SetThumbnail(manga.CoverImage.ExtraLarge).
						SetDescription(description).
						AddField("Format", manga.MediaFormat).
						AddField("Volumes", strconv.Itoa(manga.Volumes)).
						AddField("Status", manga.Status).
						AddField("Start Date", startDate).
						AddField("End Date", endDate).
						AddField("Average Score", strconv.Itoa(manga.AverageScore)+"%").
						AddField("Mean Score", strconv.Itoa(manga.MeanScore)+"%").
						AddField("Popularity", strconv.Itoa(manga.Popularity)).
						AddField("Favourites", strconv.Itoa(manga.Favourites)).
						AddField("Source", manga.Source).
						AddField("Genres", genres).
						SetFooter(manga.Title.Romaji, manga.CoverImage.ExtraLarge).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				} else {
					// query media by title
					manga := anilistgo.NewMediaQuery()
					_, err := manga.FilterMangaByID(mangaID)
					if err != nil {
						log.Error("Failed to filter media by ID: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Media not found!\n Maybe try using title?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						return
					}

					colorHex, err := convertStringHexColorToInt(manga.CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					genres := strings.Join(manga.Genres, ",")
					description, startDate, endDate := anilistAnimeData(manga)

					// start embed
					embed := NewEmbed().
						SetTitle(manga.Title.English).
						SetURL(manga.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(manga.BannerImage).
						SetThumbnail(manga.CoverImage.ExtraLarge).
						SetDescription(description).
						AddField("Format", manga.MediaFormat).
						AddField("Volumes", strconv.Itoa(manga.Volumes)).
						AddField("Status", manga.Status).
						AddField("Start Date", startDate).
						AddField("End Date", endDate).
						AddField("Average Score", strconv.Itoa(manga.AverageScore)+"%").
						AddField("Mean Score", strconv.Itoa(manga.MeanScore)+"%").
						AddField("Popularity", strconv.Itoa(manga.Popularity)).
						AddField("Favourites", strconv.Itoa(manga.Favourites)).
						AddField("Source", manga.Source).
						AddField("Genres", genres).
						SetFooter(manga.Title.Romaji, manga.CoverImage.ExtraLarge).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				}
			}
		}
	}
}

// character query media from anilist by id or name.
func character(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	args := getArguments(messageContent)
	characterArgs := args[0]

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if characterArgs == guild.GuildPrefix+"character" || characterArgs == guild.GuildPrefix+"char" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't character args.")
					return
				}

				characterQuery := args[1]

				// check if it's a number or a strings.
				characterID, err := strconv.Atoi(characterQuery)
				if err != nil {
					// query media by title
					character := anilistgo.NewCharacterQuery()
					// join multiple search query
					characterQuery = strings.Join(args[1:], " ")
					_, err := character.FilterCharacterByName(characterQuery)
					if err != nil {
						log.Error("Failed to filter character by Name: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Character not found!\n Maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						return
					}

					colorHex, err := convertStringHexColorToInt(character.Media.Nodes[0].CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					characterStartMonth := strconv.Itoa(character.DateOfBirth.Month)
					characterStartDay := strconv.Itoa(character.DateOfBirth.Day)
					convMonth := convMonthIntToStr(characterStartMonth) + " "
					dateOfBirth := convMonth + characterStartDay

					description := character.Description

					checkSpoiler := strings.Replace(description, "~!", "||", -1)
					cleanEndSpoiler := strings.Replace(checkSpoiler, "!~", "||", -1)

					description = cutDescription(cleanEndSpoiler)

					if character.Age == "" {
						character.Age = "\u200b"
					}

					if character.Gender == "" {
						character.Gender = "\u200b"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(character.Name.Full).
						SetURL(character.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(character.Media.Nodes[0].BannerImage).
						SetThumbnail(character.Image.Large).
						SetDescription(description).
						AddField("Age", character.Age).
						AddField("Gender", character.Gender).
						AddField("Date Of Birth", dateOfBirth).
						AddField("Favorites", strconv.Itoa(character.Favourites)).
						SetFooter(character.Name.UserPreferred, character.Image.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				} else {
					// query media by title
					character := anilistgo.NewCharacterQuery()
					_, err := character.FilterCharacterID(characterID)
					if err != nil {
						log.Error("Failed to filter character by Name: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Character not found!\n Maybe try using name?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}

						return
					}

					colorHex, err := convertStringHexColorToInt(character.Media.Nodes[0].CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					characterStartMonth := strconv.Itoa(character.DateOfBirth.Month)
					characterStartDay := strconv.Itoa(character.DateOfBirth.Day)
					convMonth := convMonthIntToStr(characterStartMonth) + " "
					dateOfBirth := convMonth + characterStartDay

					description := character.Description

					checkSpoiler := strings.Replace(description, "~!", "||", -1)
					cleanEndSpoiler := strings.Replace(checkSpoiler, "!~", "||", -1)

					description = cutDescription(cleanEndSpoiler)

					if character.Age == "" {
						character.Age = "\u200b"
					}

					if character.Gender == "" {
						character.Gender = "\u200b"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(character.Name.Full).
						SetURL(character.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(character.Media.Nodes[0].BannerImage).
						SetThumbnail(character.Image.Large).
						SetDescription(description).
						AddField("Age", character.Age).
						AddField("Gender", character.Gender).
						AddField("Date Of Birth", dateOfBirth).
						AddField("Favorites", strconv.Itoa(character.Favourites)).
						SetFooter(character.Name.UserPreferred, character.Image.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				}
			}
		}
	}
}

// staff query staff from anilist by id or name.
func staff(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	args := getArguments(messageContent)
	staffArgs := args[0]

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if staffArgs == guild.GuildPrefix+"staff" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				staffQuery := args[1]

				// check if it's a number or a strings.
				staffID, err := strconv.Atoi(staffQuery)
				if err != nil {
					// query media by title
					staff := anilistgo.NewStaffQuery()
					// join multiple search query
					staffQuery = strings.Join(args[1:], " ")
					_, err := staff.FilterStaffByName(staffQuery)
					if err != nil {
						log.Error("Failed to filter staff by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Staff not found!\n maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					colorHex, err := convertStringHexColorToInt(staff.StaffMedia.Nodes[0].CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					// start date
					staffBirthMonth := strconv.Itoa(staff.DateOfBirth.Month)
					staffBirthDay := strconv.Itoa(staff.DateOfBirth.Day) + ","
					staffBirthYear := strconv.Itoa(staff.DateOfBirth.Year)
					animeStartMonthString := convMonthIntToStr(staffBirthMonth) + " "
					dateOfBirth := animeStartMonthString + staffBirthDay + staffBirthYear

					// end date
					staffDeathMonth := strconv.Itoa(staff.DateOfDeath.Month)
					staffDeathDay := strconv.Itoa(staff.DateOfDeath.Day) + ","
					staffDeathYear := strconv.Itoa(staff.DateOfDeath.Year)
					animeEndMonthString := convMonthIntToStr(staffDeathMonth) + " "
					dateOfDeath := animeEndMonthString + staffDeathDay + staffDeathYear

					// check if years active length.
					var yearsActive string
					if len(staff.YearsActive) == 0 {
						yearsActive = "\u200b"
					} else if len(staff.YearsActive) == 1 {
						yearsActive = strconv.Itoa(staff.YearsActive[0]) + " - "
					} else {
						yearsActive = strconv.Itoa(staff.YearsActive[0]) + " - " + strconv.Itoa(staff.YearsActive[1])
					}

					var primaryOccupation string
					if len(staff.PrimaryOccupations) <= 0 {
						primaryOccupation = "\u200b"
					} else {
						primaryOccupation = strings.Join(staff.PrimaryOccupations, ", ")
					}

					// if staff.HomeTown is empty then just send empty string.
					if staff.HomeTown == "" {
						staff.HomeTown = "\u200b"
					}

					if staff.Gender == "" {
						staff.Gender = "\u200b"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(staff.Name.Full).
						SetURL(staff.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(staff.StaffMedia.Nodes[0].BannerImage).
						SetThumbnail(staff.Image.Large).
						SetDescription(staff.Description).
						AddField("Language", staff.LanguageV2).
						AddField("Primary Occupations", primaryOccupation).
						AddField("Gender", staff.Gender).
						AddField("Date Of Birth", dateOfBirth).
						AddField("Date Of Death", dateOfDeath).
						AddField("Age", strconv.Itoa(staff.Age)).
						AddField("Years Active", yearsActive).
						AddField("Hometown", staff.HomeTown).
						AddField("Favourites", strconv.Itoa(staff.Favourites)).
						SetFooter(staff.Name.UserPreferred, staff.Image.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				} else {
					// query media by title
					staff := anilistgo.NewStaffQuery()
					_, err := staff.FilterStaffByID(staffID)
					if err != nil {
						log.Error("Failed to filter staff by title: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("Staff not found!\n maybe try using name?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					colorHex, err := convertStringHexColorToInt(staff.StaffMedia.Nodes[0].CoverImage.Color)
					if err != nil {
						log.Error("Failed to get media Color hex: ", err)
						return
					}

					// start date
					staffBirthMonth := strconv.Itoa(staff.DateOfBirth.Month)
					staffBirthDay := strconv.Itoa(staff.DateOfBirth.Day) + ","
					staffBirthYear := strconv.Itoa(staff.DateOfBirth.Year)
					animeStartMonthString := convMonthIntToStr(staffBirthMonth) + " "
					dateOfBirth := animeStartMonthString + staffBirthDay + staffBirthYear

					// end date
					staffDeathMonth := strconv.Itoa(staff.DateOfDeath.Month)
					staffDeathDay := strconv.Itoa(staff.DateOfDeath.Day) + ","
					staffDeathYear := strconv.Itoa(staff.DateOfDeath.Year)
					animeEndMonthString := convMonthIntToStr(staffDeathMonth) + " "
					dateOfDeath := animeEndMonthString + staffDeathDay + staffDeathYear

					// check if years active length.
					var yearsActive string
					if len(staff.YearsActive) == 0 {
						yearsActive = "\u200b"
					} else if len(staff.YearsActive) == 1 {
						yearsActive = strconv.Itoa(staff.YearsActive[0]) + " - "
					} else {
						yearsActive = strconv.Itoa(staff.YearsActive[0]) + " - " + strconv.Itoa(staff.YearsActive[1])
					}

					var primaryOccupation string
					if len(staff.PrimaryOccupations) <= 0 {
						primaryOccupation = "\u200b"
					} else {
						primaryOccupation = strings.Join(staff.PrimaryOccupations, ", ")
					}

					// if staff.HomeTown is empty then just send empty string.
					if staff.HomeTown == "" {
						staff.HomeTown = "\u200b"
					}

					if staff.Gender == "" {
						staff.Gender = "\u200b"
					}

					// start embed
					embed := NewEmbed().
						SetTitle(staff.Name.Full).
						SetURL(staff.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(staff.StaffMedia.Nodes[0].BannerImage).
						SetThumbnail(staff.Image.Large).
						SetDescription(staff.Description).
						AddField("Language", staff.LanguageV2).
						AddField("Primary Occupations", primaryOccupation).
						AddField("Gender", staff.Gender).
						AddField("Date Of Birth", dateOfBirth).
						AddField("Date Of Death", dateOfDeath).
						AddField("Age", strconv.Itoa(staff.Age)).
						AddField("Years Active", yearsActive).
						AddField("Hometown", staff.HomeTown).
						AddField("Favourites", strconv.Itoa(staff.Favourites)).
						SetFooter(staff.Name.UserPreferred, staff.Image.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				}
			}
		}
	}
}

// user query user from anilist by id or name.
func user(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checks if the message has prefix from the database file.
	guild, err := db.FindGuildByID(m.GuildID)
	if err != nil {
		log.Error("Finding Guild: ", err)
		return
	}
	messageContent := strings.ToLower(m.Content)
	args := getArguments(messageContent)
	staffArgs := args[0]

	if strings.HasPrefix(messageContent, guild.GuildPrefix) {
		// check if the channel is bot channel or allowed channel.
		allowedChannels := checkAllowedChannel(m.ChannelID, guild)
		if allowedChannels {
			if m.Author.ID == s.State.User.ID {
				return
			}
			if staffArgs == guild.GuildPrefix+"user" {
				// reset messageContent and arguments
				messageContent = m.Content
				args = getArguments(messageContent)

				if len(args) < 2 {
					log.Info("Doing nothing since there wasn't media args.")
					return
				}

				userQuery := args[1]

				// check if it's a number or a strings.
				userID, err := strconv.Atoi(userQuery)
				if err != nil {
					// query media by title
					user := anilistgo.NewUserQuery()
					// join multiple search query
					userQuery = strings.Join(args[1:], " ")
					_, err := user.FilterUserByName(userQuery)
					if err != nil {
						log.Error("Failed to filter user by user name: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("User not found!\n maybe try using id?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					var colorHex int
					var userAnimeFavourites []string
					var joinedAnimeFav string

					if len(user.Favourites.Anime.Edges) > 0 {
						colorHex, err = convertStringHexColorToInt(user.Favourites.Anime.Edges[0].Node.CoverImage.Color)
						if err != nil {
							log.Error("Failed to get media Color hex: ", err)
							return
						}

						for _, anime := range user.Favourites.Anime.Edges {
							userAnimeFavourites = append(userAnimeFavourites, anime.Node.Title.English)
						}
					} else {
						colorHex = green
						userAnimeFavourites = append(userAnimeFavourites, "\u200b")
					}

					// check the length of animeFavourites

					if len(userAnimeFavourites) >= 4 {
						joinedAnimeFav = strings.Join(userAnimeFavourites[0:4], "\n")
					} else {
						joinedAnimeFav = strings.Join(userAnimeFavourites, "\n")
					}

					// start embed
					embed := NewEmbed().
						SetTitle(user.Name).
						SetURL(user.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(user.BannerImage).
						SetThumbnail(user.Avatar.Large).
						SetDescription(user.About).
						AddField("ID", strconv.FormatInt(user.ID, 10)).
						AddField("Total Anime", strconv.Itoa(user.Statistics.Anime.Count)).
						AddField("Days Watched", strconv.Itoa(user.Statistics.Anime.MinutesWatched/1440)).
						AddField("Mean Score Anime", strconv.FormatFloat(user.Statistics.Anime.MeanScore, 'f', 1, 64)).
						AddField("Total Manga", strconv.Itoa(user.Statistics.Manga.Count)).
						AddField("Chapters Read", strconv.Itoa(user.Statistics.Manga.ChaptersRead)).
						AddField("Mean Score Manga", strconv.FormatFloat(user.Statistics.Manga.MeanScore, 'f', 1, 64)).
						AddField("Created At", strconv.Itoa(user.CreatedAt)).
						AddField("Anime Favourites", joinedAnimeFav).
						SetFooter(user.Name, user.Avatar.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				} else {
					// query media by title
					user := anilistgo.NewUserQuery()
					_, err := user.FilterUserByID(userID)
					if err != nil {
						log.Error("Failed to filter user by ID: ", err)
						// start embed
						embed := NewEmbed().
							SetDescription("User not found!\n maybe try using username?").
							SetColor(red).
							MessageEmbed

						_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
						if err != nil {
							log.Error("Failed to send embed to the channel: ", err)
							return
						}
						return
					}

					var colorHex int
					var userAnimeFavourites []string
					var joinedAnimeFav string

					if len(user.Favourites.Anime.Edges) > 0 {
						colorHex, err = convertStringHexColorToInt(user.Favourites.Anime.Edges[0].Node.CoverImage.Color)
						if err != nil {
							log.Error("Failed to get media Color hex: ", err)
							return
						}

						for _, anime := range user.Favourites.Anime.Edges {
							userAnimeFavourites = append(userAnimeFavourites, anime.Node.Title.English)
						}
					} else {
						colorHex = green
						userAnimeFavourites = append(userAnimeFavourites, "\u200b")
					}

					// check the length of animeFavourites

					if len(userAnimeFavourites) >= 4 {
						joinedAnimeFav = strings.Join(userAnimeFavourites[0:4], "\n")
					} else {
						joinedAnimeFav = strings.Join(userAnimeFavourites, "\n")
					}

					// start embed
					embed := NewEmbed().
						SetTitle(user.Name).
						SetURL(user.SiteURL).
						SetAuthor("Anilist", "https://anilist.co/img/logo_al.png").
						SetImage(user.BannerImage).
						SetThumbnail(user.Avatar.Large).
						SetDescription(user.About).
						AddField("ID", strconv.FormatInt(user.ID, 10)).
						AddField("Total Anime", strconv.Itoa(user.Statistics.Anime.Count)).
						AddField("Days Watched", strconv.Itoa(user.Statistics.Anime.MinutesWatched/1440)).
						AddField("Mean Score Anime", strconv.FormatFloat(user.Statistics.Anime.MeanScore, 'f', 1, 64)).
						AddField("Total Manga", strconv.Itoa(user.Statistics.Manga.Count)).
						AddField("Chapters Read", strconv.Itoa(user.Statistics.Manga.ChaptersRead)).
						AddField("Mean Score Manga", strconv.FormatFloat(user.Statistics.Manga.MeanScore, 'f', 1, 64)).
						AddField("Created At", strconv.Itoa(user.CreatedAt)).
						AddField("Anime Favourites", joinedAnimeFav).
						SetFooter(user.Name, user.Avatar.Large).
						InlineAllFields().
						SetColor(colorHex).MessageEmbed

					_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
					if err != nil {
						log.Error("Failed to send embed to the channel: ", err)
						return
					}

					botMessageID, err := getBotMessageID(s, m)
					if err != nil {
						log.Error("Failed to get bot message ID: ", err)
						return
					}
					// check timer and reaction
					checkAnilistTimer(s, m.ChannelID, botMessageID, m.Author.ID)
				}
			}
		}
	}
}
