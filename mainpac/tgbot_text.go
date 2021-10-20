package mainpac

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"math/rand"
	"strings"
)

func (s *Service) TgTextReverse(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	x.Send(textReverse(x.Message().ReplyTo.Text))
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func textReverse(s string) (res string) {
	for _, v := range []rune(s) {
		res = string(v) + res
	}
	return
}

func (s *Service) TgTextToUpper(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	x.Send(strings.ToUpper(x.Message().ReplyTo.Text))
	c := x.Callback()
	_, err := x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, c.Data)})
	fmt.Println(err)
	return
}

func (s *Service) TgTextRandom(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	x.Send(textRandom(x.Message().ReplyTo.Text, s.Rand))
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func textRandom(s string, r *rand.Rand) (res string) {
	for _, v := range []rune(s) {
		if r.Intn(2) == 1 {
			res += strings.ToUpper(string(v))
		} else {
			res += strings.ToLower(string(v))
		}
	}
	return
}
