package bot

import (
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/config"
	"github.com/kaiserbh/gin-bot-go/database"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
)

var db = database.Connect()
var (
	red                = 0xff0000
	green              = 0x11ff00
	previousAuthor     string
	nickCoolDownAuthor []string
)
var Uptime = time.Now()

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	// Logging Method Name
	//log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("Couldn't initiate bot:  ", err)
		return
	}

	_, err = goBot.User("@me")
	if err != nil {
		log.Fatal("Couldn't get botID:  ", err)
	}

	// intent or what to store for bot?
	goBot.Identify.Intents = discordgo.IntentsAll

	// Register handlers here.
	goBot.AddHandler(guildJoinInit)

	// configurationCommands
	goBot.AddHandler(setPrefixHandler)
	goBot.AddHandler(setBotChannelHandler)
	goBot.AddHandler(setNicknameCooldown)

	// generalCommands
	go goBot.AddHandler(helpMessageHandler)
	goBot.AddHandler(pingMessageHandler)
	goBot.AddHandler(stats)
	goBot.AddHandler(setNick)
	goBot.AddHandler(resetNickHandler)
	goBot.AddHandler(botPing)
	goBot.AddHandler(invite)

	// adminCommands
	goBot.AddHandler(ban)

	//TODO:Support get support invite linko
	// anilistCommands
	//TODO:anime Query anime from Anilist
	//TODO:manga Query manga from Anilist
	//TODO:character Query character from Anilist
	//TODO:staff Query person/staff from Anilist
	//TODO:studio Query studio from Anilist
	//TODO:user Query user from Anilist

	//miscellaneousCommands
	//TODO:permissions Show your permissions or the member specified.
	//TODO:userinfo Show some information about yourself or the member specified.
	//TODO:serverinfo Get some information about this server.

	// Start bot with chan.
	err = goBot.Open()
	if err != nil {
		log.Fatal("Couldn't Connect bot:  ", err)
		return
	}

	log.Info("Bot is running")
}

// guildJoinInit runs whenever it joins a new guild or gets online.
func guildJoinInit(s *discordgo.Session, g *discordgo.GuildCreate) {
	guild, err := s.Guild(g.ID)
	if err != nil {
		log.Error("Getting guild information from Session: ", err)
		return
	}

	guildChannels := g.Channels
	var guildIDs []string
	for _, guild := range guildChannels {
		guildIDs = append(guildIDs, guild.ID)
	}
	_, err = db.FindGuildByID(guild.ID)
	if err != nil {
		log.Error("Guild not found in DB creating one with default values... ", err)
		guildSetting := model.GuildSettings{
			GuildID:               guild.ID,
			GuildName:             guild.Name,
			GuildPrefix:           config.BotPrefix,
			GuildBotChannelsID:    guildIDs,
			GuildNicknameDuration: "30",
			TimeStamp:             time.Now().UTC(),
		}
		err := db.InsertOrUpdateGuild(&guildSetting)
		if err != nil {
			log.Error("Error inserting default values into DB", err)
			return
		}
	}
	log.Info("init successful")
}
