package config

type Config struct {
	ClientId       string `env:"CLIENT_ID,notEmpty"`
	ClientSecret   string `env:"CLIENT_SECRET,notEmpty"`
	Username       string `env:"BOT_USERNAME,notEmpty"`
	Password       string `env:"BOT_PASSWORD,notEmpty"`
	TPUHost        string `env:"TPU_HOST,notEmpty"`
	TPUSecret      string `env:"TPU_SECRET"`
	EsAccessKey    string `env:"ES_ACCESS_KEY"`
	EsAccessSecret string `env:"ES_ACCESS_SECRET"`

	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
}
