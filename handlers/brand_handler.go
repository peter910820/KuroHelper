package handlers

import (
	"errors"

	"github.com/bwmarrin/discordgo"

	kurohelpererrors "kurohelper/errors"
	"kurohelper/utils"
)

func FuzzySearchBrand(s *discordgo.Session, i *discordgo.InteractionCreate, cid *CustomID) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	opt, err := utils.GetOptions(i, "查詢資料庫選項")
	if err != nil && errors.Is(err, kurohelpererrors.ErrOptionTranslateFail) {
		utils.HandleError(err, s, i)
		return
	}
	if opt == "" {
		ErogsFuzzySearchBrand(s, i, cid)
	} else {
		VndbFuzzySearchBrand(s, i, cid)
	}
}
