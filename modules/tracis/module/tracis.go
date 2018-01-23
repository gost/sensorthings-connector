package tracis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gost/sensorthings-connector/module"
)

const (
	// ISO8601 defines the ISO8601 date time format
	ISO8601 = "2006-01-02T15:04:05.999Z07:00"
	// TRACISTIME defines the date time returned for a tracis measurement
	TRACISTIME = "2006-01-02T15:04:05.000"
)

var (
	minFetchInterval = 60
	location         *time.Location
	equipmentIds     []string
)

// Setup initialised the module by setting some default values
func (m *Module) Setup() error {
	location, _ = time.LoadLocation("Europe/Amsterdam")
	m.ModuleName = "Tracis"
	m.ModuleDescription = "Publish Tracis readings to a SensorThings server"
	m.settings = Settings{}

	err := m.GetSettings(&m.settings)
	if err != nil {
		return err
	}

	if len(m.settings.APIKey) == 0 || len(m.settings.TracisHost) == 0 {
		m.SendError(fmt.Errorf("missing config parameters"), true)
	}

	equipmentIds = make([]string, 0)
	for _, mapping := range m.settings.Mappings {
		if !stringInSlice(mapping.EquipmentID, equipmentIds) {
			equipmentIds = append(equipmentIds, mapping.EquipmentID)
		}
	}

	return nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// Start receiving Netatmo readings and publish it to a SensorThings server
func (m *Module) Start(initStartup bool) error {
	interval := m.settings.FetchInterval
	if interval == 0 || interval < minFetchInterval {
		interval = minFetchInterval
	}

	// Get some readings at start
	for _, e := range equipmentIds {
		m.requestAPI(e)
	}

	m.ticker = time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		for range m.ticker.C {
			for _, e := range equipmentIds {
				m.requestAPI(e)
			}
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

func (m *Module) requestAPI(equipmentID string) {
	// GetData from tracis
	equipmentItems, err := GetData(m.settings.TracisHost, m.settings.APIKey, equipmentID, 1)
	if err != nil {
		m.SendError(err, false)
		return
	}

	// There should only be one since we requested &count=1
	for _, item := range equipmentItems {
		if item.Sensors == nil {
			continue
		}

		for _, mapping := range m.settings.Mappings {
			if mapping.EquipmentID == equipmentID {
				for _, sensor := range item.Sensors {
					sID := strconv.Itoa(sensor.ChannelNumber)
					t, _ := time.ParseInLocation(TRACISTIME, fmt.Sprintf("%s.000", sensor.DateTime), location)
					t2 := t.Format(ISO8601)
					obs := module.Observation{
						PhenomenonTime: t2,
						Result:         sensor.Value,
					}

					for _, stream := range mapping.Streams {
						if sID == stream.ChannelNumber {
							m.SendObservation(mapping.Server, stream.StreamID, obs)
						}
					}
				}
			}
		}
	}
}
