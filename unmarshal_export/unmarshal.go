package unmarshal_export

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"
)

func UnmarshalChan(file *os.File, n int) chan *Message {
	ch := make(chan *Message, 100)

	go func(ch chan *Message) {
		var lines []string

		limiter := false
		if n > 0 {
			limiter = true
		}

		i := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if i > n && limiter {
				break
			}

			if scanner.Text() == `  {` {
				lines = []string{}
			} else if scanner.Text() == `  },` {
				lines = append(lines, `  }`)
				i++
				ch <- unmarshalJson(strings.Join(lines, ""))
				continue
			}
			lines = append(lines, scanner.Text())
		}

		close(ch)
	}(ch)

	return ch
}

func unmarshalJson(str string) *Message {
	var (
		err error
		ok  bool
	)

	msgRaw := messageForUnmarshal{}
	err = json.Unmarshal([]byte(str), &msgRaw)
	if err != nil {
		return nil
	}

	msg := Message{
		ID:        int64(msgRaw.ID),
		Type:      msgRaw.Type,
		Date:      time.Time{},
		FromID:    msgRaw.FromID,
		ReplyToID: int64(msgRaw.ReplyToID),
		Text:      "",
		Hashtags:  make([]string, 0),
		Raw:       str,
	}

	msg.Date, err = time.Parse("2006-01-02T15:04:05", msgRaw.Date)

	if msgRaw.MediaType != "" {
		msg.Type = msgRaw.MediaType
	}
	if msgRaw.Photo != "" {
		msg.Type = "photo"
	}

	if _, ok = msgRaw.Text.(string); ok {
		msg.Text = msgRaw.Text.(string)
	} else {
		for i := 0; i < len(msgRaw.Text.([]interface{})); i++ {
			if _, ok = msgRaw.Text.([]interface{})[i].(string); ok {
				msg.Text += msgRaw.Text.([]interface{})[i].(string)
			} else if _, ok = msgRaw.Text.([]interface{})[i].(map[string]interface{}); ok {
				if _, ok = msgRaw.Text.([]interface{})[i].(map[string]interface{})["text"].(string); ok {
					msg.Text += msgRaw.Text.([]interface{})[i].(map[string]interface{})["text"].(string)
				} else {
					return nil
				}
				if _, ok = msgRaw.Text.([]interface{})[i].(map[string]interface{})["type"].(string); ok {
					if msgRaw.Text.([]interface{})[i].(map[string]interface{})["type"].(string) == "hashtag" {
						msg.Hashtags = append(msg.Hashtags, msgRaw.Text.([]interface{})[i].(map[string]interface{})["text"].(string))
					}
				}
			}
		}
	}

	return &msg
}
