package xtpl

import (
	"bytes"
	"html"
	"regexp"
	"strings"
)

var (
	regSections = regexp.MustCompile(`(?Us)@section\([^)]+\).*@endsection`)
	regExtends = regexp.MustCompile(`(?U)@extends\([^)]+\)`)
	regIncludes = regexp.MustCompile(`(?U)@include\([^)]+\)`)
	regYields = regexp.MustCompile(`(?U)@yield\([^)]+\)`)
	regCutSpaces = regexp.MustCompile(`@[a-z]+[\s]+\(`)
)

func (x *xtpl) preBuild(src string) []rune {
	var sections = make(map[string]string)
	var cutSections = func(src string) string {
		return strings.TrimSpace(regSections.ReplaceAllStringFunc(src, func(s string) string {
			var sectionName = strings.Trim(strings.SplitN(strings.SplitN(s, ")", 2)[0], "(", 2)[1], `"'`)
			var sectionData = strings.TrimSpace(strings.SplitN(strings.SplitN(s, "@endsection", 2)[0], ")", 2)[1])
			sections[sectionName] += sectionData
			return ""
		}))
	}
	src = strings.TrimSpace(src)

	if strings.HasPrefix(src, "@extends(") {
		src = regExtends.ReplaceAllStringFunc(src, func(s string) string {
			s = strings.Split(s, "(")[1]
			s = strings.Split(s, ")")[0]
			s = strings.Trim(s, `"'`)
			return cutSections(x.tplSource(s))
		})
	}

	src = regIncludes.ReplaceAllStringFunc(src, func(s string) string {
		s = strings.Split(s, "(")[1]
		s = strings.Split(s, ")")[0]
		s = strings.Trim(s, `"'`)
		return cutSections(x.tplSource(s))
	})

	src = cutSections(src)

	src = regYields.ReplaceAllStringFunc(src, func(s string) string {
		var arr = strings.Split(strings.SplitN(strings.SplitN(s, "(", 2)[1], ")", 2)[0], ",")
		if len(arr) == 0 {
			return ""
		}
		var sectionName = strings.Trim(strings.TrimSpace(arr[0]), `"'`)
		if data, ok := sections[sectionName]; ok {
			return data
		}
		if len(arr) < 2 {
			return ""
		}
		return strings.Trim(strings.TrimSpace(arr[1]), `"'`)
	})

	src = regCutSpaces.ReplaceAllStringFunc(src, func(s string) string {
		return strings.Join(strings.Fields(s), "")
	})

	return []rune(strings.TrimSpace(src))
}

func (x *xtpl) buildTree(src []rune, witchPrefixDog bool) (tree []treeNode) {
	var plainEcho = func(buff *bytes.Buffer) treeNode {
		var s = buff.String()
		buff.Reset()
		return func(vars *xVarCollection) []byte {
			return []byte(s)
		}
	}
	var buff = &bytes.Buffer{}

	for i := 0; i < len(src); {
		if function, offset := x.node(src[i:], witchPrefixDog); offset > 0 {
			if buff.Len() > 0 {
				tree = append(tree, plainEcho(buff))
			}
			i += offset
			tree = append(tree, function)
		} else {
			buff.WriteRune(src[i])
			i++
		}
	}

	if buff.Len() > 0 {
		tree = append(tree, plainEcho(buff))
	}

	return
}

