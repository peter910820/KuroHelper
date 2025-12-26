package utils

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	kurohelpererrors "kurohelper/errors"

	kurohelpercore "github.com/peter910820/kurohelper-core"
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
)

// 錯誤統一處理方法
func HandleError(err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	logrus.Error(err)
	switch {
	case errors.Is(err, kurohelperdb.ErrUniqueViolation):
		InteractionEmbedErrorRespond(s, i, "資料已存在，此次操作無效", true)
	case errors.Is(err, gorm.ErrRecordNotFound):
		InteractionEmbedErrorRespond(s, i, "使用者尚未建檔", true)
	case errors.Is(err, kurohelpercore.ErrRateLimit):
		InteractionEmbedErrorRespond(s, i, "速率限制，請過約1分鐘後再試", true)
	case errors.Is(err, kurohelpercore.ErrSearchNoContent):
	case errors.Is(err, kurohelpercore.ErrSearchNoContent):
		InteractionEmbedErrorRespond(s, i, "找不到任何結果喔", true)
	case errors.Is(err, kurohelpererrors.ErrTimeWrongFormat):
		InteractionEmbedErrorRespond(s, i, "日期格式錯誤，格式為YYYYMMDD", true)
	case errors.Is(err, kurohelpererrors.ErrDateExceedsTomorrow):
		InteractionEmbedErrorRespond(s, i, "日期格式錯誤，完成日期不得超過今日加一天", true)
	case errors.Is(err, kurohelpercore.ErrBangumiCharacterListSearchNotSupported):
		InteractionEmbedErrorRespond(s, i, "目前不支援對Bangumi使用角色列表搜尋", true)
	case errors.Is(err, kurohelpererrors.ErrCIDGetParameterFailed):
		fallthrough
	case errors.Is(err, kurohelpercore.ErrCacheLost):
		EmbedErrorRespond(s, i, "快取過期，請重新查詢")
	case errors.Is(err, kurohelpercore.ErrStatusCodeAbnormal):
		fallthrough
	case errors.Is(err, kurohelpererrors.ErrOptionTranslateFail):
		fallthrough
	case errors.Is(err, kurohelpererrors.ErrOptionNotFound):
		fallthrough
	default:
		InteractionEmbedErrorRespond(s, i, "該功能目前異常，請稍後再嘗試", true)
	}
}
