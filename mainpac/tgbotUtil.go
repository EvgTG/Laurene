package mainpac

import (
	"Laurene/util"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func (tg *TG) isAdmin(user *tb.User, chat int64) bool {
	if chat >= 0 && util.IntInSlice(tg.AdminList, user.ID) {
		return true
	}
	return false
}

func (tg *TG) sendToSlice(slice []int64, mesText string) {
	for _, chatID := range slice {
		tg.Bot.Send(&tb.User{ID: int(chatID)}, mesText, tb.ModeHTML)
	}
}

// 4d7h6m34s
func (tg *TG) uptimeString(timestamp time.Time) string {
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

func (tg *TG) addBtn(btn tb.Btn, key string, handler interface{}) {
	tg.Buttons[key] = &btn
	tg.Bot.Handle(&btn, handler)
}
