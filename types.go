package main

// Reading stores data for each air quality reading
type Reading struct {
	Coordinates      []float64 `json:"location>coordinates"`
	ObservationID    string    `json:"observationid"`
	ObservationTime  string    `json:"observationtime"`
	ThingID          string    `json:"thingid"`
	LocationName     string    `json:"locationname"`
	MeasurementType  string    `json:"measurementtype"`
	MeasurementValue string    `json:"measurementvalue"`
	MeasurementUnit  string    `json:"measurementunit"`
}
