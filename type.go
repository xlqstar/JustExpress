//各种类的定义汇集地

package just

import (
	"just/pinyin"
)

//相册
type Album []map[string]string

//文章
type Article string

//标签
type Category struct {
	Name  string
	Href  string
	Alias string
	//IsCategory bool
}

//标签
type Tag struct {
	Name  string
	Alias string
}

//友情链接
type Link struct {
	Name string
	Href string
}

type LogInfo struct {
	Title       string
	Date        int
	MetaData    map[string]string
	Src         string
	Type        string
	LastModTime int
	// CreateTime  string
	Summary interface{}
	Log     interface{}
}

type SiteInfo struct {
	Categorys []Category
	Links     []Link
	Tags      []Tag
	Email     string
	Author    string
}

type IndexPage struct {
	Category Category
	RelPath  string
	SiteInfo SiteInfo

	Page     int
	LogList  []LogInfo
	NextPage int
	PrevPage int
	PageSize int
}

type LogPage struct {
	RelPath  string
	LogInfo  LogInfo
	SiteInfo SiteInfo
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
	Archives map[string]LogList
}

//索引模版中需要
func (logInfo LogInfo) IsArticle() bool {
	if logInfo.Type == "article" {
		return true
	}
	return false
}

type LogList []LogInfo

func (logList *LogList) Contain(title string) bool {
	for _, v := range *logList {
		if pinyin.Convert(v.Title) == title {
			return true
		}
	}
	return false
}
