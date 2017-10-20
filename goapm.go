package goapm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/elastic/beats/libbeat/common"
)

const (
	header = "application/json"
)

type APM struct {
	data         common.MapStr
	app          *app
	transactions map[uuid.UUID]*Transaction
	host         string
}

func NewAPM(appName, appVersion string) *APM {
	apm := &APM{
		host:         "http://localhost:8200/",
		transactions: map[uuid.UUID]*Transaction{},
	}
	apm.app = newApp(appName, appVersion)

	return apm
}

func (apm *APM) StartTransaction() *Transaction {
	t := NewTransaction()
	t.Start()
	apm.transactions[t.id] = t
	return t
}

func (apm *APM) extractFinishedTransactions() []*Transaction {
	var finished []*Transaction

	for _, t := range apm.transactions {
		if t.Finished() {
			finished = append(finished, t)
			// Remove from list
			delete(apm.transactions, t.id)
		}
	}

	return finished
}

func (a *APM) send() error {

	payload := payload{
		app:          *a.app,
		Transactions: a.extractFinishedTransactions(),
	}

	json, err := json.Marshal(payload.toMapStr())
	if err != nil {
		return err
	}

	resp, err := http.Post(a.host+"v1/transactions", header, bytes.NewReader(json))
	if err != nil {
		return err
	}

	fmt.Printf("%+v", resp)

	return nil
}
