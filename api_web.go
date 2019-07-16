package main

import (
	"net/http"
	"time"
)

type RespRank struct {
	RankInfo RespRankInfo   `json:"current_rank_info"`
	RankList []RespRankList `json:"rank_list_info"`
}

type RespRankInfo struct {
	TeamName string `json:"team_name"`
	StatusA  string `json:"a_scores"`
	StatusB  string `json:"b_scores"`
}

type RespRankList struct {
	RankNum    int            `json:"rank_num"`
	ScoresInfo RespScoresInfo `json:"scores_info"`
}

type RespScoresInfo struct {
	TeamName   string `json:"team_name"`
	StatusA    string `json:"a_scores"`
	CreateTime int64  `json:"create_time"`
	ID         int    `json:"id"`
}

func QueryRankInfo(w http.ResponseWriter, r *http.Request) {
	resp := defaultResp()
	Response(w, http.StatusOK, StatusSuccess, "", resp)
}

func defaultResp() RespRank {
	return RespRank{
		RankInfo: RespRankInfo{
			TeamName: "超级大西瓜",
			StatusA:  "已完成",
			StatusB:  "未完成",
		},
		RankList: []RespRankList{
			RespRankList{
				RankNum: 0,
				ScoresInfo: RespScoresInfo{
					TeamName:   "小东瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().Unix() - 10,
					ID:         2,
				},
			},
			RespRankList{
				RankNum: 1,
				ScoresInfo: RespScoresInfo{
					TeamName:   "超级大西瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().Unix() - 100,
					ID:         3,
				},
			},
			RespRankList{
				RankNum: 2,
				ScoresInfo: RespScoresInfo{
					TeamName:   "大南瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().Unix() - 1000,
					ID:         4,
				},
			},
		},
	}
}
