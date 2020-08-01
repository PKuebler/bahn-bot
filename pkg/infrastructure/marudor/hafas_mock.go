package marudor

import (
	"context"
	"errors"
)

// HafasMock to testing
type HafasMock struct {
	TrainExists bool
}

// GetTrainByStation by name, station id and date
func (h *HafasMock) GetTrainByStation(ctx context.Context, trainName string, stationEVA int, stationDate int64) (*HafasTrain, error) {
	if !h.TrainExists {
		return &HafasTrain{}, errors.New("not found")
	}

	return &HafasTrain{
		Arrival: &HafasArrival{
			Platform:      "9",
			ScheduledTime: 1595797980000,
			Time:          1595797980000,
			Reihung:       true,
		},
		CurrentStop: &HafasStop{
			Station: &HafasStation{
				ID:    "8503000",
				Title: "Zürich HB",
				Coordinates: HafasCoordinates{
					Lng: 8.540211,
					Lat: 47.378177,
				},
			},
			Departure: &HafasDeparture{
				Platform:      "15",
				ScheduledTime: 1595782800000,
				Time:          1595782800000,
				Reihung:       true,
			},
		},
		Departure: &HafasDeparture{
			Platform:      "15",
			ScheduledTime: 1595782800000,
			Time:          1595782800000,
			Reihung:       true,
		},
		Train: HafasTrainInfo{
			Name:   "ICE 4",
			Line:   "4",
			Admin:  "85____",
			Number: "4",
			Type:   "ICE",
			Operator: HafasTrainOperator{
				Name: "SBB",
				IcoX: 1,
			},
		},
		Stops: []HafasStop{
			{
				Station: &HafasStation{
					ID:    "8503000",
					Title: "Zürich HB",
					Coordinates: HafasCoordinates{
						Lng: 8.540211,
						Lat: 47.378177,
					},
				},
				Departure: &HafasDeparture{
					Platform:      "15",
					ScheduledTime: 1595782800000,
					Time:          1595782800000,
					Reihung:       true,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8500010",
					Title: "Basel SBB",
					Coordinates: HafasCoordinates{
						Lng: 7.589548,
						Lat: 47.547408,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "11",
					ScheduledTime: 1595785980000,
					Time:          1595785980000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "",
					ScheduledTime: 1595787180000,
					Time:          1595787180000,
					Reihung:       true,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000026",
					Title: "Basel Bad Bf",
					Coordinates: HafasCoordinates{
						Lng: 7.607886,
						Lat: 47.567679,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "",
					ScheduledTime: 1595787540000,
					Time:          1595787540000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "",
					ScheduledTime: 1595787780000,
					Time:          1595787780000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000107",
					Title: "Freiburg(Breisgau) Hbf",
					Coordinates: HafasCoordinates{
						Lng: 7.841344,
						Lat: 47.997974,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "1",
					ScheduledTime: 1595789580000,
					Time:          1595789580000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "1",
					ScheduledTime: 1595789700000,
					Time:          1595789700000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000290",
					Title: "Offenburg",
					Coordinates: HafasCoordinates{
						Lng: 7.94677,
						Lat: 48.476506,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "3",
					ScheduledTime: 1595791500000,
					Time:          1595791500000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "3",
					ScheduledTime: 1595791620000,
					Time:          1595791620000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000774",
					Title: "Baden-Baden",
					Coordinates: HafasCoordinates{
						Lng: 8.190773,
						Lat: 48.790392,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "4",
					ScheduledTime: 1595792460000,
					Time:          1595792460000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "4",
					ScheduledTime: 1595792520000,
					Time:          1595792520000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000191",
					Title: "Karlsruhe Hbf",
					Coordinates: HafasCoordinates{
						Lng: 8.401939,
						Lat: 48.99353,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "3",
					ScheduledTime: 1595793540000,
					Time:          1595793540000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "3",
					ScheduledTime: 1595793660000,
					Time:          1595793660000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000244",
					Title: "Mannheim Hbf",
					Coordinates: HafasCoordinates{
						Lng: 8.469268,
						Lat: 49.479181,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "2",
					ScheduledTime: 1595795400000,
					Time:          1595795400000,
					Reihung:       true,
				},
				Departure: &HafasDeparture{
					Platform:      "2",
					ScheduledTime: 1595795580000,
					Time:          1595795580000,
					Reihung:       true,
				},
				TrainLoad: &HafasTrainLoad{
					FirstClass:  1,
					SecondClass: 1,
				},
			},
			{
				Station: &HafasStation{
					ID:    "8000105",
					Title: "Frankfurt(Main)Hbf",
					Coordinates: HafasCoordinates{
						Lng: 8.663003,
						Lat: 50.106817,
					},
				},
				Arrival: &HafasArrival{
					Platform:      "9",
					ScheduledTime: 1595797980000,
					Time:          1595797980000,
					Reihung:       true,
				},
			},
		},
		FinalDestination: "Frankfurt(Main)Hbf",
		JID:              "1|351228|0|80|26072020",
		TrainLoad: &HafasTrainLoad{
			FirstClass:  1,
			SecondClass: 2,
		},
		Type: "JNY",
	}, nil
}
