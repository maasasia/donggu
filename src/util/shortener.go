package util

import "fmt"

type Shortener interface {
	Shorten(key string) string
}

type PassthroughShortener struct{}

func (s PassthroughShortener) Shorten(key string) string {
	return key
}

type CountingShortener struct {
	keys map[string]uint32
}

func (s *CountingShortener) Shorten(key string) string {
	if _, ok := s.keys[key]; !ok {
		s.keys[key] = uint32(len(s.keys) + 1)
	}
	return fmt.Sprintf("%d", s.keys[key])
}

func NewCountingShortener() *CountingShortener {
	return &CountingShortener{keys: map[string]uint32{}}
}
