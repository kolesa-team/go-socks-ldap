package configuration

import (
	"github.com/kolesa-team/go-socks-ldap/logger"
	"github.com/kolesa-team/go-socks-ldap/socks"
	"github.com/kolesa-team/go-socks-ldap/storage"
)

// config struct
type Config struct {
	// path to config file
	file string

	// socks server options
	Server socks.Options

	// storage options
	Storage storage.Options

	// logger options
	Logger logger.Options
}
