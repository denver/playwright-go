package playwright

import (
	"fmt"
	"log"
	"reflect"
)

type Channel struct {
	EventEmitter
	guid       string
	connection *Connection
	object     interface{}
}

func (c *Channel) Send(method string, options ...interface{}) (interface{}, error) {
	params := transformOptions(options...)
	result, err := c.connection.SendMessageToServer(c.guid, method, params)
	if err != nil {
		return nil, fmt.Errorf("could not send message to server: %w", err)
	}
	if result == nil {
		return nil, nil
	}
	if reflect.TypeOf(result).Kind() == reflect.Map {
		mapV := result.(map[string]interface{})
		if len(mapV) == 0 {
			return nil, nil
		}
		for key := range mapV {
			return mapV[key], nil
		}
	}
	return result, nil
}

func (c *Channel) SendNoReply(method string, options ...interface{}) {
	params := transformOptions(options...)
	_, err := c.connection.SendMessageToServer(c.guid, method, params)
	if err != nil {
		log.Printf("could not send message to server from noreply: %v", err)
	}
}

func newChannel(connection *Connection, guid string) *Channel {
	channel := &Channel{
		connection: connection,
		guid:       guid,
	}
	channel.initEventEmitter()
	return channel
}
