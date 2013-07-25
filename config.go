//配置文件解析函数集

package just

import (
	"github.com/xlqstar/Just/pinyin"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Config map[string]string

func (config Config) GetInt(key string) int {
	key = strings.ToLower(key)
	value, err := strconv.Atoi(config[key])
	if err != nil {
		log.Fatal(key + "值未填写或填写不正确，请确认为整数")
	}
	return value
}

func (config Config) GetStr(key string) string {
	key = strings.ToLower(key)
	if config[key] == "" {
		log.Fatal(key + "值未填写或填写不正确")
	}
	return config[key]
}

func (config Config) GetArray(key string) []string {
	key = strings.ToLower(key)
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
		v = strings.TrimSpace(v)
		vArray := strings.SplitN(v, ":", 2)
		if len(vArray) == 2 && !strings.HasPrefix(v, "#") {
			key := strings.ToLower(strings.TrimSpace(vArray[0]))
			value := strings.TrimSpace(strings.Trim(strings.TrimSpace(vArray[1]), "|"))
			configMap[key] = value
		}

	}

	return configMap
}

func GetCategorys(categorys []string) []Category {
	var categorysSet []Category

	for _, categoryStr := range categorys {
		categoryArry := strings.Split(strings.TrimSpace(categoryStr), "@")
		var category Category
		name := strings.Split(strings.TrimSpace(categoryArry[0]), "(")
		category.Name = strings.TrimSpace(name[0])
		if len(name) > 2 {
			log.Fatal(categoryStr + "配置格式有误(多个()符号)，请检查确认！")
		} else if len(name) > 1 {
			category.Alias = url.QueryEscape(pinyin.Convert(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(name[1]), ")")), ""))
		} else {
			category.Alias = url.QueryEscape(pinyin.Convert(category.Name, ""))
		}
		if len(categoryArry) > 2 {
			log.Fatal(categoryStr + "配置格式有误(多个@符号)，请检查确认")
		} else if len(categoryArry) > 1 {
			category.Href = strings.TrimSpace(categoryArry[1])
		}
		categorysSet = append(categorysSet, category)
	}

	return categorysSet
}

func GetTags(tags []string) []Tag {
	var tagSet []Tag

	for _, tagStr := range tags {
		var tag Tag
		name := strings.Split(strings.TrimSpace(tagStr), "(")
		tag.Name = strings.TrimSpace(name[0])

		if len(name) > 2 {
			log.Fatal(tagStr + "配置格式有误(多个()符号)，请检查确认！")
		} else if len(name) > 1 {
			tag.Alias = url.QueryEscape(pinyin.Convert(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(name[1]), ")")), ""))
		} else {
			tag.Alias = url.QueryEscape(pinyin.Convert(tag.Name, ""))
		}
		tagSet = append(tagSet, tag)
	}

	return tagSet
}

func GetLinks(links []string) []Link {
	var linkSet []Link
	for _, linkStr := range links {
		var link Link
		linkArry := strings.SplitN(strings.TrimSpace(linkStr), "@", 2)
		link.Name = strings.TrimSpace(linkArry[0])

		if len(linkArry) < 2 {
			log.Fatal(linkStr + "友情连接配置不完整，可能未配置相应连接！")
		} else {
			link.Href = strings.TrimSpace(linkArry[1])
		}
		linkSet = append(linkSet, link)
	}

	return linkSet
}

func GetSocials(socials []string) map[string]string {
	socialSet := map[string]string{}
	for _, socialStr := range socials {
		socialArry := strings.SplitN(strings.TrimSpace(socialStr), "@", 2)
		if len(socialArry) != 2 {
			log.Fatal(socialStr + "友情连接配置不完整，可能未配置相应连接！")
		}
		socialSet[strings.ToLower(strings.TrimSpace(socialArry[0]))] = strings.TrimSpace(socialArry[1])
	}

	return socialSet
}
