package foobot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gost/sensorthings-connector/module"
)

var (
	minFetchInterval = 500 // 200 req day = 432 seconds
)

// Setup initialised the module by setting some default values
func (m *Module) Setup() error {
	m.ModuleName = "Foobot"
	m.ModuleDescription = "Publish Foobot sensor readings to a SensorThings server"
	m.Endpoints = m.getEndpoints()

	m.settings = Settings{}
	err := m.GetSettings(&m.settings)
	if err != nil {
		return err
	}

	if len(m.settings.SecretKey) == 0 {
		m.SendError(fmt.Errorf("missing config parameters"), true)
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
	for _, ma := range m.settings.Mappings {
		url := fmt.Sprintf("https://api.foobot.io/v2/device/%s/datapoint/0/last/0/", ma.UUID)

		foobotClient := http.Client{
			Timeout: time.Second * 10,
		}

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("X-API-KEY-TOKEN", m.settings.SecretKey)

		res, err := foobotClient.Do(req)
		if err != nil {
			m.SendError(err, false)
		}

		if res.StatusCode == 401 {
			// by setting fatal to true, module will stop running
			m.SendError(fmt.Errorf("incorrect api key"), true)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			m.SendError(err, false)
		}

		fj := FoobotJSON{}
		err = json.Unmarshal(body, &fj)
		if err != nil {
			m.SendError(err, false)
		}

		m.handleReadings(ma, fj)
	}
}

func (m *Module) handleReadings(ma Mapping, response FoobotJSON) {
	kvp := make(map[string]float64, 0)
	for i, s := range response.Sensors {
		kvp[s] = response.Datapoints[0][i]
	}

	for k, v := range kvp {
		for _, s := range ma.Streams {
			if s.Sensor == k {
				obs := module.Observation{
					Result:         v,
					PhenomenonTime: time.Unix(int64(response.End), 0).Format(time.RFC3339Nano),
				}

				m.SendObservation(ma.Server, s.StreamID, obs)
			}
		}
	}
}
