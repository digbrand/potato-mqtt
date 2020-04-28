package broker

//
//type topicManager struct {
//	subTree sync.Map
//}
//
//func newTopicManager() *topicManager {
//	return &topicManager{
//		subTree: sync.Map{},
//	}
//}
//
//type subTreeNode struct {
//	subList []*Client
//}
//
//func newNode() *subTreeNode {
//	return &subTreeNode{
//		subList: make([]*Client, 0),
//	}
//}
//
//func (s *subTreeNode) appendChild(c *Client) {
//	s.subList = append(s.subList, c)
//}
//
//func (s *topicManager) sub(topic string, c *Client) {
//	n, ok := s.subTree.Load(topic)
//	if !ok {
//		n = newNode()
//	}
//	n.(*subTreeNode).appendChild(c)
//	s.subTree.Store(topic, n)
//}
//
//func (s *topicManager) getSubs(topic string) []*Client {
//	n, ok := s.subTree.Load(topic)
//	if !ok {
//		return nil
//	}
//	return n.(*subTreeNode).subList
//}
