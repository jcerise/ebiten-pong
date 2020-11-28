package main

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRandomRangeGenerator(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	number := rangeNegative(-4, 4)
	assert.Equal(t, 1, number)
}
