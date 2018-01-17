package homecoach

import (
	"fmt"
	"time"

	"github.com/gost/sensorthings-connector/module"
)

var (
	minFetchInterval = 300
)

// Setup initialised the module by setting some default values
func (m *Module) Setup() error {
	m.ModuleName = "Netatmo Homecoach"
	m.ModuleDescription = "Publish Netatmo Homecoach readings to a SensorThings server"
	m.Endpoints = m.getEndpoints()

	m.settings = Settings{}
	err := m.GetSettings(&m.settings)
	if err != nil {
		return err
	}

	if len(m.settings.ClientID) == 0 || len(m.settings.ClientSecret) == 0 || len(m.settings.Username) == 0 || len(m.settings.Password) == 0 {
		m.SendError(fmt.Errorf("missing config parameters"), true)
	}

	m.client, err = NewClient(Config{
		ClientID:     m.settings.ClientID,
		ClientSecret: m.settings.ClientSecret,
		Username:     m.settings.Username,
		Password:     m.settings.Password,
	})
	if err != nil {
		m.SendError(fmt.Errorf("unable to create Netatmo Homecoach client"), true)
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
	r, err := m.client.Read()
	if err != nil {
		m.SendError(fmt.Errorf("unable to get netatmo homecoach sensor values: %v", err), false)
	} else {
		m.handleReadings(r)
	}
}

// ToDo: Lesser for loops -> create mappings?
func (m *Module) handleReadings(response *Response) {
	for _, mod := range response.Body.Devices {
		for _, mapping := range m.settings.Mappings {
			if mapping.ModuleID == mod.ID {
				data := dashboardDataToMap(mod.DashboardData)
				for dataType, value := range data {
					for _, s := range mapping.Streams {
						if s.Type == dataType {
							obs := module.Observation{
								Result:         value,
								PhenomenonTime: time.Unix(int64(mod.DashboardData.TimeUTC), 0).Format(time.RFC3339Nano),
							}

							m.SendObservation(mapping.Server, s.StreamID, obs)
						}
					}
				}
			}
		}
	}
}

func dashboardDataToMap(data DashboardData) map[string]interface{} {
	dataMap := make(map[string]interface{}, 14)
	dataMap["AbsolutePressure"] = data.AbsolutePressure
	dataMap["TimeUTC"] = data.TimeUTC
	dataMap["HealthIndex"] = data.HealthIndex
	dataMap["Noise"] = data.Noise
	dataMap["Temperature"] = data.Temperature
	dataMap["TempTrend"] = data.TempTrend
	dataMap["Humidity"] = data.Humidity
	dataMap["Pressure"] = data.Pressure
	dataMap["PressureTrend"] = data.PressureTrend
	dataMap["CO2"] = data.CO2
	dataMap["DateMaxTemp"] = data.DateMaxTemp
	dataMap["DateMinTemp"] = data.DateMinTemp
	dataMap["MinTemp"] = data.MinTemp
	dataMap["MaxTemp"] = data.MaxTemp

	return dataMap
}
