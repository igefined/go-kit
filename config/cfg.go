package config

import (
	"flag"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var envs = []*EnvVar{
	DefaultNamespace,
	DefaultEnvironment,

	DefaultHost,
	DefaultPort,
	DefaultMonitorHost,
	DefaultMonitorPort,

	DefaultRedisAddr,
	DefaultRedisDatabase,
	DefaultRedisPassword,

	DefaultRateLimitLuaScriptPath,
	DefaultRateLimitAuthRate,
	DefaultRateLimitAuthMaxTokens,
	DefaultRateLimitCommonRate,
	DefaultRateLimitCommonMaxTokens,

	DefaultDBURL,
	DefaultDBAutoCreate,
	DefaultMigrationPath,

	DefaultJWTSecretKey,
	DefaultExpirationTime,

	DefaultChainRPCtUrl,
	DefaultContractorAddr,
	DefaultOperatorAddr,
	DefaultOperatorPrivateKey,

	DefaultAWSAccessKeyID,
	DefaultAWSSecretKey,
	DefaultAWSRegion,
	DefaultAWSEndpoint,

	DefaultS3BucketName,
}

type (
	EnvVar struct {
		DefaultValue           interface{}
		Flag, Env, Description string
	}

	MainCfg struct {
		Namespace   string `mapstructure:"namespace"`
		Environment string `mapstructure:"environment"`

		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		MonitorHost string `mapstructure:"monitor_host"`
		MonitorPort string `mapstructure:"monitor_port"`
	}

	RedisCfg struct {
		Addr     string `mapstructure:"redis_address"`
		Password string `mapstructure:"redis_password"`
		Database int    `mapstructure:"redis_database"`
	}

	RateLimitCfg struct {
		RedisLuaScriptPath string `mapstructure:"rate_limit_lua_script_path"`

		AuthRate      float64 `mapstructure:"rate_limit_auth_rate"`
		AuthMaxTokens float64 `mapstructure:"rate_limit_auth_max_tokens"`

		CommonRate      float64 `mapstructure:"rate_limit_common_rate"`
		CommonMaxTokens float64 `mapstructure:"rate_limit_common_max_tokens"`
	}

	DBCfg struct {
		URL                string `mapstructure:"db_url"`
		MigrationsPath     string `mapstructure:"db_migrations_path"`
		AutoCreateDatabase bool   `mapstructure:"db_auto_create_database"`
	}

	JWTCfg struct {
		SecretKey      string        `mapstructure:"jwt_secret_key"`
		ExpirationTime time.Duration `mapstructure:"jwt_expiration_time"`
	}

	ETHCfg struct {
		RPCUrl             string `mapstructure:"chain_rpc_url"`
		ContractAddress    string `mapstructure:"contract_address"`
		OperatorAddr       string `mapstructure:"operator_address"`
		OperatorPrivateKey string `mapstructure:"operator_private_key"`
	}

	AWSCfg struct {
		AWSAccessKeyID string `mapstructure:"aws_access_key_id"`
		AWSSecretKey   string `mapstructure:"aws_secret_key"`
		AWSRegion      string `mapstructure:"aws_region"`
		AWSEndpoint    string `mapstructure:"aws_endpoint"`
	}

	S3 struct {
		S3BucketName string `mapstructure:"s3_bucket_name"`
	}

	TraceCfg struct {
		TracingServiceName  string `mapstructure:"tracing_service_name"`
		TracingCollectorURI string `mapstructure:"tracing_collector_uri"`
	}
)

func (c *DBCfg) GetDatabaseName() string {
	re := regexp.MustCompile(`(([0-9]+\/)([a-z_]+)+)`)
	out := strings.Split(re.FindString(c.URL), "/")

	if len(out) == 2 {
		return out[1]
	}

	return ""
}

func (c *DBCfg) GetDatabaseUser() string {
	startIndex := strings.Index(c.URL, "//")
	endIndex := strings.Index(c.URL, "@")

	if startIndex != -1 && endIndex != -1 {
		return strings.Split(c.URL[startIndex+2:endIndex], ":")[0]
	}

	return ""
}

func (c *DBCfg) GetDatabasePassword() string {
	startIndex := strings.Index(c.URL, "//")
	endIndex := strings.Index(c.URL, "@")

	if startIndex != -1 && endIndex != -1 {
		return strings.Split(c.URL[startIndex+2:endIndex], ":")[1]
	}

	return ""
}

func (c *DBCfg) GetDatabasePort() string {
	u, err := url.Parse(c.URL)
	if err != nil {
		return ""
	}

	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return ""
	}

	return port
}

func NewEnvVar(flag, env string, defaultValue interface{}, description string) *EnvVar {
	return &EnvVar{Flag: flag, Env: env, Description: description, DefaultValue: defaultValue}
}

func AddEnvs(customEnvs []*EnvVar) {
	var tmpEnvs []*EnvVar
	tmpEnvs = append(tmpEnvs, customEnvs...)
	for _, defaultEnv := range envs {
		check := true
		for _, customEnv := range customEnvs {
			if customEnv.Flag == defaultEnv.Flag {
				check = false
				break
			}
		}

		if check {
			tmpEnvs = append(tmpEnvs, defaultEnv)
		}
	}

	envs = tmpEnvs
}

func BindConfig() {
	for _, e := range envs {
		switch val := e.DefaultValue.(type) {
		case string:
			flag.String(e.Flag, val, e.Description)
		case int:
			flag.Int(e.Flag, val, e.Description)
		case bool:
			flag.Bool(e.Flag, val, e.Description)
		case time.Duration:
			flag.Duration(e.Flag, val, e.Description)
		default:
			continue
		}
		if e.DefaultValue != nil {
			viper.SetDefault(e.Env, e.DefaultValue)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	for _, e := range envs {
		_ = viper.BindEnv(e.Env)
	}
}

func GetConfig(cfg interface{}, customEnvs []*EnvVar) error {
	AddEnvs(customEnvs)
	BindConfig()
	return viper.Unmarshal(cfg)
}
