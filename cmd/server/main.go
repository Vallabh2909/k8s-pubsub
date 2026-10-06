package main

import (
	"fmt"
	"net/http"
	. "github.com/vallabh2909/k8s-pubsub/internal/api"
	"encoding/json"
)

func main() {
	Topics:=make(map[string]*Topic)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "k8s-pubsub server")
	})
	
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var req TopicRequest
		if r.Method != http.MethodPost {
			http.Error(w,"method not allowed",http.StatusMethodNotAllowed)
			return
		}
		err:=json.NewDecoder(r.Body).Decode(&req)

		if err!=nil {
			http.Error(w,"Invalid JSON",http.StatusBadRequest)
		}
		
		Topics[req.Name]=&Topic{
			Name: req.Name,
			Channel: make(chan PublishRequest, 1000),
			Subscribers: nil,
		}
		go Topics[req.Name].Start()
		w.Header().Set("Content-type","application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TopicResponse{Message:"Topic Registered"})
	})

	http.HandleFunc("/publish",func(w http.ResponseWriter,r *http.Request){
		var req PublishRequest
		if r.Method != http.MethodPost {
			http.Error(w,"method not allow",http.StatusMethodNotAllowed)
			return 
		}
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w,"Invalid JSON",http.StatusBadRequest)
			return
		}

		if topic,ok:=Topics[req.Topic]; ok {	
			topic.Channel<-req
		}else{
			http.Error(w,"topic not found",http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type","application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TopicResponse{Message:"Message Published"})
	})

	http.HandleFunc("/subscribe",func(w http.ResponseWriter, r *http.Request){
		var req SubscribeRequest
		if r.Method != http.MethodGet {
			http.Error(w,"method not allow",http.StatusMethodNotAllowed)
			return 
		}
		err:=json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w,"Invalid JSON",http.StatusBadRequest)
			return
		}
		if topic,ok:= Topics[req.Topic]; ok {
			channel:=make(chan PublishRequest,1000)
			topic.Subscribers=append(topic.Subscribers,channel)
			w.Header().Set("Content-type","application/json")
			w.WriteHeader(http.StatusCreated)
			flusher := w.(http.Flusher)
			for {
				
				val:=<-channel
				json.NewEncoder(w).Encode(SubscribeResponse{
					Event: val.Event,
					Message: val.Message,
				})
				flusher.Flush()
			}
		}else{
			http.Error(w,"topic not found",http.StatusBadRequest)
			return
		}
		

	})
	
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
