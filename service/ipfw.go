package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type IPFilter struct {
	config *IPFilterConfig
}

type IPFilterConfig struct {
	// DBIPToken - db-ip.com private api token
	DBIPToken string
}

type DBIPResult struct {
	IPAddress     string   `json:"ipAddress"`
	ContinentCode string   `json:"continentCode"`
	ContinentName string   `json:"continentName"`
	CountryCode   string   `json:"countryCode"`
	CountryName   string   `json:"countryName"`
	IsEuMember    bool     `json:"isEuMember"`
	CurrencyCode  string   `json:"currencyCode"`
	CurrencyName  string   `json:"currencyName"`
	PhonePrefix   string   `json:"phonePrefix"`
	Languages     []string `json:"languages"`
	StateProvCode string   `json:"stateProvCode"`
	StateProv     string   `json:"stateProv"`
	City          string   `json:"city"`
	GeonameId     int      `json:"geonameId"`
	ZipCode       string   `json:"zipCode"`
	Latitude      float32  `json:"latitude"`
	Longitude     float32  `json:"longitude"`
	GmtOffset     int      `json:"gmtOffset"`
	TimeZone      string   `json:"timeZone"`
	WeatherCode   string   `json:"weatherCode"`
	AsNumber      int      `json:"asNumber"`
	AsName        string   `json:"asName"`
	Isp           string   `json:"isp"`
	UsageType     string   `json:"usageType"`
	IsCrawler     bool     `json:"isCrawler"`
	IsProxy       bool     `json:"isProxy"`
	ThreatLevel   string   `json:"threatLevel"`
}

type IPLookupResult struct {
	IPAddress    string
	ISP          string
	IsSuspicious bool
	Reason       string
}

func NewIPFilter(config *IPFilterConfig) (*IPFilter, error) {
	logrus.Printf("[IPFilter] Initialize IP Filter...")

	return &IPFilter{
		config: config,
	}, nil
}

// Check - check user ip
func (ipfw *IPFilter) Check(ip string) (*IPLookupResult, error) {
	// Check DBIP
	dbipRes, err := ipfw.LookupDBIP(ip)
	if err != nil {
		return nil, err
	}

	if dbipRes.ThreatLevel == "medium" || dbipRes.ThreatLevel == "high" {
		return &IPLookupResult{
			IPAddress:    dbipRes.IPAddress,
			ISP:          dbipRes.Isp,
			IsSuspicious: true,
			Reason:       "THREAT_LEVEL_ABOVE_MEDIUM",
		}, nil
	} else if dbipRes.UsageType == "hosting" {
		return &IPLookupResult{
			IPAddress:    dbipRes.IPAddress,
			ISP:          dbipRes.Isp,
			IsSuspicious: true,
			Reason:       "HOSTING",
		}, nil
	}

	return &IPLookupResult{
		IPAddress:    ip,
		IsSuspicious: false,
	}, nil
}

// LookupDBIP - check user ip
func (ipfw *IPFilter) LookupDBIP(ip string) (*DBIPResult, error) {
	if len(ipfw.config.DBIPToken) == 0 {
		logrus.Printf("ERR: DB_IP_TOKEN is empty?")
	}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.db-ip.com/v2/%s/%s", ipfw.config.DBIPToken, ip),
		nil,
	)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse result
	b, _ := ioutil.ReadAll(res.Body)
	var parsed *DBIPResult

	if err := json.Unmarshal(b, &parsed); err != nil {
		return nil, err
	}

	//:wlogrus.Debugf("[IPFilter] Lookup result")

	return parsed, nil
}
