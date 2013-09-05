//各种类的定义汇集地

package just

import (
	"strconv"
	"time"
)

//图片
type Photo struct {
	Src              string
	PhotoFileName    string
	BigPhotoFileName string
	Comment          string
	Width            int
	Height           int
}

//相册
type Album []Photo

//文章
type Article string

//分类
type Category struct {
	Name  string
	Href  string
	Alias string
	Count int //关联到该标签的文章数
	//IsCategory bool
}

//标签
type Tag struct {
	Name  string
	Alias string
	Count int //关联到该标签的文章数
}

//友情链接
type Link struct {
	Name string
	Href string
}

//存档
type Archive struct {
	YearMonth string
	LogList   LogList
	Count     int //关联到该标签的文章数
}

/*type CategoryStatis struct {
	Name  string
	Alias string
	Href  string
	Count int
}

type TagStatis struct {
	Name  string
	Alias string
	Count int
}
*/
/*//统计对象
type ArchiveStatis struct {
	YearMonth string
	Count     int
}*/

type LogInfo struct {
	Title       string
	Date        TimeStamp
	Tags        []Tag
	Categorys   []Category
	Permalink   string
	MetaData    map[string]string
	Src         string
	Type        string
	LastModTime TimeStamp
	// CreateTime  string
	Summary interface{}
	Log     interface{}
}

type SiteInfo struct {
	Site      string
	SitePath  string
	Domain    string
	Categorys []Category
	Links     []Link
	Tags      []Tag
	Archives  []Archive
	Email     string
	Author    string
	Socials   map[string]string
	PageSize  int

	ImgWidth    int
	BigImgWidth int
	ThemeName   string

	GlobalTpl map[string]string
}

type IndexPage struct {
	Category Category
	RelPath  string
	SiteInfo SiteInfo

	Page      Page
	LogList   []LogInfo
	NextPage  Page
	PrevPage  Page
	TotalPage Page
	PageSize  int
}

type LogPage struct {
	RelPath  string
	LogInfo  LogInfo
	SiteInfo SiteInfo
	NextLog  LogInfo
	PrevLog  LogInfo
}

type TagPage struct {
	RelPath  string
	SiteInfo SiteInfo
	LogList  []LogInfo
	Tag      Tag
}

type ArchivePage struct {
	RelPath  string
	SiteInfo SiteInfo
	// LogList   []LogInfo
	Archives []Archive
}

func (photo Photo) HasComment() bool {
	if len(photo.Comment) > 0 {
		return true
	}
	return false
}

//索引模版中需要
func (logInfo LogInfo) IsArticle() bool {
	if logInfo.Type == "article" {
		return true
	}
	return false
}

type LogList []LogInfo

func (logList *LogList) Contain(permalink string) bool {
	for _, v := range *logList {
		if v.Permalink == permalink {
			return true
		}
	}
	return false
}

type TimeStamp int64

func (timeStamp TimeStamp) Format(layout string) string {
	t := time.Unix(int64(timeStamp), 0)
	return t.Format(layout)
}

/*func (timeStamp TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format(time.RFC3339))), nil
}
*/

//页码
type Page int

func (page Page) PageName() string {
	pageInt := int(page)
	if pageInt == 1 {
		return "index"
	}
	return "index_" + strconv.Itoa(pageInt)
}

func (page Page) PageList() []Page {
	var pageList []Page
	pageInt := int(page)
	for i := 1; i <= pageInt; i++ {
		pageList = append(pageList, Page(i))
	}
	return pageList
}
