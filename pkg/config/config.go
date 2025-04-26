package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Config struct {
	// Stellar Network Configuration
	Network struct {
		NetworkPassphrase string `json:"network_passphrase"`
		HorizonURL        string `json:"horizon_url"`
		NetworkURL        string `json:"network_url"`
	} `json:"network"`

	// Token Configuration
	Token struct {
		MaxSupply    string `json:"max_supply"` // 100,000,000,000
		TokenName    string `json:"token_name"`
		TokenCode    string `json:"token_code"`
		IssuerSecret string `json:"issuer_secret"`
	} `json:"token"`

	// Service Categories
	Services struct {
		FreightForwarding bool `json:"freight_forwarding"`
		CustomsBrokerage  bool `json:"customs_brokerage"`
		Shipping          bool `json:"shipping"`
		AirFreight        bool `json:"air_freight"`
	} `json:"services"`

	// Development Settings
	Development struct {
		Debug       bool   `json:"debug"`
		Port        string `json:"port"`
		DatabaseURL string `json:"database_url"`
	} `json:"development"`
}

var globalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	// Override with environment variables if set
	overrideWithEnv(config)

	globalConfig = config
	return config, nil
}

func GetConfig() *Config {
	return globalConfig
}

func overrideWithEnv(config *Config) {
	if val := os.Getenv("NETWORK_PASSPHRASE"); val != "" {
		config.Network.NetworkPassphrase = val
	}
	if val := os.Getenv("HORIZON_URL"); val != "" {
		config.Network.HorizonURL = val
	}
	if val := os.Getenv("NETWORK_URL"); val != "" {
		config.Network.NetworkURL = val
	}
	if val := os.Getenv("TOKEN_MAX_SUPPLY"); val != "" {
		config.Token.MaxSupply = val
	}
	if val := os.Getenv("TOKEN_NAME"); val != "" {
		config.Token.TokenName = val
	}
	if val := os.Getenv("TOKEN_CODE"); val != "" {
		config.Token.TokenCode = val
	}
	if val := os.Getenv("ISSUER_SECRET"); val != "" {
		config.Token.IssuerSecret = val
	}
	if val := os.Getenv("SERVICE_FREIGHT_FORWARDING"); val != "" {
		config.Services.FreightForwarding, _ = strconv.ParseBool(val)
	}
	if val := os.Getenv("SERVICE_CUSTOMS_BROKERAGE"); val != "" {
		config.Services.CustomsBrokerage, _ = strconv.ParseBool(val)
	}
	if val := os.Getenv("SERVICE_SHIPPING"); val != "" {
		config.Services.Shipping, _ = strconv.ParseBool(val)
	}
	if val := os.Getenv("SERVICE_AIR_FREIGHT"); val != "" {
		config.Services.AirFreight, _ = strconv.ParseBool(val)
	}
	if val := os.Getenv("DEV_DEBUG"); val != "" {
		config.Development.Debug, _ = strconv.ParseBool(val)
	}
	if val := os.Getenv("DEV_PORT"); val != "" {
		config.Development.Port = val
	}
	if val := os.Getenv("DEV_DATABASE_URL"); val != "" {
		config.Development.DatabaseURL = val
	}
}
