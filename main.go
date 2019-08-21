package main

import (
	"flag"
	config "github.com/kolesa-team/go-socks-ldap/configuration"
	"github.com/kolesa-team/go-socks-ldap/logger"
	"github.com/kolesa-team/go-socks-ldap/socks"
	"github.com/kolesa-team/go-socks-ldap/storage"
	"github.com/kolesa-team/go-socks-ldap/validator"
	"github.com/sirupsen/logrus"
	"log"
)

const DefaultConfig string = "config/development/config.toml"

func main() {
	configPath := flag.String("config", DefaultConfig, "Путь до файла конфига")
	debug := flag.Bool("debug", false, "Режим отладки")
	flag.Parse()

	cfg := config.NewConfig(*configPath)
	if err := cfg.Init(); err != nil {
		log.Fatalln("unable to parse config", err)
	}

	logger.Init(cfg.Logger, *debug)

	vld := validator.NewCryptMd5Validator()
	stg := storage.NewClient(cfg.Storage, vld)

	if err := stg.Init(); err != nil {
		logrus.Error(err)
	}

	srv := socks.NewClient(cfg.Server, stg)

	if err := srv.Init(); err != nil {
		logrus.Fatal(err)
	}

	logrus.Error(srv.Serve())
}
