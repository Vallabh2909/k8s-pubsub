package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	. "github.com/vallabh2909/k8s-pubsub/internal/api"
	"github.com/vallabh2909/k8s-pubsub/internal/broker"
)

func GenerateID(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return "sub-" + string(b)
}
func main() {
	broker := broker.NewBroker()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "k8s-pubsub server")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var req TopicRequest
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&req)

		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = broker.Register(req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TopicResponse{Message: "Topic Registered"})
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var req PublishRequest
		if r.Method != http.MethodPost {
			http.Error(w, "method not allow", http.StatusMethodNotAllowed)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = broker.Publish(req.Topic, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TopicResponse{Message: "Message Published"})
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {

		var req SubscribeRequest
		if r.Method != http.MethodGet {
			http.Error(w, "method not allow", http.StatusMethodNotAllowed)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		user := GenerateID(10)
		channel, err := broker.Subscribe(req.Topic, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(user, " subscribed to topic", req.Topic)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		flusher := w.(http.Flusher)
		for {
			select {
			case <-ctx.Done():
				broker.Unsubscribe(req.Topic, user)
				close(channel)
				fmt.Println(user, " unsubscribed to topic", req.Topic)
				fmt.Println()
				return
			case val := <-channel:
				json.NewEncoder(w).Encode(SubscribeResponse{
					Event:   val.Event,
					Message: val.Message,
				})

				flusher.Flush()
			}

		}

	})

	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
