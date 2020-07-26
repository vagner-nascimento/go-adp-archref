package config

import (
	"encoding/json"
	"fmt"
	"github.com/vagner-nascimento/go-adp-bridge/src/infra/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ConnectionRetry struct {
	Sleep    time.Duration `json:"sleep"`
	MaxTries *int          `json:"maxTries"`
}

type AmqpDataConfig struct {
	ConnStr              string          `json:"connStr"`
	ConnRetry            ConnectionRetry `json:"connectionRetry"`
	ExitOnLostConnection bool            `json:"exitOnLostConnection"`
}

type DataConfig struct {
	Amqp AmqpDataConfig `json:"amqp"`
}

type PresentationWebConfig struct {
	Port int `json:"port"`
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
	Fails      TopicConfig `json:"fails"`
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
	Affiliations     ResClientConfig `json:"affiliations"`
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

func Load() (err error) {
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "LOCAL"

		logger.Info("GO_ENV not informed, using LOCAL", nil)
	}

	if config == nil {
		path, _ := filepath.Abs(fmt.Sprintf("config/%s.json", strings.ToLower(env)))

		var confData []byte
		confData, err = ioutil.ReadFile(path)

		if err == nil {
			if err = json.Unmarshal(confData, &config); err == nil {
				logger.Info("**CONFIGS**", string(confData))

				config.Env = env
			}
		}
	}

	return
}

func Get() Config {
	if config == nil {
		panic("config not loaded")
	}
	return *config
}
