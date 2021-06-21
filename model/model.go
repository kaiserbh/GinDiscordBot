package model

import (
	"time"
)

type User struct {
	UserID       string    `bson:"user_id"`
	NickName     string    `bson:"nick_name"`
	Date         time.Time `bson:"date"`
	DurationLeft time.Time `bson:"duration_left"`
	TimeStamp    time.Time `bson:"time_stamp"`
}

type GuildSettings struct {
	GuildID            string    `bson:"guild_id"`
	GuildName          string    `bson:"guild_name"`
	GuildPrefix        string    `bson:"guild_prefix"`
	GuildBotChannelsID []string  `bson:"guild_bot_channels_id"`
	TimeStamp          time.Time `bson:"time_stamp"`
}
