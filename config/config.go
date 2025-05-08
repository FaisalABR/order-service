package config

import (
	"order-service/common/util"
	"os"

	"github.com/sirupsen/logrus"
)

var Config AppConfig

type AppConfig struct {
	Port                  int             `json:"port"`
	AppName               string          `json:"appName"`
	AppEnv                string          `json:"appEnv"`
	SignatureKey          string          `json:"signatureKey"`
	Database              Database        `json:"database"`
	RateLimiterMaxRequest float64         `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int             `json:"rateLimiterTimeSecond"`
	InternalService       InternalService `json:"internalService"`
	Kafka                 Kafka           `json:"kafka"`
}

type Database struct {
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Name                  string `json:"name"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxLifeTimeConnection int    `json:"maxLifeTimeConnection"`
	MaxIdleTime           int    `json:"maxIdleTime"`
}

type InternalService struct {
	User    User    `json:"user"`
	Field   Field   `json:"field"`
	Payment Payment `json:"payment"`
}

type User struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Field struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Payment struct {
	Host         string `json:"host"`
	SignatureKey string `json:"signatureKey"`
}

type Kafka struct {
	Brokers               []string `json:"brokers"`
	TimeoutInMS           int      `json:"timeoutInMs"`
	MaxRetry              int      `json:"maxRetry"`
	MaxWaitTimeInMs       int      `json:"maxWaitTimeInMs"`
	MaxProcessingTimeInMs int      `json:"maxProcessingTimeInMs"`
	BackoffTimeInMs       int      `json:"backoffTimeInMs"`
	Topics                []string `json:"topics"`
	GroupID               string   `json:"groupID"`
}

func Init() {
	err := util.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		logrus.Infof("failed to bind config: %v", err)
		err = util.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_PATH"))
		if err != nil {
			panic(err)
		}
	}
}
