package repository

import (
	"sync"
	"time"
)

type Repository struct {
	mu    *sync.RWMutex
	links map[int]map[string]string // словарь сайтов. ключ - уникальный номер запроса, значение - словарь с сайтами в виде "https://google.com": "(not) available"
}

func NewRepository() *Repository {
	return &Repository{
		links: make(map[int]map[string]string, 0),
		mu:    &sync.RWMutex{},
	}
}

// AddLinks - добавление ссылок в кэш
func (r *Repository) AddLinks(newLinks map[string]string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	linksID := int(time.Now().Unix())

	r.links[linksID] = newLinks

	return linksID
}

// GetLinksByLinksID - получение ссылок из кэша по id
func (r *Repository) GetLinksByLinksID(linksID int) map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.links[linksID]
}
