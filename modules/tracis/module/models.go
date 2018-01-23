package tracis

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
	APIKey        string    `json:"apiKey"`
	TracisHost    string    `json:"tracisHost"`
	FetchInterval int       `json:"fetchIntervalSeconds"`
	Mappings      []Mapping `json:"mappings"`
}

// Mapping contains information about the link between the netatmo stations, sensors and datastreams
type Mapping struct {
	Name        string   `json:"name"`
	Server      string   `json:"server"`
	EquipmentID string   `json:"equipmentId"`
	Streams     []Stream `json:"streams"`
}

// Stream Netatmo type to SensorThings stream
type Stream struct {
	ChannelNumber string `json:"channelNumber"`
	StreamID      string `json:"streamId"`
}

type Equipment struct {
	EquipmentID string       `json:"equipmentID"`
	Sensors     []SensorData `json:"sensorData"`
}

type SensorData struct {
	ChannelNumber int     `json:"channelNumber"`
	PortNumber    string  `json:"portNumber"`
	SensorNumber  int     `json:"sensorNumber"`
	SensorType    string  `json:"sensorType"`
	DateTime      string  `json:"dateTime"`
	Units         string  `json:"units"`
	Value         float64 `json:"value"`
}
