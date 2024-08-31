package programapi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type bootOption struct {
	//必須値
	DiscordToken string
	GuildID      string
	ProgramName  string

	//Option
	ProgamEmbedColor int //通常の埋め込みメッセージの色
}

type bootOptionFunc func(*bootOption)

func WithProgramName(name string) bootOptionFunc {
	return func(option *bootOption) {
		option.ProgramName = name
	}
}

func WithProgramEmbedColor(color int) bootOptionFunc {
	return func(option *bootOption) {
		option.ProgamEmbedColor = color
	}
}

func NewBootOption(discordToken string, guildID string, ProgramName string, opt ...bootOptionFunc) *bootOption {
	discordToken = "Bot " + discordToken
	ans := &bootOption{
		DiscordToken: discordToken,
		GuildID:      guildID,
		ProgramName:  ProgramName,

		ProgamEmbedColor: defaultProgramColor,
	}

	for _, o := range opt {
		o(ans)
	}
	return ans
}

type BotManager struct {
	session     *discordgo.Session
	guildID     string
	botUserInfo *discordgo.User
	programName string

	programColor int
}

func newBotManager(sess *discordgo.Session, guildID string, programName string, opt *bootOption) *BotManager {
	return &BotManager{
		session:     sess,
		guildID:     guildID,
		botUserInfo: sess.State.User,
		programName: programName,

		programColor: opt.ProgamEmbedColor,
	}

}

// ()->(manager, closeFunc, error)
func Boot(option *bootOption) (manager *BotManager, closeFunc func() error, err error) {
	discordSess, err := discordgo.New(option.DiscordToken)

	if err != nil {
		return nil, nil, fmt.Errorf("Error creating Discord session: %w", err)
	}

	discordSess.Open()

	manager = newBotManager(discordSess, option.GuildID, option.ProgramName, option)
	closeFunc = func() error {
		err := discordSess.Close()
		if err != nil {
			return fmt.Errorf("Error closing Discord session: %w", err)
		}
		return nil
	}

	return manager, closeFunc, nil
}

func BootManager() {}
