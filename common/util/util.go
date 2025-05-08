package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaginationParam struct {
	Page  int         `json:"page"`
	Count int64       `json:"count"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

type PaginationResult struct {
	TotalPage    int         `json:"totalPage"`
	Page         int         `json:"page"`
	NextPage     *int        `json:"nextPage"`
	PreviousPage *int        `json:"previousPage"`
	TotalData    int64       `json:"totalData"`
	Limit        int         `json:"limit"`
	Data         interface{} `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var (
		nextPage int
		prevPage int
	)

	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	if params.Page > 1 {
		prevPage = params.Page - 1
	}

	result := PaginationResult{
		TotalPage:    totalPage,
		Page:         params.Page,
		NextPage:     &nextPage,
		PreviousPage: &prevPage,
		TotalData:    params.Count,
		Limit:        params.Limit,
		Data:         params.Data}

	return result
}

func GenerateSHA256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func FormatRupiah(amount *float64) string {
	stringValue := "0"
	if amount != nil {
		humanizeValue := humanize.CommafWithDigits(*amount, 0)
		stringValue = strings.ReplaceAll(humanizeValue, ",", ".")
	}

	return fmt.Sprintf("Rp %s", stringValue)

}

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		logrus.Errorf("Error reading config file, %s", err)
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		logrus.Errorf("Error unmarshalling config file, %s", err)
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
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		default:
			panic("unsupported type")

		}

		err = os.Setenv(k, val)
		if err != nil {
			logrus.Errorf("failed to set env: %v", err)
			return err
		}
	}

	return nil
}

func BindFromConsul(dest any, endpoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")

	err := v.AddRemoteProvider("consul", endpoint, path)
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

func Recover() {
	if r := recover(); r != nil {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Errorf("recovered from panic: %v", r)
	}
}
