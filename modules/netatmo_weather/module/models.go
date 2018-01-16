package weather

import (
	"time"

	"github.com/gost/sensorthings-connector/module"
	netatmo "github.com/tebben/netatmo-api-go"
)

// Module adds support for publishing Netatmo weather module readings
// to a SensorThings server.
type Module struct {
	module.ConnectorModuleBase
	settings Settings
	client   *netatmo.Client
	ticker   *time.Ticker
}

// Settings contains information on Netatmo login and sensor reading to datastream mappings
type Settings struct {
	ClientID      string    `json:"clientId"`
	ClientSecret  string    `json:"clientSecret"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	FetchInterval int       `json:"fetchIntervalSeconds"`
	Mappings      []Mapping `json:"mappings"`
}

// Mapping contains information about the link between the netatmo stations, sensors and datastreams
type Mapping struct {
	ModuleID string   `json:"moduleId"`
	Name     string   `json:"name"`
	Server   string   `json:"server"`
	Streams  []Stream `json:"streams"`
}

// Stream Netatmo type to SensorThings stream
type Stream struct {
	Type     string `json:"type"`
	StreamID string `json:"streamId"`
}
