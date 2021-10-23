package mainpac

import (
	"Laurene/util"
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

	x.Send("<pre>"+textReverse(x.Message().ReplyTo.Text)+"</pre>", tb.ModeHTML)
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

	x.Send("<pre>"+strings.ToUpper(x.Message().ReplyTo.Text)+"</pre>", tb.ModeHTML)
	c := x.Callback()
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, c.Data)})
	return
}

func (s *Service) TgTextRandom(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	x.Send("<pre>"+textRandom(x.Message().ReplyTo.Text, s.Rand)+"</pre>", tb.ModeHTML)
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

func (s *Service) TgTextAtbash(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	text := s.Other.AtbashAlphabet.Replace(x.Message().ReplyTo.Text)
	key := util.CreateKey(8)
	s.Other.AtbashCache.Add(key, x.Message().ReplyTo.Text)
	rm := *s.TG.menu.atbashBtns
	rm.InlineKeyboard[0][0].Data = key

	x.Send("<pre>"+text+"</pre>", &rm, tb.ModeHTML)
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func (s *Service) TgTextAtbashBtn(x tb.Context) (errReturn error) {
	if x.Data() == "" {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	key := x.Data()
	text, ok := s.Other.AtbashCache.Get(key)
	if !ok {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Сообщение слишком старое. Перешли его в бота, чтобы расшифровать.", ShowAlert: true})
		return
	}

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: text.(string), ShowAlert: true})
	return
}
