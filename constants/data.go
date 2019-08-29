package constants

import (
	"os"
	"path"
)

const (
	SlackURLKey           = "slack_url"
	OpenvpnSelected       = "openvpn_config_selected"
	OpenvpnDir            = "openvpn_config_dir"
	OpenvpnCredentialFile = "credentials"
)

// SitePrefix returns the uri prefix of the site
func SitePrefix() string {
	return path.Join("/", os.Getenv("SITE_PREFIX"))
}
