package model

import (
	"fmt"
	"github.com/wesleybits/go-simple-globs/repo"
	"testing"
)

// make a simple in-memory repo that tracks the methods called on it
type Funcall struct {
	Function string
	Args     map[string]string
}

type TestRepo[T any] struct {
	memory map[string]T
	calls  []Funcall
}

// implement the repo interface
func (r *TestRepo[T]) Get(id string) *repo.IdWrap[T] {
	r.calls = append(r.calls, Funcall{
		Function: "Get",
		Args:     map[string]string{"id": id},
	})

	if elt, ok := r.memory[id]; ok {
		return &repo.IdWrap[T]{Id: id, Data: elt}
	} else {
		return nil
	}
}

func (r *TestRepo[T]) GetAll() []repo.IdWrap[T] {
	r.calls = append(r.calls, Funcall{
		Function: "GetAll",
	})

	result := make([]repo.IdWrap[T], 0, len(r.memory))
	for id, elt := range r.memory {
		result = append(result, repo.IdWrap[T]{Id: id, Data: elt})
	}
	return result
}

func (r *TestRepo[T]) Put(id string, elt T) {
	r.calls = append(r.calls, Funcall{
		Function: "Put",
		Args: map[string]string{
			"id":  id,
			"elt": fmt.Sprintf("%v", elt),
		},
	})

	r.memory[id] = elt
}

func (r *TestRepo[T]) Delete(id string) {
	r.calls = append(r.calls, Funcall{
		Function: "Delete",
		Args: map[string]string{
			"id": id,
		},
	})

	delete(r.memory, id)
}

// some test utilities
func (r *TestRepo[T]) MatchGet(id string) bool {
	top := r.calls[0]
	r.calls = r.calls[1:]

	aid, ok := top.Args["id"]
	return top.Function == "Get" && ok && aid == id
}

func (r *TestRepo[T]) MatchGetAll() bool {
	top := r.calls[0]
	r.calls = r.calls[1:]

	return top.Function == "GetAll"
}

func (r *TestRepo[T]) MatchPut(id string, elt T) bool {
	top := r.calls[0]
	r.calls = r.calls[1:]

	aid, idok := top.Args["id"]
	aelt, eltok := top.Args["elt"]
	return top.Function == "Put" && idok && eltok && aid == id && fmt.Sprintf("%v", elt) == aelt
}

func (r *TestRepo[T]) MatchDelete(id string) bool {
	top := r.calls[0]
	r.calls = r.calls[1:]

	aid, idok := top.Args["id"]
	return top.Function == "Delete" && idok && id == aid
}

func (r *TestRepo[T]) Exists(id string) bool {
	_, ok := r.memory[id]
	return ok
}

func (r *TestRepo[T]) ExistsAs(id string, elt T) bool {
	x, ok := r.memory[id]
	return ok && fmt.Sprintf("%v", elt) == fmt.Sprintf("%v", x)
}

// blank out the repo between tests
func (r *TestRepo[T]) Reset() {
	r.memory = map[string]T{}
	r.calls = make([]Funcall, 0, 10)
}

type Data struct {
	Message string
}

func NewData(mesg string) Data {
	return Data{Message: mesg}
}

// tests start here
var testrepo *TestRepo[Data]
var model *Model[Data]

func TestMain(m *testing.M) {
	testrepo = &TestRepo[Data]{calls: make([]Funcall, 0, 10), memory: map[string]Data{}}
	model = NewModel[Data](testrepo)

	m.Run()
}

func TestCreate(t *testing.T) {
	defer testrepo.Reset()

	m1 := NewData("message 1")

	id1 := model.Create(m1)

	if !testrepo.ExistsAs(id1, m1) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m1) {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	defer testrepo.Reset()

	m1 := NewData("message 1")
	id1 := model.Create(m1)
	model.Delete(id1)

	if testrepo.Exists(id1) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m1) {
		t.Fail()
	}

	if !testrepo.MatchDelete(id1) {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	defer testrepo.Reset()

	m1 := NewData("message 1")
	id1 := model.Create(m1)
	m2 := NewData("message 2")
	model.Update(id1, m2)

	if testrepo.ExistsAs(id1, m1) {
		t.Fail()
	}

	if !testrepo.ExistsAs(id1, m2) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m1) {
		t.Fail()
	}

	if !testrepo.MatchDelete(id1) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m2) {
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	defer testrepo.Reset()

	m1 := NewData("message 1")
	m2 := NewData("message 2")
	m3 := NewData("message 3")
	id1 := model.Create(m1)
	id2 := model.Create(m2)
	id3 := model.Create(m3)

	mesgs := model.GetAll()

	if len(mesgs) != 3 {
		t.Fail()
	}

	if !testrepo.ExistsAs(id1, m1) || !testrepo.ExistsAs(id2, m2) || !testrepo.ExistsAs(id3, m3) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m1) || !testrepo.MatchPut(id2, m2) || !testrepo.MatchPut(id3, m3) {
		t.Fail()
	}

	if !testrepo.MatchGetAll() {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	defer testrepo.Reset()

	m1 := NewData("message 1")
	id1 := model.Create(m1)
	res := model.Get(id1)

	if m1 != res.Data {
		t.Fail()
	}

	if !testrepo.ExistsAs(id1, m1) {
		t.Fail()
	}

	if !testrepo.MatchPut(id1, m1) {
		t.Fail()
	}

	if !testrepo.MatchGet(id1) {
		t.Fail()
	}
}
