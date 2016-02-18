package service

import "io"

type Service interface {
	PostText(key, text string, cacheID int)
	MessageFromJSON(reader io.Reader) Message
	ServiceMonitor() *Monitor
}

type Message interface {
	GroupID() string
	UserName() string
	UserID() string
	Text() string
	UserType() string
}

type Monitor interface {
	ValueFor(cachedID int) int
}

var serviceMap = map[string]Service{}

func NewService(name string) *Service {
	serv, ok := serviceMap[name]

	if ok {
		return &serv
	} else {
		return nil
	}
}

func AddService(name string, serv Service) {
	serviceMap[name] = serv
}
