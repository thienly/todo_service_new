package config

import "testing"

func TestLoadConfig(t *testing.T) {
	config := LoadConfig()
	if config.Db.Conn == "" {
		t.Fail()
	}
}