package main

import "net/http"

type Controller struct {
	//
	module *Module
}

func NewController (m *Module) *Controller {
	return &Controller{
		module:m,
	}
}

func (c *Controller) QueryRankInfo(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		Response(w, http.StatusBadRequest, StatusFail, "miss team name", nil)
		return
	}
	resp := defaultResp()
	Response(w, http.StatusOK, StatusSuccess, "", resp)
	return
}
