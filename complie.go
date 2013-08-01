//just主要入口函数，负责编译日志
package just

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var siteInfo SiteInfo //全局变量：站点信息

func Complie(siteDirPath string, onlyRebuildHtml bool) {
	var logList LogList             //全部日志列表
	var oldLogList map[string]int64 //已生成的日志列表
	var updatedLogList []string     //有修改变动的日志列表

	//路径参数配置列表
	var compliedDirPath = siteDirPath + "\\complied"
	var themeDirPath = compliedDirPath + "\\style"
	var postDirPath = compliedDirPath + "\\posts"
	var archiveDirPath = compliedDirPath + "\\archives"
	var tagDirPath = compliedDirPath + "\\tags"
	// var indexDirPath = compliedDirPath
	siteCfg := Configure(compliedDirPath + "\\setting")
	themeCfg := Configure(themeDirPath + "\\meta")

	siteInfo.Site = siteCfg.GetStr("site")
	siteInfo.Domain = siteCfg.GetStr("domain")
	siteInfo.Author = siteCfg.GetStr("author")
	siteInfo.PageSize = siteCfg.GetInt("pageSize")
	siteInfo.Categorys = GetCategorys(siteCfg.GetArray("categorys"))
	siteInfo.Tags = GetTags(siteCfg.GetArray("tags"))
	siteInfo.Links = GetLinks(siteCfg.GetArray("links"))
	siteInfo.Socials = GetSocials(siteCfg.GetArray("socials"))

	siteInfo.ThemeName = themeCfg.GetStr("ThemeName")
	siteInfo.ImgWidth = themeCfg.GetInt("ImgWidth")
	siteInfo.OriginImgWidth = themeCfg.GetInt("OriginImgWidth")

	siteInfo.SitePath = siteDirPath

	oldLogList = getOldLogList(postDirPath)

	//解析日志、统计站点信息============================

	//解析log数据
	logDirList, _ := filepath.Glob(siteDirPath + "\\*")
	os.Mkdir(postDirPath, os.ModePerm)
	for k := range logDirList {
		if filepath.Base(logDirList[k]) == "complied" {
			continue
		}
		logInfo := Decode_log(logDirList[k], siteInfo)

		if logList.Contain(logInfo.Permalink) {
			log.Fatal("《" + logInfo.Title + "》文章转拼音后跟其他文章冲突,请指定alias别名")
		}

		logList = append(logList, logInfo)
	}

	//检查变动
	for k := range logList {
		//如果没有创建过 或者 如果有更新
		_, ok := oldLogList[logList[k].Permalink]
		if !ok || int64(logList[k].LastModTime) > oldLogList[logList[k].Permalink] || onlyRebuildHtml {
			updatedLogList = append(updatedLogList, logList[k].Title)
		}
	}

	// deletedLogList := oldLogList
	//检查是否有删除的日志，获取被删除日志列表
	for _, logInfo := range logList {
		delete(oldLogList, logInfo.Permalink)
	}

	//若无改动则直接返回
	if len(updatedLogList) == 0 && len(oldLogList) == 0 {
		//清理生成日志目录（删除未预料到的文件或者已经删除、改名的日志）
		cleanDestPostDir(postDirPath, oldLogList)
		//复制主题
		SyncTheme(themeDirPath, siteInfo.ThemeName)
		return
	}

	//日志列表排序
	logList = LogSort(logList)

	//=======================统计===========================

	//统计tag数据
	for k := range siteInfo.Tags {
		for kk := range logList {
			if strings.Contains(logList[kk].MetaData["tag"], siteInfo.Tags[k].Name) {
				siteInfo.Tags[k].Count += 1
			}
		}
	}

	//统计category数据
	for k := range siteInfo.Categorys {
		//统计每个分类底下的日志个数
		for kk := range logList {
			if strings.Contains(logList[kk].MetaData["category"], siteInfo.Categorys[k].Name) || logList[kk].Type == siteInfo.Categorys[k].Alias {
				siteInfo.Categorys[k].Count += 1
			}
		}
	}

	//统计archive数据
	for _, log := range logList {
		year_month := log.Date.Format("2006-1")
		length := len(siteInfo.Archives)
		if length == 0 || siteInfo.Archives[length-1].YearMonth != year_month {
			archive := Archive{YearMonth: year_month, Count: 1}
			siteInfo.Archives = append(siteInfo.Archives, archive)
		} else {
			siteInfo.Archives[length-1].Count += 1
		}
	}

	//生成全局引用模版
	siteInfo.GlobalTpl = make(map[string]string)
	tplList, _ := filepath.Glob(themeDirPath + "\\*")
	for _, tplFilePath := range tplList {
		fileName := filepath.Base(tplFilePath)
		if strings.HasPrefix(fileName, "_") && strings.HasSuffix(fileName, ".html") {
			key := strings.TrimPrefix(strings.TrimSuffix(fileName, ".html"), "_")
			template_byte, err := ioutil.ReadFile(tplFilePath)
			if err != nil {
				log.Fatal(tplFilePath + "模版文件读取失败，请检查该模版是否存在？")
			}
			siteInfo.GlobalTpl[key] = string(renderTpl(siteInfo, string(template_byte), strings.TrimSuffix(filepath.Base(tplFilePath), filepath.Ext(tplFilePath))))
		}
	}
	//========================生成==============================

	//日志生成
	last := len(logList) - 1
	for k, logInfo := range logList {
		if In_array(logInfo.Title, updatedLogList) {
			var logPage = LogPage{LogInfo: logList[k], SiteInfo: siteInfo, RelPath: "../../"}
			if (k - 1) < 0 {
				logPage.PrevLog = logList[last]
			} else {
				logPage.PrevLog = logList[k-1]
			}
			if (k + 1) > last {
				logPage.NextLog = logList[0]
			} else {
				logPage.NextLog = logList[k+1]
			}
			Build_log(logPage, themeDirPath, postDirPath, onlyRebuildHtml)
		} else {
			//TODO
			update_log(postDirPath+"\\"+logInfo.Permalink+"\\index.html", "../../")
		}

	}

	//总索引生成============================
	// os.Mkdir(indexDirPath, os.ModePerm)
	indexPage := IndexPage{Category: Category{Name: "首页", Alias: "index"}, SiteInfo: siteInfo, PageSize: siteInfo.PageSize, LogList: logList, RelPath: "./"}
	// if len(updatedLogList) > 0 || !Exist(compliedDirPath+"\\index.html") {
	Build_index(indexPage, themeDirPath, compliedDirPath)
	// }

	//按分类生成索引
	for _, category := range siteInfo.Categorys {
		haveUpdated := false
		_logList := LogList{}
		categoryDirPath := compliedDirPath + "\\" + category.Alias

		//判断该分类下的日志是否有变动(添加日志或修改日志)
		if category.Count == 0 {
			if !Exist(categoryDirPath + "\\index.html") {
				haveUpdated = true
			}
		} else if category.Count > 0 {
			//剔出该分类下的日志列表
			for _, logInfo := range logList {
				if strings.Contains(logInfo.MetaData["category"], category.Name) || logInfo.Type == category.Alias {
					_logList = append(_logList, logInfo)
				}
			}
			for _, logInfo := range _logList {
				if In_array(logInfo.Title, updatedLogList) /* || !Exist(categoryDirPath) */ {
					haveUpdated = true
					break
				}
			}
		} else {
			continue
		}

		indexPage.LogList = _logList
		indexPage.Category = category
		indexPage.RelPath = "../"

		if haveUpdated {
			Remkdir(categoryDirPath)
			Build_index(indexPage, themeDirPath, categoryDirPath)
		} else {
			//TODO
			update_index(indexPage, categoryDirPath)
		}

	}

	//标签页生成============================
	os.Mkdir(tagDirPath, os.ModePerm)

	tagPage := TagPage{SiteInfo: siteInfo, RelPath: "../"}
	//按标签生成索引
	for _, tag := range siteInfo.Tags {
		haveUpdated := false
		_logList := LogList{}
		tagPagePath := tagDirPath + "\\" + tag.Alias + ".html"

		//如果没有日志关联到该分类下
		if tag.Count == 0 {
			if !Exist(tagPagePath) {
				haveUpdated = true
			}
		} else if tag.Count > 0 {
			//选出关联至该标签的日志列表
			for _, logInfo := range logList {
				//如果关联到该标签并且有修改或是新创建日志
				if strings.Contains(logInfo.MetaData["tag"], tag.Name) {
					_logList = append(_logList, logInfo)
				}
			}
			//判断该分类下的日志是否有变动(添加日志或修改日志)
			for _, logInfo := range _logList {
				if In_array(logInfo.Title, updatedLogList) /* || !Exist(tagPagePath)*/ {
					haveUpdated = true
					break
				}
			}
		} else {
			continue
		}

		tagPage.LogList = _logList
		tagPage.Tag = tag

		if haveUpdated {
			Build_tagpage(tagPage, themeDirPath, tagDirPath)
		} else {
			//TODO
			update_tag(tagDirPath+"\\"+tag.Alias+".html", "../")
		}

	}

	//归档生成============================
	os.Mkdir(archiveDirPath, os.ModePerm)
	archivePage := ArchivePage{SiteInfo: siteInfo, RelPath: "../"}
	archives := []Archive{}
	for _, log := range logList {
		year_month := log.Date.Format("2006-1")
		length := len(archives)

		if length == 0 || siteInfo.Archives[length-1].YearMonth != year_month {
			var logList LogList
			logList = append(logList, log)
			archive := Archive{YearMonth: year_month, LogList: logList}
			archives = append(archives, archive)
		} else {
			archives[length-1].LogList = append(archives[length-1].LogList, log)
		}
	}
	if len(archives) > 0 {
		archivePage.Archives = archives
		Build_archive(archivePage, themeDirPath, archiveDirPath)
	}

	//清理生成日志目录（删除未预料到的文件或者已经删除、改名的日志）
	cleanDestPostDir(postDirPath, oldLogList)
	//复制主题
	SyncTheme(themeDirPath, siteInfo.ThemeName)
	build_loglistdata(compliedDirPath, logList)

}

