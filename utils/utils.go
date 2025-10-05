package utils

import (
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	kurohelpererrors "kurohelper/errors"
)

// handle interaction command common respond
//
// 這邊用來當作如果嵌入式訊息發送失敗的最後發送手段
func InteractionRespond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

// handle interaction command embed respond
//
// editFlag參數為有無需要修改因為defer而產生的interaction訊息(機器人正在思考...)
func InteractionEmbedRespond(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components *discordgo.ActionsRow, editFlag bool) {
	var comps []discordgo.MessageComponent
	if components != nil {
		comps = []discordgo.MessageComponent{*components}
	} else {
		comps = []discordgo.MessageComponent{}
	}

	if editFlag {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{embed},
			Components: &comps,
		})
		if err != nil {
			logrus.Error(err)
			InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
		}
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: comps,
			},
		})
	}
}

// handle interaction command embed respond
// 管理員專用版本
//
// editFlag參數為有無需要修改因為defer而產生的interaction訊息(機器人正在思考...)
func InteractionEmbedRespondForSelf(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, editFlag bool) {
	if editFlag {
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			logrus.Error(err)
			InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
		}
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

// 傳送嵌入訊息內建包裝錯誤的版本
func InteractionEmbedErrorRespond(s *discordgo.Session, i *discordgo.InteractionCreate, errString string, editFlag bool) {
	embed := &discordgo.MessageEmbed{
		Title: "❌錯誤",
		Color: 0xcc543a,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "說明",
				Value:  errString,
				Inline: false,
			},
		},
	}
	InteractionEmbedRespond(s, i, embed, nil, editFlag)
}

func EmbedErrorRespond(s *discordgo.Session, i *discordgo.InteractionCreate, errString string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	embed := &discordgo.MessageEmbed{
		Title: "❌錯誤",
		Color: 0xcc543a,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "說明",
				Value:  errString,
				Inline: false,
			},
		},
	}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &[]discordgo.MessageComponent{},
	})
}

func EditEmbedRespond(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components *discordgo.ActionsRow) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	var comps []discordgo.MessageComponent
	if components != nil {
		comps = []discordgo.MessageComponent{*components}
	} else {
		comps = []discordgo.MessageComponent{}
	}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &comps,
	})
}

// get slash command options
func GetOptions(i *discordgo.InteractionCreate, name string) (string, error) {
	for _, v := range i.ApplicationCommandData().Options {
		if v.Name == name {
			value, ok := v.Value.(string) // type assertion
			if !ok {
				return "", kurohelpererrors.ErrOptionTranslateFail
			} else {
				return value, nil
			}
		}
	}
	return "", kurohelpererrors.ErrOptionNotFound
}

func IsEnglish(r rune) bool {
	if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
		return false
	}
	return true
}

func GetEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}
	return def
}
