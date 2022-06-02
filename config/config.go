package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

func ExitsFile(path string) bool {
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func Read(path string, data interface{}) error {
	if !ExitsFile(path) {
		return fmt.Errorf("open %s error: File does not exist", path)
	}
	return loadConfig(reflect.ValueOf(data), path)
}

func loadConfig(vTarge reflect.Value, path string) error {
	oTarge := vTarge.Type()
	if oTarge.Elem().Kind() != reflect.Struct {
		return errors.New("type of received parameter is not struct")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, vTarge.Interface())
	if err != nil {
		return err
	}
	return nil
}

func Write(path string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, body, 0666)
	return err
}
