package leakybucket

import (
	"sync"
	"time"
)

// Bucket contains the counters and timers to implement your rate limit.
type Bucket struct {

	// Capacity of the bucket is how many total spots are available in the bucket.
	// It is the burst limit. A capacity of 0 will reject all Add()s
	Capacity int

	// DrainAmount is the number of spots that open up each drain period. It is
	// the volume portion of the throughput limit. A drain amount of 0 means the
	// bucket doesn't leak, and once capacity is reached, will reject all
	// subsequent events.
	DrainAmount int

	// DrainPeriod is the frequency the bucket drains. It is the time portion of
	// the throughput limit. A drain period of 0 means always drain - every Add()
	// also drains, effectively meaning no throughput limit (so long as Capacity
	// and DrainAmount are greater than zero).
	DrainPeriod time.Duration

	level     int        // level is the current
	lock      sync.Mutex // used for making it all threadsafe
	lastDrain time.Time  // the time the bucket was last drained. used for calculating whether to drain

	testNower nower // to make testing easier
}

// BucketOverflow error will be returned when you fail to Add() to a bucket
// because the rate limit has been reached.
type BucketOverflow struct{}

// Error implements the error interface
func (b BucketOverflow) Error() string {
	return "Bucket Overflow"
}

// Add increment the counter in the bucket. If the bucket overflows, the rate
// limit has been exceeded and Add() will return a BucketOverflow error.
func (b *Bucket) Add() error {
	b.lock.Lock()
	defer b.lock.Unlock()
	durSinceDrain := b.now().Sub(b.lastDrain)
	if durSinceDrain >= b.DrainPeriod || durSinceDrain < 0 {
		b.drain()
	}
	if b.level < b.Capacity {
		b.level++
		return nil
	}
	return BucketOverflow{}
}

// drain reduces the bucket level as much as it needs to be reduced
// if multiple drain periods have elapsed since last drain, we'll drain
// multiple times, flooring the result to 0
// the caller should be holding a lock to protect b.level and b.lastDrain
func (b *Bucket) drain() {
	timeSince := b.now().Sub(b.lastDrain)
	if timeSince < 0 {
		timeSince = 0
	}
	// how many periods has it been since your last drain?
	var numTimes int
	if b.DrainPeriod == 0 {
		// drain period of 0 means always drain, aka no limit.
		numTimes = 1
	} else {
		numTimes = int(timeSince / b.DrainPeriod)
	}
	drainAmount := numTimes * b.DrainAmount
	b.level -= drainAmount
	if b.level < 0 {
		b.level = 0
	}
	b.lastDrain = b.now()
}

func (b *Bucket) now() time.Time {
	if b.testNower != nil {
		return b.testNower.Now()
	}
	return time.Now()
}

// makes testing easier
type nower interface {
	Now() time.Time
}
