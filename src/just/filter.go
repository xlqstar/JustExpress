package just
import(
	"os"
)

//数组删除
func array_delete(array *[]string, k int) {

	*array = append((*array)[:k], (*array)[k+1:]...)

}

//过滤非目录文件
func filter_dir(logDirList *[]string) {
	for key := range *logDirList {

		logDirFI, _ := os.Stat((*logDirList)[key])

		if !logDirFI.IsDir() {

			array_delete(logDirList, key)

		}
	}
}