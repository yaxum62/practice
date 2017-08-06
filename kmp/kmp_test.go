package kmp

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"time"
)

func TestCompile(t *testing.T) {
	pattern := Compile("")
	assert.Equal(t, "", pattern.pattern)
	assert.Equal(t, []int(nil), pattern.next)

	pattern = Compile("abcde")
	assert.Equal(t, "abcde", pattern.pattern)
	assert.Equal(t, []int{-1, 0, 0, 0, 0}, pattern.next)

	pattern = Compile("aabaabaaaa")
	assert.Equal(t, []int{-1, 0, 1, 0, 1, 2, 3, 4, 5, 2}, pattern.next)
}

func genString(rand *rand.Rand, size int) string {
	charSet := []byte{'a', 'b', 'c'}

	result := make([]byte, size)

	for idx := range result {
		result[idx] = charSet[rand.Intn(len(charSet))]
	}

	return string(result)
}

type testLongPattern string

func (p testLongPattern) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(testLongPattern(genString(rand, size)))
}

type testShortPattern string

func (p testShortPattern) Generate(rand *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(testShortPattern(genString(rand, 4)))
}

type testText string

func (p testText) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(testText(genString(rand, size)))
}

func TestCompile2(t *testing.T) {
	err := quick.Check(func(pattern testLongPattern) bool {
		pat := Compile(string(pattern))
		for i := 1; i < len(pat.next); i++ {
			sub := pattern[:i]
			k := pat.next[i]
			if !assert.Equal(t, sub[:k], sub[len(sub)-k:]) {
				return false
			}
			for bigger_k := k + 1; k < len(sub); k++ {
				if !assert.NotEqual(t, sub[:bigger_k], sub[len(sub)-bigger_k]) {
					return false
				}
			}
		}
		return true
	}, &quick.Config{
		MaxCountScale: 10,
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
	})
	assert.Nil(t, err)
}

func TestPattern_FindIn(t *testing.T) {
	pattern := Compile("")
	assert.Equal(t, 0, pattern.FindIn("abc"))
	assert.Equal(t, 0, pattern.FindIn(""))

	pattern = Compile("abcdabd")
	assert.Equal(t, 15, pattern.FindIn("abc abcdab abcdabcdabde"))
	assert.Equal(t, -1, pattern.FindIn("abc abcdab abcdabcdabcde"))
}

func TestPattern_FindIn2(t *testing.T) {
	err := quick.Check(func(pattern testShortPattern, text testText) bool {
		return assert.Equal(t,
			strings.Index(string(text), string(pattern)),
			Compile(string(pattern)).FindIn(string(text)),
		)
	}, &quick.Config{
		MaxCountScale: 10,
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
	})
	assert.Nil(t, err)
}
