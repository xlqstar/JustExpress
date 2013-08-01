//just工具链

package just

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetSitePath(siteName string) string {
	sitesRoot := SitesRoot("")
	if len(siteName) == 0 {
		var sitesArry []string
		fileList, _ := filepath.Glob(sitesRoot + "\\*")
		for k := range fileList {
			if Exist(fileList[k] + "\\complied\\setting") {
				sitesArry = append(sitesArry, fileList[k])
			}
		}
		if len(sitesArry) == 0 {
			log.Fatal("站点根目录不存在任何站点")
		} else if len(sitesArry) > 1 {
			var sites string
			for _, v := range sitesArry {
				if sites != "" {
					sites += ","
				}
				sites += filepath.Base(v)
			}
			log.Fatal("有多个站点(" + sites + ")，请通过-site参数明确指定当前操作对象")
		}
		return sitesArry[0]
	}

	sitePath := sitesRoot + "\\" + siteName
	if !Exist(sitePath + "\\complied\\setting") {
		log.Fatal(sitePath + "\\complied\\setting" + "站点目录结构不符合预期，有异常！确定站点标识输入是否错误！？")
	}

	return sitePath
}

func Delete(sitePath string, title string) {
	fileList, _ := filepath.Glob(sitePath + "\\*")
	for k := range fileList {
		fileName := filepath.Base(fileList[k])
		pos := strings.LastIndex(fileName, "@")
		var fileTitle string
		if pos >= 0 {
			fileTitle = fileName[0:pos]
		}
		if fileTitle == title {
			err := os.RemoveAll(fileList[k])
			if err != nil {
				log.Fatal("删除《" + title + "》操作失败")
			}
			log.Println("删除《" + title + "》操作成功")
			return
		}
	}
	log.Fatal("未找到您输入的日志")
}

func ImgResize(sitePath string) {
	fileList, _ := filepath.Glob(sitePath + "\\*")
	for _, path := range fileList {
		if Parse_logType(path) == "album" {
			var siteInfo SiteInfo
			logInfo := Decode_log(path, siteInfo)
			os.RemoveAll(sitePath + "\\complied\\posts\\" + logInfo.Permalink)
		}
	}
	Complie(sitePath, false)
}

func SwitchTheme(sitePath string, themeName string) {
	SyncTheme(sitePath+"\\complied\\style", themeName)
	os.RemoveAll(sitePath + "\\complied\\style")
	CopyDir(".\\themes\\"+themeName, sitePath+"\\complied\\style")
	// Complie(sitePath, true)
}

//重新构建(只构建html部分)
func Rebuild(sitePath string) {
	Complie(sitePath, true)
}

//重新构建(彻底重新构建，包括图片及其所有附件)
func RebuildAll(sitePath string) {
	fileList, _ := filepath.Glob(sitePath + "\\complied\\*")
	for k := range fileList {
		if !(filepath.Base(fileList[k]) == "style" || filepath.Base(fileList[k]) == "setting") {
			os.RemoveAll(fileList[k])
		}
	}
	Complie(sitePath, false)
}

func NewSite(siteName string) {

	if !Exist(siteName) {
		dataBytes, _ := ioutil.ReadFile("data")
		reg, _ := regexp.Compile(`#===站点默认配置[\s\S]*#===//站点默认配置`)

		content := string(reg.Find(dataBytes))
		content = strings.TrimPrefix(content, "#===站点默认配置")
		content = strings.TrimSuffix(content, "#===//站点默认配置")
		content = strings.TrimSpace(content)
		content = strings.Replace(content, "{{sitename}}", siteName, -1)

		ioutil.WriteFile(siteName, []byte(content), os.ModePerm)
		log.Println("已生成配置文件（" + siteName + "），请自定义调整站点配置参数！")
	} else {
		fi, _ := os.Stat(siteName)
		if fi.IsDir() {
			if Exist(siteName + "\\complied\\setting") { //是站点的话
				log.Fatal("已经存在该站点，请勿重复创建")
			} else {
				log.Fatal("该目录存在于站点同名的文件或文件夹，请检查")
			}
		} else {
			fileData, _ := ioutil.ReadFile(".\\" + siteName)
			os.Remove(".\\" + siteName)
			sitesRoot := SitesRoot("")
			sitePath := sitesRoot + "\\" + siteName
			err := os.MkdirAll(sitePath+"\\complied", os.ModePerm)
			if err != nil {
				ioutil.WriteFile(".\\"+siteName, fileData, os.ModePerm)
				panic(err)
				log.Fatal(sitePath + "目录创建失败。")
			}
			ioutil.WriteFile(".\\"+siteName+"\\complied\\setting", fileData, os.ModePerm)
			// ioutil.WriteFile(".\\"+siteName, fileData, os.ModePerm)
			CopyDir(".\\themes\\default", sitePath+"\\complied\\style")
		}
	}

}

