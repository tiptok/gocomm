package limit

import (
	"sync"
	"time"
)

type (
	CounterLimiter struct {
		rate    int           //计数周期内最多允许的请求数
		cycle   time.Duration //计数周期
		couters sync.Map
	}
	counterLimitItem struct {
		key   string
		rate  int           //计数周期内最多允许的请求数
		begin time.Time     //计数开始时间
		cycle time.Duration //计数周期
		count int           //计数周期内累计收到的请求数
	}
)

const (
	Unknown = iota
	Allowed
	HitQuota
	OverQuota
)

func NewCounterLimitItem(key string, rate int, cycle time.Duration) *counterLimitItem {
	return &counterLimitItem{
		key:   key,
		rate:  rate,
		cycle: cycle,
		begin: time.Now(),
		count: 0,
	}
}
func (l *counterLimitItem) Take() (int, error) {
	if l.count == l.rate-1 {
		now := time.Now()
		if now.Sub(l.begin) >= l.cycle {
			//速度允许范围内， 重置计数器
			l.Reset(now)
			return Allowed, nil
		} else {
			return OverQuota, nil
		}
	} else {
		//没有达到速率限制，计数加1
		l.count++
		return Allowed, nil
	}
	return Allowed, nil
}
func (l *counterLimitItem) Set(r int, cycle time.Duration) {
	l.rate = r
	l.begin = time.Now()
	l.cycle = cycle
	l.count = 0
}
func (l *counterLimitItem) Reset(t time.Time) {
	l.begin = t
	l.count = 0
}

// NewCounterLimiter  实例化一个计数器制器
// defaultRate   默认速率   eg:10
// defaultCycle  默认周期   eg:1 second
// counterItem   自定义限速项
func NewCounterLimiter(defaultRate int, defaultCycle time.Duration, counterItem ...*counterLimitItem) *CounterLimiter {
	limit := new(CounterLimiter)
	limit.rate = defaultRate
	limit.cycle = defaultCycle
	for i := range counterItem {
		item := counterItem[i]
		limit.couters.Store(item.key, item)
	}
	return limit
}

// Allow  从键值的限速器获取权限
// true：Allowed
func (l *CounterLimiter) Allow(key string) bool {
	if value, ok := l.couters.Load(key); ok {
		code, _ := value.(*counterLimitItem).Take()
		return code == Allowed
	}
	l.couters.Store(key, NewCounterLimitItem(key, l.rate, l.cycle))
	return true
}

// Set	设置限速器 速率、周期
func (l *CounterLimiter) Set(key string, r int, cycle time.Duration) {
	if value, ok := l.couters.Load(key); ok {
		value.(*counterLimitItem).Set(r, cycle)
		return
	}
	l.couters.Store(key, NewCounterLimitItem(key, r, cycle))
}
