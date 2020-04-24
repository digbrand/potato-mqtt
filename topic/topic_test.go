package topic

//func TestGetMatchs(t *testing.T) {
//
//	tree := NewTree()
//	createTopicNode("a/b/#", "client1", tree)
//	createTopicNode("a/+/b/x", "client2", tree)
//	createTopicNode("/a/b", "client5", tree)
//
//	fmt.Println(tree.String())
//
//	arr:=make([]*TreeNode,0)
//
//	tt,_:= getTopicArray("/a/b")
//	getMatches(tt,tree.root.childs,&arr)
//
//	js,_:=json.Marshal(arr)
//	fmt.Println(string(js))
//}
//
//func getMatches(ts []string, childs cmap.ConcurrentMap,arr *[]*TreeNode){
//	for item:=range childs.IterBuffered(){
//		if item.Key==ts[0] || item.Key=="+"{
//			if len(ts)==1{
//				*arr=append(*arr,item.Val.(*TreeNode))
//			}else{
//				getMatches(ts[1:],item.Val.(*TreeNode).childs,arr)
//			}
//		}else if item.Key=="#"{
//			*arr=append(*arr,item.Val.(*TreeNode))
//		}
//	}
//}

//
//func BenchmarkTestTree(b *testing.B) {
//	tree := NewTree()
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		topic := randomdata.StringNumber(12, "/")
//		ts := make([]string, 0)
//		if err := recursionFindTopic(topic, &ts); err != nil {
//			log.Println(err)
//		}
//		tree.root.addNodes(ts, "client1", &SubscribeThing{
//			clientId: "client1",
//		})
//	}
//}

func (t *ConcurrentTree) GetMatchItems(topic string) ([]interface{}, error) {
	arr, err := getTopicArray(topic)
	if err != nil {
		return nil, err
	}
	if arr == nil || len(arr) == 0 {
		return nil, nil
	}

	//nodes := t.root.getChildrenWildCard(arr[0])
	//t.getMatchItems(arr[1:], &nodes)
	return nil, nil
}

func (t *ConcurrentTree) getMatchItems(arr []string, nodes *[]*TreeNode) {
	for _, a := range arr {
		for _, n := range *nodes {
			if n.Id == a {

			}
		}
	}
}
