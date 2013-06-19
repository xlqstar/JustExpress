package just

import (
	"os"
	"path/filepath"
	"strings"
)

func Parse_logType(dir string) string {

	file, _ := os.Open(dir)
	fi, _ := file.Stat()
	if !fi.IsDir() {
		return "article"
	}

	srcDir := dir + "\\*"
	fileList, _ := filepath.Glob(srcDir)

	if is_album(fileList) {
		return "album"
		// fmt.Println("is album")
	} else if is_article(fileList) {
		return "article"
		//fmt.Println("is article")
	} else {
		return "unknow logtype"
		// fmt.Println("unknow logtype")
	}

}

//是否是相册型
func is_album(fileList []string) bool {
	for k := range fileList {
		ext := filepath.Ext(fileList[k])
		ext = strings.ToLower(ext)
		imgExtArray := [...]string{".jpg", ".jpeg", ".gif", ".bmp", ".png"}
		if !in_array(strings.ToLower(ext), imgExtArray[0:]) {
			return false
		}
	}
	return true
}

//是否是文章型
func is_article(fileList []string) bool {
	for k := range fileList {
		ext := filepath.Ext(fileList[k])
		ext = strings.ToLower(ext)
		imgExtArray := [...]string{".md", ".markdown", ".html", ".htm"}
		if in_array(strings.ToLower(ext), imgExtArray[0:]) {
			return true
		}
	}
	return false
}

func in_array(v string, array []string) bool {
	for k := range array {
		if array[k] == v {
			return true
		}
	}
	return false
}

/*
func InArray(obj interface{}, slice interface{}) (bool, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(slice))
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	if sliceValue.Kind() != reflect.Slice {
		return false, errors.New("expected a slice")
	}

	if sliceValue.Len() < 1 {
		return false, nil
	}

	for i := 0; i < sliceValue.Len(); i++ {
		if sliceValue.Index(i).Interface() == objValue.Interface() {
			return true, nil
		}
		fmt.Println(sliceValue.Index(i).Interface(), objValue.Interface())
	}
	return false, nil
}
*/
