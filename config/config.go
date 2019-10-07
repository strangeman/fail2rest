package config

import (
	"encoding/json"

	"github.com/Strum355/log"
	"github.com/spf13/viper"
)

func Load() {
	loadDefaults()
	viper.AutomaticEnv()
}

func PrintSettings() {
	settings := viper.AllSettings()
	settings["fail2rest"].(map[string]interface{})["secret"] = "[secret]"

	out, _ := json.MarshalIndent(settings, "", "\t")
	log.Debug("config:\n" + string(out))
}
