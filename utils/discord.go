package utils

import (
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
func InteractionEmbedRespondForSelf(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed, components *discordgo.ActionsRow, editFlag bool) {
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
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		})
	}
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

// Use discordgo.MessageComponent slice to make ActionsRow
func MakeActionsRow(messageComponent []discordgo.MessageComponent) *discordgo.ActionsRow {
	if len(messageComponent) != 0 {
		return &discordgo.ActionsRow{
			Components: messageComponent,
		}
	} else {
		return nil
	}

}

func MakeErrorEmbedMsg(errString string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
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
}

// get user discord ID
func GetUserID(i *discordgo.InteractionCreate) string {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	return userID
}

// get user discord name
func GetUsername(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.Username
	} else if i.User != nil {
		return i.User.Username
	}
	return ""
}

func GetAvatarURL(user *discordgo.User) string {
	if user.Avatar != "" {
		// 自訂大頭貼
		return discordgo.EndpointUserAvatar(user.ID, user.Avatar)
	}

	// 沒有自訂大頭貼 → 使用預設頭貼
	discriminator, _ := strconv.Atoi(user.Discriminator)
	return discordgo.EndpointDefaultUserAvatar(discriminator % 5)
}

//
//以下為新版API架構Utils
//

func InteractionRespondV2(s *discordgo.Session, i *discordgo.InteractionCreate, components []discordgo.MessageComponent) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: components,
		},
	})
	if err != nil {
		logrus.Error(err)
		InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}
}

func InteractionRespondEditComplex(s *discordgo.Session, i *discordgo.InteractionCreate, components []discordgo.MessageComponent) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Components: &components,
	})
	if err != nil {
		logrus.Error(err)
		InteractionRespond(s, i, "該功能目前異常，請稍後再嘗試")
	}
}

func MakeErrorComponentV2(errMsg string) []discordgo.MessageComponent {
	color := 0xcc543a
	divider := true
	containerComponents := []discordgo.MessageComponent{
		discordgo.TextDisplay{
			Content: "# ❌錯誤 \n## " + errMsg,
		},
		discordgo.Separator{Divider: &divider},
		discordgo.TextDisplay{
			Content: "聯絡我們: https://discord.gg/6rkrm7tsXr",
		},
	}

	return []discordgo.MessageComponent{
		discordgo.Container{
			AccentColor: &color,
			Components:  containerComponents,
		},
	}
}
