package util

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"os"
	"reflect"
	"strconv"
)

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}
	return nil
}
func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)
	err := v.Unmarshal(&env)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}
	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)
		switch valOf.Kind() {
		case reflect.String:
			val = valOf.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// This will catch iota constants defined as int types
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = strconv.FormatUint(valOf.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			val = strconv.FormatFloat(valOf.Float(), 'f', -1, 64)
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			// For complex types or custom types based on iota
			// Consider converting to string representation
			val = reflect.Indirect(valOf).String()
		}
		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}
	return nil
}
func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()
	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		logrus.Errorf("failed to add remote provider: %v", err)
		return err
	}
	err = v.ReadRemoteConfig()
	if err != nil {
		logrus.Errorf("failed to read remote config: %v", err)
		return err
	}
	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("failed to unmarshal: %v", err)
		return err
	}
	err = SetEnvFromConsulKV(v)
	if err != nil {
		logrus.Errorf("failed to set env from consul kv: %v", err)
		return err
	}
	return nil
}
