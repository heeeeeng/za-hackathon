package main

import "net/http"

func InitRouter() {
	http.HandleFunc("/icrsmp/api/hackathon/rank", QueryRankInfo)
}
