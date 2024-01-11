package xtpl_test

import (
	"bytes"
	"embed"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/PavelVershinin/xtpl"
)

var ptr = struct {
	Value string
}{
	Value: "Test the pointer",
}

var testData1 = map[string]interface{}{
	"page_title": "Тестовая страница",
	"strings":    []string{"Первый", "Второй", "Третий", "048465450"},
	"numbers":    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
	"structs": []struct {
		ID       int
		Name     string
		Value    string
		Function func(s string, i int) string
		Pointer  interface{}
	}{
		{
			1,
			"my name",
			"my value",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " дней"
			},
			&ptr,
		},
		{
			2,
			"my name 2",
			"my value 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " дней"
			},
			nil,
		},
	},
	"pointer": &ptr,
}

var testData2 = map[string]interface{}{
	"page_title": "Тестовая страница 2",
	"strings":    []string{"Первый 2", "Второй 2", "Третий 2", "07869786"},
	"numbers":    []int{0, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	"structs": []struct {
		ID       int
		Name     string
		Value    string
		Function func(s string, i int) string
		Pointer  interface{}
	}{
		{
			3,
			"my name 2",
			"my value 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " часов"
			},
			&ptr,
		},
		{
			4,
			"my name 2 2",
			"my value 2 2",
			func(s string, i int) string {
				return s + " " + strconv.Itoa(i) + " часов"
			},
			nil,
		},
	},
}

var update = flag.Bool("update", false, "update .golden files")

//go:embed all:testdata
var embeddedFS embed.FS

func init() {
	xtpl.ViewsPath(filepath.Join(".", "testdata", "templates"))
	xtpl.ViewExtension("tpl")
	xtpl.CycleLimit(100)
	xtpl.Debug(false)
	xtpl.Functions(map[string]interface{}{
		"date": func(timestamp int64, layout string) string {
			t := time.Unix(timestamp, 0)
			return t.UTC().Format(layout)
		},
	})

	xtpl.EmbeddedFs(embeddedFS)
}

func TestXtplCollection_View(t *testing.T) {
	t.Run("test data 1", func(t *testing.T) {
		t.Parallel()
		buff := bytes.Buffer{}
		goldenFile := filepath.Join(".", "testdata", "result_1.gold")
		if err := xtpl.View("index", testData1, &buff); err != nil {
			t.Fatal(err)
		}
		if *update {
			t.Log("update golden file")
			if err := os.WriteFile(goldenFile, buff.Bytes(), os.ModePerm); err != nil {
				t.Fatalf("error writing golden file: %s", err)
			}
		}
		res, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatalf("error reading golden file: %s", err)
		}
		if !bytes.Equal(res, buff.Bytes()) {
			t.Error("result not equal golden file")
		}
	})

	t.Run("test data 2", func(t *testing.T) {
		t.Parallel()
		buff := bytes.Buffer{}
		goldenFile := filepath.Join(".", "testdata", "result_2.gold")
		if err := xtpl.View("index", testData2, &buff); err != nil {
			t.Fatal(err)
		}
		if *update {
			t.Log("update golden file")
			if err := os.WriteFile(goldenFile, buff.Bytes(), os.ModePerm); err != nil {
				t.Fatalf("error writing golden file: %s", err)
			}
		}
		res, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatalf("error reading golden file: %s", err)
		}
		if !bytes.Equal(res, buff.Bytes()) {
			t.Error("result not equal golden file")
		}
	})
}

func TestXtplCollection_String(t *testing.T) {
	if *update {
		t.Skip()
	}

	template, err := os.ReadFile(filepath.Join(".", "testdata", "templates", "index.tpl"))
	if err != nil {
		t.Fatalf("error reading template file: %s", err)
	}

	t.Run("test data 1", func(t *testing.T) {
		t.Parallel()
		buff := bytes.Buffer{}
		goldenFile := filepath.Join(".", "testdata", "result_1.gold")
		if err := xtpl.String(string(template), testData1, &buff); err != nil {
			t.Fatal(err)
		}
		res, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatalf("error reading golden file: %s", err)
		}
		if !bytes.Equal(res, buff.Bytes()) {
			t.Error("result not equal golden file")
		}
	})

	t.Run("test data 2", func(t *testing.T) {
		t.Parallel()
		buff := bytes.Buffer{}
		goldenFile := filepath.Join(".", "testdata", "result_2.gold")
		if err := xtpl.String(string(template), testData2, &buff); err != nil {
			t.Fatal(err)
		}
		res, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatalf("error reading golden file: %s", err)
		}
		if !bytes.Equal(res, buff.Bytes()) {
			t.Error("result not equal golden file")
		}
	})
}

func BenchmarkView(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buff = &bytes.Buffer{}
		if err := xtpl.View("index", testData1, buff); err != nil {
			b.Fatal(err)
		}
	}
}
