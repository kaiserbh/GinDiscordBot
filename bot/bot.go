package bot

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/config"
	"github.com/kaiserbh/gin-bot-go/database"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
)

var db = database.Connect()
var (
	red            = 0xff0000
	green          = 0x11ff00
	previousAuthor string
)
var Uptime = time.Now()

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func init() { flag.Parse() }

var s *discordgo.Session

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
	//goBot, err := discordgo.New("Bot " + config.Token)
	//if err != nil {
	//	log.Fatal("Couldn't initiate bot:  ", err)
	//	return
	//}

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	//s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
	//		h(s, i)
	//	}
	//})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := malCommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := anilistCommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.Identify.Presence = discordgo.GatewayStatusUpdate{
		Game: discordgo.Activity{
			Name:          "Using Slash Command Now",
			Type:          0,
			URL:           "",
			CreatedAt:     time.Now(),
			ApplicationID: "",
			State:         "",
			Details:       "Okay I don't know anymore",
			Timestamps:    discordgo.TimeStamps{},
			Emoji:         discordgo.Emoji{},
			Party:         discordgo.Party{},
			Assets:        discordgo.Assets{},
			Secrets:       discordgo.Secrets{},
			Instance:      false,
			Flags:         1,
		},
		Status: "Just testing ya know",
		AFK:    true,
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Info("Adding commands")

	// MAl
	malRegisteredCommands := make([]*discordgo.ApplicationCommand, len(malCommands))
	for i, v := range malCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		malRegisteredCommands[i] = cmd
	}

	//TODO:GeneralSlashCommandsAdd (Maybe)
	anilistRegisteredCommands := make([]*discordgo.ApplicationCommand, len(alCommands))
	for i, v := range alCommands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		anilistRegisteredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		//anilistRegisteredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		//if err != nil {
		//	log.Fatalf("Could not fetch registered commands: %v", err)
		//}
		//
		//malRegisteredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		//if err != nil {
		//	log.Fatalf("Could not fetch registered commands: %v", err)
		//}

		for _, v := range anilistRegisteredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}

		for _, v := range malRegisteredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
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
