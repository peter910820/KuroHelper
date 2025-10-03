package utils

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"

	kurohelpererrors "kurohelper/errors"
)

func HandleError(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	logrus.Error(err)
	switch {
	case errors.Is(err, kurohelpererrors.ErrRateLimit):
		InteractionEmbedErrorRespond(s, i, "速率限制，請過約1分鐘後再試", true)
	case errors.Is(err, kurohelpererrors.ErrSearchNoContent):
		InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
	case errors.Is(err, kurohelpererrors.ErrCacheLost):
		EmbedErrorRespond(s, i, "快取過期，請重新查詢")
	case errors.Is(err, kurohelpererrors.ErrStatusCodeAbnormal):
		fallthrough
	case errors.Is(err, kurohelpererrors.ErrOptionTranslateFail):
		fallthrough
	case errors.Is(err, kurohelpererrors.ErrOptionNotFound):
		fallthrough
	default:
		InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
	}
}
