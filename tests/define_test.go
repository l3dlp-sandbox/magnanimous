package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"github.com/renatoathaydes/magnanimous/mg/expression"
	"strings"
	"testing"
	"time"
)

func TestDefineNumber(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define a 2 }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(2)

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineString(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My Site\" }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["title"] = "My Site"

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineDate(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define date1 date[\"2017-11-23T22:12:21\"][\"2016-01-02\"] }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)
	check(err)

	date1, err := time.Parse("2006-01-02T15:04:05", "2017-11-23T22:12:21")
	check(err)
	expectedCtx := make(map[string]interface{})
	expectedCtx["date1"] = &expression.DateTime{Time: date1, Format: "2016-01-02"}

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineDateNow(t *testing.T) {
	now := time.Now()

	r := bufio.NewReader(strings.NewReader("{{ define date1 date[\"now\"] }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)
	check(err)

	ctx := processed.ResolveContext(mg.WebFilesMap{}, mg.ContextStack{})
	if actualDate1, ok := ctx.Get("date1"); ok {
		if d1, ok := actualDate1.(*expression.DateTime); ok {
			if d1.Time.Unix()-now.Unix() > 1 {
				t.Errorf("Time difference between now and evaluated date[now] is too big: %d",
					d1.Time.Unix()-now.Unix())
			}
		} else {
			t.Errorf("Expected DateTime but got %v", actualDate1)
		}
	} else {
		t.Errorf("Expected map[]{date1: now} but got %v", ctx)
	}
}

func TestDefinePath(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define f1 path[\"/processed/f1.txt\"] }}" +
		"Greeting: {{ eval f1.hello }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.txt", 43, nil)
	check(err)

	r2 := bufio.NewReader(strings.NewReader("Defining hello: {{ define hello \"Hey Joe\" }}"))
	f1, err := mg.ProcessReader(r2, "source/processed/f1.txt", 11, nil)
	check(err)

	expectedCtx := make(map[string]interface{})
	expectedCtx["f1"] = &expression.Path{Value: "/processed/f1.txt"}

	files := mg.WebFilesMap{WebFiles: make(map[string]mg.WebFile, 1)}
	files.WebFiles["/processed/f1.txt"] = mg.WebFile{Processed: f1}

	checkParsing(t, files, processed, expectedCtx, []string{"", "Greeting: ", "Hey Joe"})
}

func TestDefineStringConcat(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define title \"My\" + \" Site\" }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["title"] = "My Site"

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineFromExpression(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define n 10 + 10 * (2 + 4) }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["n"] = float64(70)

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineFromOrExpression(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(`{{ define n not_defined || "alternative" }}`))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["n"] = "alternative"

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{""})
}

func TestDefineBasedOnPreviousDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ define a 10 }}" +
			"{{ define b 4 }}" +
			"{{ define c a * b }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["a"] = float64(10)
	expectedCtx["b"] = float64(4)
	expectedCtx["c"] = float64(40)

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{"", "", ""})
}

func TestDefineBasedOnPreviousEmptyStringDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader(
		"{{ define baseURL \"\" }}" +
			"{{ define INSTRUCTIONS_PATH baseURL + \"/hello\" }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedCtx := make(map[string]interface{})
	expectedCtx["baseURL"] = ""
	expectedCtx["INSTRUCTIONS_PATH"] = "/hello"

	checkParsing(t, emptyFilesMap, processed, expectedCtx, []string{"", ""})
}

func TestMalformedDefine(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("{{ define }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"{{ define }}"})

	r = bufio.NewReader(strings.NewReader("{{ define abc }}"))
	processed, err = mg.ProcessReader(r, "source/processed/hi.html", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, emptyFilesMap, processed, emptyContext, []string{"{{ define abc }}"})
}
