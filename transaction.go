package goapm

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/elastic/beats/libbeat/common"
)

type Transaction struct {
	id        uuid.UUID
	name      string
	typ       string
	timestamp time.Time
	duration  time.Duration // us
}

func NewTransaction() *Transaction {
	return &Transaction{
		id: uuid.NewV4(),
	}
}

func (t *Transaction) toMapStr() common.MapStr {
	return common.MapStr{
		"timestamp": t.timestamp.UTC(),
		"duration":  int64(t.duration / time.Microsecond),
		"id":        t.id,
		"name":      "GET /api/types",
		"type":      "measurement",
	}
}

func (t *Transaction) Start() {
	t.timestamp = time.Now()
}

func (t *Transaction) Stop() {
	t.duration = time.Now().Sub(t.timestamp)
}

func (t *Transaction) Finished() bool {
	return t.duration != 0
}
