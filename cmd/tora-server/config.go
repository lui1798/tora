package main

import (
	"bytes"
	"github.com/leizongmin/tora/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Listen string       `json:"listen"` // 监听地址
	Log    ConfigLog    `json:"log"`    // 日志配置
	Enable []string     `json:"enable"` // 开启的模块
	Module ConfigModule `json:"module"` // 模块对应的配置
	Auth   ConfigAuth   `json:"auth"`   // 授权规则
}

type ConfigLog struct {
	Format string `json:"format"` // 日志格式，text或json
	Level  string `json:"level"`  // 日志等级，可选：debug，info，warning，error，fatal，panic
}

type ConfigAuth struct {
	Token map[string]ConfigAuthItem `json:"token"` // 允许指定token
	IP    map[string]ConfigAuthItem `json:"ip"`    // 允许指定ip
}

type ConfigAuthItem struct {
	Allow   bool     `json:"allow"`   // 是否允许访问
	Modules []string `json:"modules"` // 允许访问的模块
}

type ConfigModule struct {
	File  ConfigModuleFile  `json:"file"`  // file 模块配置
	Shell ConfigModuleShell `json:"shell"` // shell 模块配置
	Log   ConfigModuleLog   `json:"log"`   // log 模块配置
}

type ConfigModuleFile struct {
	Root         string `json:"root"`         // 根目录
	AllowPut     bool   `json:"allowPut"`     // 允许上传文件
	AllowDelete  bool   `json:"allowDelete"`  // 允许删除文件
	AllowListDir bool   `json:"allowListDir"` // 允许列出目录
}

type ConfigModuleShell struct{}

type ConfigModuleLog struct{}

func GetDefaultConfig() Config {
	c := Config{
		Listen: server.DefaultListenAddr,
		Log: ConfigLog{
			Format: "text",
			Level:  "info",
		},
		Enable: []string{},
		Module: ConfigModule{
			File:  ConfigModuleFile{},
			Shell: ConfigModuleShell{},
			Log:   ConfigModuleLog{},
		},
		Auth: ConfigAuth{
			IP:    make(map[string]ConfigAuthItem),
			Token: make(map[string]ConfigAuthItem),
		},
	}
	c.Auth.IP["127.0.0.1"] = ConfigAuthItem{Allow: true, Modules: []string{}}
	return c
}

func LoadConfigFile(filename string) (Config, error) {
	c := GetDefaultConfig()
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	d := yaml.NewDecoder(bytes.NewReader(b))
	err = d.Decode(&c)
	return c, err
}