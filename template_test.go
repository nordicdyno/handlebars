package handlebars

import . "github.com/smartystreets/goconvey/convey"
import "testing"
import "path"
import "os"

type Test struct {
	tmpl     string
	context  interface{}
	expected string
}

type Data struct {
	A bool
	B string
}

type User struct {
	Name string
	Age  int
}

type settings struct {
	Allow bool
}

func (u User) Func1() string {
	return u.Name
}

func (u *User) Func2() string {
	return u.Name
}

func (u *User) Func3() (map[string]string, error) {
	return map[string]string{"name": u.Name}, nil
}

func (u *User) Func4() (map[string]string, error) {
	return nil, nil
}

func (u *User) Func5() (*settings, error) {
	return &settings{true}, nil
}

func (u *User) Func6() ([]interface{}, error) {
	var v []interface{}
	v = append(v, &settings{true})
	return v, nil
}

func (u User) Truefunc1() bool {
	return true
}

func (u *User) Truefunc2() bool {
	return true
}

func makeVector(n int) []interface{} {
	var v []interface{}
	for i := 0; i < n; i++ {
		v = append(v, &User{"haoxin", 1})
	}
	return v
}

type Category struct {
	Tag         string
	Description string
}

func (c Category) DisplayName() string {
	return c.Tag + " - " + c.Description
}

