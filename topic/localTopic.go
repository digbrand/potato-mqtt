package topic

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	ShareGroupCompile = `^\$share/([0-9a-zA-Z_-]+)/(.+)$`
)

var (
	groupRegExp = regexp.MustCompile(ShareGroupCompile)
)

/**
struct a Subscribe tree,store Subscribe Client and path
*/
type LocalTree struct {
	tree *ConcurrentTree
	lock sync.RWMutex
}

func newLocalTree() *LocalTree {
	return &LocalTree{
		tree: NewTree("root"),
		lock: sync.RWMutex{},
	}
}

/**
convert a topic to string array
*/
func getTopicArray(topic string) ([]string, error) {
	ts := make([]string, 0)
	if err := recursionFindTopic(topic, &ts); err != nil {
		return nil, err
	}
	return ts, nil
}

/**
recursion find topic,until topic string reach last character
*/
func recursionFindTopic(topic string, result *[]string) error {
	node, other, err := findFirstTopicPart(topic)
	if err != nil {
		return err
	}
	*result = append(*result, node)
	if other != "" {
		return recursionFindTopic(other, result)
	}
	return nil
}

/**
given a string,find first mathing topic part
for example:
fmt.println(findFirstTopicPart("/a/b/c"))
output:
a,b/c,nil
*/
func findFirstTopicPart(str string) (node string, other string, err error) {
	for i, s := range str {
		switch s {
		case '/':
			if i == 0 {
				return "^", str[1:], nil
			}
			if i == len(str)-1 {
				return str, str, errors.New("/ must not last character")
			}
			return str[0:i], str[i+1:], nil
		case '#':
			if i != len(str)-1 {
				return str, str, errors.New("# wildcard must be last character")
			}
			if i != 0 {
				return str, str, errors.New("# wildcard previous character must be /")
			}
			return "#", "", nil
		case '+':
			if i != 0 {
				return str, str, errors.New("# wildcard previous character must be /")
			}
			if i != len(str)-1 {
				if str[i+1] != '/' {
					return str, str, errors.New("+ wildcard next character must be /")
				} else {
					return "+", str[i+2:], nil
				}
			}
			return "+", "", nil
		default:
			if i == len(str)-1 {
				return str, "", nil
			}
		}
	}
	return "", str, nil
}

/**
give a topic array,tree will add all nodes
*/
func (st *LocalTree) AddNodes(ts []string, dataKey string, dataVal *SubscribeThing) {
	st.lock.Lock()
	defer st.lock.Unlock()

	current := st.tree.root
	for _, v := range ts {
		if current.hasChild(v) {
			current = current.getChild(v)
			continue
		}
		newNode := NewTreeNode(v)
		current.addNode(newNode)
		current = newNode
	}
	current.addDataItem(dataKey, dataVal)
}

func (st *LocalTree) GetSubscribers(topic string) (map[string]*SubscribeThing, error) {
	arr, err := getTopicArray(topic)
	if err != nil {
		return nil, err
	}

	nodes := make([]*TreeNode, 0)
	getMatches(arr, st.tree.root.childs, &nodes)

	mp := make(map[string]*SubscribeThing)

	for _, v := range nodes {
		if v.Data.IsEmpty() {
			continue
		}
		for item := range v.Data.IterBuffered() {
			mp[item.Key] = item.Val.(*SubscribeThing)
		}
	}

	if mp != nil {
		return clearSubscribers(mp), nil
	}
	return nil, nil
}

func getMatches(ts []string, childs cmap.ConcurrentMap, arr *[]*TreeNode) {
	for item := range childs.IterBuffered() {
		if item.Key == ts[0] || item.Key == "+" {
			if len(ts) == 1 {
				*arr = append(*arr, item.Val.(*TreeNode))
			} else {
				getMatches(ts[1:], item.Val.(*TreeNode).childs, arr)
			}
		} else if item.Key == "#" {
			*arr = append(*arr, item.Val.(*TreeNode))
		}
	}
}

/**
Client Subscribe topic
*/
func (st *LocalTree) Subscribe(topic string, thing *SubscribeThing) error {
	if strings.HasPrefix(topic, "$share/") {
		parts := groupRegExp.FindStringSubmatch(topic)
		if parts == nil {
			return errors.New("topic is share group,but format failure")
		}
		thing.share = true
		thing.groupName = parts[1]
		thing.Id = thing.Id + "|" + thing.groupName
		topic = parts[2]
	}

	arr, err := getTopicArray(topic)
	if err != nil {
		return err
	}

	st.AddNodes(arr, thing.Id, thing)
	return nil
}

/**
handle subscribers share group,and handle repetition thingId,prevent send message
twice or more
*/
func clearSubscribers(mp map[string]*SubscribeThing) map[string]*SubscribeThing {
	rx := regexp.MustCompile("^(.+)\\|(.+)$")
	grp := make(map[string]map[string]*SubscribeThing)
	normal := make(map[string]*SubscribeThing)

	for k, v := range mp {
		arr := rx.FindStringSubmatch(k)
		//this is share group thing
		if arr != nil {
			a, ok := grp[arr[2]]
			if !ok {
				a = make(map[string]*SubscribeThing)
				grp[arr[2]] = a
			}
			a[arr[1]] = v
			continue
		}
		normal[k] = v
	}
	//if thing exists share and normal same time,remove from share
	for _, v := range grp {
		for x, _ := range v {
			_, ok := normal[x]
			if ok {
				delete(v, x)
			}
		}
	}

	//remove null group and merge group thing to normal thing
	for k, v := range grp {
		if len(v) == 0 {
			delete(grp, k)
			continue
		}
		ln := len(v)
		rand.Seed(time.Now().UnixNano())
		t, o := getMapValueByIndex(v, rand.Intn(ln))
		if o != nil {
			normal[t] = o
		}
	}

	return normal
}

func getMapValueByIndex(mp map[string]*SubscribeThing, index int) (string, *SubscribeThing) {
	i := 0
	for k, v := range mp {
		if i == index {
			return k, v
		}
		i++
	}
	return "", nil
}
