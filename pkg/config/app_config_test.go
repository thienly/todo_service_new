package config_test

import (
	"new_todo_project/pkg/config"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestLoadFromJsonOrPanic(t *testing.T){
	dir, _ := os.Getwd()
	rootPath:= path.Dir(path.Dir(dir))
	data, err:= config.LoadFromJsonOrPanic(filepath.Join(rootPath,"mocks","config","config.json"))
	if err != nil {
		t.Fail()
	}
	if data.Email == nil {
		t.Fail()
	}
}