package marudor

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// HafasService serve the marudor hafas api
type HafasService struct {
	client *APIClient
}

// HafasArrival ressource
type HafasArrival struct {
	Platform      string `json:"platform"`
	ScheduledTime int64  `json:"scheduledTime"`
	Time          int64  `json:"time"`
	Reihung       bool   `json:"reihung"`
}

// HafasDeparture ressource
type HafasDeparture struct {
	Platform      string `json:"platform"`
	ScheduledTime int64  `json:"scheduledTime"`
	Time          int64  `json:"time"`
	Reihung       bool   `json:"reihung"`
}

// HafasCoordinates ressource
type HafasCoordinates struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// HafasTrainLoad ressource
type HafasTrainLoad struct {
	FirstClass  int `json:"first"`
	SecondClass int `json:"second"`
}

// HafasStation ressource
type HafasStation struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Coordinates HafasCoordinates `json:"coordinates"`
}

// HafasStop ressource
type HafasStop struct {
	Station   *HafasStation   `json:"station"`
	Departure *HafasDeparture `json:"departure"`
	Arrival   *HafasArrival   `json:"arrival"`
	TrainLoad *HafasTrainLoad `json:"auslastung"`
}

// HafasTrainOperator resource
type HafasTrainOperator struct {
	Name string `json:"name"`
	IcoX int64  `json:"icoX"`
}

// HafasTrainInfo ressource
type HafasTrainInfo struct {
	Name     string             `json:"name"`
	Line     string             `json:"line"`
	Admin    string             `json:"admin"`
	Number   string             `json:"number"`
	Type     string             `json:"type"`
	Operator HafasTrainOperator `json:"operator"`
}

// HafasTrain ressource
type HafasTrain struct {
	JID       string          `json:"jid"`
	Train     HafasTrainInfo  `json:"train"`
	Departure *HafasDeparture `json:"departure"`
	Arrival   *HafasArrival   `json:"arrival"`
	TrainLoad *HafasTrainLoad `json:"auslastung"`
	Type      string          `json:"type"`
	// duration
	// train
	// segmentStart
	// segmentDestination
	CurrentStop      *HafasStop  `json:"currentStop"`
	Stops            []HafasStop `json:"stops"`
	FinalDestination string      `json:"finalDestination"`
	// messages
}

// HafasTrainResult from find train search
type HafasTrainResult struct {
	Train     HafasTrainInfo `json:"train"`
	JID       string         `json:"jid"`
	FirstStop *HafasStop     `json:"firstStop"`
	LastStop  *HafasStop     `json:"lastStop"`
	Stops     []HafasStop    `json:"stops"`
}

// FindTrain by name and date
func (h *HafasService) FindTrain(ctx context.Context, trainName string, date time.Time) (*[]HafasTrainResult, error) {
	path := "hafas/v1/enrichedJourneyMatch"

	requestBody := struct {
		TrainName            string `json:"trainName"`
		InitialDepartureDate string `json:"initialDepartureDate"`
	}{
		TrainName:            trainName,
		InitialDepartureDate: strconv.FormatInt(date.UnixNano()/int64(time.Millisecond), 10),
	}

	req, err := h.client.newAPIRequest("POST", path, "", requestBody)
	if err != nil {
		return nil, err
	}

	var results []HafasTrainResult
	_, err = h.client.do(req, &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// GetTrainByStation by name, station id and date
func (h *HafasService) GetTrainByStation(ctx context.Context, trainName string, stationEVA int, stationDate int64) (*HafasTrain, error) {
	path := fmt.Sprintf("hafas/v1/details/%s", trainName)

	v := url.Values{}
	v.Add("date", strconv.FormatInt(stationDate, 10))
	v.Add("station", strconv.Itoa(stationEVA))

	req, err := h.client.newAPIRequest("GET", path, v.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var train HafasTrain
	_, err = h.client.do(req, &train)
	if err != nil {
		return nil, err
	}

	return &train, nil
}
