package mainpack

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"io/fs"
	"os"
	"sort"
	"time"

	"Laurene/go-log"
	"Laurene/util"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	tb "gopkg.in/telebot.v3"
)

func (s *Service) TgPic(x tb.Context) (errReturn error) {
	if x.Message().AlbumID != "" {
		if s.Bot.AlbumsManager.AddPhotoMes(x.Message().AlbumID, x.Message()) {
			return
		}
		time.Sleep(time.Millisecond * 500)
	}

	if x.Message().AlbumID == "" {
		x.Send(s.Bot.Text(x, "pic"), &tb.SendOptions{ReplyTo: x.Message(), AllowWithoutReply: true}, s.Bot.Markup(x, "pic"))
		return
	}

	x.Send(s.Bot.Text(x, "pic_album"), &tb.SendOptions{ReplyTo: x.Message(), AllowWithoutReply: true}, s.Bot.Markup(x, "pic_album"))
	return
}

func (s *Service) TgAlbumToPic(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}
	defer s.Bot.EditReplyMarkup(x.Message(), delBtn(x.Message().ReplyMarkup, x.Callback().Data))

	albumID := x.Message().ReplyTo.AlbumID
	if albumID == "" {
		return
	}

	okLock := s.Bot.AlbumsManager.LockAlbum(albumID)
	if !okLock {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_lock"), ShowAlert: true})
		return
	}
	defer s.Bot.AlbumsManager.UnLockAlbum(albumID)
	photosMes := s.Bot.AlbumsManager.GetAlbum(albumID)
	if len(photosMes) == 0 {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_len0"), ShowAlert: true})
		return
	}

	sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })

	pathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, fs.ModePerm)
	defer os.RemoveAll(dir)
	for _, mes := range photosMes {
		photo := mes.Photo
		path := dir + photo.FileID + ".jpg"
		pathes = append(pathes, path)
		err := s.Bot.Download(photo.MediaFile(), path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic Bot.Download"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgAlbumToPic image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
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
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err")})
		return
	}

	outPath := dir + util.CreateKey(12) + ".jpg"
	out, err := os.Create(outPath)
	if err != nil {
		log.Error(errors.Wrap(err, "TgAlbumToPic os.Create(outPath)"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}
	defer out.Close()

	var opt jpeg.Options
	opt.Quality = 85
	err = jpeg.Encode(out, rgba, &opt)
	if err != nil {
		log.Error(errors.Wrap(err, "TgAlbumToPic jpeg.Encode"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}

	err = x.Send(&tb.Document{File: tb.FromDisk(outPath), FileName: "pic.jpg", Caption: photosMes[0].Caption}, s.Bot.Markup(x, "picfile_to_pic"))
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}
	x.Respond()
	return
}

func (s *Service) TgFilePicToPic(x tb.Context) (errReturn error) {
	if x.Message() == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err")})
		return
	}
	if x.Message().Document == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err")})
		return
	}

	outPath := "files/temp/" + util.CreateKey(12) + ".jpg"
	defer os.Remove(outPath)
	err := s.Bot.Download(x.Message().Document.MediaFile(), outPath)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err")})
		return
	}

	x.Send(&tb.Photo{File: tb.FromDisk(outPath), Caption: x.Message().Caption})
	s.Bot.EditReplyMarkup(x.Message(), nil)
	x.Respond()
	return
}

