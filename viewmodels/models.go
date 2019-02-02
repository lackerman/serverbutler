package viewmodels

import (
	"path"

	"github.com/lackerman/serverbutler/constants"
)

// Site represents the basic variables for a page
type Site struct {
	Page string
}

// Prefix returns the uri prefix of the site
func (Site) Prefix(uri string) string {
	return path.Join(constants.SitePrefix(), uri)
}

// Home represents the variables for the home page
type Home struct {
	Site
}

// Config represents the variables for the config page
type Config struct {
	Site
	OpenVPN
	Slack
}

// OpenVPN represents the variables for openvpn section
type OpenVPN struct {
	Notification string
	Username     string
	Password     string
	Configs      []string
	Selected     string
	ConfigDir    string
}

// Slack represents the variables for slack section
type Slack struct {
	URL string
}

// IpInfo represents the variables for ipinfo section
type IPInfo struct {
	IP        string
	City      string
	Region    string
	Country   string
	Postal    string
	Latitude  float32
	Longitude float32
	Timezone  string
	Asn       string
	Org       string
}

// ErrorMessage is the definition of a JSON error message
type ErrorMessage struct {
	Message string `json:"message"`
}
