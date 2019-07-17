package main

import "net/http"

func (c *Controller) InitRouter() {
	http.HandleFunc("/icrsmp/api/hackathon/rank", c.QueryRankInfo)
}
