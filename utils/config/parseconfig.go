package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	App App `json:"app"`

	MySql MySql `json:"mysql"`

	Redis Redis `json:"redis"`

	LogConfig LogConfig `json:"log"`
}

type App struct {
	Appname    string `json:"app_name"`
	Apphost    string `json:"app_host"`
	Appport    string `json:"app_port"`
	StaticPath string `json:"static_path"`
	TmpPath    string `json:"tmp_path"`
	AuthPath   string `json:"persvali_path"`
}

type MySql struct {
	HostName string `json:"hostname"`
	Port     string `json:"port"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	DbName   string `json:"dbname"`
}

type Redis struct {
	Redishost string `json:"redis_host"`
	Redisport string `json:"redis_port"`
	Redisauth string `json:"redis_auth"`
	Redisdb   int    `json:"redis_db"`
}

type LogConfig struct {
	Level      string `json:"level"`       // 记录等级
	Filename   string `json:"filename"`    // 保存的文件名
	MaxSize    int    `json:"maxsize"`     // 在进行切割之前,日志文件的最大大小
	MaxAge     int    `json:"max_age"`     // 保留旧文件的最大天数
	MaxBackups int    `json:"max_backups"` // 保留旧文件的最大个数
	Compress   bool   `json:"compress"`    // 是否压缩/归档旧文件
}

var Cfg *Config = nil

// 解析Config 文件
// path: config 文件的路径
// return：*Config Config结构体；error：错误信息
func ParseConfig(path string) (*Config, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// 从io 中读取文件数据
	buff := bufio.NewReader(file)
	// 解码成json
	decode := json.NewDecoder(buff)
	if err = decode.Decode(&Cfg); err != nil {
		return nil, err
	}
	return Cfg, nil
}
