package programapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (mng *BotManager) createEmbedMsg(color int, title, msg string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: msg,
		Color:       color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    mng.programName,
			IconURL: mng.botUserInfo.AvatarURL("24"),
		},
	}
	return embed
}

func (mng *BotManager) attachProgramName(msg string) string {
	return fmt.Sprintf("**[%s]**\n", mng.programName) + msg
}

func (mng *BotManager) SendMsg(channelID, content string) {
	content = mng.attachProgramName(content)

	mng.session.ChannelMessageSend(channelID, content)
}

func (mng *BotManager) SendNormalEmbedMsg(channelID, title, content string) error {

	embed := mng.createEmbedMsg(defaultNormalEmbedColor, title, content)
	_, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return err
	}

	return nil
}

func (mng *BotManager) SendErrorEmbedMsg(channelID, title, content string) error {
	color := 0xff0000

	embed := mng.createEmbedMsg(color, title, content)
	_, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return err
	}

	return nil
}

func (mng *BotManager) SendProgramEmbedMsg(channelID, title, content string) error {
	embed := mng.createEmbedMsg(mng.programColor, title, content)
	_, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return err
	}

	return nil
}
