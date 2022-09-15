package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	alCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "anilist-anime",
			Description: "Search Anilist by anime title or id",
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
					Description: "Enter the anime id from anilist.",
					Required:    false,
				},
			},
		},

		{
			Name:        "anilist-manga",
			Description: "Search Anilist by manga title or id",
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
					Description: "Enter the manga id from anilist.",
					Required:    false,
				},
			},
		},

		{
			Name:        "anilist-character",
			Description: "Search Anilist by character name or id",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Enter the character name.",
					Required:    false,
				},
				// optional
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Enter the character id from anilist.",
					Required:    false,
				},
			},
		},

		{
			Name:        "anilist-staff",
			Description: "Search Anilist by staff name or id",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Enter the staff name.",
					Required:    false,
				},
				// optional
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Enter the staff id from anilist.",
					Required:    false,
				},
			},
		},

		{
			Name:        "anilist-users",
			Description: "Search users on anilist by their username or id",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Enter the username.",
					Required:    false,
				},
				// optional
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "id",
					Description: "Enter the user id from anilist.",
					Required:    false,
				},
			},
		},
	}

	anilistCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"anilist-anime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				log.Debug("Entered String for Anime Anilist: ", option.StringValue())
				anilistAnime, err := anilistQueryByTitle(option.StringValue())
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
						Embeds: anilistAnime,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer for Anime Anilist: ", option.StringValue())
				anilistAnime, err := anilistQueryAnimeByID(int(option.IntValue()))
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
						Embeds: anilistAnime,
					},
				})
			}
		},

		"anilist-manga": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				log.Debug("Entered String for Manga: ", option.StringValue())
				anilistManga, err := anilistQueryMangaByTitle(option.StringValue())
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
						Embeds: anilistManga,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer For Manga:", option.IntValue())
				anilistMangaID, err := anilistQueryMangaByID(option.IntValue())
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
						Embeds: anilistMangaID,
					},
				})
			}
		},

		"anilist-character": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			margs := make([]interface{}, 0, len(options))

			if option, ok := optionMap["name"]; ok {
				margs = append(margs, option.StringValue())
				log.Debug("Entered String for character: ", option.StringValue())
				anilistCharacter, err := anilistQueryCharacterByName(option.StringValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get character By Title, Maybe try using ID?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistCharacter,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer For character:", option.IntValue())
				anilistCharacterID, err := anilistQueryCharacterByID(option.IntValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get character By ID, Maybe try using name?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistCharacterID,
					},
				})
			}
		},

		"anilist-staff": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			margs := make([]interface{}, 0, len(options))

			if option, ok := optionMap["name"]; ok {
				margs = append(margs, option.StringValue())
				log.Debug("Entered String for staff: ", option.StringValue())
				anilistStaff, err := anilistQueryStaffByName(option.StringValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get staff By Name, Maybe try using ID?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistStaff,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer For staff:", option.IntValue())
				anilistStaffID, err := anilistQueryStaffByID(option.IntValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get staff By ID, Maybe try using name?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistStaffID,
					},
				})
			}
		},

		"anilist-users": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			margs := make([]interface{}, 0, len(options))

			if option, ok := optionMap["username"]; ok {
				margs = append(margs, option.StringValue())
				log.Debug("Entered String for username: ", option.StringValue())
				anilistUserName, err := anilistQueryUserByUserName(option.StringValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get user by username, Maybe try using ID or check the spelling agen?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistUserName,
					},
				})
			} else if option, ok := optionMap["id"]; ok {
				margs = append(margs, option.IntValue())
				log.Debug("Entered Integer For username:", option.IntValue())
				anilistUsersID, err := anilistQueryUserByID(option.IntValue())
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("%s", "Failed to get user By ID, Maybe try using username?"),
						},
					})
					return
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: anilistUsersID,
					},
				})
			}
		},
	}
)
