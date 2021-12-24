package config

import (
	"github.com/spf13/viper"
	"path"
	"runtime"
	"time"
)

type DB struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"name" yaml:"name"`
}

type HTTP struct {
	Port string `json:"port" yaml:"port"`
}

type GRPC struct {
	Port string `json:"port" yaml:"port"`
}

type Metrics struct {
	Port string `json:"port" yaml:"port"`
}

type Server struct {
	Name    string `json:"name" yaml:"name"`
	Host    string `json:"host" yaml:"host"`
	GRPC    GRPC
	HTTP    HTTP
	Metrics Metrics
}

type Blog struct {
	Server Server
}

type User struct {
	DB     DB
	Server Server
}

type Post struct {
	DB     DB
	Server Server
}

type Comment struct {
	DB     DB
	Server Server
}

type Auth struct {
	Server Server
}

type DTM struct {
	Server Server
}

type JWT struct {
	Secret  string        `json:"secret" yaml:"secret"`
	Expires time.Duration `json:"expires" yaml:"expires"`
}

type Config struct {
	Blog    Blog
	User    User
	Post    Post
	Comment Comment
	Auth    Auth
	DTM     DTM
	JWT     JWT
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = GetPath()
	}
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.JWT.Expires = conf.JWT.Expires * time.Second

	return conf, nil
}

func GetPath() string {
	dir := getSourcePath()
	return dir + "/../../../configs/config.yaml"
}

func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
