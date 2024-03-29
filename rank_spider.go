package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

const (
	ScoreInfoCollectionName = "register_status"
	TeamInfoCollectionName  = "team_info"
)

type RankSpider struct {
	module              *Module
	client              *http.Client
	OgUri               string
	webUri              string
	quit                chan bool
	contractAddr        string
	ScoreInfoCollection *mgo.Collection
	TeamInfoCollection  *mgo.Collection
}

func NewRankSpider(m *Module, ogrul string, contractAddr string, webUri string) *RankSpider {
	return &RankSpider{
		module: m,
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: http.DefaultTransport,
		},
		OgUri:        ogrul,
		quit:         make(chan bool),
		contractAddr: contractAddr,
		webUri:       webUri,
	}
}

func (r *RankSpider) Start() {
	r.ScoreInfoCollection = r.module.GetCollection(ScoreInfoCollectionName)
	r.TeamInfoCollection = r.module.GetCollection(TeamInfoCollectionName)
	go r.start()
}

func (r *RankSpider) Stop() {
	close(r.quit)
}

func (r *RankSpider) start() {
	//r.module.db.DropDatabase()
	//return
	r.fetchTeamInfo()
	// TODO
	for {
		select {
		case <-time.After(time.Second * 30):
			r.fetchDataFromOg()
		case <-time.After(time.Second * 19 * 19 * 2):
			r.fetchTeamInfo()
		case <-r.quit:
			logrus.Info("got quit signal , stopping")
		}

	}
}

type NewQueryContractReq struct {
	Address string `json:"address"`
	Data    string `json:"data"`
}

func (a *RankSpider) fetchDataFromOg() {
	//todo get phone list from db
	var TeamInfos []TeamInfo
	for _, team := range TeamInfos {
		//todo get result from db
		var score ScoreInfo
		err := a.ScoreInfoCollection.Find(bson.M{"phone": team.Phone}).One(&score)
		if err != nil {
			logrus.WithField("not found score fore team ", team).WithError(err).Error("should never happen")
			continue
		}
		if score.StatusA =="已完成" {
			//already registered
			continue
		}
		//if not found
		ok, err := a.getRegisterStatusFromOg(team.Phone)
		if err != nil {
			logrus.WithError(err).Warn("get response error")
		}
		if err == nil && ok {
			//wrie ok
			score.StatusA = "已完成"
			err = a.ScoreInfoCollection.Update(bson.M{"phone": score.Phone}, &score)
			if err != nil {
				logrus.WithField("score ", score).WithError(err).Error("update data err ")
			}
		}
	}
	return
}

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func (a *RankSpider) calculateDate(phone string) string {
	code := "hello hacker: "
	hash := Keccak256([]byte(code + phone))

	parameter := fmt.Sprintf("%x", hash)

	finalData := "8e7d4b1d" + parameter
	return finalData
}

func (a *RankSpider) fetchTeamInfo() {
	teams, err := a.getTeamInfoFromAPI()
	if err != nil {
		logrus.WithError(err).Warn("get response error")
		return
	}
	for _, team := range teams {
		if !team.isBlockChainTeam() {
			continue
		}
		var teamInfo TeamInfo
		//fmt.Println(a.TeamInfoCollection,"139")
		err := a.TeamInfoCollection.Find(bson.M{"phone": team.CaptainPhone}).One(&teamInfo)
		if err == nil {
			logrus.WithField("team ", team).Debug("already have this team info")
			continue
		}else {
			logrus.WithField("team ", teamInfo.Phone).WithError(err).Debug("not found team data")
			//var teaminfos []TeamInfo
			//a.TeamInfoCollection.Find(bson.M{}).All(&teaminfos)
			//logrus.Debug("all ",  teaminfos)
		}
		teamInfo = team.TeamInfo()
		//insert score first
		newScore := &ScoreInfo{
			Phone:      teamInfo.Phone,
			StatusA:    "未完成",
			//ID:         bson.NewObjectId(),
			CreateTime: time.Now().UnixNano(),
			UpdateTime: time.Now().Unix(),
		}
		err = a.ScoreInfoCollection.Insert(newScore)
		if err != nil {
			logrus.WithField("score ", newScore).WithError(err).Error("insert data err ")
			continue
		}
		teamInfo.CreateTime = time.Now().Unix()
		teamInfo.UpdateTime = time.Now().Unix()
		err = a.TeamInfoCollection.Insert(&teamInfo)
		if err != nil {
			logrus.WithField("team", team).WithError(err).Error("insert failed")
			continue
		}
		logrus.WithField("request score ", newScore).Debug("inserted data")
	}
	return

}

