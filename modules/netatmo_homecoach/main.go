package main

import (
	"github.com/gost/sensorthings-connector/module"
	homecoach "github.com/gost/sensorthings-connector/modules/netatmo_homecoach/module"
)

// Module is a mandatory var which is used by the connector
var Module module.IConnectorModule = &homecoach.Module{}

func main() {}
