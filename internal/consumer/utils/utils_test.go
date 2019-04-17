package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestViperConfig(t *testing.T) {
	Convey("Given yaml config ", t, func() {
		Convey("When config is read", func() {
			v := GetConfig()
			Convey("Values must be extracted", func() {
				So(v.GetString("sequelize2.database"), ShouldEqual, "hiveos_eth")
				So(v.GetString("influx.database"), ShouldEqual, "minerdash")
			})
		})
	})
}

func TestParseJSON(t *testing.T) {
	var str string
	var res map[string]interface{}
	Convey("Given json string ", t, func() {
		str = `{"num":6.13,"str":"aaa"}`
		Convey("When string is parsed", func() {
			res = ParseJSON(str,false)
			Convey("Values must be extracted", func() {
				So(res["num"].(float64), ShouldEqual, 6.13)
				So(res["str"].(string), ShouldEqual, "aaa")
			})

		})
	})
}

func TestParseTimestampToUnix(t *testing.T) {
	Convey("Given timestamp string ", t, func() {
		str:= "2019-01-10T19:09:13.754410804"
		Convey("When string is parsed", func() {
			res, _ := ParseTimestampToUnix(str)
			Convey("Values must be equal to ", func() {
				So(res, ShouldEqual, 1547118793754410804)
			})
		})
	})
}