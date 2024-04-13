package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewEnvVar(t *testing.T) {
	flag := "test_flag"
	env := "test_env"
	description := "test description"

	t.Run("string default", func(t *testing.T) {
		defaultValue := "test_default_value"
		envVar := NewEnvVar(flag, env, defaultValue, description)
		assert.NotNil(t, envVar)
		switch envVar.DefaultValue.(type) {
		case string:
			return
		default:
			t.Fatalf("default value should be string")
		}
	})

	t.Run("int default", func(t *testing.T) {
		defaultValue := 1
		envVar := NewEnvVar(flag, env, defaultValue, description)
		assert.NotNil(t, envVar)
		switch envVar.DefaultValue.(type) {
		case int:
			return
		default:
			t.Fatalf("default value should be string")
		}
	})
}

func TestAddEnvs(t *testing.T) {
	customEnvs := []*EnvVar{
		{
			Flag:         "custom_flag1",
			Env:          "CUSTOM_ENV_1",
			DefaultValue: "default_value_1",
			Description:  "Custom environment variable 1",
		},
		{
			Flag:         "custom_flag2",
			Env:          "CUSTOM_ENV_2",
			DefaultValue: "default_value_2",
			Description:  "Custom environment variable 2",
		},
	}

	AddEnvs(customEnvs)

	for _, customEnv := range customEnvs {
		found := false
		for _, env := range envs {
			if customEnv.Flag == env.Flag {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected custom environment variable %s not found in envs", customEnv.Flag)
		}
	}

	for i, env := range envs {
		for j := i + 1; j < len(envs); j++ {
			if env.Flag == envs[j].Flag {
				t.Errorf("Duplicate environment variable %s found in envs", env.Flag)
			}
		}
	}
}

func TestBindConfig(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
		envs = nil
	})

	testArgs := []string{
		"--host=test_host",
		"--port=8080",
		"--monitor_host=test_monitor_host",
		"--custom_flag=test_value",
	}

	os.Args = append([]string{"cmd"}, testArgs...)

	pflag.String("host", "", "nftique host")
	pflag.String("port", "", "nftique port")
	pflag.String("monitor_host", "", "nftique monitor port")
	pflag.String("custom_flag", "", "Custom flag")

	BindConfig()

	assert.Equal(t, "test_host", viper.GetString("HOST"))
	assert.Equal(t, "8080", viper.GetString("PORT"))
	assert.Equal(t, "test_monitor_host", viper.GetString("MONITOR_HOST"))
	assert.Equal(t, "test_value", viper.GetString("CUSTOM_FLAG"))
}

func TestGetConfig(t *testing.T) {
	type C struct {
		MainCfg `mapstructure:",squash"`
		Test    string
	}

	testVars := []*EnvVar{
		NewEnvVar("test", "TEST", "test", "custom test env"),
	}

	var expectedPort = "9090"
	_ = os.Setenv("MONITOR_PORT", expectedPort)

	var c C
	err := GetConfig(&c, testVars)
	assert.NoError(t, err)

	assert.Equal(t, c.MonitorPort, expectedPort)
	assert.Equal(t, c.Test, "test")
}

func TestGetDatabaseName(t *testing.T) {
	const (
		url    = "postgres://postgres:postgres@localhost:5466/test_clients?sslmode=disable"
		expect = "test_clients"
	)

	cfg := DBCfg{
		URL: url,
	}

	assert.Equal(t, cfg.GetDatabaseName(), expect)
}

func TestGetDatabaseUser(t *testing.T) {
	tCases := []struct {
		url      string
		expected string
	}{
		{
			url:      "postgres://postgres:postgres@localhost:5466/test?sslmode=disable",
			expected: "postgres",
		},
		{
			url:      "postgres://root:postgres@localhost:5466/test",
			expected: "root",
		},
		{
			url:      "postgres://docker:12345@localhost:5432/common?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			expected: "docker",
		},
	}

	for _, c := range tCases {
		cfg := DBCfg{
			URL: c.url,
		}
		assert.Equal(t, c.expected, cfg.GetDatabaseUser())
	}
}

func TestGetDatabasePassword(t *testing.T) {
	tCases := []struct {
		url      string
		expected string
	}{
		{
			url:      "postgres://postgres:postgres@localhost:5466/test?sslmode=disable",
			expected: "postgres",
		},
		{
			url:      "postgres://root:12345@localhost:5466/test",
			expected: "12345",
		},
		{
			url:      "postgres://docker:pass@localhost:5432/common?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			expected: "pass",
		},
	}

	for _, c := range tCases {
		cfg := DBCfg{
			URL: c.url,
		}
		assert.Equal(t, c.expected, cfg.GetDatabasePassword())
	}
}

func TestGetDatabasePort(t *testing.T) {
	tCases := []struct {
		url      string
		expected string
	}{
		{
			url:      "postgres://postgres:postgres@localhost:5466/test?sslmode=disable",
			expected: "5466",
		},
		{
			url:      "postgres://root:12345@localhost:5432/test",
			expected: "5432",
		},
		{
			url:      "postgres://docker:pass@localhost:9090/common?sslmode=disable&pool_max_conns=16&pool_max_conn_idle_time=30m&pool_max_conn_lifetime=1h&pool_health_check_period=1m", //nolint:lll
			expected: "9090",
		},
	}

	for _, c := range tCases {
		cfg := DBCfg{
			URL: c.url,
		}
		assert.Equal(t, c.expected, cfg.GetDatabasePort())
	}
}
