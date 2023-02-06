package mainpac

import (
	"Laurene/util"
	tb "gopkg.in/telebot.v3"
	"math/rand"
	"strings"
)

func (s *Service) TgTextReverse(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
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
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
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
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
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

func (s *Service) TgEmojiAlphabet(x tb.Context) (errReturn error) {
	text := ""
	n := 0

	for i := 0; i < len(alphabetRus); i++ {
		text += string(alphabetRus[i]) + "-" + emojiAlphabetRus[i]
		n++
		if n == 7 {
			text += "\n"
			n = 0
		} else {
			text += " "
		}
	}

	x.Send(text)
	return
}

func (s *Service) TgTextEmoji(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
		return
	}

	x.Send("<pre>"+textEmoji(x.Message().ReplyTo.Text)+"</pre>", tb.ModeHTML)
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

var (
	alphabetRus      = []rune{'а', 'б', 'в', 'г', 'д', 'е', 'ё', 'ж', 'з', 'и', 'й', 'к', 'л', 'м', 'н', 'о', 'п', 'р', 'с', 'т', 'у', 'ф', 'х', 'ц', 'ч', 'ш', 'щ', 'ъ', 'ы', 'ь', 'э', 'ю', 'я'}
	emojiAlphabetRus = []string{"🍍", "🔩", "🚁", "👁", "🏠", "🌲", "🎄", "🦒", "🦷", "🪡", "🪡", "🐳", "🌝", "⚽", "🛸", "🦅", "🕷", "🌹", "🧃", "🌮", "🦆", "🏁", "🐹", "⛓", "🕑", "🎱", "🛡", "🪨🪧", "🤣🚀", "🛏🪧", "🧝‍♂", "🇿🇦", "⚓"}
)

func textEmoji(s string) (res string) {
	s = strings.ToLower(s)

	for _, v := range []rune(s) {
		if v == ' ' {
			res += "  "
			continue
		}

		ok := false
		for i, char := range alphabetRus {
			if v == char {
				res += emojiAlphabetRus[i]
				ok = true
				break
			}
		}
		if !ok {
			res += string(v)
		}
	}

	return
}

func (s *Service) TgTextAtbash(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
		return
	}

	text := s.Other.AtbashAlphabet.Replace(x.Message().ReplyTo.Text)
	key := util.CreateKey(8)
	s.Other.AtbashCache.Add(key, x.Message().ReplyTo.Text)
	rm := *s.Bot.Markup(x, "atbash")
	rm.InlineKeyboard[0][0].Data = key

	x.Send("<pre>"+text+"</pre>", &rm, tb.ModeHTML)
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func (s *Service) TgTextAtbashBtn(x tb.Context) (errReturn error) {
	if x.Data() == "" {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "err"), ShowAlert: true})
		return
	}

	key := x.Data()
	text, ok := s.Other.AtbashCache.Get(key)
	if !ok {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "atbash_old"), ShowAlert: true})
		return
	}

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: text.(string), ShowAlert: true})
	return
}
