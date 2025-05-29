package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type AppConfig struct {
	ServiceName string                  `yaml:"service_name"`
	Language    string                  `yaml:"language"`
	Mode        string                  `yaml:"mode"`
	DB          DB                      `mapstructure:"db"`
	Services    Services                `yaml:"services"`
	Core        Core                    `yaml:"core"`
	Monitoring  Monitoring              `yaml:"monitoring"`
	Kafka       Kafka                   `mapstructure:"kafka"`
	Schedulers  map[string]SchedulerJob `mapstructure:"scheduler"`
}

type Core struct {
	Http Http `mapstructure:"http"`
}

type Http struct {
	Address string `yaml:"address"`
	Port    uint16 `yaml:"port"`
	Url     string `yaml:"url"`
}

type DB struct {
	Postgres  Postgres `mapstructure:"postgres"`
	Redis     Redis    `yaml:"redis"`
	RunSeeder bool     `mapstructure:"run_seeder"`
}

type Postgres struct {
	Host               string        `yaml:"host"`
	AppName            string        `yaml:"appName"`
	Port               int           `yaml:"port"`
	Username           string        `yaml:"username"`
	Password           string        `mask:"filled" yaml:"password"`
	Name               string        `yaml:"name"`
	Driver             string        `yaml:"driver"`
	AutoMigration      bool          `mapstructure:"auto_migration"`
	Ssl                string        `yaml:"ssl"`
	MigrationsURL      string        `mapstructure:"migrations_url"`
	TransactionTimeout time.Duration `yaml:"transaction_timeout" mapstructure:"transaction_timeout"`
	MaxIdleConnection  int           `yaml:"max_idle_connection" mapstructure:"max_idle_connection"`
	MaxOpenConnection  int           `yaml:"max_open_connection" mapstructure:"max_open_connection"`
	ConnMaxLifetime    time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	TraceStacks        bool          `yaml:"trace_stacks" mapstructure:"trace_stacks"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `mask:"filled" yaml:"password"`
	DB       int    `yaml:"db"`
}

type Services struct{}

type Monitoring struct {
	OpenTelemetry OpenTelemetry `yaml:"open_telemetry"`
}

type OpenTelemetry struct {
	ServiceName  string `yaml:"service_name"`
	CollectorURL string `yaml:"collector_url"`
}

type KafkaRouterConfig struct {
	CloseTimeout                   string `yaml:"close_timeout"`
	Retries                        int    `yaml:"retries"`
	RetriesIntervalMs              int    `yaml:"retries_interval_ms"`
	SigTermAndSigIntHandledAlready bool   `yaml:"sig_term_and_sig_int_handled_already"`
	InitialInterval                int64  `yaml:"initial_interval"`
}

type KafkaRetries struct {
	Retries int64 `yaml:"retries"`
}

type InboxRetry struct {
	MaxRetries  int `yaml:"max_retries"`
	BaseDelay   int `yaml:"base_delay"`
	MaxDelay    int `yaml:"max_delay"`
	ScaleFactor int `yaml:"scale_factor"`
}

type Kafka struct {
	Brokers       []string          `yaml:"brokers"`
	DebugLog      bool              `yaml:"debug_log"`
	InboxRetry    InboxRetry        `mapstructure:"inbox_retry"`
	Retry         KafkaRetries      `mapstructure:"retry"`
	Router        KafkaRouterConfig `mapstructure:"router"`
	GroupId       string            `yaml:"group_id" mapstructure:"group_id"`
	FromBeginning bool              `yaml:"from_beginning" mapstructure:"from_beginning"`
	Topics        map[string]string `yaml:"topics" mapstructure:"topics"`
}

type SchedulerJob struct {
	Time    string `yaml:"time" mapstructure:"time"`
	Enabled bool   `yaml:"enabled" mapstructure:"enabled"`
}

func LoadConfig(configPath string) *AppConfig {
	conf := NewConfig(configPath, &AppConfig{})
	configJson, err := json.Marshal(conf.Internal.(*AppConfig))
	if err != nil {
		panic(fmt.Sprintf("Can't make the json the config file:%v", err))
	}
	fmt.Printf("config value: %v", string(configJson))

	return conf.Internal.(*AppConfig)
}
