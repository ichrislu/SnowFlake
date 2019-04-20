// Package snowflake provides a very simple Twitter snowflake generator and parser.
package snowflake

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1288834974657

	// NodeBits holds the number of bits to use for Node
	// Remember, you have a total 22 bits to share between Node/Step
	NodeBits uint8 = 10

	// StepBits holds the number of bits to use for Step
	// Remember, you have a total 22 bits to share between Node/Step
	StepBits uint8 = 12

	// DEPRECATED: the below four variables will be removed in a future release.
	mu        sync.Mutex
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// A Node struct holds the basic information needed for a snowflake generator
// node
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  time.Duration
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID int64

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewNode(node int64) (*Node, error) {

	// re-calc in case custom NodeBits or StepBits were set
	// DEPRECATED: the below block will be removed in a future release.
	mu.Lock()
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits
	mu.Unlock()

	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// Generate creates and returns a unique snowflake ID
func (n *Node) Generate() ID {

	n.mu.Lock()

	now := time.Since(n.epoch)

	// 解决跨毫秒时请求量不大，全生成偶数的问题
	// https://blog.csdn.net/wangzhanzheng/article/details/84937021
	n.step = (n.step + 1) & n.stepMask
	// if now-n.time < time.Millisecond {
	// 	n.step = (n.step + 1) & n.stepMask

	// 	if n.step == 0 {
	// 		for now-n.time < time.Millisecond {
	// 			now = time.Since(n.epoch)
	// 		}
	// 	}
	// } else {
	// 	n.step = 0
	// }

	n.time = now

	r := ID((now.Nanoseconds()/1000000)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step),
	)

	n.mu.Unlock()
	return r
}
