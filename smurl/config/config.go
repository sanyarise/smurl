package config

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Структура конфигурации
type Config struct {
	DNS                string `toml:"dns" default:"postgres://postgres:1110@localhost/test?sslmode=disable"`
	Port               string `toml:"port" default:":8000"`
	ServerURL          string `toml:"server_url" default:"http://localhost/"`
	ReadTimeout        int    `toml:"read_timeout" default:"30"`
	WriteTimeout       int    `toml:"write_timeout" default:"30"`
	WriteHeaderTimeout int    `toml:"write_header_timeout" default:"30"`
	Logger             *zap.Logger
	LogLevel           int `toml:"log_level" default:"2"`
}

//Функция для инициализации конфигурации
func InitConfig() (*Config, error) {
	var configPath string

	//При запуске без флага с указанием пути с файлом конфигурации
	//по умолчанию путь "./config/config.toml"
	flag.StringVar(&configPath, "config-path", "./config/config.toml", "path to file in .toml format")
	flag.Parse()

	var cfg = Config{}
	//Декодирование файла конфигурации
	s, err := toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("can't load configuration file: %s", err)
	}
	log.Printf("load config successful %v", s)
	atomicLevel := zap.NewAtomicLevel()
	//Установка уровня логирования на основании данных из
	//файла конфигурации
	switch cfg.LogLevel {
	case 0:
		{
			atomicLevel.SetLevel(zap.InfoLevel)
		}
	case 1:
		{
			atomicLevel.SetLevel(zap.WarnLevel)
		}
	case 2:
		{
			atomicLevel.SetLevel(zap.DebugLevel)
		}
	case 3:
		{
			atomicLevel.SetLevel(zap.ErrorLevel)
		}
	case 4:
		{
			atomicLevel.SetLevel(zap.PanicLevel)
		}
	case 5:
		{
			atomicLevel.SetLevel(zap.FatalLevel)
		}
	}
	//Установка параметров логгера
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	), zap.AddCaller())
	cfg.Logger = logger
	return &cfg, err
}
