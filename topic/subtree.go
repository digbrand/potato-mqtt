package topic

import (
	"errors"
	"github.com/digbrand/potato-mqtt/broker"
	cmap "github.com/orcaman/concurrent-map"
	"sync"
)

const (
	_GroupTopicRegexp = `^\$share/([0-9a-zA-Z_-]+)/(.*)$`
)

/**
struct a subscribe tree,store subscribe client and path
*/
type SubscribeTree struct {
	tree *ConcurrentTree
	lock sync.RWMutex
}

type SubscribeThing struct {
	Id        string
	client    *broker.Client
	share     bool   //是否共享订阅
	groupName string //订阅组名称
}

func newSubscribeTree() *SubscribeTree {
	return &SubscribeTree{
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
func (st *SubscribeTree) addNodes(ts []string, dataKey string, dataVal *SubscribeThing) {
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

func (st *SubscribeTree) GetSubscribers(topic string) (map[string]*SubscribeThing, error) {
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
	return mp, nil
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
client subscribe topic
*/
func (st *SubscribeTree) subscribe(topic string, thing *SubscribeThing) error {
	arr, err := getTopicArray(topic)
	if err != nil {
		return err
	}
	st.addNodes(arr, thing.Id, thing)
	return nil
}
