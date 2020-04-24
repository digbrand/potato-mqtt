package topic

import (
	"github.com/disiqueira/gotree"
	cmap "github.com/orcaman/concurrent-map"
)

/**
thread safe tree,depends on "github.com/orcaman/concurrent-map"
*/
type ConcurrentTree struct {
	root *TreeNode
}

/**
build a thread safe tree,root as root node clientId
*/
func NewTree(root string) *ConcurrentTree {
	return &ConcurrentTree{
		root: NewTreeNode("root"),
	}
}

func (t *ConcurrentTree) String() string {
	r := gotree.New(t.root.Id)
	t.stringChild(t.root, r)

	return r.Print()
}

func (t *ConcurrentTree) stringChild(n *TreeNode, r gotree.Tree) {
	for item := range n.childs.IterBuffered() {
		pn := r.Add(item.Key)

		nd := item.Val.(*TreeNode)
		if !nd.Data.IsEmpty() {
			for it := range nd.Data.IterBuffered() {
				pn.Add(it.Key)
			}
		}
		if nd.childs != nil {
			t.stringChild(nd, pn)
		}
	}
}

type TreeNode struct {
	Id     string
	childs cmap.ConcurrentMap
	Data   cmap.ConcurrentMap
}

func NewTreeNode(id string) *TreeNode {
	return &TreeNode{
		Id:     id,
		childs: cmap.New(),
		Data:   cmap.New(),
	}
}

func (n *TreeNode) hasChild(id string) bool {
	return n.childs.Has(id)
}

func (n *TreeNode) getChild(id string) *TreeNode {
	result, ok := n.childs.Get(id)
	if ok {
		return result.(*TreeNode)
	}
	return nil
}

func (n *TreeNode) addNode(node *TreeNode) (nd *TreeNode, success bool) {
	if n.hasChild(node.Id) {
		return node, false
	}
	n.childs.Set(node.Id, node)
	return node, true
}

func (n *TreeNode) addDataItem(dataKey string, dataVal interface{}) {
	if n.Data.Has(dataKey) {
		return
	}
	n.Data.Set(dataKey, dataVal)
}
