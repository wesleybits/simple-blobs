package repo

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

type Payload string

var repo *InMemoryRepo[Payload]

func TestMain(m *testing.M) {
	wg := &sync.WaitGroup{}
	ctx, stop := context.WithCancel(context.Background())
	defer func() {
		stop()
		wg.Wait()
	}()

	repo = NewInMemoryRepo[Payload](ctx)
	repo.Start(wg)

	m.Run()
}

func TestMessages(t *testing.T) {
	message1, id1 := Payload("message 1"), "1"
	message2, id2 := Payload("message 2"), "2"
	fmt.Printf("Putting %s: %v\n", id1, message1)
	repo.Put(id1, message1)
	fmt.Printf("Putting %s: %v\n", id2, message2)
	repo.Put(id2, message2)

	fmt.Printf("Getting %s\n", id1)
	storedMessage1 := repo.Get(id1)
	if storedMessage1 == nil {
		t.Log("expected stored message to be retrivable")
		t.Fail()
	}

	if *storedMessage1 != message1 {
		t.Logf("expected \"%s\", found \"%s\"", message1, *storedMessage1)
		t.Fail()
	}

	fmt.Printf("Deleting %s\n", id2)
	repo.Delete(id2)
	storedMessage2 := repo.Get(id2)
	if storedMessage2 != nil {
		t.Log("failed to delete message")
		t.Fail()
	}

	fmt.Printf("Getting %s\n", id1)
	storedMessage1 = repo.Get(id1)
	if storedMessage1 == nil {
		t.Log("non-deleted message expected to be retrivable")
		t.Fail()
	}

	if *storedMessage1 != message1 {
		t.Logf("expected \"%s\", found \"%s\"", message1, *storedMessage1)
		t.Fail()
	}
}
