package xtpl

import (
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
	serialID      *serial
)

func init() {
	collection = make(map[string]*xtpl)
	cyclesLimit = 10000
	viewsPath = "."
	viewExtension = "tpl"
	serialID = &serial{}
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
func View(tplPath string, data map[string]interface{}, writer io.Writer) error {
	var xTpl *xtpl

	if debug {
		xTpl = xtplInit(tplPath)
		xTpl.run(data, writer)
	} else {
		var ok bool
		m.RLock()
		xTpl, ok = collection[tplPath]
		m.RUnlock()

		if !ok {
			m.Lock()
			xTpl = xtplInit(tplPath)
			collection[tplPath] = xTpl
			m.Unlock()
		}

		xTpl.run(data, writer)
	}

	return xTpl.errors.Error()
}

// ParseString Обработает строку как шаблон, вернёт строку с результатом обработки
func ParseString(source string, data map[string]interface{}, writer io.Writer) error {
	var xTpl *xtpl

	if debug {
		xTpl = xtplInitFromSource(source)
		xTpl.run(data, writer)
	} else {
		var h = md5.New()
		if _, err := h.Write([]byte(source)); err != nil {
			return err
		}
		var tplKey = "xtpl_" + fmt.Sprintf("%x", h.Sum(nil))
		var ok bool

		m.RLock()
		xTpl, ok = collection[tplKey]
		m.RUnlock()

		if !ok {
			m.Lock()
			xTpl = xtplInitFromSource(source)
			collection[tplKey] = xTpl
			m.Unlock()
		}

		xTpl.run(data, writer)
	}

	return xTpl.errors.Error()
}
