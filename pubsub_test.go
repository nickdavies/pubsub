package pubsub

import (
	check "launchpad.net/gocheck"
	"runtime"
	"testing"
	"time"
)

var _ = check.Suite(new(Suite))

func Test(t *testing.T) {
	check.TestingT(t)
}

type Suite struct{}

func (s *Suite) TestSub(c *check.C) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t1")
	ch3 := ps.Sub("t2")

	ps.Pub("hi", "t1")
	c.Check(<-ch1, check.Equals, "hi")
	c.Check(<-ch2, check.Equals, "hi")

	ps.Pub("hello", "t2")
	c.Check(<-ch3, check.Equals, "hello")

	ps.Shutdown()
	_, ok := <-ch1
	c.Check(ok, check.Equals, false)
	_, ok = <-ch2
	c.Check(ok, check.Equals, false)
	_, ok = <-ch3
	c.Check(ok, check.Equals, false)
}

func (s *Suite) TestSubOnce(c *check.C) {
	ps := New(1)
	ch := ps.SubOnce("t1")

	ps.Pub("hi", "t1")
	c.Check(<-ch, check.Equals, "hi")

	_, ok := <-ch
	c.Check(ok, check.Equals, false)
	ps.Shutdown()
}

func (s *Suite) TestUnsub(c *check.C) {
	ps := New(1)
	ch := ps.Sub("t1")

	ps.Pub("hi", "t1")
	c.Check(<-ch, check.Equals, "hi")

	ps.Unsub(ch, "t1")
	_, ok := <-ch
	c.Check(ok, check.Equals, false)
	ps.Shutdown()
}

func (s *Suite) TestClose(c *check.C) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t1")
	ch3 := ps.Sub("t2")
	ch4 := ps.Sub("t3")

	ps.Pub("hi", "t1")
	ps.Pub("hello", "t2")
	c.Check(<-ch1, check.Equals, "hi")
	c.Check(<-ch2, check.Equals, "hi")
	c.Check(<-ch3, check.Equals, "hello")

	ps.Close("t1", "t2")
	_, ok := <-ch1
	c.Check(ok, check.Equals, false)
	_, ok = <-ch2
	c.Check(ok, check.Equals, false)
	_, ok = <-ch3
	c.Check(ok, check.Equals, false)

	ps.Pub("welcome", "t3")
	c.Check(<-ch4, check.Equals, "welcome")

	ps.Shutdown()
}

func (s *Suite) TestShutdown(c *check.C) {
	start := runtime.NumGoroutine()
	New(10).Shutdown()
	time.Sleep(1)
	c.Check(runtime.NumGoroutine()-start, check.Equals, 1)
}

func (s *Suite) TestMultiSub(c *check.C) {
	ps := New(1)
	ch := ps.Sub("t1", "t2")

	ps.Pub("hi", "t1")
	c.Check(<-ch, check.Equals, "hi")

	ps.Pub("hello", "t2")
	c.Check(<-ch, check.Equals, "hello")

	ps.Shutdown()
	_, ok := <-ch
	c.Check(ok, check.Equals, false)
}

func (s *Suite) TestMultiSubOnce(c *check.C) {
	ps := New(1)
	ch := ps.SubOnce("t1", "t2")

	ps.Pub("hi", "t1")
	c.Check(<-ch, check.Equals, "hi")

	ps.Pub("hello", "t2")

	_, ok := <-ch
	c.Check(ok, check.Equals, false)
	ps.Shutdown()
}

func (s *Suite) TestMultiPub(c *check.C) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t2")

	ps.Pub("hi", "t1", "t2")
	c.Check(<-ch1, check.Equals, "hi")
	c.Check(<-ch2, check.Equals, "hi")

	ps.Shutdown()
}

func (s *Suite) TestMultiUnsub(c *check.C) {
	ps := New(1)
	ch := ps.Sub("t1", "t2", "t3")

	ps.Unsub(ch, "t1")

	ps.Pub("hi", "t1")

	ps.Pub("hello", "t2")
	c.Check(<-ch, check.Equals, "hello")

	ps.Unsub(ch, "t2", "t3")
	_, ok := <-ch
	c.Check(ok, check.Equals, false)

	ps.Shutdown()
}

func (s *Suite) TestMultiClose(c *check.C) {
	ps := New(1)
	ch := ps.Sub("t1", "t2")

	ps.Pub("hi", "t1")
	c.Check(<-ch, check.Equals, "hi")

	ps.Close("t1")
	ps.Pub("hello", "t2")
	c.Check(<-ch, check.Equals, "hello")

	ps.Close("t2")
	_, ok := <-ch
	c.Check(ok, check.Equals, false)

	ps.Shutdown()
}
