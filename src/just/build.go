package just

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wendal/gor"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

//更新loglist历史记录数据
func Update_loglistdata(loglist map[string]LogInfo) {
	logList_json, _ := json.Marshal(loglist)
	ioutil.WriteFile("loglist.json", logList_json, os.ModePerm)
}

//生成索引
func Build_index(logList map[string]LogInfo, tplDirPath string, destDirPath string, pagesize int) {
	//索引处理
	logList = LogSort(logList)
	var page = 1
	var pageloglist = map[string]LogInfo{}
	var listPage ListPage
	var totalCount = len(logList) //总数
	var totalPage = totalCount / pagesize
	for k, v := range logList {

		pageloglist[k] = v

		if len(pageloglist) == pagesize || len(pageloglist) == totalCount {
			listPage.Loglist = pageloglist
			listPage.Page = page
			if page == totalPage {
				listPage.NextPage = 0
			} else {
				listPage.NextPage = page + 1
			}

			if page == 1 {
				listPage.PrevPage = 0
			} else {
				listPage.PrevPage = page - 1
			}

			build_index_page(listPage, tplDirPath, destDirPath)

			page++
			pageloglist = map[string]LogInfo{}
		}
	}
}

func build_index_page(listPage ListPage, tplDirPath string, destDirPath string) {
	page := ""
	if listPage.Page > 1 {
		page = "_" + strconv.Itoa(listPage.Page)
	}
	dest_filename := destDirPath + "\\index" + page + ".html"
	makeHTML(dest_filename, listPage, tplDirPath+"\\index.html")
	log.Println(filepath.Base(dest_filename) + " 成功生成！")
}

func Build_log(logInfo LogInfo, tplDirPath string, destDirPath string, smallImgWidth uint, bigImgWidth uint) {

	if logInfo.Type == "article" { //文章

		build_article(logInfo, tplDirPath, destDirPath)

	} else if logInfo.Type == "album" { //相册

		build_album(logInfo, tplDirPath, destDirPath, smallImgWidth, bigImgWidth)

	}

	log.Println("《" + logInfo.Title + "》成功生成！")

}

func build_album(logInfo LogInfo, tplDirPath string, destDirPath string, smallImgWidth uint, bigImgWidth uint) {
	destAlbumDir := destDirPath + "\\" + logInfo.Title
	if _, err := os.Stat(destAlbumDir); os.IsNotExist(err) {
		os.Mkdir(destAlbumDir, os.ModePerm)
	}
	photoList, _ := filepath.Glob(logInfo.Src + "\\*")
	for key := range photoList {
		srcPhotoFullFileName := photoList[key]
		destPhotoFullFileName := strings.Replace(srcPhotoFullFileName, logInfo.Src, destAlbumDir, -1)
		photoFileName := filepath.Base(srcPhotoFullFileName)
		photoWidth, _ := strconv.Atoi(logInfo.Log.(Album)[photoFileName]["width"])
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
		logInfo.Log.(Album)[photoFileName]["smallPhotoFileName"] = filepath.Base(smallPhotoFileName)
		logInfo.Log.(Album)[photoFileName]["bigPhotoFileName"] = filepath.Base(bigPhotoFileName)
	}

	makeHTML(destAlbumDir+"\\index.html", logInfo, tplDirPath+"\\album.html")
}

func build_article(logInfo LogInfo, tplDirPath string, destDirPath string) {
	destArticleDir := destDirPath + "\\" + logInfo.Title
	if filepath.Base(logInfo.Src) == "article.md" {
		srcArticleDir := filepath.Dir(logInfo.Src)
		CopyDir(srcArticleDir, destArticleDir)
	} else {
		if _, err := os.Stat(destArticleDir); os.IsNotExist(err) {
			os.Mkdir(destArticleDir, os.ModePerm)
		}
		CopyFile(logInfo.Src, destArticleDir+"\\article.md")
	}
	logInfo.Log = Article(gor.MarkdownToHtml(string(logInfo.Log.(Article))))
	makeHTML(destArticleDir+"\\index.html", logInfo, tplDirPath+"\\article.html")
}

func makeHTML(dest string, data interface{}, templatePath string) {
	template_byte, _ := ioutil.ReadFile(templatePath)
	out := bytes.NewBuffer([]byte{})
	tplName := strings.Replace(filepath.Base(templatePath), path.Ext(templatePath), "", -1)
	t, _ := template.New(tplName).Parse(string(template_byte))
	if err := t.Execute(out, data); err != nil {
		fmt.Println(err)
	}
	html := out.Bytes()
	ioutil.WriteFile(dest, html, os.ModePerm)
}
