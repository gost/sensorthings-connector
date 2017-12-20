package connector

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/gost/sensorthings-connector/module"
)

var (
	modulePaths map[string]string
)

// loadModules searches for *.so files and tries to load it as a ConnectorModule
func loadModules(modulePath string, obsChannel *chan module.ObservationMessage, locChannel *chan module.LocationMessage, errorChannel *chan module.ErrorMessage) []*module.IConnectorModule {
	flag.Parse()
	modulePaths = map[string]string{}
	modules := make([]*module.IConnectorModule, 0)
	dir, _ := filepath.Abs(filepath.Dir(modulePath))
	filepath.Walk(dir, visit)

	for k, v := range modulePaths {
		// try loading module
		loaded, err := tryLoadModule(k, v)
		d := module.NewConnectorModuleData(VERSION, k, v, obsChannel, locChannel, errorChannel)

		if err != nil {
			// Unable to load, create dummy for logging purpose
			dummy := createDummy(k, d, err)
			modules = append(modules, toPointerInterface(dummy))
		} else {
			(*loaded).SetConnectorModuleData(d)
			err = (*loaded).Setup()
			if err != nil {
				dummy := createDummy(k, d, err)
				modules = append(modules, toPointerInterface(dummy))
			} else {
				modules = append(modules, loaded)
			}
		}
	}

	return modules
}

func createDummy(moduleFileName string, d *module.ConnectorModuleData, err error) *module.ConnectorModuleBase {
	d.AddError(fmt.Errorf("error loading module %s: %v", moduleFileName, err))
	d.Status.Running = false
	d.Status.Fatal = true
	dummy := &module.ConnectorModuleBase{}
	dummy.SetConnectorModuleData(d)

	return dummy
}

func toPointerInterface(i interface{}) *module.IConnectorModule {
	cm, _ := i.(module.IConnectorModule)
	return &cm
}

func visit(path string, f os.FileInfo, err error) error {
	if !strings.HasSuffix(f.Name(), ".so") {
		return nil
	}

	modulePaths[f.Name()] = path
	return nil
}

func tryLoadModule(name, path string) (*module.IConnectorModule, error) {
	lib, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file - %v", err)
	}

	m, err := lib.Lookup("Module")
	if err != nil {
		return nil, fmt.Errorf("not exported properly - %v", err)
	}

	m2, ok := m.(*module.IConnectorModule)
	if !ok {
		return nil, fmt.Errorf("module does not implement ConnectorModule properly")
	}

	return m2, nil
}
