package just

import (
	// "fmt"
	"sort"
	"strconv"
)

type MapSorter []Item

type Item struct {
	Key string
	Val LogInfo
}

func LogSort(loglist map[string]LogInfo) map[string]LogInfo {
	/*	m := map[string]map[string]string{
			"e": {"title": "你好", "date": "12334"},
			"a": {"title": "你好", "date": "5"},
			"d": {"title": "你好", "date": "3"},
			"c": {"title": "你好", "date": "1"},
			"f": {"title": "你好", "date": "4"},
			"b": {"title": "你好", "date": "9"},
		}
	*/
	loglist_array := NewMapSorter(loglist)
	sort.Sort(loglist_array)
	loglist_new := map[string]LogInfo{}
	for _, v := range loglist_array {
		loglist_new[v.Key] = v.Val
	}
	/*	for _, item := range loglist {
		fmt.Printf("%s:%d\n", item.Key, item.Val)
	}*/
	return loglist_new

}

func NewMapSorter(m map[string]LogInfo) MapSorter {
	ms := make(MapSorter, 0, len(m))

	for k, v := range m {
		ms = append(ms, Item{k, v})
	}

	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	i_date_int, _ := strconv.ParseInt(ms[i].Val.Date, 10, 64)
	j_date_int, _ := strconv.ParseInt(ms[j].Val.Date, 10, 64)
	return i_date_int < j_date_int // 按值排序
	// return ms[i].Val < ms[j].Val // 按值排序
	//return ms[i].Key < ms[j].Key // 按键排序
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
