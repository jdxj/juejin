package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var client = newHTTPClient()

func newHTTPClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := new(http.Client)
	client.Jar = jar
	return client
}

func SetCookies(domain *url.URL, cookies []*http.Cookie) {
	client.Jar.SetCookies(domain, cookies)
}

func Do(req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

func Get(url string) (*http.Response, error) {
	time.Sleep(500 * time.Millisecond)
	return client.Get(url)
}

// todo: 实现
func ParseCookies() []*http.Cookie {
	return nil
}