func SitesRoot(sitesRoot string) string {
	if len(sitesRoot) == 0 {
		config := Configure(".\\data")
		sitesRoot = config.GetStr("SitesRoot")
		return sitesRoot
	} else {
		if !Exist(sitesRoot) {
			log.Fatal("站点根目录不存在")
		}
		data := "SitesRoot:\t" + sitesRoot
		err := ioutil.WriteFile("data", []byte(data), os.ModePerm)
		if err != nil {
			return "true"
		}
		return "false"
	}
}

func Post(sitePath string, title string, logType string, categorys string, tags string) {
	var meta_categorysArry []string
	var meta_tagsArry []string
	var meta_categorys string
	var meta_tags string
	var siteCategroysArry []string
	var siteTagsArry []string
	var metadata string
	siteCfg := Configure(sitePath + "\\complied\\setting")
	siteCategroys := GetCategorys(siteCfg.GetArray("categorys"))
	categorysArry := strings.Split(categorys, ",")
	for _, v := range siteCategroys {
		siteCategroysArry = append(siteCategroysArry, v.Name)
	}
	for k := range categorysArry {
		if categorysArry[k] == "index" {
			continue
		}
		if !In_array(categorysArry[k], siteCategroysArry) {
			log.Fatal(categorysArry[k] + " 该分类不存在于站点分类中！")
		} else {
			meta_categorysArry = append(meta_categorysArry, categorysArry[k])
		}
	}

	if len(tags) > 0 {
		siteTags := GetTags(siteCfg.GetArray("tags"))
		for _, v := range siteTags {
			siteTagsArry = append(siteTagsArry, v.Name)
		}

		tagsArry := strings.Split(tags, ",")
		for k := range tagsArry {
			if !In_array(tagsArry[k], siteTagsArry) {
				log.Fatal(tagsArry[k] + " 该标签不存在于站点标签库中！")
			} else {
				meta_tagsArry = append(meta_tagsArry, tagsArry[k])
			}
		}
	}
	for k := range meta_categorysArry {
		if meta_categorys != "" {
			meta_categorys += ","
		}
		meta_categorys += meta_categorysArry[k]
	}
	for k := range meta_tagsArry {
		if meta_tags != "" {
			meta_tags += ","
		}
		meta_tags += meta_tagsArry[k]
	}

	if len(meta_categorys) > 0 {
		metadata += "category:" + meta_categorys + "\n"
	}
	if len(meta_tags) > 0 {
		metadata += "tag:" + meta_tags + "\n"
	}
	if len(metadata) > 0 {
		metadata = "---\n" + metadata + "---\n"
	}
	logPath := sitePath + "\\" + title + "@" + time.Now().Format("2006-1-2")
	os.Mkdir(logPath, os.ModePerm)
	if logType == "article" {
		ioutil.WriteFile(logPath+"\\article.md", []byte(metadata), os.ModePerm)
	} else {
		ioutil.WriteFile(logPath+"\\meta", []byte(metadata), os.ModePerm)
	}
}

func QuickPost(sitePath string, logType string, title string) {
	var metadata string
	siteCfg := Configure(sitePath + "\\complied\\setting")
	siteCategroys := GetCategorys(siteCfg.GetArray("categorys"))
	siteTags := GetTags(siteCfg.GetArray("tags"))
	var meta_siteCategroys, meta_siteTags string
	for _, category := range siteCategroys {
		if meta_siteCategroys != "" {
			meta_siteCategroys += ","
		}
		meta_siteCategroys += category.Name
	}
	for _, tag := range siteTags {
		if meta_siteTags != "" {
			meta_siteTags += ","
		}
		meta_siteTags += tag.Name
	}

	if len(meta_siteCategroys) > 0 {
		metadata += "category:" + meta_siteCategroys + "\n"
	}
	if len(meta_siteTags) > 0 {
		metadata += "tag:" + meta_siteTags + "\n"
	}
	if len(metadata) > 0 {
		metadata = "---\n" + metadata + "---\n"
	}
	logPath := sitePath + "\\" + title + "@" + time.Now().Format("2006-1-2")
	err := os.Mkdir(logPath, os.ModePerm)
	if err == nil {
		var file string
		if logType == "album" {
			file = logPath + "\\meta"
		} else {
			file = logPath + "\\article.md"
		}
		err = ioutil.WriteFile(file, []byte(metadata), os.ModePerm)
		if err != nil {
			log.Fatal(logPath + "写入元数据失败")
		}
		cmd := exec.Command("explorer.exe", "/select,"+file)
		cmd.Run()
	} else {
		log.Fatal(logPath + "日志创建失败")
	}

}
