package cron

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	// 2020-07-06T07:58:34+00:00
	// or 2020-07-06T15:58:34+08:00
	var sec int64 = 1594022314
	baseTime := time.Unix(sec, 0)
	rand.Seed(time.Now().UnixNano())

	t1 := Trigger{
		Every: Period{
			Every: 1,
			Unit:  Day,
		},
		At: Moment{
			Hour: 23,
		},
	}
	t1Base := baseTime
	fmt.Println("t1: every day at 23:00")
	for i := 0; i < 10; i++ {
		nextBase := t1.nextWakeup(t1Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t1Base) < 0 ||
			nextBase.Sub(t1Base) > 24*time.Hour {
			t.Error("bad wakeup time", nextBase)
		}
		if nextBase.Hour() != 23 {
			t.Error("bad hour", nextBase)
		}
		t1Base = nextBase
	}

	t2 := Trigger{
		Every: Period{
			Every: 1,
			Unit:  Hour,
		},
		At: Moment{
			Minute: 58,
		},
	}
	fmt.Println("t2: every hour at :58")
	t2Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t2.nextWakeup(t2Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t2Base) < 0 ||
			nextBase.Sub(t2Base) > time.Hour {
			t.Error("bad wakeup time", nextBase)
		}
		if nextBase.Minute() != 58 {
			t.Error("bad minute", nextBase)
		}
		t2Base = nextBase
	}

	t3 := Trigger{
		Every: Period{
			Every: 17,
			Unit:  Minute,
		},
		At: Moment{
			Any: true,
		},
	}
	fmt.Println("t3: every 17 min, at any moment")
	t3Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t3.nextWakeup(t3Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t3Base) < 0 ||
			nextBase.Sub(t3Base) >= 18*time.Minute {
			t.Error("bad wakeup time", nextBase)
		}
		if i != 0 && (60+nextBase.Minute()-t3Base.Minute())%60 != 17 {
			t.Error("bad minute", nextBase)
		}
		t3Base = nextBase
	}

	t4 := Trigger{
		Every: Period{
			Every: 7,
			Unit:  Second,
		},
		At: Moment{},
	}
	fmt.Println("t4: every 7 sec")
	t4Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t4.nextWakeup(t4Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t4Base) < 0 ||
			nextBase.Sub(t4Base) > 7*time.Second {
			t.Error("bad wakeup time", nextBase)
		}
		if i != 0 && (60+nextBase.Second()-t4Base.Second())%60 != 7 {
			t.Error("bad second", nextBase)
		}
		t4Base = nextBase
	}

	t5 := Trigger{
		Every: Period{
			Every: 1,
			Unit:  Day,
		},
		At: Moment{
			Any: true,
		},
	}
	fmt.Println("t5: every day, at any moment")
	t5Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t5.nextWakeup(t5Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t5Base) < 0 ||
			nextBase.Sub(t5Base) >= 36*time.Hour {
			t.Error("bad wakeup time", nextBase)
		}
		if nextBase.Day()-t5Base.Day() != 1 {
			t.Error("bad day", nextBase)
		}
		t5Base = nextBase
	}

	t6 := Trigger{
		Every: Period{
			Every: 1,
			Unit:  Second,
		},
		At: Moment{},
	}
	fmt.Println("t6: every sec")
	t6Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t6.nextWakeup(t6Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t6Base) < 0 ||
			nextBase.Sub(t6Base) > time.Second {
			t.Error("bad wakeup time", nextBase)
		}
		if nextBase.Second()-t6Base.Second() != 1 {
			t.Error("bad sec", nextBase)
		}
		t6Base = nextBase
	}

	t7 := Trigger{
		Every: Period{
			Every: 3600,
			Unit:  Second,
		},
	}
	fmt.Println("t7: every 3600 sec")
	t7Base := baseTime
	for i := 0; i < 10; i++ {
		nextBase := t7.nextWakeup(t7Base)
		fmt.Println(nextBase, nextBase.Unix())
		if nextBase.Sub(t7Base) > 3600*time.Second {
			t.Error("bad wakeup time", nextBase)
		}
		t7Base = nextBase
	}
}
