package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// LogConfig 日志配置（保留你参考项目字段）
type LogConfig struct {
	Level string `yaml:"level"`
}

// AiServer AI 服务配置：Go 转发到 Python FastAPI
type AiServer struct {
	ApiUrl string `yaml:"apiUrl"`
}

// StorageConfig 存储配置：无数据库，全部落地到本地目录树
type StorageConfig struct {
	RootDir      string `yaml:"rootDir"`      // 教学文件存储库根目录
	MaxUploadMB  int    `yaml:"maxUploadMB"`  // 上传大小限制（MB）
	EnableLock   bool   `yaml:"enableLock"`   // 是否启用分析锁
	EnableCache  bool   `yaml:"enableCache"`  // 是否缓存 viz.json
	AllowExtList []string `yaml:"allowExtList"` // 允许上传扩展名（小写，不带点）
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
}

// Config 全局配置
type Config struct {
	App     AppConfig     `yaml:"app"`
	Log     LogConfig     `yaml:"log"`
	AiServer AiServer     `yaml:"aiServer"`
	Storage StorageConfig `yaml:"storage"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) error {
	if configPath == "" {
		_, filename, _, _ := runtime.Caller(0)
		configDir := filepath.Dir(filename)
		configPath = filepath.Join(configDir, "config.yaml")
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("app.name", "adm_bkd")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.host", "0.0.0.0")
	viper.SetDefault("app.port", 6677)

	viper.SetDefault("log.level", "debug")

	// 你的根目录要求
	viper.SetDefault("storage.rootDir", "/data/teaching_repo")
	viper.SetDefault("storage.maxUploadMB", 500)
	viper.SetDefault("storage.enableLock", true)
	viper.SetDefault("storage.enableCache", true)
	viper.SetDefault("storage.allowExtList", []string{
		"pdf", "doc", "docx", "xls", "xlsx",
		"png", "jpg", "jpeg", "webp",
		"mp4", "mov", "avi",
		"txt", "md",
	})
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}

func (c *Config) GetLogLevel() string { return c.Log.Level }

func (c *Config) GetAiServer() AiServer { return c.AiServer }

func (c *Config) GetStorage() StorageConfig { return c.Storage }
