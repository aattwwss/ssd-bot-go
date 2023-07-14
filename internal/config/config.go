package config

type Config struct {
	// reddit config
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Username     string `env:"BOT_USERNAME,notEmpty"`
	Password     string `env:"BOT_PASSWORD,notEmpty"`
	Subreddit    string `env:"SUBREDDIT,notEmpty"`

	// techpowerup config
	TPUHost     string `env:"TPU_HOST,notEmpty"`
	TPUUsername string `env:"TPU_USERNAME,notEmpty"`
	TPUSecret   string `env:"TPU_SECRET,notEmpty"`

	// elasticsearch config
	EsAddress string `env:"ES_ADDRESS,notEmpty"`

	// application config
	OverrideOldBot bool `env:"OVERRIDE_OLD_BOT,notEmpty"`

	//debugging config
	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
}
