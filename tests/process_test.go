package tests

import (
	"bufio"
	"github.com/renatoathaydes/magnanimous/mg"
	"strings"
	"testing"
)

func TestProcessSimple(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello world"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello world"})
}

func TestProcessIncludeSimple(t *testing.T) {

	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.StringContent{Text: "from another file!"})

	m := mg.WebFilesMap{}
	m["source/processed/example.html"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: m}

	r := bufio.NewReader(strings.NewReader("hello {{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hello.html", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"hello ", "from another file!"})
}

func TestMarkdownToHtml(t *testing.T) {
	file := mg.ProcessedFile{}
	file.AppendContent(&mg.StringContent{Text: "<html><body class=\"hello\">"})
	file.AppendContent(&mg.StringContent{Text: "# Hello\n", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "## Mag", MarkDown: true})
	file.AppendContent(&mg.StringContent{Text: "</body></html>"})

	html := mg.MarkdownToHtml(file)

	m := mg.WebFilesMap{}
	result, err := html.Bytes(m, nil)

	if err != nil {
		t.Fatal(err)
	}

	expectedHtml := "<html><body class=\"hello\"><h1>Hello</h1>\n<h2>Mag</h2>\n</body></html>"
	if string(result) != expectedHtml {
		t.Errorf("Expected '%s', but was '%s'", expectedHtml, result)
	}
}

func TestProcessIncludeMarkDown(t *testing.T) {
	// setup a md file to be included
	exampleFile := mg.ProcessedFile{}
	exampleFile.AppendContent(&mg.HtmlFromMarkdownContent{
		MarkDownContent: &mg.StringContent{Text: "# header", MarkDown: true},
	})

	m := mg.WebFilesMap{}
	m["source/example.md"] = mg.WebFile{Processed: &exampleFile}

	resolver := mg.DefaultFileResolver{BasePath: "source", Files: m}

	// read a file that includes the previous file
	r := bufio.NewReader(strings.NewReader("## hello {{ include /example.md }}"))
	processed, err := mg.ProcessReader(r, "source/processed/hi.md", 11, &resolver)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), m, processed, emptyContext, []string{"<h2>hello</h2>\n", "<h1>header</h1>\n"})
}

func TestProcessIgnoreEscapedBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("hello \\{{ include example.html }}"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{"hello {{ include example.html }}"})
}

func TestProcessIgnoreEscapedClosingBrackets(t *testing.T) {
	r := bufio.NewReader(strings.NewReader("Hello {{\n  bad-instruction \"contains \\}} ignored\"\n}}. How are you?"))
	processed, err := mg.ProcessReader(r, "", 11, nil)

	if err != nil {
		t.Fatal(err)
	}

	checkParsing(t, processed.Context(), emptyFilesMap, processed, emptyContext, []string{
		"Hello ",
		"{{\n  bad-instruction \"contains }} ignored\"\n}}",
		". How are you?",
	})
}
