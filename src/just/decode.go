package just

import (
	"fmt"
	"io/ioutil"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"os"
	"path/filepath"
	"github.com/axgle/mahonia"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path"
	"regexp"
	"strings"
	"syscall"
)

//解析日志信息
func Decode_log(src string) LogInfo { //decode_log

	logType := Parse_logType(src)

	var logInfo LogInfo

	if logType == "article" { //文章

		logInfo = decode_article(src)

	} else if logType == "album" { //相册

		logInfo = decode_album(src)

	} else {

		log.Fatal("解析日志：未知类型的日志目录，目录结构不符合预期！")

	}

	return logInfo
}

//文章型博客解析
//两种格式：F:\\kuaipan\\blogMaker\\src\\哇哈哈@2013-12-03
//			F:\\kuaipan\\blogMaker\\src\\哇哈哈@2013-12-03.md
func decode_article(src string) LogInfo {

	var articleInfo = LogInfo{MetaData: make(map[string]string)}

	fileName := filepath.Base(src)
	articleInfo.Title = strings.Split(strings.Replace(fileName, ".md", "", -1), "@")[0]
	articleInfo.Date = strings.Split(strings.Replace(fileName, ".md", "", -1), "@")[1]

	fileInfo, _ := os.Stat(src)
	articleInfo.LastModTime = getLastModTime(fileInfo)
	if fileInfo.IsDir() {
		src = src + "\\article.md"
	}

	//获取元数据
	fi, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	content := string(fi)
	//re := regexp.MustCompile("-{3,}([\\S\\s]*?)-{3,}")
	metaDataBlock := GetMetaDataBlock(content)
	//metaDataBlock = strings.Replace(metaDataBlock, "-", "", -1)
	array := strings.Split(metaDataBlock, "\n")
	for _, metaDataLine := range array {
		if strings.Contains(metaDataLine, "---") {
			continue
		}
		articleInfo.MetaData[strings.Split(metaDataLine, ":")[0]] = strings.Split(metaDataLine, ":")[1]
	}

	//获取summary
	content = GetContent(content, metaDataBlock)

	index := strings.Index(content, "\n<!--more-->")
	summary := content[0:index]
	articleInfo.Summary = summary
	articleInfo.Src = src
	articleInfo.Log = Article(content)
	articleInfo.Type = "article"
	return articleInfo
}

func decode_album(src string) LogInfo {
	var albumInfo = LogInfo{Log: Album{}}
	fileInfo, _ := os.Stat(src)
	albumInfo.LastModTime = getLastModTime(fileInfo)
	fileName := filepath.Base(src)
	albumInfo.Title = strings.Split(strings.Replace(fileName, ".md", "", -1), "@")[0]
	albumInfo.Date = strings.Split(strings.Replace(fileName, ".md", "", -1), "@")[1]
	photoList, _ := filepath.Glob(src + "\\*")
	decoder := mahonia.NewDecoder("UTF-16LE")
	var comment_str string
	var imgInfo image.Config
	for key := range photoList {
		fullFileName := photoList[key]
		photo_fi, err := os.Open(fullFileName)
		if err != nil {
			log.Fatal(err)
		}
		photo_exif, err := exif.Decode(photo_fi)
		if err != nil {
			comment_str = ""
		}else{
			comment, err := photo_exif.Get(exif.UserComment)
			if err != nil {
				comment_str = ""
			}else{
				comment_str = decoder.CConvertString(comment.Val)
			}
		}
		photo_fi.Close()
		photo_fi, err = os.Open(fullFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer photo_fi.Close()
		switch strings.ToLower(path.Ext(fullFileName)) {
		case ".jpg", ".jpeg":
			imgInfo, err = jpeg.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误："+fullFileName)
				os.Exit(0)
			}

		case ".png":
			imgInfo, err = png.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误："+fullFileName)
				os.Exit(0)
			}

		case ".gif":
			imgInfo, err = gif.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误："+fullFileName)
				os.Exit(0)
			}
		}
		//albumInfo.Log = append(albumInfo.Log.(Album), map[string]string{"src": photoList[key], "comment": comment_str, "width": fmt.Sprintf("%d", imgInfo.Width), "height": fmt.Sprintf("%d", imgInfo.Height)})
		albumInfo.Log.(Album)[filepath.Base(fullFileName)] = map[string]string{"src": fullFileName, "comment": comment_str, "width": fmt.Sprintf("%d", imgInfo.Width), "height": fmt.Sprintf("%d", imgInfo.Height)}
	}
	albumInfo.Src = src
	albumInfo.Type = "album"
	return albumInfo
}

func getLastModTime(fi os.FileInfo) string {
	ModTime := fi.Sys().(*syscall.Win32FileAttributeData).LastWriteTime.Nanoseconds() / (1000*1000*1000)
	ModTimeStr := fmt.Sprintf("%d", ModTime)
	return ModTimeStr
}

func GetMetaDataBlock(content string) string {
	re := regexp.MustCompile("-{3,}([\\S\\s]*?)-{3,}")
	metaDataBlock := re.FindString(content)
	return metaDataBlock
}

func GetContent(content string, metaDataBlock string) string {
	content = strings.Replace(content, metaDataBlock, "", -1)
	return content
}