func TestTemplate(t *testing.T) {
	Convey("template", t, func() {
		Convey("basic", func() {
			tests := []Test{
				{`hello world`, nil, "hello world"},
				{`hello {{name}}`, map[string]string{"name": "world"}, "hello world"},
				{`{{var}}`, map[string]string{"var": "5 > 2"}, "5 &gt; 2"},
				{`{{{var}}}`, map[string]string{"var": "5 > 2"}, "5 > 2"},
				{`{{a}}{{b}}{{c}}{{d}}`, map[string]string{"a": "a", "b": "b", "c": "c", "d": "d"}, "abcd"},
				{`0{{a}}1{{b}}23{{c}}456{{d}}89`, map[string]string{"a": "a", "b": "b", "c": "c", "d": "d"}, "0a1b23c456d89"},
				{`hello {{! comment }}world`, map[string]string{}, "hello world"},
				{`{{ a }}{{=<% %>=}}<%b %><%={{ }}=%>{{ c }}`, map[string]string{"a": "a", "b": "b", "c": "c"}, "abc"},
				{`{{ a }}{{= <% %> =}}<%b %><%= {{ }}=%>{{c}}`, map[string]string{"a": "a", "b": "b", "c": "c"}, "abc"},

				// does not exist
				{`{{dne}}`, map[string]string{"name": "world"}, ""},
				{`{{dne}}`, User{"haoxin", 1}, ""},
				{`{{dne}}`, &User{"haoxin", 1}, ""},
				{`{{#has}}{{/has}}`, &User{"haoxin", 1}, ""},

				// section tests
				{`{{#A}}{{B}}{{/A}}`, Data{true, "hello"}, "hello"},
				{`{{#A}}{{{B}}}{{/A}}`, Data{true, "5 > 2"}, "5 > 2"},
				{`{{#A}}{{B}}{{/A}}`, Data{true, "5 > 2"}, "5 &gt; 2"},
				{`{{#A}}{{B}}{{/A}}`, Data{false, "hello"}, ""},
				{`{{a}}{{#b}}{{b}}{{/b}}{{c}}`, map[string]string{"a": "a", "b": "b", "c": "c"}, "abc"},
				{`{{#A}}{{B}}{{/A}}`, struct {
					A []struct {
						B string
					}
				}{[]struct {
					B string
				}{{"a"}, {"b"}, {"c"}}},
					"abc",
				},
				{`{{#A}}{{b}}{{/A}}`, struct{ A []map[string]string }{[]map[string]string{{"b": "a"}, {"b": "b"}, {"b": "c"}}}, "abc"},

				{`{{#users}}{{Name}}{{/users}}`, map[string]interface{}{"users": []User{{"haoxin", 1}}}, "haoxin"},

				{`{{#users}}gone{{Name}}{{/users}}`, map[string]interface{}{"users": nil}, ""},
				{`{{#users}}gone{{Name}}{{/users}}`, map[string]interface{}{"users": (*User)(nil)}, ""},
				{`{{#users}}gone{{Name}}{{/users}}`, map[string]interface{}{"users": []User{}}, ""},

				{`{{#users}}{{Name}}{{/users}}`, map[string]interface{}{"users": []*User{{"haoxin", 1}}}, "haoxin"},
				{`{{#users}}{{Name}}{{/users}}`, map[string]interface{}{"users": []interface{}{&User{"haoxin", 12}}}, "haoxin"},
				{`{{#users}}{{Name}}{{/users}}`, map[string]interface{}{"users": makeVector(1)}, "haoxin"},
				{`{{Name}}`, User{"haoxin", 1}, "haoxin"},
				{`{{Name}}`, &User{"haoxin", 1}, "haoxin"},
				{"{{#users}}\n{{Name}}\n{{/users}}", map[string]interface{}{"users": makeVector(2)}, "haoxin\nhaoxin\n"},
				{"{{#users}}\r\n{{Name}}\r\n{{/users}}", map[string]interface{}{"users": makeVector(2)}, "haoxin\r\nhaoxin\r\n"},

				// implicit iterator tests
				{`"{{#list}}({{.}}){{/list}}"`, map[string]interface{}{"list": []string{"a", "b", "c", "d", "e"}}, "\"(a)(b)(c)(d)(e)\""},
				{`"{{#list}}({{.}}){{/list}}"`, map[string]interface{}{"list": []int{1, 2, 3, 4, 5}}, "\"(1)(2)(3)(4)(5)\""},
				{`"{{#list}}({{.}}){{/list}}"`, map[string]interface{}{"list": []float64{1.10, 2.20, 3.30, 4.40, 5.50}}, "\"(1.1)(2.2)(3.3)(4.4)(5.5)\""},

				// inverted section tests
				{`{{a}}{{^b}}b{{/b}}{{c}}`, map[string]string{"a": "a", "c": "c"}, "abc"},
				{`{{a}}{{^b}}b{{/b}}{{c}}`, map[string]interface{}{"a": "a", "b": false, "c": "c"}, "abc"},
				{`{{^a}}b{{/a}}`, map[string]interface{}{"a": false}, "b"},
				{`{{^a}}b{{/a}}`, map[string]interface{}{"a": true}, ""},
				{`{{^a}}b{{/a}}`, map[string]interface{}{"a": "nonempty string"}, ""},
				{`{{^a}}b{{/a}}`, map[string]interface{}{"a": []string{}}, "b"},

				// function tests
				{`{{#users}}{{Func1}}{{/users}}`, map[string]interface{}{"users": []User{{"haoxin", 1}}}, "haoxin"},
				{`{{#users}}{{Func1}}{{/users}}`, map[string]interface{}{"users": []*User{{"haoxin", 1}}}, "haoxin"},
				{`{{#users}}{{Func2}}{{/users}}`, map[string]interface{}{"users": []*User{{"haoxin", 1}}}, "haoxin"},

				{`{{#users}}{{#Func3}}{{name}}{{/Func3}}{{/users}}`, map[string]interface{}{"users": []*User{{"haoxin", 1}}}, "haoxin"},
				{`{{#users}}{{#Func4}}{{name}}{{/Func4}}{{/users}}`, map[string]interface{}{"users": []*User{{"haoxin", 1}}}, ""},
				{`{{#Truefunc1}}abcd{{/Truefunc1}}`, User{"haoxin", 1}, "abcd"},
				{`{{#Truefunc1}}abcd{{/Truefunc1}}`, &User{"haoxin", 1}, "abcd"},
				{`{{#Truefunc2}}abcd{{/Truefunc2}}`, &User{"haoxin", 1}, "abcd"},
				{`{{#Func5}}{{#Allow}}abcd{{/Allow}}{{/Func5}}`, &User{"haoxin", 1}, "abcd"},
				{`{{#user}}{{#Func5}}{{#Allow}}abcd{{/Allow}}{{/Func5}}{{/user}}`, map[string]interface{}{"user": &User{"haoxin", 1}}, "abcd"},
				{`{{#user}}{{#Func6}}{{#Allow}}abcd{{/Allow}}{{/Func6}}{{/user}}`, map[string]interface{}{"user": &User{"haoxin", 1}}, "abcd"},

				// context chaining
				{`hello {{#section}}{{name}}{{/section}}`, map[string]interface{}{"section": map[string]string{"name": "world"}}, "hello world"},
				{`hello {{#section}}{{name}}{{/section}}`, map[string]interface{}{"name": "bob", "section": map[string]string{"name": "world"}}, "hello world"},
				{`hello {{#bool}}{{#section}}{{name}}{{/section}}{{/bool}}`, map[string]interface{}{"bool": true, "section": map[string]string{"name": "world"}}, "hello world"},
				{`{{#users}}{{canvas}}{{/users}}`, map[string]interface{}{"canvas": "hello", "users": []User{{"haoxin", 1}}}, "hello"},
				{`{{#categories}}{{DisplayName}}{{/categories}}`, map[string][]*Category{
					"categories": {&Category{"a", "b"}},
				}, "a - b"},

				// invalid syntax - https://github.com/hoisie/mustache/issues/10
				{`{{#a}}{{#b}}{{/a}}{{/b}}}`, map[string]interface{}{}, "line 1: interleaved closing tag: a"},

				// dotted names(dot notation)
				{`"{{person.name}}" == "{{#person}}{{name}}{{/person}}"`, map[string]interface{}{"person": map[string]string{"name": "hx"}}, `"hx" == "hx"`},
				{`"{{{person.name}}}" == "{{#person}}{{{name}}}{{/person}}"`, map[string]interface{}{"person": map[string]string{"name": "hx"}}, `"hx" == "hx"`},
				{`"{{a.b.c.d.e.name}}" == "Phil"`, map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": map[string]interface{}{"d": map[string]interface{}{"e": map[string]string{"name": "Phil"}}}}}}, `"Phil" == "Phil"`},
				{`"{{a.b.c}}" == ""`, map[string]interface{}{}, `"" == ""`},
				{`"{{a.b.c.name}}" == ""`, map[string]interface{}{"a": map[string]interface{}{"b": map[string]string{}}, "c": map[string]string{"name": "Jim"}}, `"" == ""`},
				{`"{{#a}}{{b.c.d.e.name}}{{/a}}" == "Phil"`, map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": map[string]interface{}{"d": map[string]interface{}{"e": map[string]string{"name": "Phil"}}}}}, "b": map[string]interface{}{"c": map[string]interface{}{"d": map[string]interface{}{"e": map[string]string{"name": "Wrong"}}}}}, `"Phil" == "Phil"`},
				{`{{#a}}{{b.c}}{{/a}}`, map[string]interface{}{"a": map[string]interface{}{"b": map[string]string{}}, "b": map[string]string{"c": "ERROR"}}, ""},
			}
			// basic
			for _, test := range tests {
				output := RenderString(test.tmpl, test.context)
				// debug("test: %+v, output: %s", test, output)
				So(output, ShouldEqual, test.expected)
			}
		})

		Convey("render file", func() {
			filename := path.Join(path.Join(os.Getenv("PWD"), "template"), "index.hbs")
			output := RenderFile(filename, map[string]interface{}{
				"title": "test",
				"body":  "<pre>hello world</pre>",
				"users": []User{{
					Name: "haoxin",
					Age:  2,
				}, {
					Name: "xin",
					Age:  1,
				}},
			})
			expected := readTestFile("index")
			So(output, ShouldEqual, expected)
		})

		Convey("partial", func() {
			filename := path.Join(path.Join(os.Getenv("PWD"), "template"), "use-partial.hbs")
			output := RenderFile(filename, map[string]string{
				"name": "haoxin",
			})
			expected := readTestFile("use-partial")
			So(output, ShouldEqual, expected)
		})

		Convey("multi context", func() {
			output1 := RenderString(`{{hello}} {{World}}`, map[string]string{"hello": "hello"}, struct{ World string }{"world"})
			output2 := RenderString(`{{hello}} {{World}}`, struct{ World string }{"world"}, map[string]string{"hello": "hello"})
			// debug("output1: %s, output2: %s", output1, output2)
			So(output1, ShouldEqual, "hello world")
			So(output2, ShouldEqual, "hello world")
		})

		Convey("invalid", func() {
			tests := []Test{
				{`{{#a}}{{}}{{/a}}`, Data{true, "hello"}, "empty tag"},
				{`{{}}`, nil, "empty tag"},
				{`{{}`, nil, "unmatched open tag"},
				{`{{`, nil, "unmatched open tag"},
			}

			for _, test := range tests {
				output := RenderString(test.tmpl, test.context)
				// debug("output: %s", output)
				So(output, ShouldContainSubstring, test.expected)
			}
		})
	})
}
