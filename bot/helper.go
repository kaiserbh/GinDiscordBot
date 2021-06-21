package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kaiserbh/gin-bot-go/model"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// getArguments split the "parameter from a command
func getArguments(message string) []string {
	trimMessage := strings.TrimSpace(message)
	splitMessage := strings.Split(trimMessage, " ")
	return splitMessage
}

func checkPrefix(message string) bool {
	// check if the prefix or string is greater than 1 character
	if len(message) > 10 {
		log.Warn("Checking prefix; prefix must be less than 10 character")
		return false
	}
	return true
}

func sendEmbed(name, value string, inline bool) discordgo.MessageEmbedField {
	embedField := discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	return embedField
}

// MemberHasPermission checks if a member has the given permission
// for example, If you would like to check if user has the administrator
// permission you would use
// --- MemberHasPermission(s, guildID, userID, discordgo.PermissionAdministrator)
// If you want to check for multiple permissions you would use the bitwise OR
// operator to pack more bits in. (e.g): PermissionAdministrator|PermissionAddReactions
// =================================================================================
//     s          :  discordgo session
//     guildID    :  guildID of the member you wish to check the roles of
//     userID     :  userID of the member you wish to retrieve
//     permission :  the permission you wish to check for
func memberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}
	return false, nil
}

// check latency or how long a function takes to execute
func measureTime(funcName string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Time taken by %s function is %v \n", funcName, time.Since(start))
	}
}

// check allowed channel to use bot command.
func checkAllowedChannel(id string, settings *model.GuildSettings) bool {
	channelIDs := settings.GuildBotChannelsID
	for _, ID := range channelIDs {
		if id == ID {
			return true
		}
	}
	return false
}
