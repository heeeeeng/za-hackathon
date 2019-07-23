package main

import (
	"time"
)

type RespRank struct {
	RankInfo RespRankInfo   `json:"current_rank_info"`
	RankList []RespRankList `json:"rank_list_info"`
}

type RespRankInfo struct {
	MissionB   string             `json:"mission_b"`
	RankNum    string             `json:"rank_num"`
	ScoresInfo RespRankInfoScores `json:"scores_info"`
}

type RespRankInfoScores struct {
	TeamName    string `json:"team_name"`
	StatusA     string `json:"a_scores"`
	StatusB     string `json:"b_scores"`
	StatusTotal string `json:"total_scores"`
}

type RespRankList struct {
	RankNum    int            `json:"rank_num"`
	ScoresInfo RespScoresInfo `json:"scores_info"`
}

type RespScoresInfo struct {
	TeamName    string `json:"team_name"`
	StatusA     string `json:"a_scores"`
	StatusB     string `json:"b_scores"`
	StatusTotal string `json:"total_scores"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
	ID          int    `json:"id"`
}

func defaultResp() RespRank {
	return RespRank{
		RankInfo: RespRankInfo{
			ScoresInfo: RespRankInfoScores{
				TeamName: "超级大西瓜",
				StatusA:  "已完成",
				StatusB:  "未完成",
			},
		},
		RankList: []RespRankList{
			RespRankList{
				RankNum: 0,
				ScoresInfo: RespScoresInfo{
					TeamName:   "小东瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().UnixNano() - 10000,
					ID:         2,
				},
			},
			RespRankList{
				RankNum: 1,
				ScoresInfo: RespScoresInfo{
					TeamName:   "超级大西瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().UnixNano() - 100000,
					ID:         3,
				},
			},
			RespRankList{
				RankNum: 2,
				ScoresInfo: RespScoresInfo{
					TeamName:   "大南瓜",
					StatusA:    "已完成",
					CreateTime: time.Now().UnixNano() - 1000000,
					ID:         4,
				},
			},
		},
	}
}

func notFoundResp() RespRank {
	return RespRank{
		RankInfo: RespRankInfo{
			ScoresInfo: RespRankInfoScores{
				TeamName: "暂无",
				StatusA:  "暂无",
				StatusB:  "暂无",
			},
		},
	}
}
