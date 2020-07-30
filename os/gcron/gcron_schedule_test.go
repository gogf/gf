package gcron

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/test/gtest"
)

func Test_cronSchedule_meet(t *testing.T) {
	const TIMEFORMAT = "20060102"
	type ts struct {
		T string
		M bool
	}
	testdata := map[string][]ts{
		"0 0 0 L * ?": {
			{"20200731", true},
			{"20200229", true},
			{"20200701", false},
			{"20200830", false},
			{"20200630", true},
		},
		"0 0 0 LW * ?": {
			{"20200229", false},
			{"20200228", true},
			{"20200731", true},
			{"20200829", false},
			{"20200831", true},
			{"20200830", false},
			{"20200630", true},
		},
		"0 0 0 4W * ?": {
			{"20200704", false},
			{"20200703", true},
			{"20200705", false},
			{"20200706", false},
			{"20200903", false},
			{"20200904", true},
			{"20200905", false},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		cron := New()
		for p, data := range testdata {
			entry, err := cron.Add(p, func() {

			})
			t.Assert(err, nil)
			for _, v := range data {
				t1, err := time.Parse(TIMEFORMAT, v.T)
				t.Assert(err, nil)
				test := entry.schedule.meet(t1)
				if test != v.M {
					t.Fatal(fmt.Sprintf("EXPECT %s %v == %v", v.T, test, v.M))
				}
			}
		}
	})

}
