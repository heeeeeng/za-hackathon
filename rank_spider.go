package main

import (
	"github.com/annchain/OG/client/tx_client"
	"github.com/annchain/OG/common"
	"github.com/annchain/OG/common/crypto"
	"github.com/annchain/OG/rpc"

	"time"
)

type RankSpider struct {
	module *Module
	client *tx_client.TxClient
	OgUri  string
	quit   chan bool
}

func NewRankSpider(m *Module, ogrul string, sk crypto.PrivateKey) *RankSpider {
	client := tx_client.NewTxClient(ogrul, false)
	return &RankSpider{
		module: m,
		client: &client,
		OgUri:  ogrul,
		quit:   make(chan bool),
	}
}

func (r *RankSpider) Start() {
	go r.start()
}

func (r *RankSpider) stop() {
	close(r.quit)
}

func (r *RankSpider) start() {
	// TODO
	for {
		select {
		case <-time.After(time.Second * 30):
			r.fetchDataFromOg()

		}
	}
	time.Sleep(time.Second * 5)
}

func (a *RankSpider) fetchDataFromOg() {
	addr := common.RandomAddress()
	req := &rpc.NewQueryContractReq{
		Address: addr.String(),
		Data:    "testdata",
	}
	_ = req
}
