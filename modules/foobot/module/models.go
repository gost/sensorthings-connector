package foobot

import (
	"time"

	"github.com/gost/sensorthings-connector/module"
)

// Module adds support for publishing Foobot air quality readings
// to a SensorThings server
type Module struct {
	module.ConnectorModuleBase
	settings Settings
	ticker   *time.Ticker
}

// Settings contains information on Netatmo login and sensor reading to datastream mappings
type Settings struct {
	SecretKey     string    `json:"secretKey"`
	FetchInterval int       `json:"fetchIntervalSeconds"`
	Mappings      []Mapping `json:"mappings"`
}

// Mapping contains information about the link between the netatmo stations, sensors and datastreams
type Mapping struct {
	UUID    string   `json:"uuid"`
	Name    string   `json:"name"`
	Server  string   `json:"server"`
	Streams []Stream `json:"streams"`
}

// Stream Netatmo type to SensorThings stream
type Stream struct {
	Sensor   string `json:"sensor"`
	StreamID string `json:"streamId"`
}

// FoobotJSON response from Foobot api
type FoobotJSON struct {
	UUID       string      `json:"uuid"`
	Start      int64       `json:"start"`
	End        int64       `json:"end"`
	Sensors    []string    `json:"sensors"`
	Units      []string    `json:"units"`
	Datapoints [][]float64 `json:"datapoints"`
}
