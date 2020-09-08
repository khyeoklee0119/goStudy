package model

import "sync"

type repository map[string]Model

var mutex = &sync.Mutex{}
var repo = repository{}

func CreateRepository() Repository {
	return repo
}

func (repo repository) Store(key string, model Model) error {
	mutex.Lock()
	repo[key] = model
	mutex.Unlock()
	return nil
}

func (repo repository) Get(key string) Model {
	return repo[key]
}

func (repo repository) Size() int {
	return len(repo)
}

func (repo repository) Remove(key string) {
	delete(repo, key)
}

func (repo repository) HasModel(key string) bool {
	_, ok := repo[key]
	return ok
}

func (repo repository) IsEmpty() bool {
	return len(repo) == 0
}

type Repository interface {
	Store(key string, model Model) error
	Get(key string) Model
	Size() int
	Remove(key string)
	IsEmpty() bool
	HasModel(key string) bool
}
