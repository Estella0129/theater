package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TMDB struct {
		APIToken string `yaml:"api_token"`
	} `yaml:"tmdb"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

var AppConfig Config

// JWTSecret 用于JWT token签名的密钥
var JWTSecret string

func init() {
	// 读取配置文件
	yamlFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	err = yaml.Unmarshal(yamlFile, &AppConfig)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}

// GetTMDBToken 从配置中获取TMDB API Token
func GetTMDBToken() (string, error) {
	if AppConfig.TMDB.APIToken == "" {
		return "", fmt.Errorf("TMDB API Token未配置")
	}
	return AppConfig.TMDB.APIToken, nil
}

// GetJWTSecret 从配置中获取JWT签名密钥
func GetJWTSecret() (string, error) {
	if AppConfig.JWT.Secret == "" {
		return "", fmt.Errorf("JWT签名密钥未配置")
	}
	JWTSecret = AppConfig.JWT.Secret
	return JWTSecret, nil
}
