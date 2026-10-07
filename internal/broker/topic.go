package broker

import "github.com/vallabh2909/k8s-pubsub/internal/api"

type Topic struct {
	Subscribers map[string]chan api.PublishRequest
	Channel     chan api.PublishRequest
	Name        string
}

func (t *Topic) Start() {
	for {
		msg := <-t.Channel
		for _, subscriber := range t.Subscribers {
			subscriber <- msg
		}
	}
}
func NewTopic(name string) *Topic {
	return &Topic{
		Name:        name,
		Channel:     make(chan api.PublishRequest, 1000),
		Subscribers: nil,
	}
}
