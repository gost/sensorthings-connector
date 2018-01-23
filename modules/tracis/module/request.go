package tracis

import (
	"fmt"

	"github.com/gost/sensorthings-connector/module"
)

const (
	basePath    = "/api/equipmentdata"
	getDataPath = "/getdata"
)

// GetData retrieves the measurement from tracis for given equipment id
func GetData(host, apiKey, equipmentID string, count int) ([]Equipment, error) {
	var equipmentItems []Equipment
	url := constructURL(host, getDataPath, apiKey, equipmentID, fmt.Sprintf("&count=%v", count))
	err := module.GetJSON(url, &equipmentItems)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from tracis: %v", err)
	}
	if equipmentItems == nil {
		return nil, fmt.Errorf("Unable to retrieve data from tracis for equipment %s, please check the equipmentID", equipmentID)
	}
	return equipmentItems, nil
}

// ToDo: GetDataByDate
// ToDo: GetDataInDateRange
func constructURL(host, methodPath, apiKey, equipmentID, trail string) string {
	return fmt.Sprintf("%s%s%s?apikey=%s&equipmentid=%s%s", host, basePath, methodPath, apiKey, equipmentID, trail)
}
