package config

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/rancher/go-rancher/v3"
)

const (
	kuberentesHostEnv = "KUBERNETES_SERVICE_HOST"
	kuberentesPortEnv = "KUBERNETES_SERVICE_PORT"
)

type Config struct {
	KubernetesURL   string
	CattleURL       string
	CattleAccessKey string
	CattleSecretKey string
	WorkerCount     int
	HealthCheckPort int
}

func Conf(context *cli.Context) Config {
	kubernetesURL := fmt.Sprintf("https://%s:%s", os.Getenv(kuberentesHostEnv), os.Getenv(kuberentesPortEnv))
	config := Config{
		KubernetesURL:   kubernetesURL,
		CattleURL:       context.String("cattle-url"),
		CattleAccessKey: context.String("cattle-access-key"),
		CattleSecretKey: context.String("cattle-secret-key"),
		WorkerCount:     context.Int("worker-count"),
		HealthCheckPort: context.Int("health-check-port"),
	}

	return config
}

func GetRancherClient(conf Config) (*client.RancherClient, error) {
	return client.NewRancherClient(&client.ClientOpts{
		Url:       conf.CattleURL,
		AccessKey: conf.CattleAccessKey,
		SecretKey: conf.CattleSecretKey,
	})
}
