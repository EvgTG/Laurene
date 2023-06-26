package mainpack

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

	// Все сообщения в список, остаётся один поток
	s.Bot.VideoAlbumsManager.Mutex.Lock()
	va, ok := s.Bot.VideoAlbumsManager.Map[userID]
	if !ok {
		va = &VideoAlbum{}
		s.Bot.VideoAlbumsManager.Map[userID] = va
	}
	va.Videos = append(va.Videos, x.Message())

	if s.Bot.VideoAlbumsManager.MapLock[userID] {
		s.Bot.VideoAlbumsManager.Mutex.Unlock()
		return
	} else {
		s.Bot.VideoAlbumsManager.MapLock[userID] = true
	}
	s.Bot.VideoAlbumsManager.Mutex.Unlock()

	mes, err := s.Bot.Send(x.Sender(), s.Bot.Text(x, "vid_wait"))
	if err != nil {
		return
	}

	time.Sleep(time.Second)
	defer func() {
		s.Bot.Delete(mes)
		s.Bot.VideoAlbumsManager.Mutex.Lock()
		delete(s.Bot.VideoAlbumsManager.Map, userID)
		delete(s.Bot.VideoAlbumsManager.MapLock, userID)
		s.Bot.VideoAlbumsManager.Mutex.Unlock()
	}()

	sort.Slice(va.Videos, func(i, j int) bool { return va.Videos[i].ID < va.Videos[j].ID })
	album := tb.Album{}
	for _, videoMes := range va.Videos {
		album = append(album, videoMes.Video)
	}
	if len(album) < 2 {
		x.Send(s.Bot.Text(x, "vid_len1"))
		return
	}
	x.SendAlbum(album)

	return
}
