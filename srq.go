package srq

import (
	"math"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func init() {
	fetcherMessageScript = redis.NewScript(3, `
local name = KEYS[1]
local waitWeightLimit = KEYS[2]
local limit = KEYS[3]
local results = redis.call('zrangebyscore', name, '-inf', waitWeightLimit, 'LIMIT', 0, limit)
if table.getn(results) > 0 then
redis.call('zrem', name, unpack(results))
end
return results
  `)
}

var fetcherMessageScript *redis.Script

type Queue struct {
	Name string
	c    redis.Conn
}

func NewConnection(network string, address string, options ...redis.DialOption) (redis.Conn, error) {
	return redis.Dial(network, address, options...)
}

func NewQueue(queueName string, conn redis.Conn) *Queue {
	return &Queue{c: conn, Name: queueName}
}

func (queue *Queue) Push(message string, waitWeight int64) (bool, error) {
	return queue.enqueue(message, waitWeight)
}

func (queue *Queue) Pop() (string, error) {
	messages, err := queue.PopMessages(1)
	if err != nil {
		return "", err
	}

	if len(messages) == 0 {
		return "", nil
	}

	return messages[0], nil
}

func (queue *Queue) Clear() error {
	_, err := queue.c.Do("DEL", queue.Name)
	return err
}

func (queue *Queue) Length() (int64, error) {
	return redis.Int64(queue.c.Do("ZCARD", queue.Name))
}

func (queue *Queue) Close() error {
	return queue.c.Close()
}

func (queue *Queue) enqueue(message string, waitWeight int64) (bool, error) {
	return redis.Bool(queue.c.Do("ZADD", queue.Name, waitWeight, message))
}

func (queue *Queue) PopMessages(limit int) ([]string, error) {
	return redis.Strings(fetcherMessageScript.Do(queue.c, queue.Name, math.MaxInt64, strconv.Itoa(limit)))
}
