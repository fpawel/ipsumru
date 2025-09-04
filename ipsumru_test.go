package ipsumru

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSentenceGenerator(t *testing.T) {
	gen, err := NewSentenceGenerator()
	require.NoError(t, err)
	fmt.Println(gen.NextSentences(10))
}
