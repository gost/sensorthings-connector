package homecoach

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// Urls to use the homecoach API
const (
	baseURL = "https://api.netatmo.net/"
	authURL = baseURL + "oauth2/token"
	dataURL = baseURL + "/api/gethomecoachsdata"
)

// Config is used to specify credential to Netatmo API
// ClientID : Client ID from netatmo app registration at http://dev.netatmo.com/dev/listapps
// ClientSecret : Client app secret
// Username : Your netatmo account username
// Password : Your netatmo account password
type Config struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

// Client use to make request to Netatmo API
type Client struct {
	oauth        *oauth2.Config
	httpClient   *http.Client
	httpResponse *http.Response
}

// NewClient create a handle authentication to Netamo API
func NewClient(config Config) (*Client, error) {
	oauth := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"read_homecoach"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  baseURL,
			TokenURL: authURL,
		},
	}

	token, err := oauth.PasswordCredentialsToken(oauth2.NoContext, config.Username, config.Password)

	return &Client{
		oauth:      oauth,
		httpClient: oauth.Client(oauth2.NoContext, token),
	}, err
}

// Response from homecoach endpoint
type Response struct {
	Body       Body    `json:"body"`
	Status     string  `json:"status"`
	TimeExec   float64 `json:"time_exec"`
	TimeServer int64   `json:"time_server"`
}

// Body part of response including devices
type Body struct {
	Devices []Device `json:"devices"`
	User    User     `json:"user"`
}

// Device describes the returned Netatmo device including sensor data
type Device struct {
	ID              string        `json:"_id"`
	CipherID        string        `json:"cipher_id"`
	LastStatusStore int64         `json:"last_status_store"`
	Place           Place         `json:"place"`
	Type            string        `json:"type"`
	DashboardData   DashboardData `json:"dashboard_data"`
	DataType        []string      `json:"data_type"`
	Co2Calibrating  bool          `json:"co2_calibrating"`
	DateSetup       int64         `json:"date_setup"`
	LastSetup       int64         `json:"last_setup"`
	ModuleName      string        `json:"module_name"`
	Firmware        int           `json:"firmware"`
	LastUpgrade     int64         `json:"last_upgrade"`
	WifiStatus      int           `json:"wifi_status"`
	Name            string        `json:"name"`
}

// Place of module
type Place struct {
	Altitude float64   `json:"altitude"`
	City     string    `json:"city"`
	Country  string    `json:"country"`
	Timezone string    `json:"timezone"`
	Location []float64 `json:"location"`
}

// DashboardData contains the actual sensor data from a device
type DashboardData struct {
	AbsolutePressure float32 `json:"AbsolutePressure"`
	TimeUTC          int64   `json:"time_utc"`
	HealthIndex      int     `json:"health_idx"`
	Noise            int     `json:"Noise"`
	Temperature      float32 `json:"Temperature"`
	TempTrend        string  `json:"temp_trend"`
	Humidity         float32 `json:"Humidity"`
	Pressure         float32 `json:"Pressure"`
	PressureTrend    string  `json:"pressure_trend"`
	CO2              float32 `json:"CO2"`
	DateMaxTemp      int64   `json:"date_max_temp"`
	DateMinTemp      int64   `json:"date_min_temp"`
	MinTemp          float32 `json:"min_temp"`
	MaxTemp          float32 `json:"max_temp"`
}

// User describes the device owner
type User struct {
	Mail           string         `json:"mail"`
	Administrative Administrative `json:"administrative"`
}

// Administrative infor about device owner
type Administrative struct {
	Country      string `json:"country"`
	FeelLikeAlgo int    `json:"feel_like_algo"`
	Lang         string `json:"lang"`
	Pressureunit int    `json:"pressureunit"`
	RegLocale    string `json:"reg_locale"`
	Unit         int    `json:"unit"`
	Windunit     int    `json:"windunit"`
}

// send http GET request
func (c *Client) doHTTPGet(url string, data url.Values) (*http.Response, error) {
	if data != nil {
		url = url + "?" + data.Encode()
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.doHTTP(req)
}

// do a generic HTTP request
func (c *Client) doHTTP(req *http.Request) (*http.Response, error) {
	var err error
	c.httpResponse, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return c.httpResponse, nil
}

// process HTTP response
// Unmarshall received data into holder struct
func processHTTPResponse(resp *http.Response, err error, holder interface{}) error {
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// check http return code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad HTTP return code %d", resp.StatusCode)
	}

	// Unmarshall response into given struct
	return json.NewDecoder(resp.Body).Decode(holder)
}

// GetStations returns the list of stations owned by the user, and their modules
func (c *Client) Read() (*Response, error) {
	resp, err := c.doHTTPGet(dataURL, nil)
	if err != nil {
		return nil, err
	}

	r := Response{}
	err = processHTTPResponse(resp, err, &r)

	if err != nil {
		return nil, err
	}

	return &r, err
}
