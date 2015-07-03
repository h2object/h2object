package template

import (
	"io"	
	"io/ioutil"
	"regexp"
	"strings"
	"html/template"
)

var ERROR_CLASS = "hasError"


type Template interface {
	Name() string
	Content() []string
	Render(wr io.Writer, arg interface{}) error
}

var invalidSlugPattern = regexp.MustCompile(`[^a-z0-9 _-]`)
var whiteSpacePattern = regexp.MustCompile(`\s+`)

// Adapter for Go Templates.
type GoTemplate struct {
	*template.Template
	loader *TemplateLoader
}

// return a 'revel.Template' from Go's template.
func (gotmpl GoTemplate) Render(wr io.Writer, arg interface{}) error {
	return gotmpl.Execute(wr, arg)
}

func (gotmpl GoTemplate) Content() []string {
	content, _ := ReadLines(gotmpl.loader.templatePaths[gotmpl.Name()])
	return content
}

// Reads the lines of the given file.  Panics in the case of error.
func ReadLines(filename string) ([]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\n"), nil
}
