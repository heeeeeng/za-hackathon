package main

type TeamInfo struct {
	TeamName string `json:"team_name"`
	Phone    string `json:"phone"`
	//Id     bson.ObjectId `bson:"_id",json:"id"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

type ScoreInfo struct {
	Phone   string `json:"phone",bson:"phone"`
	StatusA string `json:"status_a"`
	StatusB string `json:"status_b"`
	//ID      bson.ObjectId `bson:"_id",json:"id"`
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}

func (r *ScoreInfo) ToResponse(t *TeamInfo) RespScoresInfo {
	resp := RespScoresInfo{
		TeamName:   t.TeamName,
		CreateTime: r.CreateTime,
		UpdateTime: r.UpdateTime,
		StatusB:    r.StatusB,
		StatusA:    r.StatusA,
	}
	return resp
}

type TeamInfoAPIRet struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data []TeamInfoAPIItemRet `json:"data"`
}

type TeamInfoAPIItemRet struct {
	//CellStyleMap      interface{} `json:"cell_style_map"` //we don't need this
	Id           int    `json:"id"`
	School       string `json:"school"`
	TeamName     string `json:"teamName"`
	Type         string `json:"type"` //AI,BC
	CaptainPhone string `json:"captainPhone"`
}

func (t *TeamInfoAPIItemRet) isBlockChainTeam() bool {
	return t.Type == "BC"
}

func (t *TeamInfoAPIItemRet) TeamInfo() TeamInfo {
	return TeamInfo{
		TeamName: t.TeamName,
		//TeamId:t.Id,
		Phone: t.CaptainPhone,
	}
}

type ScoresInfo []ScoreInfo

func (s ScoresInfo) Len() int      { return len(s) }
func (s ScoresInfo) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ScoresInfo) Less(i, j int) bool {
	if s[i].StatusA == "已完成" && s[j].StatusA == "未完成" {
		return true
	}
	if s[i].StatusA == "未完成" && s[j].StatusA == "已完成" {
		return false

	}
	return s[i].UpdateTime < s[j].UpdateTime
}
