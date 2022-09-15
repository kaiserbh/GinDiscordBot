package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/anilistgo"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// anime query media from anilist by id or name.
func anilistQueryByTitle(title string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// query media by title
	anime := anilistgo.NewMediaQuery()
	_, err := anime.FilterAnimeByTitle(title)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Media not found!\n Maybe try using id?").
			SetColor(red).MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	// making sure color hex is not empty
	var animeColorHex int
	if anime.CoverImage.Color == "" {
		animeColorHex = green
	} else {
		animeColorHex, err = convertStringHexColorToInt(anime.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
	}

	averageScore := strconv.Itoa(anime.AverageScore) + "%"
	meanScore := strconv.Itoa(anime.MeanScore) + "%"
	//popularity := strconv.Itoa(anime.Popularity)

	genres := strings.Join(anime.Genres, ",")
	if genres == "" {
		genres = "\u200b"
	}

	animeStudios := anime.Studios.Edges
	var mainStudio string
	for _, studio := range animeStudios {
		if studio.IsMain {
			mainStudio = studio.Node.Name
			break
		}
	}

	if mainStudio == "" {
		mainStudio = "\u200b"
	}

	if anime.Season == "" {
		anime.Season = "\u200b"
	}

	if anime.Source == "" {
		anime.Source = "\u200b"
	}

	// making sure the title is not empty.
	var animeTitle string
	if anime.Title.English != "" {
		animeTitle = anime.Title.English
	} else if anime.Title.Romaji != "" {
		animeTitle = anime.Title.Romaji
	} else if anime.Title.Native != "" {
		animeTitle = anime.Title.Native
	} else {
		animeTitle = anime.Title.UserPreferred
	}

	description, startDate, endDate := anilistAnimeData(anime)

	// start embed
	embed := NewEmbed().
		SetTitle(animeTitle).
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
		//AddField("Season", anime.Season).
		AddField("Average Score", averageScore).
		AddField("Mean Score", meanScore).
		//AddField("Popularity", popularity).
		//AddField("Favourites", strconv.Itoa(anime.Favourites)).
		//AddField("Source", anime.Source).
		AddField("Genres", genres).
		AddField("Studio", mainStudio).
		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
		InlineAllFields().
		SetColor(animeColorHex).MessageEmbed

	arr = append(arr, embed)

	return arr, nil
}

func anilistQueryAnimeByID(id int) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	anime := anilistgo.NewMediaQuery()
	_, err := anime.FilterAnimeByID(id)
	if err != nil {
		log.Error("Failed to filter media by id: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Media not found!\n Maybe try using title?").
			SetColor(red).MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}
	// making sure color hex is not empty
	var animeColorHex int
	if anime.CoverImage.Color == "" {
		animeColorHex = green
	} else {
		animeColorHex, err = convertStringHexColorToInt(anime.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
	}

	averageScore := strconv.Itoa(anime.AverageScore) + "%"
	meanScore := strconv.Itoa(anime.MeanScore) + "%"
	//popularity := strconv.Itoa(anime.Popularity)

	genres := strings.Join(anime.Genres, ",")
	if genres == "" {
		genres = "\u200b"
	}

	animeStudios := anime.Studios.Edges
	var mainStudio string
	for _, studio := range animeStudios {
		if studio.IsMain {
			mainStudio = studio.Node.Name
			break
		}
	}

	if mainStudio == "" {
		mainStudio = "\u200b"
	}

	if anime.Season == "" {
		anime.Season = "\u200b"
	}

	if anime.Source == "" {
		anime.Source = "\u200b"
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
		//AddField("Season", anime.Season).
		AddField("Average Score", averageScore).
		AddField("Mean Score", meanScore).
		//AddField("Popularity", popularity).
		//AddField("Favourites", strconv.Itoa(anime.Favourites)).
		//AddField("Source", anime.Source).
		AddField("Genres", genres).
		AddField("Studio", mainStudio).
		SetFooter(anime.Title.Romaji, anime.CoverImage.ExtraLarge).
		InlineAllFields().
		SetColor(animeColorHex).MessageEmbed

	arr = append(arr, embed)

	return arr, nil
}

func anilistQueryMangaByTitle(title string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// query media by title
	manga := anilistgo.NewMediaQuery()
	_, err := manga.FilterMangaByTitle(title)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Media not found!\n Maybe try using id?").
			SetColor(red).MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	// making sure color hex is not empty
	var colorHex int
	if manga.CoverImage.Color == "" {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(manga.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
	}

	genres := strings.Join(manga.Genres, ",")
	if genres == "" {
		genres = "\u200b"
	}

	description, startDate, endDate := anilistAnimeData(manga)

	if manga.Source == "" {
		manga.Source = "\u200b"
	}

	// making sure the title is not empty.
	var mangaTitle string
	if manga.Title.English != "" {
		mangaTitle = manga.Title.English
	} else if manga.Title.Romaji != "" {
		mangaTitle = manga.Title.Romaji
	} else if manga.Title.Native != "" {
		mangaTitle = manga.Title.Native
	} else {
		mangaTitle = manga.Title.UserPreferred
	}

	// start embed
	embed := NewEmbed().
		SetTitle(mangaTitle).
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
		//AddField("Popularity", strconv.Itoa(manga.Popularity)).
		//AddField("Favourites", strconv.Itoa(manga.Favourites)).
		//AddField("Source", manga.Source).
		AddField("Genres", genres).
		SetFooter(manga.Title.Romaji, manga.CoverImage.ExtraLarge).
		InlineAllFields().
		SetColor(colorHex).MessageEmbed

	arr = append(arr, embed)

	return arr, nil
}

func anilistQueryMangaByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed

	manga := anilistgo.NewMediaQuery()
	_, err := manga.FilterMangaByID(int(id))
	if err != nil {
		log.Error("Failed to filter manga by ID: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Manga not found!\n Maybe try using title?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	var colorHex int
	if manga.CoverImage.Color == "" {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(manga.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
	}

	genres := strings.Join(manga.Genres, ",")
	if genres == "" {
		genres = "\u200b"
	}

	description, startDate, endDate := anilistAnimeData(manga)

	if manga.Source == "" {
		manga.Source = "\u200b"
	}

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
		//AddField("Popularity", strconv.Itoa(manga.Popularity)).
		//AddField("Favourites", strconv.Itoa(manga.Favourites)).
		//AddField("Source", manga.Source).
		AddField("Genres", genres).
		SetFooter(manga.Title.Romaji, manga.CoverImage.ExtraLarge).
		InlineAllFields().
		SetColor(colorHex).MessageEmbed

	arr = append(arr, embed)

	return arr, err
}

// character query media from anilist by id or name.
func anilistQueryCharacterByName(name string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	character := anilistgo.NewCharacterQuery()
	_, err := character.FilterCharacterByName(name)
	if err != nil {
		log.Error("Failed to filter character by Name: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Character not found!\n Maybe try using id?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	// making sure color hex is not empty
	var colorHex int
	if character.Media.Nodes[0].CoverImage.Color == "" {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(character.Media.Nodes[0].CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			colorHex = green
		}
	}

	characterStartMonth := strconv.Itoa(character.DateOfBirth.Month)
	characterStartDay := strconv.Itoa(character.DateOfBirth.Day)
	convMonth := convMonthIntToStr(characterStartMonth) + " "
	dateOfBirth := convMonth + characterStartDay

	description := character.Description
	description = cutDescription(description)

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

	arr = append(arr, embed)

	return arr, nil
}

func anilistQueryCharacterByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed

	character := anilistgo.NewCharacterQuery()
	_, err := character.FilterCharacterID(int(id))
	if err != nil {
		log.Error("Failed to filter character by Name: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Character not found!\n Maybe try using name?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	// making sure color hex is not empty
	var colorHex int
	if character.Media.Nodes[0].CoverImage.Color == "" {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(character.Media.Nodes[0].CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			colorHex = green
		}
	}

	characterStartMonth := strconv.Itoa(character.DateOfBirth.Month)
	characterStartDay := strconv.Itoa(character.DateOfBirth.Day)
	convMonth := convMonthIntToStr(characterStartMonth) + " "
	dateOfBirth := convMonth + characterStartDay

	description := character.Description

	description = cutDescription(description)

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

	arr = append(arr, embed)

	return arr, nil
}

// staff query staff from anilist by id or name.
func anilistQueryStaffByName(name string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	staff := anilistgo.NewStaffQuery()
	_, err := staff.FilterStaffByName(name)
	if err != nil {
		log.Error("Failed to filter staff by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Staff not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	var colorHex int

	if len(staff.StaffMedia.Nodes) == 0 {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(staff.StaffMedia.Nodes[0].CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
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

	// if staff HomeTown is empty then just send empty string.
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

	arr = append(arr, embed)

	return arr, nil
}

func anilistQueryStaffByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// query media by title
	staff := anilistgo.NewStaffQuery()
	_, err := staff.FilterStaffByID(int(id))
	if err != nil {
		log.Error("Failed to filter staff by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Staff not found!\n maybe try using name?").
			SetColor(red).MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	var colorHex int
	if len(staff.StaffMedia.Nodes) == 0 {
		colorHex = green
	} else {
		colorHex, err = convertStringHexColorToInt(staff.StaffMedia.Nodes[0].CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			return arr, err
		}
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

	// if staff  HomeTown is empty then just send empty string.
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

	arr = append(arr, embed)

	return arr, nil
}

// user query user from anilist by id or name.
func anilistQueryUserByUserName(username string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// query media by title
	user := anilistgo.NewUserQuery()
	_, err := user.FilterUserByName(username)
	if err != nil {
		log.Error("Failed to filter user by user name: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("User not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	var colorHex int
	var userAnimeFavourites []string
	var joinedAnimeFav string

	if len(user.Favourites.Anime.Edges) > 0 {
		colorHex, err = convertStringHexColorToInt(user.Favourites.Anime.Edges[0].Node.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: setting default colour to green: ", err)
			colorHex = green
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

	// join date seconds to unix time.
	createdAt := user.CreatedAt
	var joinedDate time.Time
	var joinDateString string

	if createdAt != 0 {
		joinedDate = time.Unix(int64(createdAt), 0).UTC()
		joinDateString = joinedDate.Format(time.RFC822)
		split := strings.Split(joinDateString, " ")
		joinDateString = strings.Join(split[0:3], "-") + "\n" + strings.Join(split[3:], " ")
	} else {
		joinDateString = "0"
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
		AddField("Created At", joinDateString).
		AddField("Anime Favourites", joinedAnimeFav).
		SetFooter(user.Name, user.Avatar.Large).
		InlineAllFields().
		SetColor(colorHex).MessageEmbed

	arr = append(arr, embed)
	return arr, nil
}

func anilistQueryUserByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// query media by title
	user := anilistgo.NewUserQuery()
	_, err := user.FilterUserByID(int(id))
	if err != nil {
		log.Error("Failed to filter user by ID: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("User not found!\n maybe try using username?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	var colorHex int
	var userAnimeFavourites []string
	var joinedAnimeFav string

	if len(user.Favourites.Anime.Edges) > 0 {
		colorHex, err = convertStringHexColorToInt(user.Favourites.Anime.Edges[0].Node.CoverImage.Color)
		if err != nil {
			log.Error("Failed to get media Color hex: ", err)
			colorHex = green
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

	// join date seconds to unix time.
	createdAt := user.CreatedAt
	var joinedDate time.Time
	var joinDateString string

	if createdAt != 0 {
		joinedDate = time.Unix(int64(createdAt), 0).UTC()
		joinDateString = joinedDate.Format(time.RFC822)
		split := strings.Split(joinDateString, " ")
		joinDateString = strings.Join(split[0:3], "-") + "\n" + strings.Join(split[3:], " ")
	} else {
		joinDateString = "0"
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
		AddField("Created At", joinDateString).
		AddField("Anime Favourites", joinedAnimeFav).
		SetFooter(user.Name, user.Avatar.Large).
		InlineAllFields().
		SetColor(colorHex).MessageEmbed

	arr = append(arr, embed)
	return arr, nil
}
