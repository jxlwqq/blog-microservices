package config

import (
	"fmt"
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
	DSN      string `json:"dsn" yaml:"dsn"`
}

func (db *DB) setDSN() {
	db.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
	)
}

type Server struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
	Addr string `json:"addr" yaml:"addr"`
}

func (s *Server) setAddr() {
	s.Addr = fmt.Sprintf("%s%s", s.Host, s.Port)
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

	conf.User.DB.setDSN()
	conf.Post.DB.setDSN()
	conf.Comment.DB.setDSN()
	conf.User.Server.setAddr()
	conf.Post.Server.setAddr()
	conf.Comment.Server.setAddr()
	conf.Auth.Server.setAddr()

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
