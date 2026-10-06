package broker

import (
	"github.com/vallabh2909/k8s-pubsub/internal/api"
)

type Broker struct{
	Topics map[string]*Topic
}

func NewBroker() *Broker{
	return &Broker{
		Topics: make(map[string]*Topic),
	}
}
func (b *Broker) Register(name string) *Topic{
	b.Topics[name] = NewTopic(name)
	return b.Topics[name]
}
func (b *Broker) Publish(topic string, message PublishRequest){
	b.Topics[topic].Channel <- message
}
func (b *Broker) Subscribe(topic string) *Topic{
	return b.Topics[topic]
}

func NewTopic(name string) *Topic{
	return &Topic{
		Name: name,
		Channel: make(chan PublishRequest, 1000),
		Subscribers: nil,
	}
}
func (t *Topic) Start(){
	for{
		msg:=<-t.Channel
		for _,subscriber:= range t.Subscribers{
			subscriber<-msg
		}
	}
}