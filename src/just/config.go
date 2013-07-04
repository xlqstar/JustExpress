package just

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Config map[string]string

func (config Config) GetInt(key string) int {
	value, err := strconv.Atoi(config[key])
	if err != nil {
		log.Fatal(key + "值未填写或填写不正确，请确认为整数")
	}
	return value
}

func (config Config) GetStr(key string) string {
	if config[key] == "" {
		log.Fatal(key + "值未填写或填写不正确")
	}
	return config[key]
}

func (config Config) GetArray(key string) []string {
	value := config.GetStr(key)
	array := strings.Split(value, "|")
	return array
}

func Configure(filePath string) Config {
	configByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("配置文件不存在！")
	}
	configMap := Config{}
	config := strings.Replace("\r\n", "\n", string(configByte), -1)
	config = strings.Replace("|\n", "|", config, -1)
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

func GetCategorys(categorys []string) []Category {
	var categorysSet []Category

	for _, categoryStr := range categorys {
		categoryArry := strings.Split(Trim(categoryStr), "@")
		var category Category
		name := strings.Split(Trim(categoryArry[0]), "(")
		category.Name = Trim(name[0])
		if len(name) > 1 {
			category.Alias = strings.TrimSuffix(Trim(name[1]), ")")
		} else {
			category.Alias = category.Name
		}
		if len(categoryArry) > 1 {
			category.Href = Trim(categoryArry[1])
		}
		categorysSet = append(categorysSet, category)
	}

	return categorysSet
}
