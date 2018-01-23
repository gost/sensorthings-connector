package main

import (
	"github.com/gost/sensorthings-connector/module"
	"github.com/gost/sensorthings-connector/modules/tracis/module"
)

// Module is a mandatory var which is used by the connector
var Module module.IConnectorModule = &tracis.Module{}

func main() {}
