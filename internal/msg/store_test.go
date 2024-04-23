package msg

import (
	"testing"
	"time"
)

func TestInMemoryStoreInsert(t *testing.T) {
	store := NewInMemoryMessageStore()

	m := Message{
		Author:     "Tester.testington",
		Title:      "Test Message",
		Content:    "This is a test message",
		Topic:      []string{"test", "unit"},
		ExpiryUnix: time.Now().Add(time.Second * 4).Unix(),
	}

	store.Insert(m)

	if len(store.Messages) != 1 {
		t.Error("Expected store to have 1 message")
	} else {
		t.Logf("Store has %d messages", len(store.Messages))
		messageJson, messageErr := store.Messages[0].ToJSON(true)
		if messageErr == nil {
			t.Logf("Message:\n %s", messageJson)
		}
	}
}

func TestInMemoryStoreList(t *testing.T) {
	store := NewInMemoryMessageStore()

	m := Message{
		Author:     "Tester.testington",
		Title:      "Test Message",
		Content:    "This is a test message",
		Topic:      []string{"test", "unit"},
		ExpiryUnix: time.Now().Add(time.Second * 4).Unix(),
	}

	store.Insert(m)

	all, err := store.List()
	if err != nil {
		t.Errorf("Error occured in store.List : %s", err)
	}

	if len(all) != 1 {
		t.Error("Expected store to list 1 message")
	}
}

func TestInMemoryStorePrune(t *testing.T) {
	store := NewInMemoryMessageStore()

	m := Message{
		Author:     "Tester.testington",
		Title:      "Test Message",
		Content:    "This is a test message",
		Topic:      []string{"test", "unit"},
		ExpiryUnix: time.Now().Add(time.Second * 4).Unix(),
	}

	store.Insert(m)

	time.Sleep(time.Second * 5)
	store.Prune()

	if len(store.Messages) != 0 {
		t.Fatal("Expected store to have 0 messages after pruning")
	}
}
