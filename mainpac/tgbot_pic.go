package mainpac

import (
	"Laurene/go-log"
	"Laurene/util"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v3"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"sort"
	"time"
)

func (s *Service) TgPic(x tb.Context) (errReturn error) {
	if x.Message().AlbumID != "" {
		if s.TG.AlbumsManager.AddPhotoMes(x.Message().AlbumID, x.Message()) {
			return
		}
		time.Sleep(time.Millisecond * 500)
	}

	if x.Message().AlbumID == "" {
		x.Send("Нет действий.", &tb.SendOptions{ReplyTo: x.Message()})
		return
	}

	text := "" +
		"Что сделать с фотографиями?" +
		"\n" +
		"\n1. Объединить фотографии"

	x.Send(text, &tb.SendOptions{ReplyTo: x.Message()}, s.TG.menu.picBtns)
	return
}

func (s *Service) TgAlbumToPic(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}

	albumID := x.Message().ReplyTo.AlbumID
	if albumID == "" {
		return
	}

	okLock := s.TG.AlbumsManager.LockAlbum(albumID)
	if !okLock {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Подождите пока другой запрос выполнится.", ShowAlert: true})
		return
	}
	defer s.TG.AlbumsManager.UnLockAlbum(albumID)
	photosMes := s.TG.AlbumsManager.GetAlbum(albumID)
	if len(photosMes) == 0 {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Фото не найдены, попробуйте переслать их в этого в бота.", ShowAlert: true})
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
		err := x.Bot().Download(photo.MediaFile(), path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic Bot.Download"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка скачивания, попробуйте позже.", ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	var rgba *image.RGBA

	c := x.Callback()
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
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка, напишите создателю."})
		return
	}

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	out, err := os.Create(outPath)
	if err != nil {
		log.Error(errors.Wrap(err, "TgAlbumToPic os.Create(outPath)"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
		return
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 85
	err = jpeg.Encode(out, rgba, &opt)
	if err != nil {
		log.Error(errors.Wrap(err, "TgAlbumToPic jpeg.Encode"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, попробуйте позже.", ShowAlert: true})
		return
	}

	rm := &tb.ReplyMarkup{}
	rm.Inline([]tb.Btn{*s.TG.Buttons["picfile_to_pic"]})
	err = x.Send(&tb.Document{File: tb.FromDisk(outPath), FileName: "pic.jpg", Caption: photosMes[0].Caption}, rm)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка отправки.", ShowAlert: true})
		return
	}
	x.Respond()
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, c.Data)})
	return
}

func (s *Service) TgFilePicToPic(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	if x.Message() == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка"})
		return
	}
	if x.Message().Document == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка"})
		return
	}

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	err := x.Bot().Download(x.Message().Document.MediaFile(), outPath)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка"})
		return
	}

	x.Send(&tb.Photo{File: tb.FromDisk(outPath), Caption: x.Message().Caption})
	x.Bot().EditReplyMarkup(x.Message(), nil)
	x.Respond()
	return
}
