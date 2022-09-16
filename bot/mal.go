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
func malQueryByAnimeTitle(title string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed

	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	anime_list, _, err := c.Anime.List(ctx, title,
		mal.Fields{
			"rank",
			"popularity",
			"synopsis",
			"mean",
			"background",
			"status",
			"source",
			"media_type",
			"num_episodes",
			"average_episode_duration",
			"start_date",
			"end_date",
			"alternative_titles",
			"start_season",
			"genres",
			"studios"},
		mal.Limit(5),
	)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Anime not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed

		arr = append(arr, embed)
		return arr, err
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
			SetDescription(strings.Join(animeList, "\n") + "Please use '\\'mal-anime-choice to choose.").SetColor(green).MessageEmbed

		arr = append(arr, embed)
		return arr, nil
	} else {
		var chosen_anime int
		chosen_anime = 0
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
		embed := NewEmbed().
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

		arr = append(arr, embed)
		return arr, nil
	}
}

// anime query media from MAL by id or name.
func malQueryByAnimeTitleAndChoice(title string, choice int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed

	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	anime_list, _, err := c.Anime.List(ctx, title,
		mal.Fields{
			"rank",
			"popularity",
			"synopsis",
			"mean",
			"background",
			"status",
			"source",
			"media_type",
			"num_episodes",
			"average_episode_duration",
			"start_date",
			"end_date",
			"alternative_titles",
			"start_season",
			"genres",
			"studios"},
		mal.Limit(5),
	)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Anime not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed

		arr = append(arr, embed)
		return arr, err
	}

	var chosen_anime int
	chosen_anime = int(choice)
	genresList := anime_list[chosen_anime].Genres
	var genres []string
	for _, genre := range genresList {
		genres = append(genres, genre.Name)
	}

	genresJoin := strings.Join(genres, ",")
	if genresJoin == "" {
		genresJoin = "\u200b"
	}

	studiosList := anime_list[chosen_anime].Studios
	var studios []string
	for _, studio := range studiosList {
		studios = append(studios, studio.Name)
	}

	studiosJoin := strings.Join(studios, ",")
	if studiosJoin == "" {
		studiosJoin = "\u200b"
	}

	if anime_list[chosen_anime].EndDate == "" {
		anime_list[chosen_anime].EndDate = "Ongoing"
	}

	// start embed
	embed := NewEmbed().
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
		AddField("Genres", genresJoin).
		AddField("Studio", studiosJoin).
		SetFooter(anime_list[chosen_anime].AlternativeTitles.Ja, anime_list[chosen_anime].MainPicture.Medium).
		InlineAllFields().
		SetColor(green).MessageEmbed

	arr = append(arr, embed)
	return arr, nil
}

func malQueryAnimeByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed

	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	animeId, _, err := c.Anime.Details(ctx, int(id),
		mal.Fields{"rank", "popularity", "synopsis", "mean", "background", "status", "source", "media_type", "num_episodes", "average_episode_duration", "start_date", "end_date", "alternative_titles", "start_season", "genres", "studios"},
	)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Anime not found!\n maybe try using Title?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	genresList := animeId.Genres
	var genres []string
	for _, genre := range genresList {
		genres = append(genres, genre.Name)
	}

	genresJoin := strings.Join(genres, ",")
	if genresJoin == "" {
		genresJoin = "\u200b"
	}

	studiosList := animeId.Studios
	var studios []string
	for _, studio := range studiosList {
		studios = append(studios, studio.Name)
	}

	studiosJoin := strings.Join(studios, ",")
	if studiosJoin == "" {
		studiosJoin = "\u200b"
	}

	if animeId.EndDate == "" {
		animeId.EndDate = "Ongoing"
	}

	// start embed
	embed := NewEmbed().
		SetTitle(animeId.Title).
		SetURL("https://myanimelist.net/anime/"+strconv.Itoa(animeId.ID)).
		SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
		SetThumbnail(animeId.MainPicture.Large).
		SetDescription(cutDescription(animeId.Synopsis)).
		AddField("Format", animeId.MediaType).
		AddField("Episodes", strconv.Itoa(animeId.NumEpisodes)).
		AddField("Episode Duration", strconv.Itoa(animeId.AverageEpisodeDuration/60)+" min").
		AddField("Status", animeId.Status).
		AddField("Start Date", animeId.StartDate).
		AddField("End Date", animeId.EndDate).
		//AddField("Season", anime_id.StartSeason.Season).
		AddField("Mean Score", strconv.FormatFloat(animeId.Mean, 'f', 3, 64)).
		//AddField("Popularity", strconv.Itoa(anime_id.Popularity)).
		//AddField("Rank", strconv.Itoa(anime_id.Rank)).
		//AddField("Source", anime_id.Source).
		AddField("Genres", genresJoin).
		AddField("Studio", studiosJoin).
		SetFooter(animeId.AlternativeTitles.Ja, animeId.MainPicture.Medium).
		InlineAllFields().
		SetColor(green).MessageEmbed
	arr = append(arr, embed)

	return arr, nil
}

