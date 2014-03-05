package main

type parser struct {
	r              Renderer
	refs           map[string]*reference
	inlineCallback [256]inlineParser
	flags          int
	nesting        int
	maxNesting     int
	insideLink     bool
}

func main() {
	content := `+ Email：xlqstar at gmail.com  
+ Weibo：[http://weibo.com/u/2867098822][7]  
+ GitHub：[https://github.com/xlqstar][8]


  [1]: ./images/dsfsC.jpg
  [2]: ./images/bibo8.jpg
  [3]: ./images/tuohuang.jpg
  [4]: ./images/baidu.jpg
  [5]: ./images/efound.jpg`

	beg, end := 0, 0
	for beg < len(content) { // iterate over lines

	}
}

func isReference(p *parser, data []byte) int {
	// up to 3 optional leading spaces
	if len(data) < 4 {
		return 0
	}
	i := 0
	for i < 3 && data[i] == ' ' {
		i++
	}

	// id part: anything but a newline between brackets
	if data[i] != '[' {
		return 0
	}
	i++
	idOffset := i
	for i < len(data) && data[i] != '\n' && data[i] != '\r' && data[i] != ']' {
		i++
	}
	if i >= len(data) || data[i] != ']' {
		return 0
	}
	idEnd := i

	// spacer: colon (space | tab)* newline? (space | tab)*
	i++
	if i >= len(data) || data[i] != ':' {
		return 0
	}
	i++
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i < len(data) && (data[i] == '\n' || data[i] == '\r') {
		i++
		if i < len(data) && data[i] == '\n' && data[i-1] == '\r' {
			i++
		}
	}
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i >= len(data) {
		return 0
	}

	// link: whitespace-free sequence, optionally between angle brackets
	if data[i] == '<' {
		i++
	}
	linkOffset := i
	for i < len(data) && data[i] != ' ' && data[i] != '\t' && data[i] != '\n' && data[i] != '\r' {
		i++
	}
	linkEnd := i
	if data[linkOffset] == '<' && data[linkEnd-1] == '>' {
		linkOffset++
		linkEnd--
	}

	// optional spacer: (space | tab)* (newline | '\'' | '"' | '(' )
	for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
		i++
	}
	if i < len(data) && data[i] != '\n' && data[i] != '\r' && data[i] != '\'' && data[i] != '"' && data[i] != '(' {
		return 0
	}

	// compute end-of-line
	lineEnd := 0
	if i >= len(data) || data[i] == '\r' || data[i] == '\n' {
		lineEnd = i
	}
	if i+1 < len(data) && data[i] == '\r' && data[i+1] == '\n' {
		lineEnd++
	}

	// optional (space|tab)* spacer after a newline
	if lineEnd > 0 {
		i = lineEnd + 1
		for i < len(data) && (data[i] == ' ' || data[i] == '\t') {
			i++
		}
	}

	// optional title: any non-newline sequence enclosed in '"() alone on its line
	titleOffset, titleEnd := 0, 0
	if i+1 < len(data) && (data[i] == '\'' || data[i] == '"' || data[i] == '(') {
		i++
		titleOffset = i

		// look for EOL
		for i < len(data) && data[i] != '\n' && data[i] != '\r' {
			i++
		}
		if i+1 < len(data) && data[i] == '\n' && data[i+1] == '\r' {
			titleEnd = i + 1
		} else {
			titleEnd = i
		}

		// step back
		i--
		for i > titleOffset && (data[i] == ' ' || data[i] == '\t') {
			i--
		}
		if i > titleOffset && (data[i] == '\'' || data[i] == '"' || data[i] == ')') {
			lineEnd = titleEnd
			titleEnd = i
		}
	}
	if lineEnd == 0 { // garbage after the link
		return 0
	}

	// a valid ref has been found

	// id matches are case-insensitive
	id := string(bytes.ToLower(data[idOffset:idEnd]))
	p.refs[id] = &reference{
		link:  data[linkOffset:linkEnd],
		title: data[titleOffset:titleEnd],
	}

	return lineEnd
}
