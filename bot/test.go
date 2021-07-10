package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/mtibben/confusables"
)

func chatFilter(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// if confusables.Confusable(m.Content, "paypal") {
	// 	s.ChannelMessageSend(m.ChannelID, "banned word")
	// }

	if strings.Contains(m.Content, confusables.Skeleton("paypal")) {
		s.ChannelMessageSend(m.ChannelID, confusables.Skeleton("paypal"))
	}
}
