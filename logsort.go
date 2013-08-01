//日志排序操作

package just

import (
	"sort"
)

func LogSort(loglist LogList) LogList {
	sort.Sort(loglist)
	return loglist
}

func (logList LogList) Len() int {
	return len(logList)
}

func (logList LogList) Less(i, j int) bool {
	return logList[i].Date > logList[j].Date // 按值排序
}

func (logList LogList) Swap(i, j int) {
	logList[i], logList[j] = logList[j], logList[i]
}
