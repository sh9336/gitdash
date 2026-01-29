package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Dashboard DashboardConfig `mapstructure:"dashboard"`
	Commits   CommitsConfig   `mapstructure:"commits"`
	Display   DisplayConfig   `mapstructure:"display"`
}

type DashboardConfig struct {
	RefreshInterval string `mapstructure:"refresh_interval"`
	DefaultView     string `mapstructure:"default_view"`
}

type CommitsConfig struct {
	ShowCount        int  `mapstructure:"show_count"`
	ShowAuthor       bool `mapstructure:"show_author"`
	ShowRelativeTime bool `mapstructure:"show_relative_time"`
}

type DisplayConfig struct {
	Colors  bool `mapstructure:"colors"`
	Unicode bool `mapstructure:"unicode"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("dashboard.refresh_interval", "30s")
	v.SetDefault("commits.show_count", 10)
	v.SetDefault("commits.show_author", true)
	v.SetDefault("display.colors", true)

	// Config file
	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigName(".gitdash")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	v.SetEnvPrefix("GITDASH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
