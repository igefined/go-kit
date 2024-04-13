package test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	cfg "github.com/igefined/go-kit/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	cfg *cfg.DBCfg
	*postgres.PostgresContainer
}

func NewPostgresContainer(ctx context.Context, cfg *cfg.DBCfg, opt *Opt) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(opt.Image),
		postgres.WithDatabase(cfg.GetDatabaseName()),
		postgres.WithUsername(cfg.GetDatabaseUser()),
		postgres.WithPassword(cfg.GetDatabasePassword()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, err
	}

	cfg.URL = strings.Replace(cfg.URL, cfg.GetDatabasePort(), port.Port(), 1)

	fmt.Printf("PostgresContainer: mapped port %s\n", port.Port())

	return &PostgresContainer{cfg: cfg, PostgresContainer: container}, nil
}

type ImageSubstitutor struct {
	proxy string
}

func (i *ImageSubstitutor) Description() string {
	format := "docker proxy substitutor %s"
	if len(i.proxy) != 0 {
		return fmt.Sprintf(format, fmt.Sprintf("to %s", i.proxy))
	}

	return fmt.Sprintf(format, "is disabled")
}

func (i *ImageSubstitutor) Substitute(image string) (string, error) {
	if len(i.proxy) != 0 && !strings.HasPrefix(image, i.proxy) {
		p, err := url.JoinPath(i.proxy, image)
		if err != nil {
			return "", nil
		}

		image = p
	}

	return image, nil
}

func NewImageSubstitutor() *ImageSubstitutor {
	return &ImageSubstitutor{
		proxy: os.Getenv("DOCKER_PROXY"),
	}
}
