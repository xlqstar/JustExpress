//日志生成函数集

package just

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wendal/gor"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//生成日志列表数据
func build_loglistdata(destDirPath string, loglist []LogInfo) {
	logList_json, _ := json.MarshalIndent(loglist, "", "\t")
	ioutil.WriteFile(destDirPath+"\\loglist.json", logList_json, os.ModePerm)
}

//生成索引
func Build_index(indexPage IndexPage, tplDirPath string, destDirPath string) {

	//如果没有内容则直接生成
	if len(indexPage.LogList) == 0 {
		indexPage.Page = Page(1)
		build_index_page(indexPage, tplDirPath, destDirPath)
		return
	}

	//索引处理
	var page = 1
	var _logList = []LogInfo{}
	var totalCount = len(indexPage.LogList) //总数
	var totalPage = int(math.Ceil(float64(totalCount) / float64(indexPage.PageSize)))

	indexPage.TotalPage = Page(totalPage)
	for _, v := range indexPage.LogList {
		_logList = append(_logList, v)

		if len(_logList) == indexPage.PageSize || len(_logList) == totalCount {
			indexPage.LogList = _logList
			indexPage.Page = Page(page)
			if page == totalPage {
				indexPage.NextPage = Page(1)
			} else {
				indexPage.NextPage = Page(page + 1)
			}

			if page == 1 {
				indexPage.PrevPage = Page(totalPage)
			} else {
				indexPage.PrevPage = Page(page - 1)
			}

			build_index_page(indexPage, tplDirPath, destDirPath)

			page++
			_logList = []LogInfo{}
		}
	}
}

func build_index_page(indexPage IndexPage, tplDirPath string, destDirPath string) {
	for k, v := range indexPage.LogList {

		if v.Summary == nil || v.Summary == "" {
			if v.Type == "article" {
				content, _, _ := _decode_article(v.Src)
				indexPage.LogList[k].Summary = gor.MarkdownToHtml(string(content))
			} else if v.Type == "album" {
				indexPage.LogList[k].Summary, _ = _decode_album(v.Src)
			}
		} else {
			if v.Type == "article" {
				indexPage.LogList[k].Summary = gor.MarkdownToHtml(string(v.Summary.(Article)))
			}
		}
	}

	page := ""
	if indexPage.Page > 1 {
		page = "_" + strconv.Itoa(int(indexPage.Page))
	}

	dest_filename := destDirPath + "\\index" + page + ".html"
	makeHTML(&indexPage, dest_filename, tplDirPath+"\\index.html")
	log.Println(ConvertPath(dest_filename) + " 成功生成！")
}

func Build_log(logPage LogPage, tplDirPath string, destDirPath string, onlyRebuildHtml bool) {
	var destLogDir string
	destLogDir = destDirPath + "\\" + logPage.LogInfo.Permalink

	if !onlyRebuildHtml {
		err := os.Mkdir(destLogDir, os.ModePerm)
		if os.IsExist(err) {
			os.RemoveAll(destLogDir)
			err = os.Mkdir(destLogDir, os.ModePerm)
			if err != nil {
				log.Fatal("重建" + destLogDir + "目录发生未预料到的错误，可能正有其他进程占用该目录。")
			}
		}
	}

	if logPage.LogInfo.Type == "article" { //文章
		build_article(logPage, tplDirPath, destLogDir, onlyRebuildHtml)
	} else if logPage.LogInfo.Type == "album" { //相册
		build_album(logPage, tplDirPath, destLogDir, onlyRebuildHtml)
	}

	log.Println("《" + logPage.LogInfo.Title + "》成功生成！")
}

func build_album(logPage LogPage, tplDirPath string, destLogDir string, onlyRebuildHtml bool) {
	destAlbumDir := destLogDir
	smallImgWidth := logPage.SiteInfo.ImgWidth
	bigImgWidth := logPage.SiteInfo.OriginImgWidth

	logPage.LogInfo.Log, _ = _decode_album(logPage.LogInfo.Src)
	for k, photo := range logPage.LogInfo.Log.(Album) {
		srcPhotoFullFileName := photo.Src
		destPhotoFullFileName := strings.Replace(srcPhotoFullFileName, logPage.LogInfo.Src, destAlbumDir, -1)
		photoFileName := photo.PhotoFileName
		photoWidth := photo.Width
		var originPhotoFullFileName string
		if strings.ToLower(path.Ext(photoFileName)) == ".gif" || photoWidth < smallImgWidth {
			originPhotoFullFileName = destPhotoFullFileName
			if !onlyRebuildHtml {
				CopyFile(srcPhotoFullFileName, destPhotoFullFileName)
			}
		} else if photoWidth > bigImgWidth {
			originPhotoFullFileName = strings.Replace(destPhotoFullFileName, photoFileName, "origin_"+photoFileName, -1)
			if !onlyRebuildHtml {
				Resize(srcPhotoFullFileName, destPhotoFullFileName, uint(smallImgWidth))
				if siteInfo.OriginImgWidth > 0 {
					Resize(srcPhotoFullFileName, originPhotoFullFileName, uint(bigImgWidth))
				}
			}
		} else if photoWidth > smallImgWidth {
			originPhotoFullFileName = strings.Replace(destPhotoFullFileName, photoFileName, "origin_"+photoFileName, -1)
			if !onlyRebuildHtml {
				Resize(srcPhotoFullFileName, destPhotoFullFileName, uint(smallImgWidth))
				CopyFile(srcPhotoFullFileName, originPhotoFullFileName)
			}
		}

		photo.OriginPhotoFileName = filepath.Base(originPhotoFullFileName)
		logPage.LogInfo.Log.(Album)[k] = photo
	}
	makeHTML(&logPage, destAlbumDir+"\\index.html", tplDirPath+"\\album.html")
}

