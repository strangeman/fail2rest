package config

import (
	"github.com/spf13/viper"
)

func loadDefaults() {
	viper.SetDefault("fail2rest.production", false)
	viper.SetDefault("fail2rest.secret", "abcd1234")
	viper.SetDefault("fail2rest.port", 8080)
	viper.SetDefault("fail2rest.fail2ban", "/var/run/fail2ban/fail2ban.sock")
	viper.SetDefault("fail2rest.consul_host", "127.0.0.1:8500")
	viper.SetDefault("fail2rest.consul_token", "")
	viper.SetDefault("fail2rest.token-location", "fail2rest-token")
}
