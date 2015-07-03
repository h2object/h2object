package commands

import (
	"os"
	"path"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
	"errors"
	"time"
)

type AuthConfig struct{
	AuthID 	string 	`json:"authid,omitempty"`
	Secret 	string	`json:"secret"`
	Token  	string 	`json:"token"`
	ExpireAt int64  `json:"expire_at"`
}

type ConfigFile struct{
	Auth 		*AuthConfig	`json:"auth"`
	filename 	string	
}

func NewConfigFile(fn string) *ConfigFile {
	return &ConfigFile{
		Auth: &AuthConfig{},
		filename:    fn,
	}
}

func LoadConfigFile(workdir string) (*ConfigFile, error) {
	if workdir == "" {
		return nil, errors.New("working directory absent")
	}

	fn := path.Join(workdir, ".h2object")

	configFile := NewConfigFile(fn)

	if _, err := os.Stat(configFile.filename); err == nil {
		file, err := os.Open(configFile.filename)
		if err != nil {
			return configFile, err
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(configFile); err != nil {
			return configFile, err
		}

		if time.Now().Unix() >= configFile.Auth.ExpireAt {
			configFile.Auth.Token = ""
			configFile.Auth.ExpireAt = 0
		}

	} else if !os.IsNotExist(err) {
		return configFile, err
	}

	return configFile, nil 
}

func (configFile *ConfigFile) Save() error {
	configFile.Auth.AuthID = ""
	data, err := json.MarshalIndent(configFile, "", "\t")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configFile.filename), 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(configFile.filename, data, 0600); err != nil {
		return err
	}
	return nil
}

func (configFile *ConfigFile) Remove() error {
	return os.Remove(configFile.filename)
}
