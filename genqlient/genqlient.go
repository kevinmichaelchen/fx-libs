package genqlient

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("genqlient",
	fx.Provide(
		NewConfig,
		NewClient,
	),
)

type Config struct {
	*NestedConfig `env:",prefix=HASURA_GRAPHQL_"`
}

type NestedConfig struct {
	SecretKey string `env:"ADMIN_SECRET,default=myadminsecretkey"`
	Endpoint  string `env:"URL,default=http://localhost:8080/v1/graphql"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewClient(cfg *Config, httpClient *http.Client) graphql.Client {
	return graphql.NewClient(cfg.Endpoint, authedClient{
		client:    httpClient,
		secretKey: cfg.SecretKey,
	})
}

type authedClient struct {
	client    graphql.Doer
	secretKey string
}

func (c authedClient) Do(req *http.Request) (*http.Response, error) {
	req.Header["x-hasura-admin-secret"] = []string{c.secretKey}
	return c.client.Do(req)
}
