package leakybucket

import (
	"sync"
	"testing"
	"time"
)

func TestUnitializedBucket(t *testing.T) {
	b := Bucket{}
	if err := b.Add(); err == nil {
		t.Errorf("An uninitialized bucket should have a capacity of zero and always error. Got nil instead.")
	}
}

func TestNoDrainAmount(t *testing.T) {
	b := Bucket{
		Capacity: 10,
		// DrainAmount: 0, // no drain amount means never drain aka this bucket doesn't leak
	}
	nower := newFakeNower()
	b.testNower = nower
	var successes, failures int
	for i := 0; i < 20; i++ {
		if err := b.Add(); err == nil {
			successes++
		} else {
			failures++
		}
	}
	// we should have hit the capacity of 10 and then rejected everything else.
	if successes != 10 || failures != 10 {
		t.Errorf("expected 10 successse and 10 failures. Instead got %d successes and %d failures.", successes, failures)
	}
	// let some time go by in an attempt to drain the bucket
	nower.incrementSec(5)

	// add in another 20, they should all fail
	for i := 0; i < 20; i++ {
		if err := b.Add(); err == nil {
			successes++
		} else {
			failures++
		}
	}
	// we should have hit the capacity of 10 and then rejected everything else.
	if successes != 10 || failures != 30 {
		t.Errorf("expected 10 successse and 30 failures. Instead got %d successes and %d failures.", successes, failures)
	}
}

func TestNoDrainPeriod(t *testing.T) {
	b := Bucket{
		Capacity:    1,
		DrainAmount: 1,
		// DrainPeriod: 0, // no drain period means drain on every Add()
	}
	for i := 0; i < 20; i++ {
		if err := b.Add(); err != nil {
			t.Errorf("Zero DrainPeriod should always be empty. Got error on iteration %d: %s", i, err)
		}
	}
}

func TestBasicBucket(t *testing.T) {
	b := Bucket{
		Capacity:    10,
		DrainAmount: 2,
		DrainPeriod: 3 * time.Second,
	}
	nower := newFakeNower()
	b.testNower = nower
	var err error
	var i int
	for i = 0; i < 20; i++ {
		err = b.Add()
		if err != nil {
			break
		}
	}
	if i != 10 {
		t.Errorf("expected to hit bucket capacity at 10, instead hit at %d. err = %s", i, err)
	}
	// make 1 second pass
	nower.incrementSec(1)
	err = b.Add()
	if err == nil {
		t.Errorf("should have errored; it's not yet time to drain the bucket")
	}
	// make 3 more seconds pass - we should now have 2 additional slots in the bucket
	nower.incrementSec(3)
	if err = b.Add(); err != nil {
		t.Errorf("should not have yet hit capacity: %s", err)
	}
	if err = b.Add(); err != nil {
		t.Errorf("should not have yet hit capacity: %s", err)
	}
	if err = b.Add(); err == nil {
		t.Errorf("should have yet hit capacity but didn't")
	}
	// make 9 more seconds pass - we should now have 6 additional slots in the bucket
	nower.incrementSec(9)
	for i = 0; i < 6; i++ {
		err = b.Add()
	}
	if err = b.Add(); err == nil {
		t.Errorf("should have yet hit capacity but didn't")
	}
	// make 20 more seconds pass - we should now have all ten slots
	nower.incrementSec(20)
	for i = 0; i < 10; i++ {
		err = b.Add()
	}
	if err = b.Add(); err == nil {
		t.Errorf("should have yet hit capacity but didn't")
	}
}

// test for raciness in Add()
func TestAddRace(t *testing.T) {
	b := Bucket{
		Capacity:    10,
		DrainAmount: 5,
		DrainPeriod: 3 * time.Second,
	}
	nower := newFakeNower()
	b.testNower = nower
	var wg sync.WaitGroup
	// fill the bucket
	for i := 0; i < 10; i++ {
		b.Add()
	}
	// wait 4 seconds so it'll be time to drain
	nower.incrementSec(4)
	// try and add eight more in parallel (drain amount, + attempted extras)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			b.Add()
			wg.Done()
		}()
	}
	wg.Wait()
	// bucket should once again be full
	if b.level != 10 {
		t.Errorf("race detected. level should be 10 and is %d", b.level)
	}

}

// for easy time manipulation during tests
type fakeNower struct {
	now time.Time
}

func (f *fakeNower) Now() time.Time {
	return f.now
}

func (f *fakeNower) incrementSec(numSec int) {
	f.now = f.now.Add(time.Duration(numSec) * time.Second)
}

func newFakeNower() *fakeNower {
	return &fakeNower{
		now: time.Date(2010, time.June, 21, 15, 4, 5, 0, time.UTC),
	}
}

func ExampleBucket(t *testing.T) {
	// Example: allow no more than one request every 10 seconds.
	_ = Bucket{
		Capacity:    1,
		DrainAmount: 1,
		DrainPeriod: 10 * time.Second,
	}
	// Example: Allow a steady stream of 100 requests per minute, with occasional
	// bursts up to 500 in one minute
	_ = Bucket{
		Capacity:    500,
		DrainAmount: 100,
		DrainPeriod: time.Minute,
	}

}
