package api

type TopicRequest struct {
	Name string `json:"name"`
}
type TopicResponse struct {
	Message string `json:"message"`
}
type PublishRequest struct {
	Topic   string `json:"topic"`
	Event   string `json:"event"`
	Message string `json:"message"`
}

type PublishResponse struct {
	Message string `json:"message"`
}

type SubscribeRequest struct {
	Topic string `json:"topic"`
}

type SubscribeResponse struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}
