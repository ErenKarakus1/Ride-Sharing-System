package events

type Event struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	Data any    `json:"data"`
}

type RideEvent struct {
	Type string        `json:"type"`
	Data RideEventData `json:"data"`
}

type RideEventData struct {
	ID string `json:"id"`
}
