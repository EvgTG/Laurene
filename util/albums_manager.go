package util

import (
	tb "gopkg.in/tucnak/telebot.v3"
	"sync"
	"time"
)

type AlbumsManager struct {
	mt           sync.Mutex
	albums       map[string][]*tb.Message
	albumsLocker []string
	albumsTime   []albumTime
}

type albumTime struct {
	id string
	tm time.Time
}

func NewAlbumsManager() *AlbumsManager {
	a := &AlbumsManager{
		mt:           sync.Mutex{},
		albums:       make(map[string][]*tb.Message, 0),
		albumsLocker: make([]string, 0),
		albumsTime:   make([]albumTime, 0),
	}

	go a.delAlbumsHour()
	return a
}

func (a *AlbumsManager) Len() int {
	return len(a.albums)
}

func (a *AlbumsManager) AddPhotoMes(albumID string, photo *tb.Message) bool {
	a.mt.Lock()
	defer a.mt.Unlock()

	_, ok := a.albums[albumID]
	if !ok {
		a.albums[albumID] = make([]*tb.Message, 0, 10)
		a.albumsTime = append(a.albumsTime, albumTime{id: albumID, tm: time.Now()})
	}
	a.albums[albumID] = append(a.albums[albumID], photo)

	return ok
}

func (a *AlbumsManager) GetAlbum(albumID string) []*tb.Message {
	a.mt.Lock()
	defer a.mt.Unlock()

	photos := make([]*tb.Message, 0)
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

func (a *AlbumsManager) LockAlbum(albumID string) bool { // true - альбом заблокирован и в использовании
	a.mt.Lock()
	defer a.mt.Unlock()

	for _, s := range a.albumsLocker {
		if s == albumID {
			return false
		}
	}

	a.albumsLocker = append(a.albumsLocker, albumID)

	return true
}

func (a *AlbumsManager) UnLockAlbum(albumID string) {
	a.mt.Lock()
	defer a.mt.Unlock()

	for i, s := range a.albumsLocker {
		if s == albumID {
			a.albumsLocker = append(a.albumsLocker[:i], a.albumsLocker[i+1:]...)
			return
		}
	}
}
