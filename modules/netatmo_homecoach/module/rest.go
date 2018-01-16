package homecoach

import (
	"encoding/json"
	"net/http"

	"github.com/gost/sensorthings-connector/module"
	"github.com/julienschmidt/httprouter"
)

func (m *Module) getEndpoints() []module.Endpoint {
	eps := make([]module.Endpoint, 0)
	ep := module.Endpoint{
		Name: "Settings",
		Operations: []module.EndpointOperation{
			module.EndpointOperation{
				OperationType: module.HTTPOperationGet,
				Path:          "/Settings",
				Handler:       m.getSettingsHandler,
			},
		},
	}

	eps = append(eps, ep)
	return eps
}

func (m *Module) getSettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := &Settings{}
	*c = m.settings
	c.ClientID = ""
	c.ClientSecret = ""
	c.Password = ""

	b, _ := json.MarshalIndent(c, "", "   ")
	w.Write(b)
}
