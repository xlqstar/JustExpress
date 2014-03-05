package justExpress

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
)

func update_index(indexPage IndexPage, categoryDirPath string) {
	total := len(indexPage.LogList)
	maxPage := int(math.Ceil(float64(total) / float64(int(indexPage.PageSize))))
	indexFileName := "index.html"
	for page := 1; page <= maxPage; page++ {
		if page > 1 {
			indexFileName = "index_" + strconv.Itoa(page) + ".html"
		}

		indexFilePath := categoryDirPath + "\\" + indexFileName
		update_log(indexFilePath, indexPage.RelPath)
	}
}

func update_tag(logFilePath string, relPath string) {
	update_log(logFilePath, relPath)
}

func update_log(logFilePath string, relPath string) {
	fileContent_bytes, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		log.Fatal("读取" + logFilePath + "出错！")
	}
	fileContent := replaceGlobalTlp(string(fileContent_bytes))
	fileContent = replaceRelPath(fileContent, relPath)
	ioutil.WriteFile(logFilePath, []byte(fileContent), os.ModePerm)
}
