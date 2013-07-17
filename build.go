//日志生成函数集

package just

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wendal/gor"
	"io/ioutil"
	"just/pinyin"
	"log"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

//生成日志列表数据
func build_loglistdata(destDirPath string, loglist []LogInfo) {
	logList_json, _ := json.Marshal(loglist)
	ioutil.WriteFile(destDirPath+"\\loglist.json", logList_json, os.ModePerm)
}

//生成索引
func Build_index(indexPage IndexPage, tplDirPath string, destDirPath string) {

	pagesize := indexPage.PageSize
	// fmt.Println(indexPage.LogList)
	//索引处理
	indexPage.LogList = LogSort(indexPage.LogList)
	var page = 1
	var _logList = []LogInfo{}
	var totalCount = len(indexPage.LogList) //总数
	var totalPage = int(math.Ceil(float64(totalCount) / float64(pagesize)))
	for _, v := range indexPage.LogList {
		_logList = append(_logList, v)

		if len(_logList) == pagesize || len(_logList) == totalCount {
			indexPage.LogList = _logList
			indexPage.Page = page
			if page == totalPage {
				indexPage.NextPage = 0
			} else {
				indexPage.NextPage = page + 1
			}

			if page == 1 {
				indexPage.PrevPage = 0
			} else {
				indexPage.PrevPage = page - 1
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
		page = "_" + strconv.Itoa(indexPage.Page)
	}
	relDirPath := ""
	dest_filename := ""
	if indexPage.Category.Alias != "index" {
		relDirPath = "../"
	}
	dest_filename = destDirPath + "\\index" + page + ".html"
	indexPage.RelPath = relDirPath
	makeHTML(&indexPage, dest_filename, tplDirPath+"\\index.html")
	log.Println(dest_filename + " 成功生成！")
}

func Build_log(logPage LogPage, tplDirPath string, destDirPath string, smallImgWidth uint, bigImgWidth uint) {
	var destLogDir string
	if logPage.LogInfo.MetaData["alias"] == "" {
		destLogDir = destDirPath + "\\" + url.QueryEscape(pinyin.Convert(logPage.LogInfo.Title))
	} else {
		destLogDir = destDirPath + "\\" + url.QueryEscape(logPage.LogInfo.MetaData["alias"])
	}

	err := os.Mkdir(destLogDir, os.ModePerm)
	if os.IsExist(err) {
		os.RemoveAll(destLogDir)
		err = os.Mkdir(destLogDir, os.ModePerm)
		if err != nil {
			log.Fatal("重建" + destLogDir + "目录发生未预料到的错误")
		}
	}

	if logPage.LogInfo.Type == "article" { //文章

		build_article(logPage, tplDirPath, destLogDir)

	} else if logPage.LogInfo.Type == "album" { //相册

		build_album(logPage, tplDirPath, destLogDir, smallImgWidth, bigImgWidth)

	}

	log.Println("《" + logPage.LogInfo.Title + "》成功生成！")

}

func build_album(logPage LogPage, tplDirPath string, destLogDir string, smallImgWidth uint, bigImgWidth uint) {
	destAlbumDir := destLogDir
	if logPage.LogInfo.Log == nil {
		log.Println("logPage.LogInfo.Log.(Album)为nil")
		return
	}

	for _, v := range logPage.LogInfo.Log.(Album) {
		srcPhotoFullFileName := v["src"]
		destPhotoFullFileName := strings.Replace(srcPhotoFullFileName, logPage.LogInfo.Src, destAlbumDir, -1)
		photoFileName := filepath.Base(srcPhotoFullFileName)
		photoWidth, _ := strconv.Atoi(v["width"])
		smallPhotoFileName, bigPhotoFileName := "", ""
		if strings.HasPrefix(photoFileName, "~") || strings.ToLower(path.Ext(photoFileName)) == ".gif" {
			smallPhotoFileName, bigPhotoFileName = destPhotoFullFileName, destPhotoFullFileName
			CopyFile(srcPhotoFullFileName, destPhotoFullFileName)
		} else if uint(photoWidth) > bigImgWidth {
			smallPhotoFileName, bigPhotoFileName = strings.Replace(destPhotoFullFileName, photoFileName, "small_"+photoFileName, -1), strings.Replace(destPhotoFullFileName, photoFileName, "big_"+photoFileName, -1)
			Resize(srcPhotoFullFileName, smallPhotoFileName, smallImgWidth)
			Resize(srcPhotoFullFileName, bigPhotoFileName, bigImgWidth)
		} else if uint(photoWidth) > smallImgWidth {
			smallPhotoFileName, bigPhotoFileName = strings.Replace(destPhotoFullFileName, photoFileName, "small_"+photoFileName, -1), destPhotoFullFileName
			Resize(srcPhotoFullFileName, smallPhotoFileName, smallImgWidth)
		} else {
			smallPhotoFileName, bigPhotoFileName = destPhotoFullFileName, destPhotoFullFileName
			CopyFile(srcPhotoFullFileName, destPhotoFullFileName)
		}
		v["smallPhotoFileName"] = filepath.Base(smallPhotoFileName)
		v["bigPhotoFileName"] = filepath.Base(bigPhotoFileName)
	}
	logPage.LogInfo.Log, _ = _decode_album(logPage.LogInfo.Src)
	makeHTML(&logPage, destAlbumDir+"\\index.html", tplDirPath+"\\album.html")
}

func build_article(logPage LogPage, tplDirPath string, destLogDir string) {
	destArticleDir := destLogDir
	if filepath.Base(logPage.LogInfo.Src) == "article.md" {
		srcArticleDir := filepath.Dir(logPage.LogInfo.Src)
		CopyDir(srcArticleDir, destArticleDir)
	} else {
		if _, err := os.Stat(destArticleDir); os.IsNotExist(err) {
			os.Mkdir(destArticleDir, os.ModePerm)
		}
		CopyFile(logPage.LogInfo.Src, destArticleDir+"\\article.md")
	}
	content, _, _ := _decode_article(logPage.LogInfo.Src)
	logPage.LogInfo.Log = gor.MarkdownToHtml(string(content))
	makeHTML(&logPage, destArticleDir+"\\index.html", tplDirPath+"\\article.html")
}

func Build_tagpage(tagPage TagPage, tplDirPath string, destDirPath string) {
	destTagDir := destDirPath + "\\tag"
	makeHTML(&tagPage, destTagDir+"\\"+tagPage.Tag.Alias+".html", tplDirPath+"\\tag.html")
}

func Build_archive(archivePage ArchivePage, tplDirPath string, destDirPath string) {
	destArchiveDir := destDirPath + "\\archive"
	makeHTML(&archivePage, destArchiveDir+"\\index.html", tplDirPath+"\\archive.html")
}

func makeHTML(data interface{}, dest string, templatePath string) {

	template_byte, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatal(templatePath + "模版文件读取失败，请检查该模版是否存在？")
	}
	template_str := strings.Replace(string(template_byte), "././././", reflect.ValueOf(data).Elem().FieldByName("RelPath").String(), -1)
	out := bytes.NewBuffer([]byte{})
	tplName := strings.Replace(filepath.Base(templatePath), path.Ext(templatePath), "", -1)
	t, _ := template.New(tplName).Parse(template_str)
	if err := t.Execute(out, data); err != nil {
		fmt.Println(err)
	}
	html := out.Bytes()
	ioutil.WriteFile(dest, html, os.ModePerm)
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
