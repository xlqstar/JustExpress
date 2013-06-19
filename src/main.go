/*

--buildtarget=log/index/both
--whichlog=from/all/auto

几种日志创建方式：

0：创建未创建的日志,以及有更新的日志
1：按日期创建日志
2：全部创建

*/


package main

import(
	"fmt"
	"path/filepath"
	"os"
	//"reflect"
)

var srcDir = "D:\\goPoject\\src\\*"
var destDir = ""

func init() {

	//--buildtarget=log/index/both
	buildtarget := flag.String("buildtarget", 'both', "choose to build log or index or both")
	
	//--whichlog=from/all/auto
	whichlog := flag.String("whichlog", 'noyet', "choose whichlog to be build")

	if whichlog == "from" {
		fromtime := flag.Args()
	}

	flag.Parse()
}


func main(){

	mode = os.Args[1]
	
	logDirList,_ := filepath.Glob(srcDir)

	var logList [][]

	//过滤
	filter_dir(&logDirList)
	
	//日志生成
	for k := range logDirList{
		
		logInfo := decode_log(logDirList[k])

		//filter log
		if whichlog == "auto" {
		


		}else if whichlog == "from" {



		}else if whichlog == "all"{



		}


		build_log(logInfo)

		append(logList,logInfo)

		rebuild_list()

	}

	//索引处理
	if buildtarget == "index" || buildtarget == "both" {
		
		build_index()
	
	}

	fmt.Println(logDirList)

}






//数组删除
func array_delete(array *[]string, k int){
	
	*array = append((*array)[:k],(*array)[k+1:]...)

}

//过滤非目录文件
func filter_dir(logDirList *[]string) {

	for key := range *logDirList{

		logDirFI,_ := os.Stat((*logDirList)[key])
		
		if !logDirFI.IsDir() {
			
			array_delete(logDirList,key)

		}
	}
}

//解析日志类型
func parse_logType() {
	
}

//解析日志信息
func decode_log(dir) {

	logType := parse_logType(dir)
	
	if logType == 'article' {//文章
		
		logInfo = decode_article()

	}else if logType == 'album' {//相册

		logInfo = decode_album()

	}

}


//生成文章型日志
func build_article() {
	
}

//生成相册型日志
func build_album() {
	
}

//生成索引
func build_index() {
	
}


func rebuild_list() {
	
}