//日志数据解析函数集

package justExpress

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/xlqstar/JustExpress/pinyin"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	// "regexp"
	"strings"
	"syscall"
	"time"
)

//解析日志信息
func Decode_log(src string, logType string, siteInfo SiteInfo) LogInfo { //decode_log

	var logInfo LogInfo

	if logType == "article" { //文章

		logInfo = decode_article(src)

	} else if logType == "album" { //相册

		logInfo = decode_album(src)

	} else {

		log.Fatal("解析日志：未知类型的日志目录，目录结构不符合预期！")

	}
	//读取标题、日期等信息
	fileName := filepath.Base(src)
	pos := strings.LastIndex(fileName, "@")
	if pos >= 0 {
		logInfo.Title = fileName[0:pos]
		_date := strings.TrimSuffix(fileName[pos+1:len(fileName)], ".md")
		// t, err := time.Parse("2006-5-5", logInfo.Date)
		loc, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation("2006-1-2", _date, loc)
		if err != nil {
			log.Println(err)
			log.Fatal("日志日期填写格式错误，请参照2006-1-2这样的格式，无前导!")
		}
		logInfo.Date = TimeStamp(t.Unix())
	} else {
		log.Fatal(fileName + "日志文件名不合法，请修改后再build")
	}

	//if logInfo.MetaData["alias"] == "" {
	logInfo.Permalink = url.QueryEscape(pinyin.Convert(logInfo.Title, "-"))
	//} else {
	//	logInfo.Permalink = logInfo.MetaData["alias"]
	//}
	/*
		for _, category := range siteInfo.Categorys {
			if strings.Contains(logInfo.MetaData["category"], category.Name) {
				logInfo.Categorys = append(logInfo.Categorys, category)
			}
		}*/
	/*	for _, tag := range siteInfo.Tags {
		if strings.Contains(logInfo.MetaData["tag"], tag.Name) {
			logInfo.Tags = append(logInfo.Tags, tag)
		}
	}*/
	return logInfo
}

//文章型博客解析
//两种格式：F:\\kuaipan\\blogMaker\\src\\哇哈哈@2013-12-03
//			F:\\kuaipan\\blogMaker\\src\\哇哈哈@2013-12-03.md
func decode_article(src string) LogInfo {

	var articleInfo LogInfo

	//获取最新修改时间
	fileInfo, _ := os.Stat(src)
	articleInfo.LastModTime = TimeStamp(getLastModTime(fileInfo))
	if fileInfo.IsDir() {
		src = src + "\\article.md"
	}
	// _, articleInfo.MetaData = _decode_article(src)
	articleInfo.Src = src
	articleInfo.Type = "article"
	return articleInfo
}

func decode_album(src string) LogInfo {
	var albumInfo LogInfo
	fileInfo, _ := os.Stat(src)
	albumInfo.LastModTime = TimeStamp(getLastModTime(fileInfo))

	//获取元数据
	/*	file, err := ioutil.ReadFile(src + "\\meta")
		if err == nil {
			albumInfo.MetaData, _ = decode_meta(string(file))
		}*/
	_, summary := _decode_album(src)
	if len(summary) > 0 {
		albumInfo.Summary = summary
	}
	albumInfo.Src = src
	albumInfo.Type = "album"
	return albumInfo
}

func getLastModTime(fi os.FileInfo) int64 {
	ModTime := fi.Sys().(*syscall.Win32FileAttributeData).LastWriteTime.Nanoseconds() / (1000 * 1000 * 1000)
	// ModTimeStr := fmt.Sprintf("%d", ModTime)
	return ModTime
}

func GetCreationTime(fi os.FileInfo) int64 {
	CreationTime := fi.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds() / (1000 * 1000 * 1000)
	// CreationTimeStr := fmt.Sprintf("%d", CreationTime)
	return CreationTime
}

func _decode_album(src string) (Album, Album) {
	photoList, _ := filepath.Glob(src + "\\*")
	var album Album
	var albumSummary Album
	var imgInfo image.Config
	var comment_str string
	decoder := mahonia.NewDecoder("UTF-16LE")

	for key := range photoList {
		fullFileName := photoList[key]

		//读取图片长宽等信息
		photo_fi, err := os.Open(fullFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer photo_fi.Close()
		switch strings.ToLower(path.Ext(fullFileName)) {
		case ".jpg", ".jpeg":
			imgInfo, err = jpeg.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误：" + fullFileName)
				os.Exit(0)
			}

		case ".png":
			imgInfo, err = png.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误：" + fullFileName)
				os.Exit(0)
			}

		case ".gif":
			imgInfo, err = gif.DecodeConfig(photo_fi)
			if err != nil {
				fmt.Println("图片解析错误：" + fullFileName)
				os.Exit(0)
			}
		default:
			continue
		}

		//读取图片评注信息
		photo_fi, err = os.Open(fullFileName)
		if err != nil {
			log.Fatal(err)
		}
		photo_exif, err := exif.Decode(photo_fi)
		if err != nil {
			comment_str = ""
		} else {
			comment, err := photo_exif.Get(exif.UserComment)
			if err != nil {
				comment_str = ""
			} else {
				comment_str = decoder.CConvertString(comment.Val)
			}
		}
		photo_fi.Close()

		var bigPhotoFileName, photoFileName string
		photoFileName = filepath.Base(fullFileName)
		if strings.ToLower(path.Ext(fullFileName)) == ".gif" || imgInfo.Width < siteInfo.ImgWidth { //global var : SiteInfo
			bigPhotoFileName = photoFileName
		} else {
			bigPhotoFileName = "big_" + photoFileName
		}

		//追加图片信息
		var photo = Photo{Src: fullFileName, Comment: comment_str, Width: imgInfo.Width, Height: imgInfo.Height, PhotoFileName: photoFileName, BigPhotoFileName: bigPhotoFileName}
		album = append(album, photo)
		if strings.HasPrefix(filepath.Base(fullFileName), "~") {
			albumSummary = append(albumSummary, photo)
		}

	}
	return album, albumSummary
}

func _decode_article(src string) Article /*, map[string]string*/ {
	//获取元数据
	fi, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	// metaData, content := decode_meta(string(fi))
	// metaData, content := decode_meta(string(fi))
	return Article(fi) /*, metaData*/
}

/*func decode_meta(content string) (map[string]string, string) {
	// re := regexp.MustCompile("-{3,}([\\S\\s]*?)-{3,}")
	metaDataMap := map[string]string{}
	re := regexp.MustCompile("^(\\s*)-{3,}([\\S\\s]*?)-{3,}")
	metaDataBlock := re.FindString(content)
	if metaDataBlock != "" {
		re := regexp.MustCompile("^" + metaDataBlock)
		content = re.ReplaceAllString(content, "")
		metaDataBlock = strings.Replace(metaDataBlock, "\r\n", "\n", -1)
		metaDataArry := strings.Split(metaDataBlock, "\n")
		for _, metaDataLine := range metaDataArry {
			// log.Println(metaDataLine)
			if strings.Contains(metaDataLine, "---") || !strings.Contains(metaDataLine, ":") || strings.TrimSpace(metaDataLine) == "" {
				continue
			}
			metaDataLineArry := strings.SplitN(metaDataLine, ":", 2)
			metaDataMap[strings.ToLower(strings.TrimSpace(metaDataLineArry[0]))] = strings.TrimSpace(metaDataLineArry[1])
		}
	}

	return metaDataMap, content
}
*/
