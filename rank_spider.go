package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/annchain/OG/common/crypto"
	_ "github.com/annchain/OG/rpc"
	"io/ioutil"
	"net/http"
	"time"
)

type RankSpider struct {
	module *Module
	client  *http.Client
	OgUri  string
	quit   chan bool
}

func NewRankSpider(m *Module , sk crypto.PrivateKey) *RankSpider {
	return &RankSpider{
		module: m,
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: http.DefaultTransport,
		},
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

}

func (a *RankSpider) queryDataFromOG(url string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	resp, err := a.client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//now := time.Now()
	defer resp.Body.Close()
	resDate, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	str := string(resDate)
	if err != nil {
		fmt.Println(str, err)
		return "", err
	}
	if resp.StatusCode != 200 {
		//panic( resp.StatusCode)
		fmt.Println(resp.StatusCode)
		return "", errors.New(resp.Status)
	}
	var respStruct struct {
		Data string `json:"data"`
	}
	err = json.Unmarshal(resDate, &respStruct)
	if err != nil {
		//fmt.Println(str, err)
		return "", err
	}
	return respStruct.Data, nil

}

