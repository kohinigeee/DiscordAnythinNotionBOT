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

func (mng *BotManager) SendMsg(channelID, content string) (*discordgo.Message, error) {
	content = mng.attachProgramName(content)

	msg, err := mng.session.ChannelMessageSend(channelID, content)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mng *BotManager) SendNormalEmbedMsg(channelID, title, content string) (*discordgo.Message, error) {

	embed := mng.createEmbedMsg(defaultNormalEmbedColor, title, content)
	msg, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mng *BotManager) SendErrorEmbedMsg(channelID, title, content string) (*discordgo.Message, error) {
	color := 0xff0000

	embed := mng.createEmbedMsg(color, title, content)
	msg, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mng *BotManager) SendProgramEmbedMsg(channelID, title, content string) (*discordgo.Message, error) {
	embed := mng.createEmbedMsg(mng.programColor, title, content)
	msg, err := mng.session.ChannelMessageSendEmbed(channelID, embed)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mng *BotManager) MessageThreadStart(channelID, messageID, threadName string, archiveDuartionMin int) (*discordgo.Channel, error) {

	thread, err := mng.session.MessageThreadStart(channelID, messageID, threadName, archiveDuartionMin)

	if err != nil {
		return nil, err
	}

	return thread, nil
}
