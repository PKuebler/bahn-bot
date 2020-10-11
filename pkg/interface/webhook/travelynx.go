package webhook

// TravelynxWebhook message container
type TravelynxWebhook struct {
	Reason string          `json:"reason"`
	Status TravelynxStatus `json:"status"`
}

// TravelynxStatus of the application
type TravelynxStatus struct {
	ActionTime        int64            `json:"actionTime"`
	CheckedIn         bool             `json:"checkedIn"`
	Deprecated        bool             `json:"deprecated"`
	FromStation       TravelynxStation `json:"fromStation"`
	IntermediateStops []TravelynxStop  `json:"intermediateStops"`
	ToStation         TravelynxStation `json:"toStation"`
	Train             TravelynxTrain   `json:"train"`
}

// TravelynxStation on start or end
type TravelynxStation struct {
	DS100         *string  `json:"ds100"`
	Latitude      *float32 `json:"latitude"`
	Longitude     *float32 `json:"longitude"`
	Name          *string  `json:"name"`
	RealTime      int64    `json:"realTime"`
	ScheduledTime int64    `json:"scheduledTime"`
	UIC           *int64   `json:"uic"`
}

// TravelynxStop on journey
type TravelynxStop struct {
	Name               string `json:"name"`
	RealArrival        *int64 `json:"realArrival"`
	RealDeparture      *int64 `json:"realDeparture"`
	ScheduledArrival   *int64 `json:"scheduledArrival"`
	ScheduledDeparture *int64 `json:"scheduledDeparture"`
}

// TravelynxTrain info
type TravelynxTrain struct {
	ID   string  `json:"id"`
	Line *string `json:"line"`
	No   string  `json:"no"`
	Type string  `json:"type"`
}
