package topic

import (
	"fmt"
	"github.com/digbrand/potato-mqtt/topic/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestSubscribe(t *testing.T) {

	stree := newLocalTree()
	stree.Subscribe("a/b", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})
	stree.Subscribe("c/b/x", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("c/b/+/t", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("c/b/2/t", &SubscribeThing{
		Id: "client4",
	})

	fmt.Println(stree.tree)
}

func TestGetMatchs(t *testing.T) {

	stree := newLocalTree()
	stree.Subscribe("c/b", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})
	stree.Subscribe("c/b/x", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("c/b/+/t", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("c/b/2/t", &SubscribeThing{
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

	stree := newLocalTree()
	stree.Subscribe("$share/group1//a/b", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("$share/group1//a/b", &SubscribeThing{
		Id: "client1",
	})
	stree.Subscribe("$share/group1//a/+", &SubscribeThing{
		Id: "client5",
	})
	stree.Subscribe("$share/group1//a/b", &SubscribeThing{
		Id: "client4",
	})

	stree.Subscribe("$share/group2//a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("/a/b", &SubscribeThing{
		Id: "client2",
	})
	stree.Subscribe("/a/b", &SubscribeThing{
		Id: "client4",
	})
	stree.Subscribe("/a/#", &SubscribeThing{
		Id: "client3",
	})

	fmt.Println(stree.tree)

	mp, err := stree.GetSubscribers("/a/b")
	if err != nil {
		panic(err)
	}

	//client1,client5 will random published,but not exactly take turns,
	//because use random algorithm ,not one by one
	fmt.Println(mp)
}

func BenchmarkBuildTopicTree(b *testing.B) {

	rand.Seed(time.Now().UnixNano())
	st := newLocalTree()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tp := createRandomTopic(random.Letters, 12, 6)
		st.Subscribe(tp, &SubscribeThing{
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

func TestShareGroupName(t *testing.T) {
	rx := regexp.MustCompile(ShareGroupCompile)
	arr := rx.FindStringSubmatch("$share/group1//cc/cbbt")
	assert.Equal(t, arr[1], "group1")
	assert.Equal(t, arr[2], "/cc/cbbt")
	fmt.Println(arr[1], arr[2])

	arr = rx.FindStringSubmatch("$share/group1/")
	assert.Nil(t, arr)

}
