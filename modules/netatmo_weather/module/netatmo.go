package weather

import (
	"fmt"
	"time"

	"github.com/gost/sensorthings-connector/module"
	netatmo "github.com/tebben/netatmo-api-go"
)

var (
	minFetchInterval = 300
)

// Setup initialised the module by setting some default values
func (m *Module) Setup() error {
	m.ModuleName = "Netatmo Weather"
	m.ModuleDescription = "Publish Netatmo Weather readings to a SensorThings server"
	m.Endpoints = m.getEndpoints()

	m.settings = Settings{}
	err := m.GetSettings(&m.settings)
	if err != nil {
		return err
	}

	if len(m.settings.ClientID) == 0 || len(m.settings.ClientSecret) == 0 || len(m.settings.Username) == 0 || len(m.settings.Password) == 0 {
		m.SendError(fmt.Errorf("missing config parameters"), true)
	}

	m.client, err = netatmo.NewClient(netatmo.Config{
		ClientID:     m.settings.ClientID,
		ClientSecret: m.settings.ClientSecret,
		Username:     m.settings.Username,
		Password:     m.settings.Password,
	})
	if err != nil {
		m.SendError(fmt.Errorf("unable to create Netatmo Weather client"), true)
		return err
	}

	return nil
}

// Start receiving Netatmo readings and publish it to a SensorThings server
func (m *Module) Start(initStartup bool) error {
	interval := m.settings.FetchInterval
	if interval == 0 || interval < minFetchInterval {
		interval = minFetchInterval
	}

	// Get some readings at start
	go m.getReadings()

	m.ticker = time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		for range m.ticker.C {
			m.getReadings()
		}
	}()

	return nil
}

// Stop receiving Netatmo readings
func (m *Module) Stop() {
	if m.ticker != nil {
		m.ticker.Stop()
	}
}

func (m *Module) getReadings() {
	dc, err := m.client.Read()
	if err != nil {
		m.SendError(fmt.Errorf("unable to get netatmo sensor values: %v", err), false)
	} else {
		for _, station := range dc.Stations() {
			go m.handleReadings(station.Modules())
		}
	}
}

// ToDo: Lesser for loops -> create mappings?
func (m *Module) handleReadings(modules []*netatmo.Device) {
	for _, mod := range modules {
		for _, mapping := range m.settings.Mappings {
			if mapping.ModuleID == mod.ID {
				ts, data := mod.Data()
				for dataType, value := range data {
					for _, s := range mapping.Streams {
						if s.Type == dataType {
							obs := module.Observation{
								Result:         value,
								PhenomenonTime: time.Unix(int64(ts), 0).Format(time.RFC3339Nano),
							}

							m.SendObservation(mapping.Server, s.StreamID, obs)
						}
					}
				}
			}
		}
	}
}
