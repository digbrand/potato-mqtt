package topic

type Manager interface {
	AddNodes(ts []string, dataKey string, dataVal *SubscribeThing)
	GetSubscribers(topic string) (map[string]*SubscribeThing, error)
	Subscribe(topic string, thing *SubscribeThing) error
}

type SubscribeThing struct {
	Id        string
	Client    interface{} //store client pointer
	share     bool        //is share group subscribe
	groupName string      //share group name
}

func NewSubscribeThing(id string, client interface{}) *SubscribeThing {
	return &SubscribeThing{
		Id:     id,
		Client: client,
	}
}

func NewManager() Manager {
	return newLocalTree()
}
