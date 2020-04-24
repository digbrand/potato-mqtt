package random

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandomString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		fmt.Println(RandomStringCustom(1, []rune("abcd")))
	}
}

func TestRandN(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		fmt.Println(rand.Intn(100))
	}
}

func TestRandom(t *testing.T) {
	fmt.Println(RandomString(12))
}

func BenchmarkRandomString(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RandomStringCustom(4, Letters)
	}

}
