package globals

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Db    *gorm.DB
	V     *viper.Viper
	Redis *redis.Client
)