// manga query media from MAL by id or name.
func malQueryMangaByTitle(title string) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	orignalMangaList, _, err := c.Manga.List(ctx, title,
		mal.Fields{
			"id",
			"title",
			"main_picture",
			"alternative_titles",
			"start_date",
			"end_date",
			"synopsis",
			"mean",
			"rank",
			"popularity",
			"num_list_users",
			"num_scoring_users",
			"nsfw",
			"genres",
			"media_type",
			"status",
			"num_volumes",
			"num_chapters",
			"authors",
		}, mal.Limit(5))
	if err != nil {
		log.Error("Failed to filter manga by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Manga not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed

		arr = append(arr, embed)
		return arr, err
	}
	// check if there is more than one result
	if len(orignalMangaList) > 1 {
		// append the title then send it as embed
		var mangaList []string
		for index, manga := range orignalMangaList {
			index += 1
			item := strconv.Itoa(index)
			mangaList = append(mangaList, item+". "+manga.Title)
		}
		// start embed
		embed := NewEmbed().
			SetTitle("More than one result found (MANGA)").
			SetDescription(strings.Join(mangaList, "\n") + "Please use '\\'mal-manga-choice to choose.").SetColor(green).MessageEmbed
		arr = append(arr, embed)
		return arr, nil
	} else {
		var chosenManga int
		chosenManga = 0

		genresList := orignalMangaList[chosenManga].Genres
		var genres []string
		for _, genre := range genresList {
			genres = append(genres, genre.Name)
		}

		genresJoin := strings.Join(genres, ",")
		if genresJoin == "" {
			genresJoin = "\u200b"
		}
		var nsfwString string
		if orignalMangaList[chosenManga].Nsfw == "white" {
			nsfwString = "white (This work is safe for work)"
		} else if orignalMangaList[chosenManga].Nsfw == "" {
			nsfwString = "Unknown"
		} else if orignalMangaList[chosenManga].Nsfw == "gray" {
			nsfwString = "grey (This work may be not safe for work)"
		} else {
			nsfwString = "black (This work is not safe for work)"
		}

		// start embed
		embed := NewEmbed().
			SetTitle(orignalMangaList[chosenManga].Title).
			SetURL("https://myanimelist.net/manga/"+strconv.Itoa(orignalMangaList[chosenManga].ID)).
			SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
			SetThumbnail(orignalMangaList[chosenManga].MainPicture.Large).
			SetDescription(cutDescription(orignalMangaList[chosenManga].Synopsis[:])).
			AddField("Type", orignalMangaList[chosenManga].MediaType).
			AddField("Volumes", strconv.Itoa(orignalMangaList[chosenManga].NumVolumes)).
			AddField("Chapters", strconv.Itoa(orignalMangaList[chosenManga].NumChapters)).
			AddField("Status", orignalMangaList[chosenManga].Status).
			AddField("Published", orignalMangaList[chosenManga].StartDate).
			AddField("Mean Score", strconv.FormatFloat(orignalMangaList[chosenManga].Mean, 'f', 3, 64)).
			//AddField("Popularity", strconv.Itoa(manga_list[chosen_manga].Popularity)).
			//AddField("Rank", strconv.Itoa(manga_list[chosen_manga].Rank)).
			AddField("Genres", genresJoin).
			AddField("NSFW", nsfwString).
			SetFooter(orignalMangaList[chosenManga].AlternativeTitles.Ja, orignalMangaList[chosenManga].MainPicture.Medium).
			InlineAllFields().
			SetColor(green).MessageEmbed

		arr = append(arr, embed)
		return arr, err
	}
}

