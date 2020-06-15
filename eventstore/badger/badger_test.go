package badger

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/mishudark/triper"
)

type TestEvent struct {
	Name string
	SKU  string
}

var (
	Aid = triper.GenerateUUID()
	cli *Client
)

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatalln(err)
	}

	reg := triper.NewEventRegister()
	reg.Set(&TestEvent{})

	cli, err = NewClient(tmpDir, reg)
	if err != nil {
		log.Fatalln(err)
	}

	result := m.Run()

	if err = os.RemoveAll(tmpDir); err != nil {
		log.Println(err)
	}
	os.Exit(result)
}

func TestClientSave(t *testing.T) {
	events := []triper.Event{
		{
			ID:            triper.GenerateUUID(),
			AggregateID:   Aid,
			AggregateType: "order",
			Version:       1,
			Type:          "test_event",
			Data: TestEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
		{
			ID:            triper.GenerateUUID(),
			AggregateID:   Aid,
			AggregateType: "order",
			Version:       1,
			Type:          "test_event",
			Data: TestEvent{
				Name: "muñeca",
				SKU:  "123",
			},
		},
	}

	err := cli.Save(events, 0)
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestClientLoad(t *testing.T) {
	reg := triper.NewEventRegister()
	reg.Set(&TestEvent{})

	events, err := cli.Load(Aid)
	if err != nil {
		t.Error("expected nil, got", err)
	}

	length := len(events)
	if length != 2 {
		t.Errorf("[events] expected: 2, got: %d", length)
	}
}
