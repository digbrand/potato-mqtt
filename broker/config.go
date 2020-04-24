package broker

type Config struct {
	TcpPort int
	Routes  []*Route
}

type Route struct {
	IP   string
	Port int
}
