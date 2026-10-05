package api

import (
	"encoding/json"
	"net/http"
)
type PublishRequest struct {
	Topic string `json:"topic"`
	Message string `json:"message"`
}

type PublishResponse struct {
	Message string `json:"message"`
}

type SubscribeRequest struct {
	Topic string `json:"topic"`
}

type SubscribeResponse struct {
	Message string `json:"message"`
}

func PublishHandler(w http.ResponseWriter,r *http.Request){
	var req PublishRequest
	err:=json.NewDecoder(r.Body).Decode(&req)

	if err!=nil{
		http.Error(w,"invalid JSON",http.StatusBadRequest)
		return
	}
	response:= PublishResponse{
		Message: "message published successfully"
	}
	json.NewEncoder(w).Encode(response)
}