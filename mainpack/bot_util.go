package mainpack

import (
	"Laurene/util"
	"fmt"
	tb "gopkg.in/telebot.v3"
	"strings"
	"time"
)

func (Bot *Bot) isNotAdmin(x tb.Context) bool {
	if x.Chat().ID >= 0 && util.IntInSlice(Bot.AdminList, x.Sender().ID) {
		return false
	}
	return true
}

func (Bot *Bot) sendToSlice(slice []int64, mesText string) {
	for _, chatID := range slice {
		Bot.Bot.Send(&tb.User{ID: chatID}, mesText, tb.ModeHTML)
	}
}

// 4d7h6m34s
func (Bot *Bot) uptimeString(timestamp time.Time) string {
	uptime := time.Since(timestamp).Round(time.Second)
	hours, hoursStr := 0, ""
	for uptime.Hours() > 24 {
		uptime -= time.Hour * 24
		hours++
	}
	if hours > 0 {
		hoursStr = fmt.Sprintf("%vd", hours)
	}
	return hoursStr + uptime.String()
}

func delBtn(rm *tb.ReplyMarkup, copyData string) *tb.ReplyMarkup {
	for i, row := range rm.InlineKeyboard {
		for i2, button := range row {
			ii := strings.Index(button.Data, "|")
			if ii < 0 {
				continue
			}
			if button.Data[ii+1:] == copyData {
				rm.InlineKeyboard[i] = append(rm.InlineKeyboard[i][:i2], rm.InlineKeyboard[i][i2+1:]...)
				return rm
			}
		}
	}
	return rm
}
