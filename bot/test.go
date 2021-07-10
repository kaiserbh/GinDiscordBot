package bot

import (
	"bytes"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/mtibben/confusables"
)

var badWords = []string{"paypal", "niger"}

func removeDups(s string) string {
	var buf bytes.Buffer
	var last rune
	for i, r := range s {
		if r != last || i == 0 {
			buf.WriteRune(r)
			last = r
		}
	}
	return buf.String()
}

func chatFilter(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	cleaned := removeDups(strings.ReplaceAll(confusables.Skeleton(m.Content), " ", ""))

	for _, filter := range badWords {
		if strings.Contains(cleaned, filter) {
			s.ChannelMessageSend(m.ChannelID, "banned word")
			return
		}
	}
}

// func testCmd(ctx *framework.CommandContext) {
// 	cleaned := removeDups(strings.ReplaceAll(confusables.Skeleton(ctx.Message.Content), " ", ""))

// 	for _, filter := range badWords {
// 		if strings.Contains(cleaned, filter) {
// 			ctx.Bot.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("bad word %s detected", filter))
// 			return
// 		}
// 	}
// }
