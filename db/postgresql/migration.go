package postgresql

import (
	"embed"
	"errors"
	"fmt"
	"regexp"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/igefined/go-kit/config"
	"github.com/igefined/go-kit/log"
)

func Migrate(logger *log.Logger, fs *embed.FS, cfg *config.DBCfg) error {
	source, err := iofs.New(fs, "files")
	if err != nil {
		return err
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, makeMigrateUrl(cfg.URL))
	if err != nil {
		return err
	}

	err = instance.Up()

	switch {
	case err == nil:
		logger.Info("The migration schema: The schema successfully upgraded!")
	case errors.Is(err, migrate.ErrNoChange):
		logger.Info("The migration schema: The schema not changed")
	default:
		logger.Error("Could not apply the migration schema: %s", zap.Error(err))
	}

	return err
}

func makeMigrateUrl(dbUrl string) string {
	urlRe := regexp.MustCompile(`^[^\\?]+`)
	url := urlRe.FindString(dbUrl)

	sslModeRe := regexp.MustCompile("(sslmode=)[a-zA-Z0-9]+")
	sslMode := sslModeRe.FindString(dbUrl)

	return fmt.Sprintf("%s?%s", url, sslMode)
}
