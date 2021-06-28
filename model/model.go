package model

import (
	"time"
)

type User struct {
	UserID            string         `bson:"user_id"`
	Guild             *GuildSettings `bson:"guild"`
	NickName          string         `bson:"nick_name"`
	Date              time.Time      `bson:"date"`
	AllowedNickChange bool           `bson:"allowed_nick_change"`
	TimeStamp         time.Time      `bson:"time_stamp"`
}

type GuildSettings struct {
	GuildID               string    `bson:"guild_id"`
	GuildName             string    `bson:"guild_name"`
	GuildPrefix           string    `bson:"guild_prefix"`
	GuildBotChannelsID    []string  `bson:"guild_bot_channels_id"`
	GuildNicknameDuration string    `bson:"guild_nickname_duration"`
	TimeStamp             time.Time `bson:"time_stamp"`
}
