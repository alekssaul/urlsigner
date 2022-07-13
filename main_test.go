package main

import (
	"encoding/base64"
	"testing"
)

func TestStringToUnix(t *testing.T) {
	_, err := stringToUnix("foo")
	if err == nil {
		t.Errorf("Function stringToUnix did not error when called with \"foo\" ")
	}

	// 1657680270 - Wednesday, July 13, 2022 2:44:30 AM
	a, err := stringToUnix("1657680270")
	if err != nil {
		t.Errorf("Function stringToUnix errored when called with \"1657680270\" ")
	}

	if a.Unix() != 1657680270 {
		t.Errorf("Function stringToUnix errored when called with \"1657680270\" ")
	}
}

func TestSignUrl(t *testing.T) {
	pKey := "AlzCNB4ySeKBhaBqKR2497AQFGBZlYZNoN9vK7lf4ZwKdxf6siUE8oAwuOQ7Rtf_oj2-E4qgMcE0MQ3M9y1xpA"
	binpKey, _ := base64.RawURLEncoding.DecodeString(pKey)

	c := AppConfig{
		KeySet:        "foo",
		PrivateKey:    pKey,
		binPrivateKey: binpKey,
	}

	expiration, _ := stringToUnix("1657680270")
	signedurl := signUrl(c, "media.m3u8", expiration)
	expectedresult := "media.m3u8?Expires=1657680270&KeyName=foo&Signature=by6D9kj2rq51WsuO-GAnwLQmB7lPtnzq3-Pq2S8BDOZdaSQ1bxV3tFtH4n5n4FynevsbtFVIIHaAqc55WKenCw"

	if signedurl != expectedresult {
		t.Errorf("Fed Static test values, however output for signURL function was not correct")
	}

}
