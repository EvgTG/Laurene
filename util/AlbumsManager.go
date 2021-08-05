package util

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"sync"
	"time"
)

type AlbumsManager struct {
	mt         sync.Mutex
	albums     map[string][]*tb.Photo
	albumsTime []albumTime
}

type albumTime struct {
	id string
	tm time.Time
}

func NewAlbumsManager() *AlbumsManager {
	a := &AlbumsManager{
		mt:         sync.Mutex{},
		albums:     make(map[string][]*tb.Photo, 0),
		albumsTime: make([]albumTime, 0),
	}

	go a.delAlbumsHour()
	return a
}

func (a *AlbumsManager) AddPhoto(albumID string, photo *tb.Photo) bool {
	a.mt.Lock()
	defer a.mt.Unlock()

	_, ok := a.albums[albumID]
	if !ok {
		a.albums[albumID] = make([]*tb.Photo, 0, 10)
		a.albumsTime = append(a.albumsTime, albumTime{id: albumID, tm: time.Now()})
	}
	a.albums[albumID] = append(a.albums[albumID], photo)

	return ok
}

func (a *AlbumsManager) GetAlbum(albumID string) []*tb.Photo {
	a.mt.Lock()
	defer a.mt.Unlock()

	photos := make([]*tb.Photo, 0)
	_, ok := a.albums[albumID]

	if ok {
		photos = append(photos, a.albums[albumID]...)
	}

	return photos
}

func (a *AlbumsManager) delAlbumsHour() {
	for range time.Tick(time.Minute * 10) {
		go func() {
			a.mt.Lock()
			defer a.mt.Unlock()

			ids := make([]string, 0)

			for _, album := range a.albumsTime {
				if time.Since(album.tm).Hours() > 1 {
					ids = append(ids, album.id)
					delete(a.albums, album.id)
				}
			}

			for _, id := range ids {
				for i, album := range a.albumsTime {
					if album.id == id {
						a.albumsTime = append(a.albumsTime[:i], a.albumsTime[i+1:]...)
						break
					}
				}
			}
		}()
	}
}
