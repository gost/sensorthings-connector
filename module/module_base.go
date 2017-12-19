package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ConnectorModuleBase is the base implementation for a module this
// can be used to easily create a new module
type ConnectorModuleBase struct {
	ID                       string
	ModuleName               string
	ModuleDescription        string
	AllowDuplicateResults    bool
	mutex                    *sync.Mutex
	ModuleData               *ConnectorModuleData
	Endpoints                []Endpoint
	LatestObservationResults map[string]map[string]string
}

// GetID returns the module id
func (c *ConnectorModuleBase) GetID() string {
	return c.ID
}

// SetID returns the module id
func (c *ConnectorModuleBase) SetID(id string) {
	c.ID = id
}

// GetName returns the module name
func (c *ConnectorModuleBase) GetName() string {
	return c.ModuleName
}

// GetDescription returns the module description
func (c *ConnectorModuleBase) GetDescription() string {
	return c.ModuleDescription
}

// GetConnectorModuleData returns the ModuleData of a module
func (c *ConnectorModuleBase) GetConnectorModuleData() *ConnectorModuleData {
	return c.ModuleData
}

// SetConnectorModuleData sets the incoming ModuleData
func (c *ConnectorModuleBase) SetConnectorModuleData(data *ConnectorModuleData) {
	c.ModuleData = data
	c.AllowDuplicateResults = true
	c.mutex = &sync.Mutex{}
	c.LatestObservationResults = make(map[string]map[string]string)
}

// GetEndpoints return the configured endpoints for the module
func (c *ConnectorModuleBase) GetEndpoints() []Endpoint {
	return c.Endpoints
}

// SendError sends an error message over the ErrorChannel to the connector
func (c *ConnectorModuleBase) SendError(err error, fatal bool) {
	msg := ErrorMessage{
		ModuleID: c.GetID(),
		Fatal:    fatal,
		Error:    err,
	}

	ch := *c.ModuleData.ErrorChannel
	ch <- msg
}

// SendObservation sends an error message over the ObservationChannel to the connector
func (c *ConnectorModuleBase) SendObservation(host, datastreamID string, observation Observation) {
	c.mutex.Lock()
	c.ModuleData.Status.LastGet = time.Now().UTC().String()

	if _, ok := c.LatestObservationResults[host]; !ok {
		c.LatestObservationResults[host] = make(map[string]string)
	}

	if latestResult, ok := c.LatestObservationResults[host][datastreamID]; ok {
		if !c.AllowDuplicateResults && latestResult == fmt.Sprintf("%v", observation.Result) {
			c.mutex.Unlock()
			return
		}
	}

	// set latest result
	c.LatestObservationResults[host][datastreamID] = fmt.Sprintf("%v", observation.Result)
	c.ModuleData.Status.LastPost = time.Now().UTC().String()
	c.mutex.Unlock()

	msg := ObservationMessage{
		Host:         host,
		DatastreamID: datastreamID,
		ModuleID:     c.GetID(),
		Observation:  observation,
		Status:       c.statusCallback,
	}

	ch := *c.ModuleData.ObservationChannel
	ch <- msg
}

// SendLocation sends a location message over the LocationChannel to the connector
func (c *ConnectorModuleBase) SendLocation(host, thingID string, location Location) {
	msg := LocationMessage{
		Host:     host,
		ThingID:  thingID,
		ModuleID: c.GetID(),
		Location: location,
		Status:   c.statusCallback,
	}

	ch := *c.ModuleData.LocationChannel
	ch <- msg
}

func (c *ConnectorModuleBase) statusCallback(resp *http.Response, err error) {
	c.mutex.Lock()
	if err != nil {
		c.ModuleData.Status.ObservationsPostedOk = c.ModuleData.Status.ObservationsPostedOk + 1
	} else {
		c.ModuleData.Status.ObservationsPostedFailed = c.ModuleData.Status.ObservationsPostedFailed + 1

		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			c.SendError(fmt.Errorf("error posting to server: %v", string(body)), false)
		} else {
			c.SendError(fmt.Errorf("error posting to server: %v", err), false)
		}
	}
	c.mutex.Unlock()
}

// GetSettings reads a (JSON) config file for the module and parses it into the given settings interface
// config files should have the name of the plugin name i.e. netatmo.so should have a config file
// named netatmo.json
func (c *ConnectorModuleBase) GetSettings(settings interface{}) error {
	errorStringBase := fmt.Sprintf("error reading settings file:")
	configLocation := fmt.Sprintf("%s", strings.Replace(c.ModuleData.ModuleFilePath, c.ModuleData.ModuleFileName, strings.Replace(c.ModuleData.ModuleFileName, ".so", ".json", 1), 1))
	source, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return fmt.Errorf("%s %v", errorStringBase, err)
	}

	err = json.Unmarshal(source, settings)
	if err != nil {
		return fmt.Errorf("%s %v", errorStringBase, err)
	}

	// Try getting an id from top level settings and set module
	dummy := &dummySettings{}
	err = json.Unmarshal(source, dummy)
	if err == nil {
		if len(dummy.ModuleID) > 0 {
			c.SetID(dummy.ModuleID)
		}

		if dummy.AllowDuplicateResultValues != nil {
			c.AllowDuplicateResults = *dummy.AllowDuplicateResultValues
		}
	}

	return nil
}
