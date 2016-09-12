package main

import (
	"errors"
	"gopkg.in/redis.v4"
	"strconv"
	"time"
)

type Datastore interface {
	GetURL(string) (string, error)
	SaveURL(string) (string, error)
	GetHits(string) (Hits, error)
	IncrementHits(string) error
}

type Redis interface {
	getHash(string) (map[string]string, error)
	incrementHash(string, string) error
	hashExists(string) (bool, error)
	getKey(string) (string, error)
	setKey(string, string) error
}

// Direct database access methods, allows for testability of business logic

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(url string) RedisClient {
	client := redis.NewClient(&redis.Options{Addr: url})
	return RedisClient{client}
}

func (r RedisClient) getHash(key string) (map[string]string, error) {
	return r.HGetAll(key).Result()
}

func (r RedisClient) incrementHash(key, field string) error {
	return r.HIncrBy(key, field, 1).Err()
}

func (r RedisClient) hashExists(key string) (bool, error) {
	len, err := r.HLen(key).Result()
	if len > 0 {
		return true, err
	}

	return false, err
}

func (r RedisClient) getKey(key string) (string, error) {
	return r.Get(key).Result()
}

func (r RedisClient) setKey(key, value string) error {
	return r.Set(key, value, 0).Err()
}

// Business logic methods, this is where the fun starts

var NilValue = errors.New("Nil value returned")

type RedisStore struct {
	Redis
	Clock
}

func NewRedisStore(url string) RedisStore {
	redisClient := NewRedisClient(url)
	clock := NewSystemClock()
	return RedisStore{redisClient, clock}
}

func (r RedisStore) GetURL(short_url string) (string, error) {
	key := "url:" + short_url
	value, err := r.getKey(key)
	if err == redis.Nil {
		return "", NilValue
	}

	return value, err
}

func (r RedisStore) SaveURL(long_url string) (string, error) {
	short_url := hashUrl(long_url)
	key := "url:" + short_url
	err := r.setKey(key, long_url)
	if err != nil {
		return "", err
	}
	return short_url, nil
}

//
// Hits -- stats about an endpoint
//

type Hits struct {
	Count int
	Days  map[time.Time]int
}

func NewHits() Hits {
	var h Hits
	h.Days = make(map[time.Time]int)
	return h
}

func (r RedisStore) GetHits(short_url string) (Hits, error) {
	key := "hits:" + short_url
	exists, err := r.hashExists(key)
	if err != nil {
		return NewHits(), err
	}

	if !exists {
		return NewHits(), NilValue
	}

	hits_map, err := r.getHash(key)
	if err != nil {
		return NewHits(), err
	}

	result := NewHits()

	total := hits_map["Total"]
	if total == "" {
		total = "0"
	}
	result.Count, err = strconv.Atoi(total)
	if err != nil {
		return NewHits(), err
	}
	delete(hits_map, "Total")

	for str_day, str_hits := range hits_map {
		day, err := strconv.Atoi(str_day)
		if err != nil {
			return NewHits(), err
		}
		hits, err := strconv.Atoi(str_hits)
		if err != nil {
			return NewHits(), err
		}

		year := r.UTCNow().Year()
		date := time.Date(year, time.January, day, 0, 0, 0, 0, time.UTC)
		if date.After(r.UTCNow()) {
			date = date.AddDate(-1, 0, 0)
		}

		result.Days[date] = hits
	}

	return result, nil
}

func (r RedisStore) IncrementHits(short_url string) error {
	key := "hits:" + short_url
	err := r.incrementHash(key, "Total")
	if err != nil {
		return err
	}

	days := r.UTCNow().YearDay()
	return r.incrementHash(key, strconv.Itoa(days))
}