func getOldLogList(postDirPath string) map[string]int64 {
	oldLogList := make(map[string]int64)
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

func cleanDestPostDir(postDirPath string, oldLogList map[string]int64) {
	for k := range oldLogList {
		_postDirPath := postDirPath + "\\" + k
		err := os.RemoveAll(_postDirPath)
		if err != nil {
			panic(err)
			log.Println("清理" + ConvertPath(_postDirPath) + "目录出现未预料到的问题")
		} else {
			log.Println(ConvertPath(_postDirPath) + "目录清理成功")
		}
	}
}

//同步主题
func SyncTheme(siteThemeDirPath string, themeName string) {

	if !Exist(siteThemeDirPath) {
		log.Fatal("模版文件缺失！可能是误删除！请使用tool添加模版文件或者直接手动复制。")
	} else /*if Exist(siteThemeDirPath + "\\meta")*/ {
		var localThemeList []string
		fileList, _ := filepath.Glob(".\\themes\\*")
		for k := range fileList {
			if len(themeName) > 0 {
				localThemeList = append(localThemeList, filepath.Base(fileList[k]))
			}
		}
		if !In_array(themeName, localThemeList) {
			CopyDir(siteThemeDirPath, ".\\themes\\"+themeName)
		}
	} /* else {
		log.Fatal("站点模板文件包不完整，缺失描述元数据文件！")
	}*/
}

func ConvertPath(path string) string {
	path = strings.TrimPrefix(path, siteInfo.SitePath)
	path = siteInfo.Site + path
	return path
}