func (x *xtpl) node(src []rune, withPrefixDog bool) (function treeNode, offset int) {
	switch {
	case hasPrefix(src, "{{--"):
		if offset = getOffset(src, "--}}", "", false, false); offset > 0 {
			return func(vars *xVarCollection) []byte {
				return nil
			}, offset + 4
		}
	case hasPrefix(src, "{*"):
		if offset = getOffset(src, "*}", "", false, false); offset > 0 {
			return func(vars *xVarCollection) []byte {
				return nil
			}, offset + 2
		}
	case hasPrefix(src, "{{"):
		if offset = getOffset(src, "}}", "", true, true); offset > 0 {
			var f = x.exec(src[2:offset])
			return func(vars *xVarCollection) []byte {
				return []byte(html.EscapeString(f(vars).toString()))
			}, offset + 2
		}
	case hasPrefix(src, "{!!"):
		if offset = getOffset(src, "!!}", "", true, true); offset > 0 {
			var f = x.exec(src[3:offset])
			return func(vars *xVarCollection) []byte {
				return f(vars).toBytes()
			}, offset + 3
		}
	case hasPrefix(src, "@foreach("):
		if offset = getOffset(src, "@endforeach", "@foreach", false, false); offset > 0 {
			return x.execForeach(src[8:offset]), offset + 11
		}
	case hasPrefix(src, "@for("):
		if offset = getOffset(src, "@endfor", "@for", false, false); offset > 0 {
			return x.execFor(src[4:offset]), offset + 7
		}
	case hasPrefix(src, "@if("):
		if offset = getOffset(src, "@endif", "@if", false, false); offset > 0 {
			return x.execIf(src[3:offset]), offset + 6
		}
	case hasPrefix(src, "@exec("):
		if offset = getOffset(src, ")", "", true, true); offset > 0 {
			var f = x.exec(src[6:offset])
			return func(vars *xVarCollection) []byte {
				f(vars)
				return nil
			}, offset + 1
		}
	}

	if f, o := x.userFunction(src, withPrefixDog); o > 0 {
		offset = o
		function = func(vars *xVarCollection) []byte {
			return f(vars).toBytes()
		}
	}

	return
}

func (x *xtpl) userFunction(src []rune, withPrefixDog bool) (function func(vars *xVarCollection) *xVar, offset int) {
	for name, f := range xtplFunctions {
		if withPrefixDog {
			name = "@" + name
		}
		if hasPrefix(src, name) && hasPrefix(src[len(name):], "(") {
			if offset = getOffset(src, ")", "", true, true); offset > 0 {
				return x.execUserFunction(src[len(name):offset+1], f), offset + 1
			}
		}
	}
	return nil, 0
}

func getOffset(src []rune, find string, skipInnerTag string, excludeTextInQuotes, excludeTextInBrackets bool) (offset int) {
	var openSingleQuotes = false
	var openDoubleQuotes = false
	var openBrackets = 0
	var openSkipTags = 0
	for offset = 0; offset < len(src); offset++ {
		if excludeTextInQuotes {
			// Проверка на одинарную кавычку
			if src[offset] == '\'' && find != "'" {
				if !openDoubleQuotes && openBrackets <= 0 && (offset == 0 || src[offset-1] != '\\') {
					openSingleQuotes = !openSingleQuotes
				}
			}
			// Проверка на двойную кавычку
			if src[offset] == '"' && find != `"` {
				if !openSingleQuotes && openBrackets <= 0 && (offset == 0 || src[offset-1] != '\\') {
					openDoubleQuotes = !openDoubleQuotes
				}
			}
		}
		// Проверка на круглую скобку
		if excludeTextInBrackets && !openSingleQuotes && !openDoubleQuotes {
			if src[offset] == '(' && find != "(" {
				openBrackets++
			} else if src[offset] == ')' {
				openBrackets--
			}
		}
		if skipInnerTag != "" && !openSingleQuotes && !openDoubleQuotes && openBrackets <= 0 {
			if hasPrefix(src[offset:], skipInnerTag) {
				openSkipTags++
			} else if hasPrefix(src[offset:], strings.Replace(skipInnerTag, "@", "@end", 1)) {
				openSkipTags--
			}
		}
		if !openSingleQuotes && !openDoubleQuotes && openBrackets <= 0 && openSkipTags <= 0 && hasPrefix(src[offset:], find) {
			return offset
		}
	}
	return -1
}

func hasPrefix(src []rune, prefix string) bool {
	if len(src) < len(prefix) {
		return false
	}
	for i, r := range prefix {
		if r != src[i] {
			return false
		}
	}
	return true
}

func trim(src []rune) []rune {
	return []rune(strings.TrimSpace(string(src)))
}
