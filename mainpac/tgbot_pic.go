package mainpac

import (
	"Laurene/go-log"
	"Laurene/util"
	"bytes"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	tb "gopkg.in/tucnak/telebot.v3"
	"image"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
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
		text := "" +
			"Что сделать с фотографиями?" +
			"\n" +
			"\n1. Сжать" +
			"\n2. Исказить (seam carving)"
		x.Send(text, &tb.SendOptions{ReplyTo: x.Message()}, s.TG.menu.picBtns)
		return
	}

	text := "" +
		"Что сделать с фотографиями?" +
		"\n" +
		"\n1. Объединить" +
		"\n2. Сжать" +
		"\n3. Исказить (seam carving)"
	x.Send(text, &tb.SendOptions{ReplyTo: x.Message()}, s.TG.menu.picAlbumsBtns)
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
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка скачивания, напишите автору.", ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	var rgba *image.RGBA

	c := x.Callback()
	if c.Data == "mesh" && len(images) == 2 {
		c.Data = "right"
	}

	var maxX, maxY int
	for _, img := range images {
		if img.Bounds().Max.X > maxX {
			maxX = img.Bounds().Max.X
		}
		if img.Bounds().Max.Y > maxY {
			maxY = img.Bounds().Max.Y
		}
	}
	switch c.Data {
	case "down":
		var sumY, sumYi int
		for _, img := range images {
			a := float64(maxX) / float64(img.Bounds().Max.X)
			sumY += int(float64(img.Bounds().Max.Y) * a)
		}

		rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{maxX, sumY}})
		for _, img := range images {
			a := float64(maxX) / float64(img.Bounds().Max.X)
			yPlus := int(float64(img.Bounds().Max.Y) * a)

			p1 := image.Point{0, sumYi}
			p2 := image.Point{maxX, sumYi + yPlus}
			r := image.Rectangle{p1, p2}

			draw.BiLinear.Scale(rgba, r, img, img.Bounds(), draw.Over, nil)
			sumYi += yPlus
		}
	case "right":
		var sumX, sumXi int
		for _, img := range images {
			a := float64(maxY) / float64(img.Bounds().Max.Y)
			sumX += int(float64(img.Bounds().Max.X) * a)
		}

		rgba = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{sumX, maxY}})
		for _, img := range images {
			a := float64(maxY) / float64(img.Bounds().Max.Y)
			xPlus := int(float64(img.Bounds().Max.X) * a)

			p1 := image.Point{sumXi, 0}
			p2 := image.Point{sumXi + xPlus, maxY}
			r := image.Rectangle{p1, p2}

			draw.BiLinear.Scale(rgba, r, img, img.Bounds(), draw.Over, nil)
			sumXi += xPlus
		}
	case "mesh":
		ln := len(images)
		var sumX, sumY int

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
			x, y, xPlus, yPlus := 0, 0, 0, 0
			if maxX-img.Bounds().Max.X <= maxY-img.Bounds().Max.Y {
				a := float64(maxX) / float64(img.Bounds().Max.X)
				yPlus = int(float64(img.Bounds().Max.Y) * a)
				xPlus = maxX
				x, y = 0, (maxY-yPlus)/2
			} else {
				a := float64(maxY) / float64(img.Bounds().Max.Y)
				xPlus = int(float64(img.Bounds().Max.X) * a)
				yPlus = maxY
				x, y = (maxY-xPlus)/2, 0
			}

			p1 := image.Point{x + sumX, y + sumY}
			p2 := image.Point{x + sumX + xPlus, y + sumY + yPlus}
			r := image.Rectangle{p1, p2}
			sumX += maxX
			if xyi == ii+1 {
				ii = -1
				sumX = 0
				sumY += maxY
			}
			ii++

			draw.BiLinear.Scale(rgba, r, img, img.Bounds(), draw.Over, nil)
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
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
		return
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 85
	err = jpeg.Encode(out, rgba, &opt)
	if err != nil {
		log.Error(errors.Wrap(err, "TgAlbumToPic jpeg.Encode"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
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

func (s *Service) TgCompress(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}

	var quality int = 10
	switch x.Callback().Data {
	case "cp1":
		quality = 10
	case "cp2":
		quality = 7
	case "cp3":
		quality = 2
	}

	var photosMes []*tb.Message
	albumID := x.Message().ReplyTo.AlbumID
	if albumID == "" {
		photosMes = []*tb.Message{x.Message().ReplyTo}
	} else {
		okLock := s.TG.AlbumsManager.LockAlbum(albumID)
		if !okLock {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Подождите пока другой запрос выполнится.", ShowAlert: true})
			return
		}
		defer s.TG.AlbumsManager.UnLockAlbum(albumID)
		photosMes = s.TG.AlbumsManager.GetAlbum(albumID)
		if len(photosMes) == 0 {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Фото не найдены, попробуйте переслать их в этого в бота.", ShowAlert: true})
			return
		}
		sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })
	}

	pathes := make([]string, 0, len(photosMes))
	newpathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, 777)
	defer os.Remove(dir)
	for _, mes := range photosMes {
		photo := mes.Photo
		path := dir + photo.FileID + ".jpg"
		newpath := dir + util.CreateKey(12) + ".jpg"
		pathes = append(pathes, path)
		newpathes = append(newpathes, newpath)
		err := x.Bot().Download(photo.MediaFile(), path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress Bot.Download"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка скачивания, напишите автору.", ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	for i, img := range images {
		err := imgCompress(img, newpathes[i], quality)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress imgCompress"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка обработки, напишите автору.", ShowAlert: true})
			return
		}
	}

	album := tb.Album{}
	for i, newpath := range newpathes {
		album = append(album, &tb.Photo{File: tb.FromDisk(newpath), Caption: photosMes[i].Caption})
	}
	err := x.SendAlbum(album)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка отправки.", ShowAlert: true})
		return
	}
	x.Respond()
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func imgCompress(img image.Image, out string, quality int) error {
	opt := jpeg.Options{
		Quality: quality,
	}
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &opt)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(out, buf.Bytes(), fs.ModePerm)
}

