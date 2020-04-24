package topic

import (
	"fmt"
	"testing"
)

func TestTreeString(t *testing.T) {
	tree := NewTree("root")

	node, _ := tree.root.addNode(NewTreeNode("first"))
	node.addNode(NewTreeNode("second"))
	fmt.Println(tree)
}
