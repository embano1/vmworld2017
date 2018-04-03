package autoscaler

import (
	"net"
	"time"

	"github.com/streadway/amqp"
)

const timeout = 5 * time.Second

// RabbitMQ initialization
func (a *AutoScaler) dial(broker string) (*amqp.Channel, error) {
	// amqp dial config with shorter timeout
	config := amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}
	conn, err := amqp.DialConfig(broker, config)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}
