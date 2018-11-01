package mg

import (
	"bytes"
	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"path/filepath"
	"strings"
)

type HtmlFromMarkdownContent struct {
	MarkDownContent []Content
}

var _ Content = (*HtmlFromMarkdownContent)(nil)

var chromaRenderer = blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.WithoutAutodetect()))

func MarkdownToHtml(file ProcessedFile) ProcessedFile {
	return ProcessedFile{
		contents:     []Content{&HtmlFromMarkdownContent{MarkDownContent: file.contents}},
		context:      file.context,
		rootScope:    file.rootScope,
		NewExtension: ".html",
	}
}

func (f *HtmlFromMarkdownContent) Write(writer io.Writer, files WebFilesMap, inclusionChain []InclusionChainItem) error {
	htmlHead, main, htmlFooter, err := readBytes(f.MarkDownContent, files, inclusionChain)
	if err != nil {
		return err
	}
	if len(htmlHead) > 0 {
		_, err = writer.Write(htmlHead)
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	}

	_, err = writer.Write(blackfriday.Run(main, chromaRenderer))
	if err != nil {
		return &MagnanimousError{Code: IOError, message: err.Error()}
	}

	if len(htmlFooter) > 0 {
		_, err = writer.Write(htmlFooter)
		if err != nil {
			return &MagnanimousError{Code: IOError, message: err.Error()}
		}
	}

	return nil
}

func readBytes(c []Content, files WebFilesMap,
	inclusionChain []InclusionChainItem) (head, body, foot []byte, err error) {
	var header, main, footer bytes.Buffer
	header.Grow(128)
	main.Grow(1024)
	footer.Grow(128)

	inHeader := true
	lastIndex := len(c) - 1

	for i, content := range c {
		var writer *bytes.Buffer = nil
		if inHeader {
			if isHtml(content) {
				writer = &header
			} else {
				inHeader = false
				writer = &main
			}
		} else {
			if i == lastIndex && isHtml(content) {
				writer = &footer
			} else {
				writer = &main
			}
		}
		err = content.Write(writer, files, inclusionChain)
		if err != nil {
			return
		}
	}

	head, body, foot = header.Bytes(), main.Bytes(), footer.Bytes()
	return
}

func isHtml(c Content) bool {
	switch inc := c.(type) {
	case *IncludeInstruction:
		return strings.ToLower(filepath.Ext(inc.Path)) == ".html"
	}
	return false
}
