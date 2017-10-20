package goapm

import "github.com/elastic/beats/libbeat/common"

type payload struct {
	Transactions []*Transaction
	app          app
}

func (p *payload) toMapStr() common.MapStr {
	// TODO: get transactions

	transactions := []common.MapStr{}

	for _, t := range p.Transactions {
		// Only completed transactions are sent
		if !t.Finished() {
			continue
		}
		transactions = append(transactions, t.toMapStr())
		// TODO: remove from transactions which are in the app queue
	}

	return common.MapStr{
		"app":          p.app.toMapStr(),
		"transactions": transactions,
	}
}

type app struct {
	Name    string
	Version string
}

func newApp(name, version string) *app {
	return &app{
		Name:    name,
		Version: version,
	}
}

func (a *app) toMapStr() common.MapStr {
	return common.MapStr{
		"name":    a.Name,
		"version": a.Version,
		"agent": common.MapStr{
			"name":    "ruflin-golang",
			"version": "0.0.1",
		},
	}
}
