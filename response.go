package main

import (
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	StatusSuccess = "success"
	StatusFail    = "failed"
)

func Response(w http.ResponseWriter, code int, status string, msg string, data interface{}) {
	resp := ResponseData{
		Code:    code,
		Status:  status,
		Message: msg,
		Data:    data,
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		Response(w, http.StatusInternalServerError, StatusFail, "", nil)
		return
	}
	w.Write(respJson)
}
