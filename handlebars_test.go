package handlebars

import . "github.com/smartystreets/goconvey/convey"
import "io/ioutil"
import "testing"
import "path"

func TestHandlebars(t *testing.T) {
	Convey("handlebars", t, func() {
		render := Create(map[string]interface{}{
			"dir": "./template",
		})

		Convey("parse", func() {
			t.Logf("dir: %s, ext: %s, cacheLimit: %d", render.Dir, render.Ext, render.CacheLimit)

			c := map[string]string{
				"name": "haoxin",
			}

			t.Logf("context: %v", c)

			s := render.Parse("name:{{name}}", c)
			// So(s, ShouldEqual, "name:haoxin")
			// todo: check this is from cache
			s = render.Parse("name:{{name}}", c)
			// So(s, ShouldEqual, "name:haoxin")
			t.Log(s)
		})

		Convey("parse in layout", func() {})

		Convey("render", func() {
			s := render.Render("index", map[string]interface{}{
				"title": "test",
				"body":  "<pre>hello</pre>",
			})

			t.Log(s)

			// So()

			s = render.Render("index", map[string]interface{}{
				"title": "test",
				"body":  "<pre>hello</pre>",
			})
			// todo: check this is from cache
		})

		Convey("render in layout", func() {})

		Convey("read file", func() {
			So(readFile("test"), ShouldEqual, "<h1>test</h1>\n")
		})
	})
}

// only for test
func readFile(filename string) string {
	var abs string

	if len(path.Ext(filename)) == 0 {
		abs = resolve("template/result", filename) + ".html"
	} else {
		abs = resolve("template/result", filename)
	}

	bytes, err := ioutil.ReadFile(abs)
	panicError(err)

	return string(bytes)
}
