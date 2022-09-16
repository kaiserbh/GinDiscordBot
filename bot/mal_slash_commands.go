package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var animeTitle string
var mangaTitleMal string

var (
	malCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "mal-anime",
			Description: "Search anime from MyAnimeList by title or id",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Enter the anime title.",
					Required:    false,
				},
				// optional
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Enter the anime id from mal.",
					Required:    false,
				},
			},
		},

		{
			Name:        "mal-anime-choice",
			Description: "Enter the chosen anime from the search result: /mal-anime before using this.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "choice",
					Description: "Enter The Choice you would like. 1-5",
					Required:    true,
					MaxValue:    5,
					MaxLength:   1,
				},
			},
		},

		{
			Name:        "mal-manga",
			Description: "Search manga from MyAnimeList by title or id",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Enter the manga title.",
					Required:    false,
				},
				// optional
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Enter the manga id from mal.",
					Required:    false,
				},
			},
		},

		{
			Name:        "mal-manga-choice",
			Description: "Enter the chosen manga from the search result: /mal-manga before using this.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "choice",
					Description: "Enter The Choice you would like. 1-5",
					Required:    true,
					MaxValue:    5,
					MaxLength:   1,
				},
			},
		},
	}

	malCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mal-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			if option, ok := optionMap["title"]; ok {

				// hold the anime title temporary so that it can be used later on with choice.
				animeTitle = optionMap["title"].StringValue()
				margs = append(margs, option.StringValue())
				log.Debug("Entered String for Anime MAL: ", option.StringValue())
				malAnime, err := malQueryByAnimeTitle(option.StringValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get anime By Title, Maybe try using ID?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malAnime,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer for Anime MAL: ", option.IntValue())
				malAnimeID, err := malQueryAnimeByID(option.IntValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get anime By ID, Maybe try using Title?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malAnimeID,
					},
				})
			}
		},

		"mal-anime-choice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			if option, ok := optionMap["choice"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered integer for Anime choice MAL: ", option.IntValue())
				value := 0
				switch option.IntValue() {
				case 1:
					value = 0
					break
				case 2:
					value = 1
					break
				case 3:
					value = 2
					break
				case 4:
					value = 3
					break
				case 5:
					value = 4
					break
				}
				if animeTitle == "" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Please use /mal-anime before using this."),
						},
					})
					return
				}
				malAnime, err := malQueryByAnimeTitleAndChoice(animeTitle, int64(value))
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get anime By Title, Maybe try using ID?"),
						},
					})
					return
				}
				// reset animeTitle by now.
				animeTitle = ""
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malAnime,
					},
				})
			}
		},

		"mal-manga": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			margs := make([]interface{}, 0, len(options))

			if option, ok := optionMap["title"]; ok {

				margs = append(margs, option.StringValue())
				log.Debug("Entered String for Manga MAL: ", option.StringValue())
				malManga, err := malQueryMangaByTitle(option.StringValue())

				// set the manga title
				mangaTitleMal = option.StringValue()
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get manga By Title, Maybe try using ID?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malManga,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer For Manga:", option.IntValue())
				malMangaID, err := malQueryMangaByID(option.IntValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get manga By ID, Maybe try using Title?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malMangaID,
					},
				})
			}
		},

		"mal-manga-choice": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			margs := make([]interface{}, 0, len(options))
			if option, ok := optionMap["choice"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered integer for Manga choice MAL: ", option.IntValue())
				value := 0
				switch option.IntValue() {
				case 1:
					value = 0
					break
				case 2:
					value = 1
					break
				case 3:
					value = 2
					break
				case 4:
					value = 3
					break
				case 5:
					value = 4
					break
				}
				if mangaTitleMal == "" {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Please use /mal-manga before using this."),
						},
					})
					return
				}
				malManga, err := malQueryMangaByTitleAndChoice(mangaTitleMal, int64(value))
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get Manga By Title, Maybe try using ID?"),
						},
					})
					return
				}
				// reset animeTitle by now.
				mangaTitleMal = ""
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: malManga,
					},
				})
			}
		},
	}
)
