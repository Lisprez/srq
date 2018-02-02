package srq_test

import (
	"srq"
	"testing"

	. "gopkg.in/check.v1"
)

type SrqSuite struct{}

var _ = Suite(&SrqSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *SrqSuite) TestConnectionRedis(c *C) {
	conn, err := srq.NewConnection("tcp", "127.0.0.1:6379")
	c.Assert(err, Equals, nil)
	queue := srq.NewQueue("test_queue", conn)
	c.Assert(queue.Name, Equals, "test_queue")
	defer queue.Clear()
}

func (s *SrqSuite) TestMessagePush(c *C) {
	conn, err := srq.NewConnection("tcp", "127.0.0.1:6379")
	c.Assert(err, Equals, nil)
	queue := srq.NewQueue("test_queue", conn)
	c.Assert(queue.Name, Equals, "test_queue")
	defer queue.Clear()
	status, err := queue.Push("test text 1", 1)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 2", 4)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 3", 3)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 4", 2)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	v1, err := queue.Pop()
	c.Assert(err, Equals, nil)
	c.Assert(v1, Equals, "test text 1")

	v2, err := queue.Pop()
	c.Assert(err, Equals, nil)
	c.Assert(v2, Equals, "test text 4")

	v3, err := queue.Pop()
	c.Assert(err, Equals, nil)
	c.Assert(v3, Equals, "test text 3")

	v4, err := queue.Pop()
	c.Assert(err, Equals, nil)
	c.Assert(v4, Equals, "test text 2")
}

func (s *SrqSuite) TestMessagesPopOrder(c *C) {
	conn, err := srq.NewConnection("tcp", "127.0.0.1:6379")
	c.Assert(err, Equals, nil)
	queue := srq.NewQueue("test_queue", conn)
	c.Assert(queue.Name, Equals, "test_queue")
	defer queue.Clear()
	status, err := queue.Push("test text 1", 1)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 2", 4)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 3", 3)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	status, err = queue.Push("test text 4", 2)
	c.Assert(err, Equals, nil)
	c.Assert(status, Equals, true)

	messages, err := queue.PopMessages(4)
	c.Assert(err, Equals, nil)
	c.Assert(len(messages), Equals, 4)
	c.Assert(messages[0], Equals, "test text 1")
	c.Assert(messages[1], Equals, "test text 4")
	c.Assert(messages[2], Equals, "test text 3")
	c.Assert(messages[3], Equals, "test text 2")
}
