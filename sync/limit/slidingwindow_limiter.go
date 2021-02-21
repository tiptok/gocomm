//copy from github.com/RussellLuo/slidingwindow
//redis store

package limit

import (
	"sync"
	"time"
)

type (
	SlidingWidowLimiter struct {
		size  time.Duration
		limit int64

		mu sync.Mutex

		curr Window
		prev Window
	}
	Window interface {
		// Start returns the start boundary.
		Start() time.Time

		// Count returns the accumulated count.
		Count() int64

		// AddCount increments the accumulated count by n.
		AddCount(n int64)

		// Reset sets the state of the window with the given settings.
		Reset(s time.Time, c int64)

		// Sync tries to exchange data between the window and the central
		// datastore at time now, to keep the window's count up-to-date.
		Sync(now time.Time)
	}
	// StopFunc stops the window's sync behaviour.
	StopFunc func()

	// NewWindow creates a new window, and returns a function to stop
	// the possible sync behaviour within it.
	NewWindow func() (Window, StopFunc)
)

// NewLimiter creates a new limiter, and returns a function to stop
// the possible sync behaviour within the current window.
func NewSlidingWidowLimiter(size time.Duration, limit int64, newWindow NewWindow) (*SlidingWidowLimiter, StopFunc) {
	currWin, currStop := newWindow()

	// The previous window is static (i.e. no add changes will happen within it),
	// so we always create it as an instance of LocalWindow.
	//
	// In this way, the whole limiter, despite containing two windows, now only
	// consumes at most one goroutine for the possible sync behaviour within
	// the current window.
	prevWin, _ := NewLocalWindow()

	lim := &SlidingWidowLimiter{
		size:  size,
		limit: limit,
		curr:  currWin,
		prev:  prevWin,
	}

	return lim, currStop
}

// Size returns the time duration of one window size. Note that the size
// is defined to be read-only, if you need to change the size,
// create a new limiter with a new size instead.
func (lim *SlidingWidowLimiter) Size() time.Duration {
	return lim.size
}

// Limit returns the maximum events permitted to happen during one window size.
func (lim *SlidingWidowLimiter) Limit() int64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

// SetLimit sets a new Limit for the limiter.
func (lim *SlidingWidowLimiter) SetLimit(newLimit int64) {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	lim.limit = newLimit
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *SlidingWidowLimiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowN reports whether n events may happen at time now.
func (lim *SlidingWidowLimiter) AllowN(now time.Time, n int64) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	lim.advance(now)

	elapsed := now.Sub(lim.curr.Start())
	weight := float64(lim.size-elapsed) / float64(lim.size)
	count := int64(weight*float64(lim.prev.Count())) + lim.curr.Count()

	// Trigger the possible sync behaviour.
	defer lim.curr.Sync(now)

	if count+n > lim.limit {
		return false
	}

	lim.curr.AddCount(n)
	return true
}

// advance updates the current/previous windows resulting from the passage of time.
func (lim *SlidingWidowLimiter) advance(now time.Time) {
	// Calculate the start boundary of the expected current-window.
	newCurrStart := now.Truncate(lim.size)

	diffSize := newCurrStart.Sub(lim.curr.Start()) / lim.size
	if diffSize >= 1 {
		// The current-window is at least one-window-size behind the expected one.

		newPrevCount := int64(0)
		if diffSize == 1 {
			// The new previous-window will overlap with the old current-window,
			// so it inherits the count.
			//
			// Note that the count here may be not accurate, since it is only a
			// SNAPSHOT of the current-window's count, which in itself tends to
			// be inaccurate due to the asynchronous nature of the sync behaviour.
			newPrevCount = lim.curr.Count()
		}
		lim.prev.Reset(newCurrStart.Add(-lim.size), newPrevCount)

		// The new current-window always has zero count.
		lim.curr.Reset(newCurrStart, 0)
	}
}

// LocalWindow represents a window that ignores sync behavior entirely
// and only stores counters in memory.
type LocalWindow struct {
	// The start boundary (timestamp in nanoseconds) of the window.
	// [start, start + size)
	start int64

	// The total count of events happened in the window.
	count int64
}

func NewLocalWindow() (*LocalWindow, StopFunc) {
	return &LocalWindow{}, func() {}
}

func (w *LocalWindow) Start() time.Time {
	return time.Unix(0, w.start)
}

func (w *LocalWindow) Count() int64 {
	return w.count
}

func (w *LocalWindow) AddCount(n int64) {
	w.count += n
}

func (w *LocalWindow) Reset(s time.Time, c int64) {
	w.start = s.UnixNano()
	w.count = c
}

func (w *LocalWindow) Sync(now time.Time) {}
