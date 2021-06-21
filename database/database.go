package database

import (
	"context"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DB struct {
	Client *mongo.Client
}

const (
	DiscordRootDatabase = "Discord"
	GuildCollection     = "guilds"
	UserCollection      = "users"
)

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error Database couldn't connect: ", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Couldn't connect to DB dying..")
		return nil
	}
	log.Info("Database Connected")

	return &DB{
		Client: client,
	}
}

// FindByID Find user by ID
func (db *DB) FindByID(ID string) *model.User {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		log.Error("Error changing primitive ID from string: ", err)
	}

	usersCollection := db.Client.Database(DiscordRootDatabase).Collection(UserCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := usersCollection.FindOne(ctx, bson.M{"_id": ObjectID})
	user := model.User{}
	err = res.Decode(&user)
	if err != nil {
		log.Error("Error failed to decoded FindByID: ", err)
	}
	return &user
}

// FindGuildByID Find guild by ID
func (db *DB) FindGuildByID(ID string) (*model.GuildSettings, error) {
	usersCollection := db.Client.Database(DiscordRootDatabase).Collection(GuildCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := usersCollection.FindOne(ctx, bson.D{
		{"guild_id", ID},
	})
	guild := model.GuildSettings{}
	err := res.Decode(&guild)
	if err != nil {
		log.Error("Error failed to decoded FindGuildByID: ", err)
		return nil, err
	}
	return &guild, nil
}

//// might not be needed at all.
//// getGuildObjectID Find guild Object ID by Guild ID
//func (db *DB) getGuildObjectID(ID string) (string, error) {
//	usersCollection := db.Client.Database("Discord").Collection("guilds")
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	res := usersCollection.FindOne(ctx, bson.M{"guild_id": ID})
//	var bsonDocument bson.D
//
//	//guild := model.GuildIDOnly{}
//	err := res.Decode(&bsonDocument)
//	if err != nil {
//		log.Println("[DB] Error failed to decoded FindByID: ", err)
//		return "", err
//	}
//	//ok := guild.GuildID
//	mappedValue := bsonDocument.Map()
//	getID := mappedValue["_id"].(primitive.ObjectID)
//
//	return getID.Hex(), nil
//}

// InsertOrUpdateGuild insert or updates document from collection if it's there if not will create it
func (db *DB) InsertOrUpdateGuild(guildSettings *model.GuildSettings) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	guildsCollection := db.Client.Database(DiscordRootDatabase).Collection(GuildCollection)

	opts := options.Update().SetUpsert(true) // if it doesn't exist create one
	filter := bson.D{
		{"guild_id", guildSettings.GuildID},
	}
	update := bson.D{{"$set", bson.D{
		{"guild_prefix", guildSettings.GuildPrefix},
		{"guild_name", guildSettings.GuildName},
		{"guild_bot_channels_id", guildSettings.GuildBotChannelsID},
		{"time_stamp", guildSettings.TimeStamp},
	}}}
	result, err := guildsCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Warn("Failed to add or update guild prefix: ", err)
		return err
	}
	if result.MatchedCount != 0 {
		log.WithFields(log.Fields{
			"matched": result.MatchedCount,
		}).Info("Matched and replaced existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		log.WithFields(log.Fields{
			"ID": result.UpsertedID,
		}).Info("Inserted a new document")
		return nil
	}
	return nil
}
