package urlcache

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type urlCache[T interface{}] struct {
	expirationTime time.Time
	cache          T
}

func (c urlCache[T]) Expired() bool {
	return c.expired(time.Now())
}

func (c urlCache[T]) expired(now time.Time) bool {
	return c.expirationTime.Before(now)
}

type Client[T interface{}] struct {
	client        *http.Client
	urlCacheTTL   time.Duration
	urlCacheMap   map[string]urlCache[T]
	urlCacheMutex sync.Mutex
}

type ClientOption[T interface{}] func(*Client[T])

func NewClient[T interface{}](client *http.Client, opts ...ClientOption[T]) *Client[T] {
	d := &Client[T]{
		client:        client,
		urlCacheTTL:   2 * time.Minute,
		urlCacheMap:   make(map[string]urlCache[T]),
		urlCacheMutex: sync.Mutex{},
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithCacheTTL[T interface{}](ttl time.Duration) ClientOption[T] {
	return func(d *Client[T]) {
		d.urlCacheTTL = ttl
	}
}

var (
	ErrNotSupportedExceptGet = errors.New("httpz: not supported except GET")
	ErrNotCacheable          = errors.New("httpz: not cacheable")
)

func (uc *Client[T]) Get(url string, cacheable func(resp *http.Response) bool, responseConvert func(resp *http.Response) (T, error)) (t T, err error) { //nolint:ireturn
	return uc.do(url, func() (*http.Response, error) { return http.Get(url) }, cacheable, responseConvert, time.Now()) //nolint:gosec,wrapcheck
}

func (uc *Client[T]) Do(req *http.Request, cacheable func(resp *http.Response) bool, responseConvert func(resp *http.Response) (T, error)) (t T, err error) { //nolint:ireturn
	var zero T

	if req.Method != http.MethodGet {
		return zero, ErrNotSupportedExceptGet
	}

	return uc.do(req.URL.String(), func() (*http.Response, error) { return uc.client.Do(req) }, cacheable, responseConvert, time.Now()) //nolint:wrapcheck
}

func (uc *Client[T]) do(url string, f func() (*http.Response, error), cacheable func(resp *http.Response) bool, responseConvert func(resp *http.Response) (T, error), now time.Time) (T, error) { //nolint:ireturn
	uc.urlCacheMutex.Lock()
	defer uc.urlCacheMutex.Unlock()
	if !uc.urlCacheMap[url].Expired() {
		return uc.urlCacheMap[url].cache, nil
	}

	var zero T

	resp, err := f() //nolint:noctx
	if err != nil {
		return zero, err //nolint:wrapcheck
	}
	defer resp.Body.Close()

	v, err := responseConvert(resp)
	if err != nil {
		return zero, err //nolint:wrapcheck
	}

	if !cacheable(resp) {
		return zero, ErrNotCacheable
	}

	uc.urlCacheMap[url] = urlCache[T]{
		cache:          v,
		expirationTime: now.Add(uc.urlCacheTTL),
	}

	return v, nil
}