func malQueryMangaByTitleAndChoice(title string, choice int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	mangaList, _, err := c.Manga.List(ctx, title,
		mal.Fields{
			"id",
			"title",
			"main_picture",
			"alternative_titles",
			"start_date",
			"end_date",
			"synopsis",
			"mean",
			"rank",
			"popularity",
			"num_list_users",
			"num_scoring_users",
			"nsfw",
			"genres",
			"media_type",
			"status",
			"num_volumes",
			"num_chapters",
			"authors",
		}, mal.Limit(5))
	if err != nil {
		log.Error("Failed to filter manga by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Manga not found!\n maybe try using id?").
			SetColor(red).
			MessageEmbed

		arr = append(arr, embed)
		return arr, err
	}

	var chosenManga int
	chosenManga = int(choice)

	genresList := mangaList[chosenManga].Genres
	var genres []string
	for _, genre := range genresList {
		genres = append(genres, genre.Name)
	}

	genresJoin := strings.Join(genres, ",")
	if genresJoin == "" {
		genresJoin = "\u200b"
	}
	var nsfwString string
	if mangaList[chosenManga].Nsfw == "white" {
		nsfwString = "white (This work is safe for work)"
	} else if mangaList[chosenManga].Nsfw == "" {
		nsfwString = "Unknown"
	} else if mangaList[chosenManga].Nsfw == "gray" {
		nsfwString = "grey (This work may be not safe for work)"
	} else {
		nsfwString = "black (This work is not safe for work)"
	}

	// start embed
	embed := NewEmbed().
		SetTitle(mangaList[chosenManga].Title).
		SetURL("https://myanimelist.net/manga/"+strconv.Itoa(mangaList[chosenManga].ID)).
		SetAuthor("MAL", "https://upload.wikimedia.org/wikipedia/commons/7/7a/MyAnimeList_Logo.png").
		SetThumbnail(mangaList[chosenManga].MainPicture.Large).
		SetDescription(cutDescription(mangaList[chosenManga].Synopsis[:])).
		AddField("Type", mangaList[chosenManga].MediaType).
		AddField("Volumes", strconv.Itoa(mangaList[chosenManga].NumVolumes)).
		AddField("Chapters", strconv.Itoa(mangaList[chosenManga].NumChapters)).
		AddField("Status", mangaList[chosenManga].Status).
		AddField("Published", mangaList[chosenManga].StartDate).
		AddField("Mean Score", strconv.FormatFloat(mangaList[chosenManga].Mean, 'f', 3, 64)).
		//AddField("Popularity", strconv.Itoa(manga_list[chosen_manga].Popularity)).
		//AddField("Rank", strconv.Itoa(manga_list[chosen_manga].Rank)).
		AddField("Genres", genresJoin).
		AddField("NSFW", nsfwString).
		SetFooter(mangaList[chosenManga].AlternativeTitles.Ja, mangaList[chosenManga].MainPicture.Medium).
		InlineAllFields().
		SetColor(green).MessageEmbed

	arr = append(arr, embed)
	return arr, err

}

func malQueryMangaByID(id int64) ([]*discordgo.MessageEmbed, error) {
	var arr []*discordgo.MessageEmbed
	// set mal
	c := mal.NewClient(publicInfoClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	manga, _, err := c.Manga.Details(ctx, int(id),
		mal.Fields{
			"id",
			"title",
			"main_picture",
			"alternative_titles",
			"start_date",
			"end_date",
			"synopsis",
			"mean",
			"rank",
			"popularity",
			"num_list_users",
			"num_scoring_users",
			"nsfw", "genres",
			"media_type",
			"status",
			"num_volumes",
			"num_chapters",
			"authors"},
	)
	if err != nil {
		log.Error("Failed to filter media by title: ", err)
		// start embed
		embed := NewEmbed().
			SetDescription("Anime not found!\n maybe try using Title?").
			SetColor(red).
			MessageEmbed
		arr = append(arr, embed)
		return arr, err
	}

	genresList := manga.Genres
	var genres []string
	for _, genre := range genresList {
		genres = append(genres, genre.Name)
	}

	genresJoin := strings.Join(genres, ",")
	if genresJoin == "" {
		genresJoin = "\u200b"
	}

	var nsfwString string
	if manga.Nsfw == "white" {
		nsfwString = "white (This work is safe for work)"
	} else if manga.Nsfw == "" {
		nsfwString = "Unknown"
	} else if manga.Nsfw == "gray" {
		nsfwString = "grey (This work may be not safe for work)"
	} else {
		nsfwString = "black (This work is not safe for work)"
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
		AddField("Genres", genresJoin).
		AddField("NSFW", nsfwString).
		SetFooter(manga.AlternativeTitles.Ja, manga.MainPicture.Medium).
		InlineAllFields().
		SetColor(green).MessageEmbed
	arr = append(arr, embed)

	return arr, nil
}
