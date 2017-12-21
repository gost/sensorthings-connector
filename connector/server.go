package connector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/gost/sensorthings-connector/module"
	log "github.com/sirupsen/logrus"
)

// State can be send to a server endpoint to start or stop a module
type State struct {
	On       bool     `json:"on"`
	ModuleID string   `json:"moduleId"`
	Errors   []string `json:"errors"`
}

// StartHTTPServer starts the HTTP server
func StartHTTPServer(host string, port int, modules map[string]*module.IConnectorModule) {
	constructModuleInfo(modules)

	log.Infof("Starting HTTP server on %s:%v", host, port)
	router := httprouter.New()
	router.GET("/Modules", moduleInfoHandler)
	router.POST("/Modules/State", stateHandler)

	// register all endpoints added by the modules
	for _, i := range moduleInfos.Modules {
		for _, e := range i.Endpoints {
			for _, o := range e.Operations {
				switch o.OperationType {
				case module.HTTPOperationGet:
					{
						router.GET(o.Path, o.Handler)
					}
				}
			}
		}
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%v", host, port), router))
}

func moduleInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, _ := json.MarshalIndent(moduleInfos, "", "   ")
	w.Write(b)
}

func stateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	state := State{}
	state.Errors = make([]string, 0)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		state.Errors = append(state.Errors, "Error reading request body")
		sendState(nil, state, w, r, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &state)
	if err != nil || len(state.ModuleID) == 0 {
		state.Errors = append(state.Errors, "POST body is not in the right format")
		sendState(nil, state, w, r, http.StatusBadRequest)
		return
	}

	var module *module.IConnectorModule
	for _, m := range Modules {
		if (*m).GetID() == state.ModuleID {
			module = m
		}
	}

	if module == nil {
		state.Errors = append(state.Errors, fmt.Sprintf("Unable to find module %s", state.ModuleID))
		sendState(nil, state, w, r, http.StatusBadRequest)
		return
	}

	if state.On {
		error := startModule(module, false)
		if error != nil {
			state.Errors = append(state.Errors, (*module).GetConnectorModuleData().Status.LastErrors...)
			sendState(module, state, w, r, http.StatusInternalServerError)
			return
		}
	} else {
		stopModule(module)
	}

	sendState(module, state, w, r, http.StatusOK)
}

func sendState(module *module.IConnectorModule, state State, w http.ResponseWriter, r *http.Request, status int) {
	js, _ := json.Marshal(state)

	if status == http.StatusOK {
		stateString := "running"
		if !state.On {
			stateString = "stopped"
		}

		log.Infof("Requested state change for module with id: %s from REST service, module is now %v", (*module).GetID(), stateString)
	} else {
		log.Errorf("Requested state change for module with id: %s from REST service, but failed: %v", state.ModuleID, state.Errors)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}
