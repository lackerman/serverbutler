package viewmodels

type Home struct {
	Title   string
	Heading string
	IpInfo  *IpInfo
}

type Config struct {
	Title   string
	Heading string
	OpenVPN OpenVPN
	Slack   Slack
}

type OpenVPN struct {
	Notification string
	Username     string
	Password     string
	Configs      []string
	Selected     string
	ConfigDir    string
}

type Slack struct {
	URL string
}

type IpInfo struct {
	Ip          string
	City        string
	Region      string
	Country     string
	Postal      string
	Latitude    float32
	Longitude   float32
	Timezone    string
	Asn         string
	Org         string
}

// ErrorMessage is the definition of a JSON error message
type ErrorMessage struct {
	Message string `json:"message"`
}