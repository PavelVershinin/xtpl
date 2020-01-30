package xtpl

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"sync"
)

var (
	m             sync.RWMutex
	collection    map[string]*xtpl
	viewsPath     string
	viewExtension string
	cyclesLimit   uint
	debug         bool
)

func init() {
	collection = make(map[string]*xtpl)
	cyclesLimit = 10000
	viewsPath = "."
	viewExtension = "tpl"
}

// ViewsPath Путь к корневой директории с шаблонами
func ViewsPath(path string) {
	m.Lock()
	viewsPath = path
	for fileName := range collection {
		collection[fileName] = xtplInit(fileName)
	}
	m.Unlock()
}

// ViewExtension Расширение файлов шаблона
func ViewExtension(extension string) {
	m.Lock()
	viewExtension = extension
	for fileName := range collection {
		collection[fileName] = xtplInit(fileName)
	}
	m.Unlock()
}

// Functions Загрузка пользовательских функций в шаблонизатор.
func Functions(functions map[string]interface{}) {
	m.Lock()
	for name, function := range functions {
		xtplFunctions[name] = function
	}
	for fileName := range collection {
		collection[fileName] = xtplInit(fileName)
	}
	m.Unlock()
}

// CycleLimit Установка ограничения, на максимальное количество итераций в циклах. По умолчанию 10000
func CycleLimit(limit uint) {
	m.Lock()
	cyclesLimit = limit
	m.Unlock()
}

// Debug Переключение в режим отладки.
// В этом режиме, все изменения в шаблонах подхватываются налету, однако обработка шаблона занимет больше времени
func Debug(on bool) {
	debug = on
}

// View Обработка шаблона
func View(tplPath string, data map[string]interface{}, writer io.Writer) {
	if debug {
		xtplInit(tplPath).run(data, writer)
		return
	}

	m.RLock()
	view, ok := collection[tplPath]
	m.RUnlock()

	if !ok {
		m.Lock()
		view = xtplInit(tplPath)
		collection[tplPath] = view
		m.Unlock()
	}
	view.run(data, writer)
}

// ParseString Обработает строку как шаблон, вернёт строку с результатом обработки
func ParseString(source string, data map[string]interface{}) string {
	var buff = &bytes.Buffer{}
	if debug {
		xtplInitFromSource(source).run(data, buff)
		return buff.String()
	}

	var h = md5.New()
	h.Write([]byte(source))
	var tplKey = "xtpl_" + fmt.Sprintf("%x", h.Sum(nil))

	m.RLock()
	view, ok := collection[tplKey]
	m.RUnlock()

	if !ok {
		m.Lock()
		view = xtplInitFromSource(source)
		collection[tplKey] = view
		m.Unlock()
	}

	view.run(data, buff)

	return buff.String()
}
