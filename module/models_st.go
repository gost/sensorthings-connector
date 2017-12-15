package module

// Note: CORE SensorThing structs cannot be used in plugin cause of @ in json definitions

// Observation in SensorThings represents a single Sensor reading of an ObservedProperty. A physical device, a Sensor, sends
// Observations to a specified Datastream. An Observation requires a FeatureOfInterest entity, if none is provided in the request,
// the Location of the Thing associated with the Datastream, will be assigned to the new Observation as the FeaturOfInterest.
type Observation struct {
	PhenomenonTime    string                 `json:"phenomenonTime,omitempty"`
	Result            interface{}            `json:"result,omitempty"`
	ResultTime        *string                `json:"resultTime,omitempty"`
	ResultQuality     string                 `json:"resultQuality,omitempty"`
	ValidTime         string                 `json:"validTime,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	FeatureOfInterest *FeatureOfInterest     `json:"featureOfInterest,omitempty"`
}

// FeatureOfInterest in SensorThings represents the phenomena an Observation is detecting. In some cases a FeatureOfInterest
// can be the Location of the Sensor and therefore of the Observation. A FeatureOfInterest is linked to a single Observation
type FeatureOfInterest struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	EncodingType string                 `json:"encodingType,omitempty"`
	Feature      map[string]interface{} `json:"feature,omitempty"`
}

// Location entity locates the Thing or the Things it associated with. A Thing’s Location entity is
// defined as the last known location of the Thing.
// A Thing’s Location may be identical to the Thing’s Observations’ FeatureOfInterest. In the context of the IoT,
// the principle location of interest is usually associated with the location of the Thing, especially for in-situ
// sensing applications. For example, the location of interest of a wifi-connected thermostat should be the building
// or the room in which the smart thermostat is located. And the FeatureOfInterest of the Observations made by the
// thermostat (e.g., room temperature readings) should also be the building or the room. In this case, the content
// of the smart thermostat’s location should be the same as the content of the temperature readings’ feature of interest.
type Location struct {
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	EncodingType string                 `json:"encodingType,omitempty"`
	Location     map[string]interface{} `json:"location,omitempty"`
}
