package leakybucket

import (
	"fmt"
	"time"
)

func Example() {
	// Create a new bucket for a rate limit of 10/sec, with a burst limit of 50
	b := Bucket{
		Capacity:    50,
		DrainAmount: 10,
		DrainPeriod: time.Second,
	}
	// override time for the sake of the example. ignore the next 2 lines.
	n := newFakeNower()
	b.testNower = n

	// count successes and failures
	var successes, failures int

	// Adding the first 50 entries will succeed, the 51st will fail.
	// This is us hitting the burst limit.
	for i := 0; i < 51; i++ {
		if err := b.Add(); err == nil {
			successes++
		} else {
			failures++
		}
	}
	fmt.Printf("Successes: %d, Failures: %d\n", successes, failures)

	// "wait" 1 second, and we should be able to add 10 more entries to the bucket
	n.incrementSec(1)

	// Add another 11 - the first 10 will succeed, the 11th will fail.
	// This shows the rate limit constraining throughput
	for i := 0; i < 11; i++ {
		if err := b.Add(); err == nil {
			successes++
		} else {
			failures++
		}
	}
	fmt.Printf("Successes: %d, Failures: %d\n", successes, failures)

	// Output:
	// Successes: 50, Failures: 1
	// Successes: 60, Failures: 2
}
