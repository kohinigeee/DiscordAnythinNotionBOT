package main

import "github.com/kohinigeee/DiscordAnythinNotionBOT/programapi"

func main() {
	token := "ExampleToken"
	guildID := "ExampleGuildID"
	chanID := "ExampleChannelID"

	opt := programapi.NewBootOption(token, guildID, "TestProgram",
		programapi.WithProgramEmbedColor(programapi.ColorCyan))

	mng, f, err := programapi.Boot(opt)

	if err != nil {
		panic(err)
	}

	defer f()

	mng.SendNormalEmbedMsg(chanID, "Normal", "test program send message")
	mng.SendErrorEmbedMsg(chanID, "Error", "test program send message")
	msg, _ := mng.SendProgramEmbedMsg(chanID, "ProgramColor", "test program send message")
	mng.SendMsg(chanID, "test program send message")

	mng.MessageThreadStart(chanID, msg.ID, "Test Trhead", 0)
}
