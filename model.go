package cron

import (
	"math/rand"
	"time"
)

type TimeUnit int

const (
	Day    TimeUnit = 1
	Hour   TimeUnit = 2
	Minute TimeUnit = 3
	Second TimeUnit = 4
)

func (u TimeUnit) Seconds() int64 {
	switch u {
	case Day:
		return 24 * 60 * 60
	case Hour:
		return 60 * 60
	case Minute:
		return 60
	case Second:
		return 1
	default:
		return 0
	}
}

type Period struct {
	Every int
	Unit  TimeUnit
}

// If `Any` is true, other conditions are ignored;
// time condition with units larger than period unit are ignored,
// e.g. if trigger condition is "every minute", only `Second` is considered
type Moment struct {
	Any    bool
	Hour   int
	Minute int
	Second int
}

// Trigger condition for cron jobs, like "every day at 12:00" etc.
type Trigger struct {
	Every Period
	At    Moment
	// if last trigger has not finished, new trigger won't fire
	OneAfterAnother bool
}

func (t Trigger) nextWakeup(now time.Time) time.Time {
	nowUnix := now.Unix()
	period := int64(t.Every.Every) * t.Every.Unit.Seconds()
	var offset, anyOffset int64
	switch t.Every.Unit {
	case Day:
		if t.At.Any {
			anyOffset = rand.Int63n(Day.Seconds() / 2)
		} else {
			offset = int64(t.At.Hour)*Hour.Seconds() +
				int64(t.At.Minute)*Minute.Seconds() +
				int64(t.At.Second)
		}
	case Hour:
		if t.At.Any {
			anyOffset = rand.Int63n(Hour.Seconds() / 2)
		} else {
			offset = int64(t.At.Minute)*Minute.Seconds() +
				int64(t.At.Second)
		}
	case Minute:
		if t.At.Any {
			anyOffset = rand.Int63n(Minute.Seconds() / 2)
		} else {
			offset = int64(t.At.Second)
		}
	default:
	}

	if t.Every.Unit == Day {
		_, zoneOffset := now.Zone()
		offset -= int64(zoneOffset)
	}
	var i, next int64
	for i = -1; i <= 2; i++ {
		next = (nowUnix/period+i)*period + offset
		if next > nowUnix {
			break
		}
	}
	return time.Unix(next+anyOffset, 0)
}

type Job struct {
	Name     string
	Trigger  Trigger
	Callback func()
}
