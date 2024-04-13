package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMigrateUrl(t *testing.T) {
	tCases := []struct {
		name   string
		val    string
		result string
	}{
		{
			name:   "successfully, with sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable",
		},
		{
			name:   "successfully, with sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?pool_max_conns=16&sslmode=disable&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?sslmode=disable",
		},
		{
			name:   "successfully, without sslmode",
			val:    "postgres://postgres:12345@localhost:5432/nh_templates?pool_max_conns=16&&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			result: "postgres://postgres:12345@localhost:5432/nh_templates?",
		},
	}

	for _, c := range tCases {
		t.Run(c.name, func(t *testing.T) {
			migrateUrl := makeMigrateUrl(c.val)
			assert.Equal(t, c.result, migrateUrl)
		})
	}
}
