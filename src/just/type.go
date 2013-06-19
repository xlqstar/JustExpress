package just

type ListPage struct {
	Page     int
	Loglist  map[string]LogInfo
	NextPage int
	PrevPage int
}

type Album map[string]map[string]string

type Article string

type LogInfo struct {
	Title         string
	Date          string
	MetaData      map[string]string
	Src           string
	Type          string
	LastModTime   string
	LastBuildTime string
	CreateTime    string
	Summary       string
	Log           interface{}
}
