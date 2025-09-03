package ipsumru

import (
	"fmt"
	"testing"
)

func TestNewSentenceGenerator(t *testing.T) {
	gen := NewSentenceGenerator()
	fmt.Println(gen.NextSentences(10))
}
