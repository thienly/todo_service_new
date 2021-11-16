package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var AppCfg = &AppConfig{}

type AppConfig struct {
	Email *Email `json:"Email"`
}
type Email struct {
	EmailFrom         string `json:"From"`
	EmailFromPassword string `json:"Password"`
}


func LoadFromJsonOrPanic(filePath string)  (*AppConfig,error){
	file, err:= os.Open(filePath)
	if err != nil {
		panic("Can not open the file")
	}
	byteValue,err := ioutil.ReadAll(file)
	if err != nil {
		panic("can not read bytes")
	}
	err = json.Unmarshal(byteValue, AppCfg)
	if err != nil {
		panic("Can not parse to object")
	}
	return AppCfg, nil
}