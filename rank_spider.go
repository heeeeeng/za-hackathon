package main

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type RankSpider struct {
	module *Module
	client *http.Client
	OgUri  string
	quit   chan bool
	contractAddr string
}

func NewRankSpider(m *Module, ogrul string,contractAddr string) *RankSpider {
	return &RankSpider{
		module: m,
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: http.DefaultTransport,
		},
		OgUri:  ogrul,
		quit:   make(chan bool),
		contractAddr:contractAddr,
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
type NewQueryContractReq struct {
	Address string `json:"address"`
	Data    string `json:"data"`
}

func (a *RankSpider) fetchDataFromOg() {
	//todo get phone list from db
	var TeamInfos []TeamInfo
	for _,team := range TeamInfos {
		//todo get result from db
		//if not found
		ok, err := a.getRegisterStatusFromOg(team.TeamName)
		if err != nil {
			//wrie ok
			_ = ok
		}
	}
}


func (a *RankSpider)calculateDate(phone string) string {
	return ""
}

func (a *RankSpider)getRegisterStatusFromOg( phone string ) (bool, error){
	a.calculateDate(phone)
	request := &NewQueryContractReq{
		Address: a.contractAddr,
		Data:  a.calculateDate(""),
	}
		data, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}
		url := a.OgUri + "/" + "query_contract"
		req, err := http.NewRequest("POST", url, bytes.NewReader(data))
		resp, err := a.client.Do(req)
		if err != nil {
			//fmt.Println(err)
			return false ,err
		}
		//now := time.Now()
		defer resp.Body.Close()
		resDate, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).Warn("got response err")
			return  false ,err
		}
		str := string(resDate)
		if err != nil {
			logrus.WithError(err).Warn(str)
			return  false ,err
		}
		if resp.StatusCode != 200 {
			//panic( resp.StatusCode)
			logrus.WithField (" code" ,resp.StatusCode).Warn("response error")
			return false ,err
		}
		var nonceResp struct {
			Data bool `json:"data"`
		}
		err = json.Unmarshal(resDate,&nonceResp)
		if err != nil {
			//fmt.Println("encode nonce errror ", err)
			return false,err
		}
		return nonceResp.Data, nil
	}