func (a *RankSpider) getTeamInfoFromAPI() ([]TeamInfoAPIItemRet, error) {
	req, err := http.NewRequest("POST", a.webUri, nil)
	resp, err := a.client.Do(req)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	//now := time.Now()
	defer resp.Body.Close()
	resDate, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Warn("got response err")
		return nil, err
	}
	str := string(resDate)
	if err != nil {
		logrus.WithError(err).Warn(str)
		return nil, err
	}
	if resp.StatusCode != 200 {
		//panic( resp.StatusCode)
		logrus.WithField(" code", resp.StatusCode).Warn("response error")
		return nil, err
	}
	var teaminfoResp TeamInfoAPIRet
	err = json.Unmarshal(resDate, &teaminfoResp)
	if err != nil {
		//fmt.Println("encode nonce errror ", err)
		return nil, err
	}
	return teaminfoResp.Data, nil
}

func (a *RankSpider) getRegisterStatusFromOg(phone string) (bool, error) {
	request := &NewQueryContractReq{
		Address: a.contractAddr,
		Data:    a.calculateDate(""),
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
		return false, err
	}
	//now := time.Now()
	defer resp.Body.Close()
	resDate, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Warn("got response err")
		return false, err
	}
	str := string(resDate)
	if err != nil {
		logrus.WithError(err).Warn(str)
		return false, err
	}
	if resp.StatusCode != 200 {
		//panic( resp.StatusCode)
		logrus.WithField(" code", resp.StatusCode).Warn("response error")
		return false, err
	}
	var stateResp struct {
		Data bool `json:"data"`
	}
	err = json.Unmarshal(resDate, &stateResp)
	if err != nil {
		//fmt.Println("encode nonce errror ", err)
		return false, err
	}
	return stateResp.Data, nil
}


func (r *RankSpider)GetRankInfo(teamName string )RespRank {
	var resp RespRank
	var team TeamInfo
	//fmt.Println(r.TeamInfoCollection,"251")
	err := r.TeamInfoCollection.Find(bson.M{"teamname":teamName}).One(&team)
	if err!=nil  {
		logrus.WithError(err).Warn("team data not found")
		resp.RankInfo =  RespRankInfo{
			ScoresInfo: RespRankInfoScores{
				TeamName: "暂无",
				StatusA:  "暂无",
				StatusB:  "暂无",
			},
		}
	}else {
		var score ScoreInfo
		err := r.ScoreInfoCollection.Find(bson.M{"phone": team.Phone}).One(&score)
		if err!=nil {
			logrus.WithError(err).Warn("score data not found")
		}else {
			resp.RankInfo = RespRankInfo{
				ScoresInfo: RespRankInfoScores{
					TeamName: teamName,
					StatusA:  score.StatusA,
					StatusB:  score.StatusB,
				},
			}
		}
	}
	var scors []ScoreInfo
	//bson.M{"statusa": "已完成"}
	err = r.ScoreInfoCollection.Find(bson.M{}).All(&scors)
	if err!=nil || len(scors) ==0 {
		return resp
	}
	sort.Sort(ScoresInfo(scors))
	for i:=0;i<len(scors)&& i<10 ;i++ {
		err := r.TeamInfoCollection.Find(bson.M{"phone":scors[i].Phone}).One(&team)
		if err!=nil {
			logrus.WithError(err).Warn("data not found")
		}
		rankList := RespRankList{
			RankNum:i,
			ScoresInfo: RespScoresInfo{
				TeamName: team.TeamName,
				StatusA:  scors[i].StatusA,
				CreateTime: scors[i].UpdateTime,

			},
		}
		resp.RankList = append(resp.RankList,rankList)
	}
	return resp
}