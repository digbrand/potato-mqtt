package topic

import (
	"fmt"
	"github.com/digbrand/potato-mqtt/topic/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {

	stree := newSubscribeTree()
	stree.subscribe("a/b", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})
	stree.subscribe("c/b/x", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("c/b/+/t", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("c/b/2/t", &SubscribeThing{
		Id: "client4",
	})

	fmt.Println(stree.tree)
}

func TestGetMatchs(t *testing.T) {

	stree := newSubscribeTree()
	stree.subscribe("c/b", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})
	stree.subscribe("c/b/x", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("c/b/+/t", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("c/b/2/t", &SubscribeThing{
		Id: "client4",
	})

	mp, err := stree.GetSubscribers("/a/b")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(mp), 2)
	assert.NotNil(t, mp["client2"], mp["client3"])

	mp, _ = stree.GetSubscribers("c/b/4/t")
	assert.Equal(t, len(mp), 1)

}

func TestShareGroup(t *testing.T) {
	stree := newSubscribeTree()
	stree.subscribe("c/b", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})
	stree.subscribe("c/b/x", &SubscribeThing{
		Id: "client1",
	})
	stree.subscribe("c/b/+/t", &SubscribeThing{
		Id: "client2",
	})
	stree.subscribe("c/b/2/t", &SubscribeThing{
		Id: "client4",
	})

	mp, err := stree.GetSubscribers("/a/b")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, len(mp), 2)
	assert.NotNil(t, mp["client2"], mp["client3"])

	mp, _ = stree.GetSubscribers("c/b/4/t")
	assert.Equal(t, len(mp), 1)

}

func BenchmarkBuildTopicTree(b *testing.B) {

	rand.Seed(time.Now().UnixNano())
	st := newSubscribeTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tp := createRandomTopic(random.Letters, 12, 6)
		st.subscribe(tp, &SubscribeThing{
			Id: random.RandomStringCustom(4, []rune("abcd")),
		})
	}
}

func TestRandomTopic(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 5; i++ {
		tp := createRandomTopic(random.Letters, 12, 6)
		fmt.Println(tp)
	}

}

func createRandomTopic(letters []rune, wordLength int, splitNumber int) string {
	w := rand.Intn(wordLength) + 1
	s := rand.Intn(splitNumber) + 1
	result := ""
	for i := 0; i < s; i++ {
		result += random.RandomStringCustom(w, letters) + "/"
	}
	return result[:len(result)-1]
}
