package mainpac

import (
	"Laurene/go-log"
	"Laurene/util"
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
		time.Sleep(time.Second)
	}

	rows := make([]tb.Row, 0, 1)

	if m.AlbumID != "" {
		rows = append(rows, []tb.Btn{*s.TG.Buttons["album_to_pic"]})
	}

	if len(rows) == 0 {
		s.TG.Bot.Send(m.Sender, "Нет действий.", &tb.SendOptions{ReplyTo: m})
		return
	}

	rm := &tb.ReplyMarkup{}
	rm.Inline(rows...)
	s.TG.Bot.Send(m.Sender, "Что сделать с фотографиями?", &tb.SendOptions{ReplyTo: m}, rm)
}

func (s *Service) TgAlbumToPic(c *tb.Callback) {
	defer s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: ""})

	if c.Message.ReplyTo == nil {
		return
	}

	albumID := c.Message.ReplyTo.AlbumID
	if albumID == "" {
		return
	}

	photosMes := s.TG.AlbumsManager.GetAlbum(albumID)
	if len(photosMes) == 0 {
		return
	}

	sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })

	pathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	for _, mes := range photosMes {
		photo := mes.Photo
		path := "files/temp/" + photo.FileID + ".jpg"
		defer os.Remove(path)
		pathes = append(pathes, path)
		err := s.TG.Bot.Download(photo.MediaFile(), path)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic Bot.Download"))
			s.TG.Bot.Send(c.Sender, "Ошибка скачивания, попробуйте позже.")
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic os.Open"))
			s.TG.Bot.Send(c.Sender, "Ошибка файлов, попробуйте позже.")
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Warn(errors.Wrap(err, "TgAlbumToPic image.Decode"))
			s.TG.Bot.Send(c.Sender, "Ошибка файлов, попробуйте позже.")
			return
		}

		images = append(images, img)
	}

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

	rgba := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxX, sumY}})
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

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	out, err := os.Create(outPath)
	if err != nil {
		log.Warn(errors.Wrap(err, "TgAlbumToPic os.Create(outPath)"))
		s.TG.Bot.Send(c.Sender, "Ошибка файлов, попробуйте позже.")
		return
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 85
	err = jpeg.Encode(out, rgba, &opt)
	if err != nil {
		log.Warn(errors.Wrap(err, "TgAlbumToPic jpeg.Encode"))
		s.TG.Bot.Send(c.Sender, "Ошибка файлов, попробуйте позже.")
		return
	}

	_, err = s.TG.Bot.Send(c.Sender, &tb.Document{File: tb.FromDisk(outPath), FileName: "pic.jpg"})
	if err != nil {
		s.TG.Bot.Send(c.Sender, "Ошибка отправки.")
		return
	}
}
