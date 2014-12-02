package handlebars

import . "github.com/tj/go-debug"
import "path/filepath"
import "strings"
import "path"
import "os"

var debug = Debug("handlebars")

type Render struct {
	dir        string
	ext        string
	cache      map[string]*Template
	cacheLimit int
	helpers    map[string]interface{}
	partials   map[string]*Template
}

func New(config map[string]interface{}) *Render {
	r := new(Render)
	r.cache = map[string]*Template{}
	// r.helpers = map[string]interface{}
	// r.partials = map[string]*Template{}

	if v, ok := config["dir"]; true {
		if ok {
			r.dir = parseConfigDir(v.(string))
		} else {
			r.dir, _ = os.Getwd()
		}
	}

	if v, ok := config["ext"]; true {
		if ok {
			s := v.(string)
			if strings.HasPrefix(s, ".") {
				r.ext = s
			} else {
				r.ext = "." + s
			}
		} else {
			r.ext = ".hbs"
		}
	}

	if v, ok := config["cacheLimit"]; true {
		if ok {
			i := v.(int)

			if i < 0 {
				r.cacheLimit = 0
			} else {
				r.cacheLimit = i
			}
		} else {
			r.cacheLimit = 50
		}
	}

	return r
}

// public
func (r *Render) Parse(data string, context ...interface{}) string {
	return r.getTemplate(data, false).Render(context...)
}

func (r *Render) Render(filename string, context ...interface{}) string {
	return r.getTemplate(filename, true).Render(context...)
}

// func (r *Render) RegisterHelper(name string, helper interface{}) {}

// func (r *Render) RegisterPartial(name string, filename string) {}

// private
func (r *Render) getAbsPath(filename string) string {
	abs := resolve(r.dir, filename)

	if len(path.Ext(abs)) == 0 {
		return abs + r.ext
	}

	return abs
}

func (r *Render) getTemplate(s string, isPath bool) *Template {
	if isPath {
		abs := r.getAbsPath(s)
		t, ok := r.cache[abs]

		if ok {
			debug("from cache")
			return t
		}

		t, err := parseFile(abs)
		debug("parse file: %s", abs)
		panicError(err)

		r.addToCache(abs, t)

		return t
	}

	key := hash(s)
	t, ok := r.cache[key]

	if ok {
		debug("from cache")
		return t
	}

	t, err := parseString(s)
	debug("parse string: %v", s)
	panicError(err)

	r.addToCache(key, t)

	return t
}

func (r *Render) addToCache(key string, t *Template) {
	if len(r.cache) > r.cacheLimit {
		return
	}

	debug("add to cache")
	r.cache[key] = t
}

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