//global var : onlyRebuildHtml
func build_article(logPage LogPage, tplDirPath string, destLogDir string, onlyRebuildHtml bool) {
	destArticleDir := destLogDir
	if !onlyRebuildHtml {
		if filepath.Base(logPage.LogInfo.Src) == "article.md" {
			srcArticleDir := filepath.Dir(logPage.LogInfo.Src)
			CopyDir(srcArticleDir, destArticleDir)
		} else {
			if _, err := os.Stat(destArticleDir); os.IsNotExist(err) {
				os.Mkdir(destArticleDir, os.ModePerm)
			}
			CopyFile(logPage.LogInfo.Src, destArticleDir+"\\article.md")
		}
	}
	content, _, _ := _decode_article(logPage.LogInfo.Src)
	logPage.LogInfo.Log = gor.MarkdownToHtml(string(content))
	makeHTML(&logPage, destArticleDir+"\\index.html", tplDirPath+"\\article.html")
}

func Build_tagpage(tagPage TagPage, tplDirPath string, destTagDir string) {
	makeHTML(&tagPage, destTagDir+"\\"+tagPage.Tag.Alias+".html", tplDirPath+"\\tag.html")
	log.Println(ConvertPath(destTagDir+"\\"+tagPage.Tag.Alias+".html") + " 标签页生成")
}

func Build_archive(archivePage ArchivePage, tplDirPath string, destArchiveDir string) {
	makeHTML(&archivePage, destArchiveDir+"\\index.html", tplDirPath+"\\archive.html")
	log.Println(ConvertPath(destArchiveDir+"\\index.html") + " 归档页生成")
}

/*#先渲染负载小 后replace负载大
renderTpl
replaceGlobalTlp
replaceRelPath

#先替换global渲染负载大 后replace，负载大
replaceGlobalTlp
renderTpl
replaceRelPath

#这种先替换relPath负载小 先替换global渲染负载大
replaceGlobalTlp
replaceRelPath
renderTpl*/
func makeHTML(data interface{}, dest string, templatePath string) {

	template_byte, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatal(templatePath + "模版文件读取失败，请检查该模版是否存在？")
	}

	relPath := reflect.ValueOf(data).Elem().FieldByName("RelPath").String()

	html := replaceGlobalTlp(string(template_byte))
	html = replaceRelPath(html, relPath)
	html_bytes := renderTpl(data, html, strings.TrimSuffix(filepath.Base(templatePath), filepath.Ext(templatePath)))

	ioutil.WriteFile(dest, html_bytes, os.ModePerm)
}

func renderTpl(data interface{}, template_str string, tplName string) []byte {
	out := bytes.NewBuffer([]byte{})
	tpl := template.New(tplName)
	tpl.Funcs(template.FuncMap{"first": first, "last": last, "eq": eq, "neq": neq, "not": not, "add": add, "minus": minus})
	tpl.Parse(template_str)
	if err := tpl.Execute(out, data); err != nil {
		fmt.Println(err)
	}
	return out.Bytes()
}

func replaceRelPath(content string, relPath string) string {
	//主题路径处理
	content = strings.Replace(content, "././js/", relPath+"style/js/", -1)
	content = strings.Replace(content, "././css/", relPath+"style/css/", -1)
	content = strings.Replace(content, "././images/", relPath+"style/images/", -1)

	content = strings.Replace(content, "././", relPath, -1)
	return content
}

//全局变量: siteInfo
//替换全局布局模版
func replaceGlobalTlp(content string) string {
	for tplName, tpl := range siteInfo.GlobalTpl {
		reg, _ := regexp.Compile(`<!--` + tplName + `-->[\s\S]*<!--//` + tplName + `-->`)
		content = reg.ReplaceAllString(content, "<!--"+tplName+"-->"+tpl+"<!--//"+tplName+"-->")
	}
	return content
}

/*
func getRelPath(subDirPath string) string {
	subDirPath = fixPath(subDirPath)
	if len(subDirPath) == 0 {
		return "./"
	}
	num := strings.Count(subDirPath, "\\")
	relPath := ""
	for i := 0; i < num+1; i++ {
		relPath = relPath + "../"
	}
	return relPath
}

func fixPath(path string) string {
	path = Trim(path)
	path = strings.Replace(path, "/", "\\", -1)
	path = strings.TrimPrefix(path, "\\")
	path = strings.TrimSuffix(path, "\\")
	return path
}*/
