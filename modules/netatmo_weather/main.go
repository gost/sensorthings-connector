package main

import (
	"github.com/gost/sensorthings-connector/module"
	weather "github.com/gost/sensorthings-connector/modules/netatmo_weather/module"
)

// Module is a mandatory var which is used by the connector
var Module module.IConnectorModule = &weather.Module{}

func main() {}
