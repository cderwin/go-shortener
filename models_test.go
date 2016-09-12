package main

import (
	"gopkg.in/redis.v4"
	"reflect"
	"strconv"
	"testing"
	"time"
)

// Mock redis client

type MockClient struct {
	values map[string]string
	hashes map[string]map[string]string
}

func CreateMockStore() (RedisStore, MockClient) {
	client := CreateMockClient()
	clock := CreateMockClock()
	return RedisStore{client, clock}, client
}

func CreateMockClient() MockClient {
	valuesMap := map[string]string{"url:blah": "google.com", "url:ghjk": "lmgtfy.com", "url:foobar": "boo.baz"}
	hashesMap := make(map[string]map[string]string)
	hashesMap["hits:blah"] = map[string]string{"Total": "1117", "1": "78", "168": "34", "296": "672"}
	hashesMap["hits:ghjk"] = map[string]string{"Total": "387", "3": "31", "204": "14", "308": "76"}
	hashesMap["hits:foobar"] = map[string]string{"Total": "7", "86": "4", "287": "1", "365": "2"}
	return MockClient{values: valuesMap, hashes: hashesMap}
}

func (r MockClient) getKey(key string) (string, error) {
	value, present := r.values[key]
	if !present {
		return "", redis.Nil
	}

	return value, nil
}

func (r MockClient) setKey(key, value string) error {
	r.values[key] = value
	return nil
}

func (r MockClient) getHash(key string) (map[string]string, error) {
	value, present := r.hashes[key]
	if !present {
		return make(map[string]string), nil
	}

	return value, nil
}

func (r MockClient) hashExists(key string) (bool, error) {
	_, present := r.hashes[key]
	return present, nil
}

func (r MockClient) incrementHash(key, field string) error {
	mapp, present := r.hashes[key]
	if !present {
		r.hashes[key] = map[string]string{field: "1"}
		return nil
	}

	value, present := mapp[field]
	if !present {
		mapp[field] = "1"
		return nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	mapp[field] = strconv.Itoa(intValue + 1)
	return nil
}

// Actual tests

func TestGetURL(t *testing.T) {
	mockStore, _ := CreateMockStore()
	testValues := map[string]string{"blah": "google.com", "ghjk": "lmgtfy.com", "foobar": "boo.baz"}

	// test links defined in the mock
	for value, expected := range testValues {
		actual, err := mockStore.GetURL(value)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		if actual != expected {
			t.Errorf("Expected: %s\nActual: %s", expected, actual)
		}
	}

	// test bogus shortlink
	_, err := mockStore.GetURL("bazang")
	if err != NilValue {
		t.Errorf("Expected: %s\nActual: %s\n", NilValue.Error(), err.Error())
	}
}

func TestSaveURL(t *testing.T) {
	mockStore, mockClient := CreateMockStore()
	testMap := map[string]string{"reddit.com": "d23wrT", "news.ycombinator.com": "bB40lN", "github.com": "d1Ymny"}

	for longUrl, expectedShortURL := range testMap {
		actualShortURL, err := mockStore.SaveURL(longUrl)
		if err != nil {
			t.Errorf("Error occurred: %s\n", err.Error())
		}

		if actualShortURL != expectedShortURL {
			t.Errorf("Expected short url: %s\nActual short url: %s\n", expectedShortURL, actualShortURL)
		}

		if _, present := mockClient.values["url:"+expectedShortURL]; !present {
			t.Errorf("Key %s not added to hash", "url:"+expectedShortURL)
		}
	}
}

func DateFromDays(yearDays int, clock Clock) time.Time {
	t := time.Date(2016, time.January, yearDays, 0, 0, 0, 0, time.UTC)
	if t.After(clock.UTCNow()) {
		t = t.AddDate(-1, 0, 0)
	}
	return t
}

func TestGetHits(t *testing.T) {
	mockStore, _ := CreateMockStore()
	c := mockStore.Clock
	expectedMap := map[string]Hits{
		"blah":   Hits{Count: 1117, Days: map[time.Time]int{DateFromDays(1, c): 78, DateFromDays(168, c): 34, DateFromDays(296, c): 672}},
		"ghjk":   Hits{Count: 387, Days: map[time.Time]int{DateFromDays(3, c): 31, DateFromDays(204, c): 14, DateFromDays(308, c): 76}},
		"foobar": Hits{Count: 7, Days: map[time.Time]int{DateFromDays(86, c): 4, DateFromDays(287, c): 1, DateFromDays(365, c): 2}},
	}

	for shortURL, expectedHit := range expectedMap {
		actualHit, err := mockStore.GetHits(shortURL)
		if err != nil {
			t.Errorf("Error occurred: %s\n", err.Error())
		}

		if !reflect.DeepEqual(actualHit, expectedHit) {
			t.Errorf("Expected: %v\nActual: %v\n", expectedHit, actualHit)
		}
	}
}

func TestIncrementHits(t *testing.T) {
	mockStore, mockClient := CreateMockStore()
	expectedMap := map[string]map[string]string{
		"blah":   map[string]string{"Total": "1118", "1": "78", "168": "35", "296": "672"},
		"ghjk":   map[string]string{"Total": "388", "3": "31", "204": "14", "308": "76", "168": "1"},
		"foobar": map[string]string{"Total": "8", "86": "4", "287": "1", "365": "2", "168": "1"},
		"baz":    map[string]string{"Total": "1", "168": "1"},
	}

	for key, expectedValue := range expectedMap {
		mockStore.IncrementHits(key)
		actualValue := mockClient.hashes["hits:"+key]
		if !reflect.DeepEqual(actualValue, expectedValue) {
			t.Errorf("Expected hash value: %#v\nActual hash value: %#v\n", expectedValue, actualValue)
		}
	}
}
