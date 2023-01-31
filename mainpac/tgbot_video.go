package mainpac

import (
	tb "gopkg.in/telebot.v3"
	"sort"
	"sync"
	"time"
)

type VideoAlbum struct {
	Videos []*tb.Message
}

type VideoAlbumsManager struct {
	Map     map[int64]*VideoAlbum // map[userID]
	MapLock map[int64]bool        // map[userID]
	Mutex   sync.Mutex
}

func (s *Service) TgVideoComb(x tb.Context) (errReturn error) {
	userID := x.Sender().ID

	s.TG.VideoAlbumsManager.Mutex.Lock()
	va, ok := s.TG.VideoAlbumsManager.Map[userID]
	if !ok {
		va = &VideoAlbum{}
		s.TG.VideoAlbumsManager.Map[userID] = va
	}
	va.Videos = append(va.Videos, x.Message())

	if s.TG.VideoAlbumsManager.MapLock[userID] {
		s.TG.VideoAlbumsManager.Mutex.Unlock()
		return
	} else {
		s.TG.VideoAlbumsManager.MapLock[userID] = true
	}
	s.TG.VideoAlbumsManager.Mutex.Unlock()

	mes, err := x.Bot().Send(x.Sender(), "Подождите...")
	if err != nil {
		return
	}

	time.Sleep(time.Second)
	s.TG.VideoAlbumsManager.Mutex.Lock()
	defer func() {
		x.Bot().Delete(mes)
		delete(s.TG.VideoAlbumsManager.Map, userID)
		delete(s.TG.VideoAlbumsManager.MapLock, userID)
		s.TG.VideoAlbumsManager.Mutex.Unlock()
	}()

	sort.Slice(va.Videos, func(i, j int) bool { return va.Videos[i].ID < va.Videos[j].ID })
	album := tb.Album{}
	for _, videoMes := range va.Videos {
		album = append(album, videoMes.Video)
	}
	if len(album) < 2 {
		x.Send("Видео только одно.")
		return
	}
	x.SendAlbum(album)

	return
}
