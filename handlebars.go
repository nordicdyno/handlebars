package handlebars

import "github.com/hoisie/mustache"
import . "github.com/tj/go-debug"
import "path/filepath"
import "strings"
import "path"
import "os"

var debug = Debug("handlebars")

type Render struct {
	Dir        string
	Ext        string
	Local      interface{}
	Cache      map[string]*mustache.Template
	CacheLimit int
}

func Create(config map[string]interface{}) *Render {
	r := new(Render)
	r.Cache = map[string]*mustache.Template{}

	if v, ok := config["dir"]; true {
		if ok {
			r.Dir = parseConfigDir(v.(string))
		} else {
			r.Dir, _ = os.Getwd()
		}
	}

	if v, ok := config["ext"]; true {
		if ok {
			s := v.(string)
			if strings.HasPrefix(s, ".") {
				r.Ext = s
			} else {
				r.Ext = "." + s
			}
		} else {
			r.Ext = ".hbs"
		}
	}

	// if v, ok := config["local"]; ok {
	// 	// todo: local
	// }

	if v, ok := config["cacheLimit"]; true {
		if ok {
			i := v.(int)

			if i < 0 {
				r.CacheLimit = 0
			} else {
				r.CacheLimit = i
			}
		} else {
			r.CacheLimit = 50
		}
	}

	return r
}

// parse string
func (r *Render) Parse(data string, context ...interface{}) string {
	debug("data: %s, context: %v", data, context)
	return r.GetTemplate(data, false).Render(context)
}

// parse string in layout
func (r *Render) ParseInLayout(data string, layoutData string, context ...interface{}) string {
	t := r.GetTemplate(data, false)
	return t.RenderInLayout(r.GetTemplate(layoutData, false), context)
}

// render file
func (r *Render) Render(filename string, context ...interface{}) string {
	debug("filename: %s, context: %v", filename, context)
	return r.GetTemplate(filename, true).Render(context)
}

// render file in layout
func (r *Render) RenderInLayout(filename string, layoutFile string, context ...interface{}) string {
	t := r.GetTemplate(filename, true)
	return t.RenderInLayout(r.GetTemplate(layoutFile, true), context)
}

// private
func (r *Render) GetAbsPath(filename string) string {
	abs := resolve(r.Dir, filename)

	if len(path.Ext(abs)) == 0 {
		return abs + r.Ext
	}

	return abs
}

// private
func (r *Render) GetTemplate(s string, isPath bool) *mustache.Template {
	if isPath {
		abs := r.GetAbsPath(s)
		t, ok := r.Cache[abs]

		if ok {
			debug("from cache")
			return t
		}

		t, err := mustache.ParseFile(abs)
		debug("parse file: %s", abs)
		panicError(err)

		r.AddToCache(abs, t)

		return t
	}

	key := hash(s)
	t, ok := r.Cache[key]

	if ok {
		debug("from cache")
		return t
	}

	t, err := mustache.ParseString(s)
	debug("parse string: %v", s)
	panicError(err)

	r.AddToCache(key, t)

	return t
}

// private
func (r *Render) AddToCache(key string, t *mustache.Template) {
	if len(r.Cache) > r.CacheLimit {
		return
	}

	debug("add to cache")
	r.Cache[key] = t
}

func AddLocal(data interface{}) {}

func RegisterHelper() {}

func RegisterPartial() {}

func RegisterPartials(dir string) {}

func parseConfigDir(s string) string {
	if filepath.IsAbs(s) {
		return s
	}

	abs, err := filepath.Abs(s)

	if err != nil {
		panic("invalid dir")
	}

	return abs
}
