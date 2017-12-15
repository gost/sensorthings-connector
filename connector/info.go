package connector

import (
	"fmt"
	"time"

	"github.com/gost/sensorthings-connector/module"
)

// info returned by the Modules endpoint
var moduleInfos = info{}

// info contains information about the loaded modules which
// can be returned by the connector HTTP server
type info struct {
	ConnectorStarted string       `json:"started"`
	Modules          []moduleInfo `json:"modules"`
}

// moduleInfo contains information about all endpoints for
// loaded modules
type moduleInfo struct {
	ID          string                        `json:"id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	FileName    string                        `json:"fileName"`
	FilePath    string                        `json:"filePath"`
	Status      *module.ConnectorModuleStatus `json:"status"`
	Endpoints   []module.Endpoint             `json:"endpoints"`
}

// constructModuleInfo creates a ModuleInfo object describing a loaded module which can
// requested by going to the /Modules endpoint
func constructModuleInfo(modules map[string]*module.IConnectorModule) {
	moduleInfos = info{
		ConnectorStarted: time.Now().UTC().String(),
		Modules:          make([]moduleInfo, 0),
	}

	for _, m := range modules {
		eps := (*m).GetEndpoints()
		if len(eps) == 0 {
			continue
		}

		mi := moduleInfo{
			ID:          (*m).GetID(),
			Name:        (*m).GetName(),
			Description: (*m).GetDescription(),
			FileName:    (*m).GetConnectorModuleData().ModuleFileName,
			FilePath:    (*m).GetConnectorModuleData().ModuleFilePath,
			Endpoints:   make([]module.Endpoint, 0),
			Status:      (*m).GetConnectorModuleData().Status,
		}

		for _, ep := range eps {
			newEp := module.Endpoint{
				Name:       ep.GetName(),
				Operations: make([]module.EndpointOperation, 0),
			}

			// Set
			for _, op := range ep.GetOperations() {
				newOp := module.EndpointOperation{
					Handler:       op.Handler,
					OperationType: op.OperationType,
					Path:          fmt.Sprintf("/%s%s", (*m).GetID(), op.Path),
				}
				newEp.Operations = append(newEp.Operations, newOp)
			}

			mi.Endpoints = append(mi.Endpoints, newEp)
		}

		moduleInfos.Modules = append(moduleInfos.Modules, mi)
	}
}
