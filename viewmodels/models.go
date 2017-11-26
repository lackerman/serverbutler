package viewmodels

type Home struct {
	Title   string
	Heading string
}

type Config struct {
	Title   string
	Heading string
	OpenVPN OpenVPN
	Slack   Slack
}

type OpenVPN struct {
	Notification string
	Configs      []string
	Selected     string
}

type Slack struct {
	URL string
}
