package module

import (
	"net/http"
)

// IConnectorModule defines the functions that need to be implemented by a module
type IConnectorModule interface {
	Setup() error
	Start(onStartup bool) error
	Stop()

	GetID() string
	GetName() string
	GetDescription() string
	GetConnectorModuleData() *ConnectorModuleData
	GetEndpoints() []Endpoint

	SetID(string)
	SetConnectorModuleData(*ConnectorModuleData)
}

// PostStatus function definition is used as a callback when posting data to a SensorThings server
type PostStatus func(response *http.Response, err error)

// ConnectorModuleStatus contains information about the status of a module
type ConnectorModuleStatus struct {
	MaxErrors                int      `json:"-"`
	Fatal                    bool     `json:"fatal"`
	Running                  bool     `json:"running"`
	LastGet                  string   `json:"lastGet"`
	LastPost                 string   `json:"lastPost"`
	ObservationsPostedOk     int64    `json:"postSuccess"`
	ObservationsPostedFailed int64    `json:"postFailed"`
	Errors                   []string `json:"errors"`
}

// ErrorMessage send over ErrorChannel, an ErrorMessage should be send from a module
// when an error occurs so it can be logged from the connector
type ErrorMessage struct {
	ModuleID string
	Fatal    bool
	Error    error
}

// ObservationMessage can be passed from a module to the connector
// observation message channel, this will be used to post observation data to a Datastream
type ObservationMessage struct {
	ModuleID     string
	Host         string
	DatastreamID string
	Observation  Observation
	Status       PostStatus
}

// LocationMessage can be passed from a module to the connector
// location message channel, this will be used to post a new location for a Thing
type LocationMessage struct {
	ModuleID string
	Host     string
	ThingID  string
	Location Location
	Status   PostStatus
}

type dummySettings struct {
	ModuleID                   string `json:"moduleId"`
	AllowDuplicateResultValues *bool  `json:"allowDuplicateResultValues"`
}
