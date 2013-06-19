/*

--buildtarget=log/index/both
--whichlog=from/all/auto

几种日志创建方式：

0：创建未创建的日志,以及有更新的日志
1：按日期创建日志
2：全部创建

*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"just"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	timeFormat = "2006-5-5"

	whichlog    = *flag.String("whichlog", "auto", "生成模式 (from/all/auto)")
	fromtime    = *flag.String("fromtime", "", "生成日志的起始时间 (2014-1-19)")
	buildtarget = *flag.String("buildtarget", "both", "生成目标 (log/index/both)")
	configFile  = *flag.String("config", "config", "配置文件")

	// whichlog    = "all"
	// fromtime    = "2014-1-19"
	// buildtarget = "log"
	// configFile  = "config"

	config        map[string]string
	srcDirPath    string
	destDirPath   string
	tplDirPath    string
	ifBuild_flag  bool
	pageSize      int
	smallImgWidth int
	bigImgWidth   int
)

func init() {

	flag.Parse()

	if whichlog == "from" {
		_, err := time.Parse(timeFormat, fromtime)
		if err != nil {
			log.Fatal("请正确填写起始时间，如：2014-1-19")
		}
	}

	config = just.Configure(configFile)
	srcDirPath = just.SetStr("srcDirPath", config)
	destDirPath = just.SetStr("destDirPath", config)
	tplDirPath = just.SetStr("tplDirPath", config)
	pageSize = just.SetInt("pageSize", config)
	smallImgWidth = just.SetInt("smallImgWidth", config)
	bigImgWidth = just.SetInt("bigImgWidth", config)

	check()

}

func main() {

	logDirList, _ := filepath.Glob(srcDirPath + "\\*")

	logList := getLogList()

	//just.filter_dir(&logDirList) //初步过滤

	if buildtarget == "log" || buildtarget == "both" {

		//日志生成
		for k := range logDirList {
			ifBuild_flag = false //置为假
			logInfo := just.Decode_log(logDirList[k])
			title := logInfo.Title
			createtime := logInfo.Date

			//filter log
			if whichlog == "auto" {

				//如果没有创建过
				if _, ok := logList[title]; !ok {
					ifBuild_flag = true
				}
				//如果有更新
				if logInfo.LastModTime < logList[title].LastBuildTime {
					ifBuild_flag = true
				}

			} else if whichlog == "from" {

				from_t, _ := time.Parse(timeFormat, fromtime)
				create_t, _ := time.Parse(timeFormat, createtime)

				if create_t.After(from_t) || create_t.Equal(from_t) {
					ifBuild_flag = true
				}

			} else if whichlog == "all" {
				ifBuild_flag = true
			}

			if ifBuild_flag {
				just.Build_log(logInfo, tplDirPath, destDirPath+"\\post", uint(smallImgWidth), uint(bigImgWidth))
				buildTime := time.Now().Unix()
				buildTimeStr := fmt.Sprintf("%d", buildTime)
				logInfo.LastBuildTime = buildTimeStr
				logList[title] = logInfo
			}
		}

	}
	if buildtarget == "index" || buildtarget == "both" {
		just.Build_index(logList, tplDirPath, destDirPath, pageSize)
	}
	just.Update_loglistdata(logList)
}

//读取loglist历史数据
func getLogList() map[string]just.LogInfo {
	logListSrc := "./loglist.json"
	logListStr, err := ioutil.ReadFile(logListSrc)
	if err != nil {
		logList := map[string]just.LogInfo{}
		return logList
	}
	var logList map[string]just.LogInfo
	json.Unmarshal(logListStr, &logList)
	return logList
}

func check() {
	if !exist(srcDirPath) {
		log.Fatal("日志目录不存在！")
	}

	if !exist(tplDirPath) {
		log.Fatal("模版目录不存在！")
	} else {
		if !exist(tplDirPath + "\\index.html") {
			log.Fatal("索引模版不存在！")
		}
		if !exist(tplDirPath + "\\article.html") {
			log.Fatal("文章模版不存在！")
		}
		if !exist(tplDirPath + "\\album.html") {
			log.Fatal("相册模版不存在！")
		}
		if !exist(tplDirPath + "\\theme") {
			log.Fatal("主题目录不存在！")
		}
	}

	if !exist(destDirPath) {
		err := os.Mkdir(destDirPath, os.ModePerm)
		if err != nil {
			log.Fatal("无法创建生成目录！")
		}
		err = os.Mkdir(destDirPath+"\\post", os.ModePerm)
		if err != nil {
			log.Fatal("无法创建日志目录！")
		}
	} else if !exist(destDirPath + "\\post") {

		err := os.Mkdir(destDirPath+"\\post", os.ModePerm)
		if err != nil {
			log.Fatal("无法创建日志目录！")
		}
	}
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
