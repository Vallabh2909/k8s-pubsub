package broker

import (
	"errors"

	"github.com/vallabh2909/k8s-pubsub/internal/api"
)

type Broker struct {
	Topics map[string]*Topic
}

func NewBroker() *Broker {
	return &Broker{
		Topics: make(map[string]*Topic),
	}
}
func (b *Broker) Register(name string) error {
	if _, ok := b.Topics[name]; ok {
		return errors.New("topic already exists")
	}
	b.Topics[name] = &Topic{
		Name:        name,
		Channel:     make(chan api.PublishRequest, 1000),
		Subscribers: make(map[string]chan api.PublishRequest),
	}
	go b.Topics[name].Start()
	return nil
}
func (b *Broker) Publish(topic string, message api.PublishRequest) error {
	if t, ok := b.Topics[topic]; ok {
		t.Channel <- message
		return nil
	}
	return errors.New("topic not found")
}
func (b *Broker) Subscribe(name string, id string) (chan api.PublishRequest, error) {
	if topic, ok := b.Topics[name]; ok {
		channel := make(chan api.PublishRequest, 1000)
		topic.Subscribers[id] = channel
		return channel, nil

	} else {
		return nil, errors.New("topic not found")
	}
}

func (b *Broker) Unsubscribe(topic, id string) {
	if t, ok := b.Topics[topic]; ok {
		delete(t.Subscribers, id)
	}
}
