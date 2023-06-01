package repo

import (
	"context"
	"sync"
	"time"
)

type IdWrap[T any] struct {
	Data T
	Id   string
}

type RepoSpec[T any] interface {
	Get(id string) *IdWrap[T]
	GetAll() []IdWrap[T]
	Put(id string, elt T)
	Delete(id string)
}

type InMemoryRepo[T any] struct {
	ctx        context.Context
	memory     map[string]T
	getchan    chan Get[T]
	getallchan chan GetAll[T]
	putchan    chan Put[T]
	deletechan chan Delete
	running    bool
}

type Get[T any] struct {
	outchan chan *IdWrap[T]
	id      string
}

type GetAll[T any] chan *IdWrap[T]

type Put[T any] struct {
	id   string
	data T
}

type Delete string

func NewInMemoryRepo[T any](ctx context.Context) *InMemoryRepo[T] {
	return &InMemoryRepo[T]{
		ctx:        ctx,
		memory:     map[string]T{},
		getchan:    make(chan Get[T]),
		getallchan: make(chan GetAll[T]),
		putchan:    make(chan Put[T]),
		deletechan: make(chan Delete),
		running:    false,
	}
}

func (repo *InMemoryRepo[T]) Get(id string) *IdWrap[T] {
	outchan := make(chan *IdWrap[T], 1)
	repo.getchan <- Get[T]{outchan, id}

	select {
	case elt := <-outchan:
		return elt
	case <-time.After(200 * time.Millisecond):
		return nil
	}
}

func (repo *InMemoryRepo[T]) GetAll() []IdWrap[T] {
	result := make([]IdWrap[T], 0, 50)
	outchan := make(chan *IdWrap[T])
	repo.getallchan <- GetAll[T](outchan)

	select {
	case firstelt := <-outchan:
		result = append(result, *firstelt)
		for restelt := range outchan {
			if restelt == nil {
				break
			}
			result = append(result, *restelt)
		}
		return result
	case <-time.After(200 * time.Millisecond):
		return result
	}
}

func (repo *InMemoryRepo[T]) Put(id string, elt T) {
	repo.putchan <- Put[T]{id, elt}
}

func (repo *InMemoryRepo[T]) Delete(id string) {
	repo.deletechan <- Delete(id)
}

func (repo *InMemoryRepo[T]) worker(wg *sync.WaitGroup) {
	wg.Add(1)

	get := func(cmd Get[T]) {
		elt, ok := repo.memory[cmd.id]
		if !ok {
			cmd.outchan <- nil
		} else {
			cmd.outchan <- &IdWrap[T]{Id: cmd.id, Data: elt}
		}
	}

	getall := func(cmd GetAll[T]) {
		for id, elt := range repo.memory {
			cmd <- &IdWrap[T]{Id: id, Data: elt}
		}
		cmd <- nil
	}

	put := func(cmd Put[T]) {
		repo.memory[cmd.id] = cmd.data
	}

	del := func(cmd Delete) {
		delete(repo.memory, string(cmd))
	}

	for {
		select {
		case <-repo.ctx.Done():
			wg.Done()
			return

		case command := <-repo.getchan:
			get(command)
		case command := <-repo.getallchan:
			getall(command)
		case command := <-repo.putchan:
			put(command)
		case command := <-repo.deletechan:
			del(command)
		}
	}
}

func (repo *InMemoryRepo[T]) Start(wg *sync.WaitGroup) {
	if repo.running {
		return
	}

	repo.running = true
	go repo.worker(wg)
}
