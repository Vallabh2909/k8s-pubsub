package api

type Topic struct{
	Subscribers []chan PublishRequest
	Channel chan PublishRequest
	Name string
}

type TopicRequest struct{
	Name string `json:"name"`
}
type TopicResponse struct{
	Message string `json:"message"`
}
type PublishRequest struct {
	Topic string `json:"topic"`
	Event string `json:"event"`
	Message string `json:"message"`
}

type PublishResponse struct {
	Message string `json:"message"`
}

type SubscribeRequest struct {
	Topic string `json:"topic"`
}

type SubscribeResponse struct {
	Event string `json:"event"`
	Message string `json:"message"`
}

func (t *Topic) Start(){
	for{
		msg:=<-t.Channel
		for _,subscriber:= range t.Subscribers{
			subscriber<-msg
		}
	}
}

