package mainpac

import (
	"Laurene/go-log"
	"Laurene/util"
	"fmt"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"sort"
	"time"
)

func (s *Service) TgPic(m *tb.Message) {
	if m.AlbumID != "" {
		if s.TG.AlbumsManager.AddPhotoMes(m.AlbumID, m) {
			return
		}
		time.Sleep(time.Millisecond * 500)
	}

	if m.AlbumID == "" {
		s.TG.Bot.Send(m.Sender, "Нет действий.", &tb.SendOptions{ReplyTo: m})
		return
	}

	text := "" +
		"Что сделать с фотографиями?" +
		"\n" +
		"\n1. Объединить фотографии"

	rm := &tb.ReplyMarkup{}
	rm.Inline([]tb.Btn{*s.TG.Buttons["album_to_pic_down"], *s.TG.Buttons["album_to_pic_right"], *s.TG.Buttons["album_to_pic_mesh"]})
	s.TG.Bot.Send(m.Sender, text, &tb.SendOptions{ReplyTo: m}, rm)
}

func (s *Service) TgAlbumToPic(c *tb.Callback) {
	if c.Message == nil || c.Message.ReplyTo == nil {
		return
	}

	albumID := c.Message.ReplyTo.AlbumID
	if albumID == "" {
		return
	}

	okLock := s.TG.AlbumsManager.LockAlbum(albumID)
	if !okLock {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Подождите пока другой запрос выполнится.", ShowAlert: true})
		return
	}
	defer s.TG.AlbumsManager.UnLockAlbum(albumID)
	photosMes := s.TG.AlbumsManager.GetAlbum(albumID)
	if len(photosMes) == 0 {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Фото не найдены, попробуйте переслать их в этого в бота.", ShowAlert: true})
		return
	}

	sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })

	pathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, 777)
	defer os.Remove(dir)
	for _, mes := range photosMes {
		photo := mes.Photo
		path := dir + photo.FileID + ".jpg"
		pathes = append(pathes, path)
		err := s.TG.Bot.Download(photo.MediaFile(), path)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic Bot.Download"))
			s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка скачивания, попробуйте позже.", ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic os.Open"))
			s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic image.Decode"))
			s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	var rgba *image.RGBA

	if c.Data == "mesh" && len(images) == 2 {
		c.Data = "right"
	}

	switch c.Data {
	case "down":
		var sumY, sumYi, maxX, maxY int
		for _, img := range images {
			if img.Bounds().Max.X > maxX {
				maxX = img.Bounds().Max.X
			}
			if img.Bounds().Max.Y > maxY {
				maxY = img.Bounds().Max.Y
			}
			sumY += img.Bounds().Max.Y
		}

		rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxX, sumY}})
		for _, img := range images {
			x := 0
			if maxX != img.Bounds().Max.X {
				x = (maxX - img.Bounds().Max.X) / 2
			}

			p1 := image.Point{x, sumYi}
			p2 := image.Point{img.Bounds().Max.X + x, sumYi + img.Bounds().Max.Y}
			r := image.Rectangle{p1, p2}

			draw.Draw(rgba, r, img, image.Point{0, 0}, draw.Src)
			sumYi += img.Bounds().Max.Y
		}
	case "right":
		var sumX, sumXi, maxX, maxY int
		for _, img := range images {
			if img.Bounds().Max.X > maxX {
				maxX = img.Bounds().Max.X
			}
			if img.Bounds().Max.Y > maxY {
				maxY = img.Bounds().Max.Y
			}
			sumX += img.Bounds().Max.X
		}

		rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sumX, maxY}})
		for _, img := range images {
			y := 0
			if maxY != img.Bounds().Max.Y {
				y = (maxY - img.Bounds().Max.Y) / 2
			}

			p1 := image.Point{sumXi, y}
			p2 := image.Point{sumXi + img.Bounds().Max.X, img.Bounds().Max.Y + y}
			r := image.Rectangle{p1, p2}

			draw.Draw(rgba, r, img, image.Point{0, 0}, draw.Src)
			sumXi += img.Bounds().Max.X
		}
	case "mesh":
		ln := len(images)
		var sumX, sumY, maxX, maxY int
		for _, img := range images {
			if img.Bounds().Max.X > maxX {
				maxX = img.Bounds().Max.X
			}
			if img.Bounds().Max.Y > maxY {
				maxY = img.Bounds().Max.Y
			}
		}

		xyi, yminus := 0, 0
		for {
			xyi++
			if xyi*xyi >= ln {
				break
			}
		}
		if xyi*xyi-xyi >= ln {
			yminus++
		}

		ii := 0
		rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{xyi * maxX, (xyi - yminus) * maxY}})
		for _, img := range images {
			x, y := 0, 0
			if maxX != img.Bounds().Max.X {
				x = (maxX - img.Bounds().Max.X) / 2
			}
			if maxY != img.Bounds().Max.Y {
				y = (maxY - img.Bounds().Max.Y) / 2
			}

			p1 := image.Point{x + sumX, y + sumY}
			p2 := image.Point{x + sumX + img.Bounds().Max.X, y + sumY + img.Bounds().Max.Y}
			r := image.Rectangle{p1, p2}
			sumX += maxX
			if xyi == ii+1 {
				ii = -1
				sumX = 0
				sumY += maxY
			}
			ii++

			draw.Draw(rgba, r, img, image.Point{0, 0}, draw.Src)
		}
	default:
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка, напишите создателю."})
		return
	}

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	out, err := os.Create(outPath)
	if err != nil {
		log.Warn(errors.Wrap(err, "TgAlbumToPic os.Create(outPath)"))
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
		return
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 85
	err = jpeg.Encode(out, rgba, &opt)
	if err != nil {
		log.Warn(errors.Wrap(err, "TgAlbumToPic jpeg.Encode"))
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
		return
	}

	rm := &tb.ReplyMarkup{}
	rm.Inline([]tb.Btn{*s.TG.Buttons["picfile_to_pic"]})
	_, err = s.TG.Bot.Send(c.Sender, &tb.Document{File: tb.FromDisk(outPath), FileName: "pic.jpg", Caption: photosMes[0].Caption}, rm)
	if err != nil {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка отправки.", ShowAlert: true})
		return
	}
	s.TG.Bot.Respond(c)
	_, err = s.TG.Bot.EditReplyMarkup(c.Message, &tb.ReplyMarkup{InlineKeyboard: delBtn(c.Message.ReplyMarkup.InlineKeyboard, c.Data)})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *Service) TgFilePicToPic(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	if c.Message == nil {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка"})
		return
	}
	if c.Message.Document == nil {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка"})
		return
	}

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	err := s.TG.Bot.Download(c.Message.Document.MediaFile(), outPath)
	if err != nil {
		s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Ошибка"})
		return
	}

	s.TG.Bot.EditReplyMarkup(c.Message, nil)
	s.TG.Bot.Send(c.Sender, &tb.Photo{File: tb.FromDisk(outPath), Caption: c.Message.Caption}, &tb.SendOptions{ReplyTo: c.Message})
	s.TG.Bot.Respond(c)
}
