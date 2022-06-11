package global

import (
	"api.frank.top/stockInfo/config"
	"api.frank.top/stockInfo/utils/timer"
	"api.frank.top/stockInfo/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB     *gorm.DB
	GVA_DBList map[string]*gorm.DB
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	GVA_LOG    *zap.Logger
	GVA_Timer  timer.Timer = timer.NewTimerTask()
	GVA_REDIS  *redis.Client
)
