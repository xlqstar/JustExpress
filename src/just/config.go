package just

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func Configure(filePath string) map[string]string {
	configByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("配置文件不存在！")
	}
	configMap := map[string]string{}
	config := strings.Replace("\r\n", "\n", string(configByte), -1)
	configArray := strings.Split(config, "\n")
	for _, v := range configArray {
		v = Trim(v)
		vArray := strings.SplitN(v, ":", 2)
		if len(vArray) == 2 && !strings.HasPrefix(v, "#") {
			key := Trim(vArray[0])
			value := Trim(vArray[1])
			configMap[key] = value
		}

	}
	return configMap
}

func SetInt(key string, config map[string]string) int {
	value, err := strconv.Atoi(config[key])
	if err != nil {
		log.Fatal(key + "值未填写或填写不正确，请确认为整数")
	}
	return value
}

func SetStr(key string, config map[string]string) string {
	if config[key] == "" {
		log.Fatal(key + "值未填写或填写不正确")
	}
	return config[key]
}

func Trim(s string) string {
	return strings.Trim(s, " \t\n\r")
}
