//日志排序操作

package just

import (
	"sort"
)

// type LogSorter LogList

/*type Item struct {
	Key int
	Val LogInfo
}*/

func LogSort(loglist LogList) LogList {
	// loglist_array := NewLogSorter(loglist)
	sort.Sort(loglist)
	/*	loglist_new := []LogInfo{}
		for _, v := range loglist_array {
			loglist_new = append(loglist_new, v.Val)
		}*/
	return loglist
}

/*func NewLogSorter(m []LogInfo) LogSorter {
	ms := make(LogSorter, 0, len(m))

	for k, v := range m {
		ms = append(ms, Item{k, v})
	}

	return ms
}
*/
func (logList LogList) Len() int {
	return len(logList)
}

func (logList LogList) Less(i, j int) bool {
	return logList[i].Date < logList[j].Date // 按值排序
}

func (logList LogList) Swap(i, j int) {
	logList[i], logList[j] = logList[j], logList[i]
}
