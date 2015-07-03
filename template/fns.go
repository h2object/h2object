package template

import (
	"fmt"
	"strings"
	"time"
	"reflect"	
	"html"
	"html/template"
)

var (
	// The functions available for use in the templates.
	TemplateFuncs = map[string]interface{}{
		"set": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			renderArgs[key] = value
			return template.JS("")
		},
		"append": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			if renderArgs[key] == nil {
				renderArgs[key] = []interface{}{value}
			} else {
				renderArgs[key] = append(renderArgs[key].([]interface{}), value)
			}
			return template.JS("")
		},
		"firstof": func(args ...interface{}) interface{} {
			for _, val := range args {
				switch val.(type) {
				case nil:
					continue
				case string:
					if val == "" {
						continue
					}
					return val
				default:
					return val
				}
			}
			return nil
		},
		// Pads the given string with &nbsp;'s up to the given width.
		"pad": func(str string, width int) template.HTML {
			if len(str) >= width {
				return template.HTML(html.EscapeString(str))
			}
			return template.HTML(html.EscapeString(str) + strings.Repeat("&nbsp;", width-len(str)))
		},

		"dump": func(v interface{}) template.HTML {
			return template.HTML(fmt.Sprintf("%v", v))
		},
		// "object": func(bucket, key string) interface{} {
		// 	ctx := GetContext()
		// 	if ctx != nil {
		// 		obj, err := NewObject(bucket, key)
		// 		if err != nil {
		// 			return nil
		// 		}
		// 		if err := ctx.app.objects.GetObject(obj, -1, -1); err != nil {
		// 			return nil
		// 		}

		// 		return obj.Value()
		// 	}
		// 	return nil
		// },
		"goto": func(uri string) template.JS {
			script := `<script language="javascript">window.location.href="%s";</script>`
			return template.JS(fmt.Sprintf(script, uri))
		},

		// Replaces newlines with <br>
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},

		// Skips sanitation on the parameter.  Do not use with dynamic data.
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},

		// Format a date according to the application's default date(time) format.
		"date": func(date time.Time) string {
			return date.Format("2006-03-23")
		},
		"datetime": func(date time.Time) string {
			return date.Format("2006-03-23 10:34:32")
		},
		"even": func(a int) bool { return (a % 2) == 0 },
	}
)

func Function(name string, fn interface{}) {
	if reflect.TypeOf(fn).Kind() == reflect.Func {
		TemplateFuncs[name] = fn	
	}
}	