func (s *Service) TgCompress(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}
	defer s.Bot.EditReplyMarkup(x.Message(), delBtn(x.Message().ReplyMarkup, x.Callback().Data))

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
		okLock := s.Bot.AlbumsManager.LockAlbum(albumID)
		if !okLock {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_lock"), ShowAlert: true})
			return
		}
		defer s.Bot.AlbumsManager.UnLockAlbum(albumID)
		photosMes = s.Bot.AlbumsManager.GetAlbum(albumID)
		if len(photosMes) == 0 {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_len0"), ShowAlert: true})
			return
		}
		sort.Slice(photosMes, func(i, j int) bool { return photosMes[i].ID < photosMes[j].ID })
	}

	pathes := make([]string, 0, len(photosMes))
	newpathes := make([]string, 0, len(photosMes))
	images := make([]image.Image, 0, len(photosMes))
	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, fs.ModePerm)
	defer os.RemoveAll(dir)
	for _, mes := range photosMes {
		photo := mes.Photo
		path := dir + util.CreateKey(12) + ".jpg"
		newpath := dir + util.CreateKey(12) + ".jpg"
		pathes = append(pathes, path)
		newpathes = append(newpathes, newpath)
		err := s.Bot.Download(photo.MediaFile(), path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress Bot.Download"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}
	}

	for _, path := range pathes {
		imgFile, err := os.Open(path)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress os.Open"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}
		defer imgFile.Close()

		img, _, err := image.Decode(imgFile)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress image.Decode"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}

		images = append(images, img)
	}

	for i, img := range images {
		err := imgCompress(img, newpathes[i], quality)
		if err != nil {
			log.Error(errors.Wrap(err, "TgCompress imgCompress"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}
	}

	album := tb.Album{}
	for i, newpath := range newpathes {
		album = append(album, &tb.Photo{File: tb.FromDisk(newpath), Caption: photosMes[i].Caption})
	}
	err := x.SendAlbum(album)
	if err != nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}
	x.Respond()
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
	return os.WriteFile(out, buf.Bytes(), fs.ModePerm)
}

func (s *Service) TgPicGif(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().ReplyTo == nil {
		return
	}

	mes, err := s.Bot.Send(x.Sender(), "Обработка...")
	if err != nil {
		return
	}
	defer s.Bot.Delete(mes)
	defer s.Bot.EditReplyMarkup(x.Message(), delBtn(x.Message().ReplyMarkup, x.Callback().Data))

	dir := "files/temp/" + util.CreateKey(5) + "/"
	os.Mkdir(dir, fs.ModePerm)
	defer os.RemoveAll(dir)
	pathPic := dir + util.CreateKey(5) + ".jpg"

	err = s.Bot.Download(x.Message().ReplyTo.Photo.MediaFile(), pathPic)
	if err != nil {
		log.Error(errors.Wrap(err, "TgPicGif Bot.Download"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}

	file, err := os.Open(pathPic)
	if err != nil {
		log.Error(errors.Wrap(err, "TgPicGif os.Open"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Error(errors.Wrap(err, "TgPicGif jpeg.Decode"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}

	img = resize.Resize(uint(img.Bounds().Max.X/2), uint(img.Bounds().Max.Y/2), img, resize.Bilinear)

	files := make([]string, 0, 20)
	quality := 25
	for quality >= 1 {
		name := dir + util.CreateKey(5) + ".jpg"
		files = append(files, name)
		err = imgCompress(img, name, quality)
		if err != nil {
			log.Error(errors.Wrap(err, "TgPicGif imgCompress"))
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
			return
		}

		quality -= 1
	}

	nameGif := dir + "outgif.gif"
	err = genGif(files, 30, nameGif)
	if err != nil {
		log.Error(errors.Wrap(err, "TgPicGif genGif"))
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: s.Bot.Text(x, "pic_err"), ShowAlert: true})
		return
	}

	x.Send(&tb.Document{File: tb.FromDisk(nameGif), FileName: "file.gif"})
	x.Respond()
	return
}

func genGif(files []string, delay int, nameOut string) error {
	var frames []*image.Paletted
	var dx = []int{}
	var dy = []int{}
	var newTempImg image.Image

	for i := range files {
		file, err := os.Open(files[i])
		if err != nil {
			return errors.Wrap(err, "os.Open")
		}
		defer file.Close()
		img, err := jpeg.Decode(file)
		if err != nil {
			return errors.Wrap(err, "jpeg.Decode1")
		}

		buf := bytes.Buffer{}
		err = gif.Encode(&buf, img, nil)
		if err != nil {
			return errors.Wrap(err, "gif.Encode1")
		}

		tmpimg, err := gif.Decode(&buf)
		if err != nil {
			return errors.Wrap(err, "gif.Decode1")
		}

		r := tmpimg.Bounds()

		var newX, newY int
		if len(dx) > 0 {
			if dx[i-1] != r.Dx() {
				newX = dx[i-1]
			}
		}

		if len(dy) > 0 {
			if dy[i-1] != r.Dy() {
				newY = dy[i-1]
			}
		}

		if newX > 0 || newY > 0 {
			newTempImg = resize.Resize(uint(newX), uint(newY), tmpimg, resize.Lanczos3)
		}

		dx = append(dx, r.Dx())
		dy = append(dy, r.Dy())

		if newTempImg != nil {
			err = gif.Encode(&buf, newTempImg, nil)
			if err != nil {
				return errors.Wrap(err, "gif.Encode2")
			}

			tempImg, err := gif.Decode(&buf)
			if err != nil {
				return errors.Wrap(err, "gif.Decode2")
			}

			frames = append(frames, tempImg.(*image.Paletted))
		} else {
			frames = append(frames, tmpimg.(*image.Paletted))
		}

	}

	delays := make([]int, len(frames))
	for j := range delays {
		delays[j] = delay
	}

	f, err := os.OpenFile(nameOut, os.O_WRONLY|os.O_CREATE, fs.ModePerm)
	if err != nil {
		errors.Wrap(err, "os.OpenFile")
	}
	defer f.Close()
	err = gif.EncodeAll(f, &gif.GIF{Image: frames, Delay: delays, LoopCount: 0})
	if err != nil {
		return errors.Wrap(err, "gif.EncodeAll")
	}

	return nil
}
