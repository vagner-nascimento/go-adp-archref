package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type ConnectionRetry struct {
	Sleep    time.Duration `json:"sleep"`
	MaxTries *int          `json:"maxTries"`
}

type AmqpDataConfig struct {
	ConnStr   string          `json:"connStr"`
	ConnRetry ConnectionRetry `json:"connectionRetry"`
}

type DataConfig struct {
	Amqp AmqpDataConfig `json:"amqp"`
}

type PresentationWebConfig struct {
	Port int16 `json:"port"`
}

type PresentationConfig struct {
	Web PresentationWebConfig `json:"web"`
}

type TopicConfig struct {
	Topic    string `json:"topic"`
	Consumer string `json:"consumer"`
}

type SubsConfig struct {
	Seller   TopicConfig `json:"seller"`
	Merchant TopicConfig `json:"merchant"`
}

type PubsConfig struct {
	CrmAccount TopicConfig `json:"crm-account"`
}

type AmqIntegrationConfig struct {
	Subs SubsConfig `json:"subs"`
	Pubs PubsConfig `json:"pubs"`
}

type ResClientConfig struct {
	BaseUrl            string        `json:"baseUrl"`
	TimeOut            time.Duration `json:"timeOut"`
	RejectUnauthorized bool          `json:"rejectUnauthorized"`
}

type RestIntegrationConfig struct {
	MerchantAccounts ResClientConfig `json:"merchantAccounts"`
	Merchants        ResClientConfig `json:"merchants"`
}

type IntegrationConfig struct {
	Amqp AmqIntegrationConfig  `json:"amqp"`
	Rest RestIntegrationConfig `json:"rest"`
}

type Config struct {
	Data         DataConfig         `json:"data"`
	Presentation PresentationConfig `json:"presentation"`
	Integration  IntegrationConfig  `json:"integration"`
	Env          string
}

var config *Config

func Load(environment string) error {
	if config != nil {
		return errors.New("config is already loaded")
	}

	path, _ := filepath.Abs(fmt.Sprintf("config/%s.json", strings.ToLower(environment)))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return err
	}

	config.Env = environment

	return nil
}

func Get() Config {
	if config == nil {
		panic("config not loaded")
	}
	return *config
}
