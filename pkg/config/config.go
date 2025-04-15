package config

import (
    "os"
    "encoding/json"
)

type Config struct {
    // Stellar Network Configuration
    Network struct {
        NetworkPassphrase string `json:"network_passphrase"`
        HorizonURL       string `json:"horizon_url"`
        NetworkURL       string `json:"network_url"`
    } `json:"network"`

    // Token Configuration
    Token struct {
        MaxSupply        string `json:"max_supply"` // 100,000,000,000
        TokenName        string `json:"token_name"`
        TokenCode        string `json:"token_code"`
        IssuerSecret     string `json:"issuer_secret"`
    } `json:"token"`

    // Service Categories
    Services struct {
        FreightForwarding bool `json:"freight_forwarding"`
        CustomsBrokerage  bool `json:"customs_brokerage"`
        Shipping         bool `json:"shipping"`
        AirFreight       bool `json:"air_freight"`
    } `json:"services"`

    // Development Settings
    Development struct {
        Debug            bool   `json:"debug"`
        Port            string `json:"port"`
        DatabaseURL     string `json:"database_url"`
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

    globalConfig = config
    return config, nil
}

func GetConfig() *Config {
    return globalConfig
}
