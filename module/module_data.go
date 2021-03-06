package module

import "fmt"

// NewConnectorModuleData creates a new ConnectorModuleData object
func NewConnectorModuleData(version, fileName, filePath string, obsChannel *chan ObservationMessage, locChannel *chan LocationMessage, errorChannel *chan ErrorMessage) *ConnectorModuleData {
	cm := ConnectorModuleData{
		ModuleFileName:     fileName,
		ModuleFilePath:     filePath,
		ConnectorVersion:   version,
		ObservationChannel: obsChannel,
		LocationChannel:    locChannel,
		ErrorChannel:       errorChannel,
		Status:             &ConnectorModuleStatus{},
	}

	cm.Status.LastErrors = make([]string, 0)

	return &cm
}

// ConnectorModuleData will be send to the init function of a ConnectorModule
// Data in here can be used for initialisation/sending data
type ConnectorModuleData struct {
	ModuleFileName     string                   `json:"fileName"`
	ModuleFilePath     string                   `json:"filePath"`
	ConnectorVersion   string                   `json:"-"`
	Status             *ConnectorModuleStatus   `json:"status"`
	ObservationChannel *chan ObservationMessage `json:"-"`
	LocationChannel    *chan LocationMessage    `json:"-"`
	ErrorChannel       *chan ErrorMessage       `json:"-"`
}

// AddError adds a new error to the list of errors for the module
func (c *ConnectorModuleData) AddError(err error) {
	maxErrors := c.Status.MaxErrors
	if maxErrors == 0 {
		maxErrors = 50
	}

	c.Status.ErrorCount = c.Status.ErrorCount + 1

	// Prepend
	c.Status.LastErrors = append([]string{fmt.Sprintf("%v", err)}, c.Status.LastErrors...)

	// Remove if more than XX errors
	if len(c.Status.LastErrors) > maxErrors {
		c.Status.LastErrors = append(c.Status.LastErrors[:maxErrors], c.Status.LastErrors[maxErrors+1:]...)
	}
}
