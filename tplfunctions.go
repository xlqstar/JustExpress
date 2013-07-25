//just主要入口函数，负责编译日志

package just

func eq(a interface{}, b interface{}) bool {
	if a == b {
		return true
	}
	return false
}

func neq(a interface{}, b interface{}) bool {
	return !eq(a, b)
}

func add(x int, y int) int {
	return x + y
}
func minus(x int, y int) int {
	return x - y
}
func not(x bool) bool {
	return !x
}

func last(key int, length int) bool {
	if length == 0 || key == length-1 {
		return true
	}
	return false
}

func first(key int) bool {
	if key == 0 {
		return true
	}
	return false
}
