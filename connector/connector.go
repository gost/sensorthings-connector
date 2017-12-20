package connector

import (
	"fmt"
	"os"
	"strings"

	"github.com/gost/sensorthings-connector/configuration"
	"github.com/gost/sensorthings-connector/module"
	log "github.com/sirupsen/logrus"
)

// constants
const (
	VERSION = "1.0"
	NAME    = "sensorthings connector"
)

var (
	// Modules holds all loaded modules including their status
	Modules      = make(map[string]*module.IConnectorModule, 0)
	observations = make(chan module.ObservationMessage)
	locations    = make(chan module.LocationMessage)
	errors       = make(chan module.ErrorMessage)
)

// Start the connector
func Start(config configuration.ConnectorConfig) {
	log.Infof("Starting %s", NAME)

	// start listening on channels
	go listenForObservations()
	go listenForLocations()
	go listenForErrors()

	// load and setup modules
	modulePath := config.ModulePath
	if len(modulePath) == 0 {
		modulePath = os.Args[0]
	}
	initModules(modulePath)

	// start the modules
	if config.StartModulesOnStartup {
		startModules(true)
	}

	// start the HTTP server (which will also keep the app running)
	StartHTTPServer(config.Host, config.Port, Modules)
}

// Stop the connector
func Stop() {
	stopModules()
}

func initModules(configPath string) {
	mod := loadModules(configPath, &observations, &locations, &errors)
	for _, m := range mod {
		data := (*m).GetConnectorModuleData()
		if data.Status.Fatal {
			log.Errorf("Error loading module %s: %v\n", data.ModuleFileName, data.Status.Errors[0])
		} else {
			log.Infof("Module %s loaded: %s - %s  ", data.ModuleFileName, (*m).GetName(), (*m).GetDescription())
		}

		addIDError := false
		if len((*m).GetID()) == 0 {
			(*m).SetID(module.RandomID(8))
			addIDError = true
		}

		Modules[(*m).GetID()] = m

		if addIDError {
			errStr := fmt.Errorf("No ID set for module %s, generated ID = %s", data.ModuleFileName, (*m).GetID())
			msg := module.ErrorMessage{
				ModuleID: (*m).GetID(),
				Error:    errStr,
			}

			errors <- msg
		}
	}
}

func startModules(isStartup bool) {
	for _, m := range Modules {
		module := m
		go startModule(module, isStartup)
	}
}

func stopModules() {
	for _, m := range Modules {
		module := m
		go stopModule(module)
	}
}

func startModule(module *module.IConnectorModule, isStartup bool) error {
	if (*module).GetConnectorModuleData().Status.Fatal {
		return fmt.Errorf("module not sarted because it is in 'Fatal' state")
	}

	error := (*module).Start(isStartup)
	if error == nil {
		(*module).GetConnectorModuleData().Status.Running = true
	} else {
		(*module).GetConnectorModuleData().Status.Running = false
		(*module).GetConnectorModuleData().AddError(error)
	}

	return error
}

func stopModule(module *module.IConnectorModule) {
	(*module).Stop()
	(*module).GetConnectorModuleData().Status.Running = false
}

func listenForObservations() {
	for {
		msg := <-observations
		go sendObservation(msg)
	}
}

func listenForLocations() {
	for {
		msg := <-locations
		go sendLocation(msg)
	}
}

func listenForErrors() {
	for {
		msg := <-errors
		m := Modules[msg.ModuleID]

		if m == nil {
			log.Errorf("incoming error from not registered module id %s: %v", msg.ModuleID, msg.Error)
			return
		}

		// Add error to the module and log the error
		(*m).GetConnectorModuleData().AddError(msg.Error)
		log.Errorf("module %s error: %v", (*m).GetConnectorModuleData().ModuleFilePath, msg.Error)
		if msg.Fatal {
			status := (*m).GetConnectorModuleData().Status
			status.Fatal = true
			status.Running = false
			(*m).Stop()
		}
	}
}

func sendObservation(msg module.ObservationMessage) {
	b, err := module.PostJSON(constructObservationURL(msg.Host, msg.DatastreamID), msg.Observation, 201)
	msg.Status(b, err)
}

func sendLocation(msg module.LocationMessage) {
	b, err := module.PostJSON(constructLocationURL(msg.Host, msg.ThingID), msg.Location, 201)
	msg.Status(b, err)
}

func constructObservationURL(host, streamID string) string {
	return fmt.Sprintf("%sDatastreams(%s)/Observations", getHostWithSuffix(host), streamID)
}

func constructLocationURL(host, thingID string) string {
	return fmt.Sprintf("%sThings(%s)/Locations", getHostWithSuffix(host), thingID)
}

func getHostWithSuffix(host string) string {
	if strings.HasSuffix(host, "/") {
		return host
	}

	return fmt.Sprintf("%s/", host)
}
