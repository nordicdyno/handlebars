package handlebars

import . "github.com/smartystreets/goconvey/convey"
import "testing"

func TestSnap(t *testing.T) {
	Convey("util", t, func() {
		Convey("resolve", func() {
			So(resolve("", "c"), ShouldEndWith, "handlebars/c")
			So(resolve("/a/b", "./c"), ShouldEqual, "/a/b/c")
			So(resolve("/a/b", "../c"), ShouldEqual, "/a/c")
			So(resolve("/a/b", "c"), ShouldEqual, "/a/b/c")
			So(resolve("/a/b", "/c"), ShouldEqual, "/c")
			So(resolve("", "/c"), ShouldEqual, "/c")
		})

		Convey("hash", func() {
			So(hash("hello world, see you!"), ShouldEqual, "ebb789525eb675a237a55c91e18dade9")
			So(hash("bye, see you!"), ShouldEqual, "119d9f0d6376efe9acd5d2547ca74fed")
		})
	})
}
