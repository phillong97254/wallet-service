package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	GRPCServer GRPCServer `mapstructure:"grpc_server"`
	API        API        `mapstructure:"api"`
	OAuth      OAuth      `mapstructure:"oauth"`
}

type GRPCServer struct {
	Name string `mapstructure:"name"`
	Addr string `mapstructure:"addr"`
}

type API struct {
	Version string `mapstructure:"version"`
	Env     string `mapstructure:"env"`
}

// OAuth: thông tin để verify access token do auth-service phát hành.
type OAuth struct {
	JWKSURL  string `mapstructure:"jwks_url"`
	Issuer   string `mapstructure:"issuer"`
	Audience string `mapstructure:"audience"`
}

// NewConfig đọc config/config.yml, biến môi trường ghi đè theo quy tắc "." -> "_"
// (ví dụ GRPC_SERVER_ADDR=:50052, OAUTH_JWKS_URL=http://auth-service:9000/.well-known/jwks.json).
func NewConfig() (*Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
