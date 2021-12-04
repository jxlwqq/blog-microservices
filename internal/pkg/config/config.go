package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type DB struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"name" yaml:"name"`
	DSN      string `json:"dsn" yaml:"dsn"`
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
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

type JWT struct {
	Secret  string        `json:"secret" yaml:"secret"`
	Expires time.Duration `json:"expires" yaml:"expires"`
}

type Config struct {
	User    User
	Post    Post
	Comment Comment
	Auth    Auth
	JWT     JWT
}

func Load(path string) (*Config, error) {
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

	conf.User.DB.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User.DB.User,
		conf.User.DB.Password,
		conf.User.DB.Host,
		conf.User.DB.Port,
		conf.User.DB.Name,
	)

	conf.Post.DB.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Post.DB.User,
		conf.Post.DB.Password,
		conf.Post.DB.Host,
		conf.Post.DB.Port,
		conf.Post.DB.Name,
	)

	conf.Comment.DB.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Comment.DB.User,
		conf.Comment.DB.Password,
		conf.Comment.DB.Host,
		conf.Comment.DB.Port,
		conf.Comment.DB.Name,
	)

	conf.JWT.Expires = conf.JWT.Expires * time.Second

	return conf, nil
}