func (s *Service) TgSeamCarving(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}

	if s.Other.SeamCarvingMutexBL {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Обработка занята, попробуйте позже (делать очередь лень).", ShowAlert: true})
		return
	}
	s.Other.SeamCarvingMutex.Lock()
	s.Other.SeamCarvingMutexBL = true
	defer func() {
		s.Other.SeamCarvingMutexBL = false
		s.Other.SeamCarvingMutex.Unlock()
	}()

	var quality int = 15
	switch x.Callback().Data {
	case "sm1":
		quality = 7
	case "sm2":
		quality = 4
	case "sm3":
		quality = 2
	}

	var photosMes []*tb.Message
	albumID := x.Message().ReplyTo.AlbumID
	if albumID == "" {
		photosMes = []*tb.Message{x.Message().ReplyTo}
	} else {
		okLock := s.TG.AlbumsManager.LockAlbum(albumID)
		if !okLock {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Подождите пока другой запрос выполнится.", ShowAlert: true})
			return
		}
		defer s.TG.AlbumsManager.UnLockAlbum(albumID)
		photosMes = s.TG.AlbumsManager.GetAlbum(albumID)
		if len(photosMes) == 0 {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Фото не найдены, попробуйте переслать их в этого в бота.", ShowAlert: true})
			return
		}
		sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })
	}

	pathes := make([]string, 0, len(photosMes))
	newpathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, 777)
	defer os.Remove(dir)
	for _, mes := range photosMes {
		photo := mes.Photo
		path := dir + photo.FileID + ".jpg"
		newpath := dir + util.CreateKey(12) + ".jpg"
		pathes = append(pathes, path)
		newpathes = append(newpathes, newpath)
		err := x.Bot().Download(photo.MediaFile(), path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgSeamCarving Bot.Download"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка скачивания, напишите автору.", ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgSeamCarving os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgSeamCarving image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка файлов, напишите автору.", ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	for i, img := range images {
		err := imgSeamCarving(img, newpathes[i], quality)
		if err != nil {
			log.Error(errors.Wrap(err, "TgSeamCarving imgCompress"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка обработки, напишите автору.", ShowAlert: true})
			return
		}
	}

	album := tb.Album{}
	for i, newpath := range newpathes {
		album = append(album, &tb.Photo{File: tb.FromDisk(newpath), Caption: photosMes[i].Caption})
	}
	err := x.SendAlbum(album)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка отправки.", ShowAlert: true})
		return
	}
	x.Respond()
	x.Bot().EditReplyMarkup(x.Message(), &tb.ReplyMarkup{InlineKeyboard: delBtn(x.Message().ReplyMarkup.InlineKeyboard, x.Callback().Data)})
	return
}

func imgSeamCarving(img image.Image, out string, quality int) error {
	var err error
	x, y := img.Bounds().Max.X, img.Bounds().Max.Y
	xn, yn := x/quality, y/quality

	img, err = util.ReduceWidth(img, xn)
	if err != nil {
		return err
	}
	img, err = util.ReduceHeight(img, yn)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return err
	}
	return ioutil.WriteFile(out, buf.Bytes(), fs.ModePerm)
}
