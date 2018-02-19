package main

import (
	"github.com/gopherjs/gopherjs/js"
	"io"
	"bytes"
	. "github.com/reusing-code/csvrewrite"
	"honnef.co/go/js/dom"
	"github.com/MJKWoolnough/gopherjs/files"
	"bufio"
	"encoding/base64"
)

type Component interface {
	Render() string
}

func main() {
	idx := IndexComponent{}
	js.Global.Get("document").Call("write", idx.Render())

	d := dom.GetWindow().Document()

	bsfilechooser := d.GetElementByID("csvinput")
	bsfilechooser.AddEventListener("change", false, func(event dom.Event) {
		input := event.Target().(*dom.HTMLInputElement)
		if len(input.Files()) > 0 {
			file := files.NewFile(input.Files()[0])
			go handleFile(file)
		}
	})


}

func handleFile(file files.File) {
	d:= dom.GetWindow().Document()

	buf := bytes.Buffer{}
	writer := bufio.NewWriter(base64.NewEncoder(base64.URLEncoding, &buf))
	reader := files.NewFileReader(file);
	defer reader.Close()

	converted, errorLog := rewrite(reader)
	outputContainer := d.GetElementByID("rewrite-output-container").(*dom.HTMLDivElement)
	if len(errorLog) > 0 {
		outputContainer.Style().SetProperty("display", "block", "")
		d.GetElementByID("rewrite-output-content").SetInnerHTML(errorLog)
	} else {
		outputContainer.Style().SetProperty("display ", "none", "")
	}
	writer.WriteString(converted)
	writer.Flush()
	elem := d.CreateElement("a")
	elem.SetAttribute("download", "converted_" + file.Name())
	elem.SetAttribute("href", "data:text/plain;base64," + buf.String())
	elem.(*dom.HTMLAnchorElement).Click()
}

func rewrite(in io.Reader) (out string, errorLog string) {
	rewriter := NewRewriter()
	errBuff := new(bytes.Buffer)
	rewriter.SetErrorWriter(errBuff)
	rewriter.SetInputProcessor(NewComdirectInput(PersonalPayees{}))
	rewriter.SetOutProcessor(&YNABOutput{})
	rewriter.ImportTransactions(in)

	outBuff := new(bytes.Buffer)
	rewriter.ExportTransactions(outBuff)

	out = outBuff.String()
	errorLog = errBuff.String()
	return
}
