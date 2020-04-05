package limit

import (
	"time"
	"github.com/go-redis/redis"
)

type Limiter struct {
	rc *redis.Client
	Limit int64
	Period time.Duration
}

func NewLimiter(redisURL string, limit int64, period time.Duration) (*Limiter, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	rc := redis.NewClient(options)

	if err := rc.Ping().Err(); err != nil {
		return nil, err
	}

	return &Limiter{rc: rc, Limit: limit, Period: period}, nil
}

func (l *Limiter) GetLimit() int64{
	return l.Limit
}

func (l *Limiter) GetPeriod() time.Duration{
	return l.Period
}

func (l *Limiter) GetRequestNum(key string) int64{
	return l.rc.LLen(key).Val()
}

func (l *Limiter) GetPipelineTTL(key string) time.Duration {
	return l.rc.TTL(key).Val()
}

func (l *Limiter) Allow(key string, limit int64, duration time.Duration) bool {
	current := l.rc.LLen(key).Val()
	if current >= limit {
		return false
	}

	if v := l.rc.Exists(key).Val(); v == 0 {
		pipe := l.rc.TxPipeline()
		pipe.RPush(key, key)
		pipe.Expire(key, duration)
		_, _ = pipe.Exec()
	} else {
		l.rc.RPushX(key, key)
	}

	return true
}