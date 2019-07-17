package main

import "net/http"

type Controller struct {
	//
	module *Module
}

func NewController(m *Module) *Controller {
	return &Controller{
		module: m,
	}
}

func (c *Controller) QueryRankInfo(w http.ResponseWriter, r *http.Request) {
	corsOrigin(&w)

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		Response(w, http.StatusBadRequest, StatusFail, "miss team name", nil)
		return
	}
	resp := defaultResp()
	Response(w, http.StatusOK, StatusSuccess, "", resp)
	return
}

func corsOrigin(wP *http.ResponseWriter) {
	(*wP).Header().Set("Access-Control-Allow-Origin", "*")
	(*wP).Header().Set("Access-Control-Allow-Methods", "*")
	(*wP).Header().Set("Access-Control-Allow-Headers",
		"Access-Control-Allow-Headers, "+
			"Access-Control-Allow-Methods, "+
			"Access-Control-Allow-Origin, "+
			"Content-Length, "+
			"Content-Type, "+
			"Date, "+
			"Access-Control-Request-Headers, "+
			"Access-Control-Request-Method, "+
			"Origin, "+
			"Referer, "+
			"User-Agent, "+
			"Accept, "+
			"X-Requested-With")
}
