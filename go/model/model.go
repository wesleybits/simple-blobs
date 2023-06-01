package model

import (
	"github.com/google/uuid"
	"github.com/wesleybits/go-simple-globs/repo"
)

type ModelSpec[T any] interface {
	Exists(id string) bool
	Get(id string) *repo.IdWrap[T]
	GetAll() []repo.IdWrap[T]
	Create(elt T) string
	Update(id string, elt T)
	Delete(id string)
}

type Model[T any] struct {
	repo repo.RepoSpec[T]
}

func NewModel[T any](repo repo.RepoSpec[T]) *Model[T] {
	return &Model[T]{repo: repo}
}

func (model *Model[T]) Exists(id string) bool {
	return model.repo.Get(id) != nil
}

func (model *Model[T]) Get(id string) *repo.IdWrap[T] {
	return model.repo.Get(id)
}

func (model *Model[T]) GetAll() []repo.IdWrap[T] {
	return model.repo.GetAll()
}

func (model *Model[T]) Create(elt T) string {
	id := uuid.New().String()
	model.repo.Put(id, elt)
	return id
}

func (model *Model[T]) Update(id string, elt T) {
	model.repo.Delete(id)
	model.repo.Put(id, elt)
}

func (model *Model[T]) Delete(id string) {
	model.repo.Delete(id)
}
