package goapm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"sync"

	"github.com/elastic/beats/libbeat/common"
)

const (
	header = "application/json"
)

type APM struct {
	data              common.MapStr
	app               *app
	transactions      map[uuid.UUID]*Transaction
	transactionsMutex *sync.Mutex
	host              string
}

func NewAPM(appName, appVersion string) *APM {
	apm := &APM{
		host:              "http://localhost:8200/",
		transactions:      map[uuid.UUID]*Transaction{},
		transactionsMutex: &sync.Mutex{},
	}
	apm.app = newApp(appName, appVersion)

	return apm
}

func (apm *APM) StartTransaction() *Transaction {
	t := NewTransaction()
	t.Start()
	apm.transactionsMutex.Lock()
	apm.transactions[t.id] = t
	apm.transactionsMutex.Unlock()

	return t
}

func (apm *APM) extractFinishedTransactions() []*Transaction {
	var finished []*Transaction

	apm.transactionsMutex.Lock()
	defer apm.transactionsMutex.Unlock()

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
