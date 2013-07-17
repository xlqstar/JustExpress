//just主要入口函数，负责编译日志
package just

import (
	"flag"
	"just/pinyin"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	config        Config
	srcDirPath    string
	destDirPath   string
	tplDirPath    string
	pageSize      int
	smallImgWidth int
	bigImgWidth   int
	siteInfo      SiteInfo

	themeDirPath   string
	postDirPath    string
	tagDirPath     string
	archiveDirPath string
	indexDirPath   string
)

func Complie(configFile string) {

	flag.Parse()

	config = Configure(configFile)
	srcDirPath = config.GetStr("srcDirPath")
	destDirPath = config.GetStr("destDirPath")
	tplDirPath = "templates\\" + config.GetStr("template")
	pageSize = config.GetInt("pageSize")
	smallImgWidth = config.GetInt("smallImgWidth")
	bigImgWidth = config.GetInt("bigImgWidth")
	siteInfo.Categorys = GetCategorys(config.GetArray("categorys"))
	siteInfo.Tags = GetTags(config.GetArray("tags"))

	//目标生成目录结构
	themeDirPath = destDirPath + "\\theme"
	postDirPath = destDirPath + "\\post"
	archiveDirPath = destDirPath + "\\archive"
	tagDirPath = destDirPath + "\\tag"
	indexDirPath = destDirPath

	logDirList, _ := filepath.Glob(srcDirPath + "\\*")

	oldLogList := getOldLogList(postDirPath)
	var logList LogList

	err := os.Mkdir(destDirPath, os.ModePerm)
	if os.IsExist(err) {
		if len(os.Args) > 1 && os.Args[1] == "rebuild" {
			Remkdir(destDirPath)
			log.Println("重新构建日志生成目录！")
		}
	} else {
		log.Fatal("日志生成目录设置错误，路径不存在！")
	}

	//日志生成============================

	os.Mkdir(postDirPath, os.ModePerm)

	for k := range logDirList {
		logInfo := Decode_log(logDirList[k])
		_title := logInfo.Title
		if logInfo.MetaData["alias"] == "" {
			_title = url.QueryEscape(pinyin.Convert(_title))
		} else {
			_title = logInfo.MetaData["alias"]
		}
		if logList.Contain(_title) {
			log.Fatal("《" + logInfo.Title + "》文章转拼音后跟其他文章冲突,请指定alias别名")
		}

		//如果没有创建过 或者 如果有更新
		_, ok := oldLogList[_title]

		if !ok || logInfo.LastModTime > oldLogList[_title] {
			var logPage = LogPage{LogInfo: logInfo, SiteInfo: siteInfo, RelPath: "../../"}
			Build_log(logPage, tplDirPath, postDirPath, uint(smallImgWidth), uint(bigImgWidth))
		}
		logList = append(logList, logInfo)
		delete(oldLogList, _title)
	}

	//索引生成============================

	os.Mkdir(indexDirPath, os.ModePerm)

	indexPage := IndexPage{Category: Category{Name: "首页", Alias: "index"}, SiteInfo: siteInfo, PageSize: pageSize, LogList: logList}
	//生成全部索引
	Build_index(indexPage, tplDirPath, destDirPath)
	//按分类生成索引
	for _, category := range siteInfo.Categorys {
		_logList := LogList{}
		for _, v := range logList {
			if strings.Contains(v.MetaData["category"], category.Name) {
				_logList = append(_logList, v)
			}
		}

		if len(_logList) > 0 {
			indexPage.LogList = _logList
			indexPage.Category = category
			categoryDirPath := destDirPath + "\\" + category.Alias
			Remkdir(categoryDirPath)
			Build_index(indexPage, tplDirPath, categoryDirPath)
		}
	}

	//标签页生成============================

	Remkdir(tagDirPath)

	tagPage := TagPage{SiteInfo: siteInfo, RelPath: "../"}
	//按标签生成索引
	for _, tag := range siteInfo.Tags {
		_logList := LogList{}
		for _, v := range logList {
			if strings.Contains(v.MetaData["tag"], tag.Name) {
				_logList = append(_logList, v)
			}
		}

		if len(_logList) > 0 {
			tagPage.LogList = _logList
			tagPage.Tag = tag
			Build_tagpage(tagPage, tplDirPath, tagDirPath)
		}
	}

	//归档生成============================

	Remkdir(archiveDirPath)

	archivePage := ArchivePage{SiteInfo: siteInfo, RelPath: "../"}
	archives := make(map[string]LogList)
	for _, log := range logList {
		// dateArry := strings.Split(log.Date, "-")
		t := time.Unix(int64(log.Date), int64(0))
		year_month := strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month()))
		if _, ok := archives[year_month]; ok {
			archives[year_month] = append(archives[year_month], log)
		} else {
			var archive LogList
			archive = append(archive, log)
			archives[year_month] = archive
		}
	}
	if len(archives) > 0 {
		archivePage.Archives = archives
		Build_archive(archivePage, tplDirPath, archiveDirPath)
	}
	//清理生成日志目录（删除未预料到的文件或者已经删除、改名的日志）
	cleanDestPostDir(postDirPath, oldLogList)
	//复制主题
	copyTheme(tplDirPath+"\\theme", themeDirPath)
	build_loglistdata(destDirPath, logList)

}

func getOldLogList(postDirPath string) map[string]int {
	oldLogList := make(map[string]int)
	dirList, _ := filepath.Glob(postDirPath + "\\*")
	for _, v := range dirList {
		fi, err := os.Stat(v)
		if err != nil {
			log.Println("读取" + v + "文件出现未预料的问题")
			continue
		} else {
			oldLogList[filepath.Base(v)] = GetCreationTime(fi)
		}
	}
	return oldLogList
}

func cleanDestPostDir(postDirPath string, oldLogList map[string]int) {
	for k := range oldLogList {
		_postDirPath := postDirPath + "\\" + k
		err := os.RemoveAll(_postDirPath)
		if err != nil {
			log.Println("清理" + _postDirPath + "目录出现未预料到的问题")
		} else {
			log.Println("清理" + _postDirPath + "目录成功")
		}
	}
}

func copyTheme(srcThemeDirPath string, destThemeDirPath string) {
	if !Exist(destThemeDirPath) {
		CopyDir(srcThemeDirPath, destThemeDirPath)
	}
}

//读取loglist历史数据
/*func getLogList() LogList {
	logListSrc := "./loglist.json"
	logListStr, err := ioutil.ReadFile(logListSrc)
	if err != nil {
		logList := just.LogList{}
		return logList
	}
	var logList just.LogList
	json.Unmarshal(logListStr, &logList)
	var logInfo just.LogInfo
	for k, _ := range logList {
		logInfo = logList[k]
		switch logInfo.Log.(type) {
		case string:
			logInfo.Log = just.Article(logInfo.Log.(string))
		case []map[string]string:
			logInfo.Log = just.Album(logInfo.Log.([]map[string]string))
		}
		logList[k] = logInfo
	}
	return logList
}*/
