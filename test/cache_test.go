package main

import (
	"fmt"
	"time"
	"testing"
	"dep/cache"
)

func Test_Add_Get(t *testing.T) {
	const interval = time.Second * 5

	cases := []struct{
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			ch := cache.New_Cache(interval)
			ch.Add(c.key, c.val)
			
			val, ok := ch.Get(c.key)
			if !ok {
				t.Errorf("Didn't find the expected key: %s", c.key)
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("The values aren't matching: '%s' expected, got '%s'", string(c.val), string(val))
				return
			}
		})
	}
}

func Test_Reap_Loop(t *testing.T) {
	const base_time = 5*time.Millisecond
	const wait_time = base_time + 5*time.Millisecond

	ch := cache.New_Cache(base_time)
	ch.Add("https://example.com", []byte("testdata"))

	_, ok := ch.Get("https://example.com")
	if !ok {
		t.Errorf("Didn't find the expected key: %s", "https://example.com")
		return
	}

	time.Sleep(wait_time)

	_, ok_two := ch.Get("https://example.com")
	if ok_two {
		t.Errorf("Expected NOT to find the key: %s", "https://example.com")
		return
	}
}
