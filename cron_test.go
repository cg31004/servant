package servant

import (
	"testing"
	"time"
)

var num1 int

const (
	answer1 = 3
)

func cronFunc1(ctx *Context) {
	time.Sleep(3 * time.Second)
	num1 = answer1
}

func TestCron1(t *testing.T) {

	c := New()
	if _, err := c.AddScheduleFunc(Every(1*time.Second), cronFunc1); err != nil {
		t.Fatal(err)
	}
	c.Start()

	time.Sleep(3 * time.Second)

	c.Stop()
	if num1 != answer1 {
		t.Log(num1)
		t.Fatal("wrong answer")
	}
}

const (
	answer2 = 1
	mw1Key  = "1"
	mw2Key  = "2"
	mw3Key  = "3"
)

var (
	num2 int
	mw   map[string]int
)

func mwFunc1(ctx *Context) {
	if val, ok := mw[mw1Key]; ok {
		ctx.Set(mw1Key, val)
	}
}
func mwFunc2(ctx *Context) {
	if val, ok := mw[mw2Key]; ok {
		ctx.Set(mw2Key, val)
	}
}
func mwFunc3(ctx *Context) {
	if val, ok := mw[mw3Key]; ok {
		ctx.Set(mw3Key, val)
	}
}
func cronFunc2(ctx *Context) {
	num2 = (ctx.Value(mw1Key).(int) + ctx.Value(mw2Key).(int)) / ctx.Value(mw3Key).(int)
}

func TestCron2(t *testing.T) {
	c := New()
	mw = map[string]int{
		mw1Key: 1,
		mw2Key: 4,
		mw3Key: 5,
	}
	c.Use(mwFunc3, mwFunc2, mwFunc1)
	if _, err := c.AddScheduleFunc(Every(1*time.Second), cronFunc2); err != nil {
		t.Fatal(err)
	}
	c.Start()

	time.Sleep(3 * time.Second)

	c.Stop()
	if num2 != answer2 {
		t.Log(num2)
		t.Fatal("wrong answer")
	}
}
