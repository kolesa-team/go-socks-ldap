package configuration

import (
	"github.com/BurntSushi/toml"
)

// Конструктор
func NewConfig(file string) *Config {
	return &Config{
		file: file,
	}
}

// Инициализация
func (c *Config) Init() error {
	_, err := toml.DecodeFile(c.file, c)
	return err
}
