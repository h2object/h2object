package template

import (
	"os"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"regexp"
	"io/ioutil"

	"github.com/h2object/h2object/log"
)

// This object handles loading and parsing of templates.
// Everything below the application's views directory is treated as a template.
type TemplateLoader struct {
	log.Logger
	// This is the set of all templates under views
	templateSet *template.Template
	// If an error was encountered parsing the templates, it is stored here.
	compileError *Error
	// template delims
	delimiters string
	// Paths to search for templates, in priority order.
	paths []string
	// Map from template name to the path from whence it was loaded.
	templatePaths map[string]string
}

func NewTemplateLoader(delimiters string, paths []string, logger log.Logger) *TemplateLoader {
	loader := &TemplateLoader{
		Logger: logger,
		delimiters: delimiters,
		paths: paths,
	}
	return loader
}

func (loader *TemplateLoader) Load() error {
	if err := loader.Refresh(); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

// This scans the views directory and parses all templates as Go Templates.
// If a template fails to parse, the error is set on the loader.
// (It's awkward to refresh a single Go Template)
func (loader *TemplateLoader) Refresh() *Error {
	loader.compileError = nil
	loader.templatePaths = map[string]string{}

	// Set the template delimiters for the project if present, then split into left
	// and right delimiters around a space character
	var splitDelims []string
	if loader.delimiters != "" {
		splitDelims = strings.Split(loader.delimiters, " ")
		if len(splitDelims) != 2 {
			loader.Error("templates loader delimiters format error")
		}
	}

	// Walk through the template loader's paths and build up a template set.
	var templateSet *template.Template = nil
	for _, basePath := range loader.paths {
		// Walk only returns an error if the template loader is completely unusable
		// (namely, if one of the TemplateFuncs does not have an acceptable signature).

		// Handling symlinked directories
		var fullSrcDir string
		f, err := os.Lstat(basePath)
		if err == nil && f.Mode()&os.ModeSymlink == os.ModeSymlink {
			fullSrcDir, err = filepath.EvalSymlinks(basePath)
			if err != nil {
				panic(err)
			}
		} else {
			fullSrcDir = basePath
		}

		var templateWalker func(path string, info os.FileInfo, err error) error
		templateWalker = func(path string, info os.FileInfo, err error) error {
			if err != nil {
				loader.Error("templates loader walking templates failed: (%s)", err.Error())
				return nil
			}

			// is it a symlinked template?
			link, err := os.Lstat(path)
			if err == nil && link.Mode()&os.ModeSymlink == os.ModeSymlink {
				loader.Trace("templates loader symlink template: (%s)", path)
				// lookup the actual target & check for goodness
				targetPath, err := filepath.EvalSymlinks(path)
				if err != nil {
					loader.Error("templates loader read symlink failed: (%s)", err.Error())
					return err
				}
				targetInfo, err := os.Stat(targetPath)
				if err != nil {
					loader.Error("templates loader stat symlink failed: (%s)", err.Error())
					return err
				}

				// set the template path to the target of the symlink
				path = targetPath
				info = targetInfo

				// need to save state and restore for recursive call to Walk on symlink
				tmp := fullSrcDir
				fullSrcDir = filepath.Dir(targetPath)
				filepath.Walk(targetPath, templateWalker)
				fullSrcDir = tmp
			}

			// Walk into watchable directories
			if info.IsDir() {
				if !loader.WatchDir(info) {
					return filepath.SkipDir
				}
				return nil
			}

			// Only add watchable
			if !loader.WatchFile(info.Name()) {
				return nil
			}

			var fileStr string

			// addTemplate loads a template file into the Go template loader so it can be rendered later
			addTemplate := func(templateName string) (err error) {
				// Convert template names to use forward slashes, even on Windows.
				if os.PathSeparator == '\\' {
					templateName = strings.Replace(templateName, `\`, `/`, -1) // `
				}

				// If we already loaded a template of this name, skip it.
				if _, ok := loader.templatePaths[templateName]; ok {
					return nil
				}

				loader.templatePaths[templateName] = path

				// Load the file if we haven't already
				if fileStr == "" {
					fileBytes, err := ioutil.ReadFile(path)
					if err != nil {
						loader.Error("templates loader reading file failed:(%s)",path)
						return nil
					}

					fileStr = string(fileBytes)
				}

				if templateSet == nil {
					// Create the template set.  This panics if any of the funcs do not
					// conform to expectations, so we wrap it in a func and handle those
					// panics by serving an error page.
					var funcError *Error
					func() {
						defer func() {
							if err := recover(); err != nil {
								funcError = &Error{
									Title:       "Panic (Template Loader)",
									Description: fmt.Sprintln(err),
								}
							}
						}()
						templateSet = template.New(templateName).Funcs(TemplateFuncs)
						// If alternate delimiters set for the project, change them for this set
						if splitDelims != nil {
							templateSet.Delims(splitDelims[0], splitDelims[1])
						} else {
							// Reset to default otherwise
							templateSet.Delims("", "")
						}
						_, err = templateSet.Parse(fileStr)
					}()

					if funcError != nil {
						return funcError
					}

				} else {
					if splitDelims != nil {
						templateSet.Delims(splitDelims[0], splitDelims[1])
					} else {
						templateSet.Delims("", "")
					}
					_, err = templateSet.New(templateName).Parse(fileStr)
				}
				return err
			}

			templateName := path[len(fullSrcDir)+1:]

			err = addTemplate(templateName)

			// Store / report the first error encountered.
			if err != nil && loader.compileError == nil {
				_, line, description := ParseTemplateError(err)
				loader.compileError = &Error{
					Title:       "Template Compilation Error",
					Path:        templateName,
					Description: description,
					Line:        line,
					SourceLines: strings.Split(fileStr, "\n"),
				}
				loader.Error("templates loader template compiling failed:(%s:%d) (%s)",
					templateName, line, description)			
			}
			return nil
		}

		funcErr := filepath.Walk(fullSrcDir, templateWalker)

		// If there was an error with the Funcs, set it and return immediately.
		if funcErr != nil {
			loader.compileError = funcErr.(*Error)
			return loader.compileError
		}
	}

	// Note: compileError may or may not be set.
	loader.templateSet = templateSet
	return loader.compileError
}

func (loader *TemplateLoader) WatchDir(info os.FileInfo) bool {
	// Watch all directories, except the ones starting with a dot.
	return !strings.HasPrefix(info.Name(), ".")
}

func (loader *TemplateLoader) WatchFile(basename string) bool {
	// Watch all files, except the ones starting with a dot.
	return !strings.HasPrefix(basename, ".")
}


// Parse the line, and description from an error message like:
// html/template:Application/Register.html:36: no such template "footer.html"
func ParseTemplateError(err error) (templateName string, line int, description string) {
	description = err.Error()
	i := regexp.MustCompile(`:\d+:`).FindStringIndex(description)
	if i != nil {
		line, err = strconv.Atoi(description[i[0]+1 : i[1]-1])
		if err == nil {
			templateName = description[:i[0]]
			if colon := strings.Index(templateName, ":"); colon != -1 {
				templateName = templateName[colon+1:]
			}
			templateName = strings.TrimSpace(templateName)
			description = description[i[1]+1:]
		}
	}
	return templateName, line, description
}

// Return the Template with the given name.  The name is the template's path
// relative to a template loader root.
//
// An Error is returned if there was any problem with any of the templates.  (In
// this case, if a template is returned, it may still be usable.)
func (loader *TemplateLoader) Template(name string) (Template, error) {
	// Case-insensitive matching of template file name
	name = strings.TrimPrefix(name, "/")
	name = strings.ToLower(name)
	for k := range loader.templatePaths {
		if name == strings.ToLower(k) {
			name = k
		}
	}
	// Look up and return the template.
	if loader.templateSet == nil {
		return nil, fmt.Errorf("template %s not found.", name)
	}

	tmpl := loader.templateSet.Lookup(name)

	// This is necessary.
	// If a nil loader.compileError is returned directly, a caller testing against
	// nil will get the wrong result.  Something to do with casting *Error to error.
	var err error
	if loader.compileError != nil {
		err = loader.compileError
	}

	if tmpl == nil && err == nil {
		return nil, fmt.Errorf("template %s not found.", name)
	}

	return GoTemplate{tmpl, loader}, err
}