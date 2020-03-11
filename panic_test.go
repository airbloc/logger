package logger

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWrapRecover(t *testing.T) {
	convey.Convey("When Recover is called", t, func() {
		convey.Convey("It should not panic", func() {
			convey.So(func() {
				defer func() {
					if err := WrapRecover(recover()); err != nil {
						fmt.Println(err.Pretty())
					}
				}()
				panicStation()
			}, convey.ShouldNotPanic)
		})
	})
}

func panicStation() {
	var empty map[string]string
	empty["a"] = "b"
}
