package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyMatch3(t *testing.T) {
	type urlMatch struct {
		router   string
		urlMatch map[string]bool
	}
	urlMatches := []urlMatch{
		{
			router: "/user/{userId}",
			urlMatch: map[string]bool{
				"/user/1":                true,
				"/user/info":             true,
				"/user/info/address":     false,
				"/user1/info":            false,
				"/user/1/info":           false, // TODO:需要支持
				"/us":                    false,
				"/user?userId=1&name=aa": false,
			},
		},
		{
			router: "/user/*",
			urlMatch: map[string]bool{
				"/user/static":           true,
				"/user/info":             true,
				"/user/info/address":     true,
				"/user?userId=1&name=aa": false,
				"/user1/info":            false,
				"/us":                    false,
			},
		},
		{
			router: "/user",
			urlMatch: map[string]bool{
				"/user/static":           false,
				"/user/info":             false,
				"/user/info/address":     false,
				"/user?userId=1&name=aa": false,
				"/user":                  true,
				"/user1/info":            false,
				"/us":                    false,
			},
		},
	}

	for _, item := range urlMatches {
		for k, v := range item.urlMatch {
			result := KeyMatch3(k, item.router)
			assert.Equal(t, v, result, fmt.Sprintf("url:%v not match router:%v", k, item.router))
		}
	}
}
