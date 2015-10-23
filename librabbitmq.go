// Copyright 2015 ZeroStack, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// librabbitmq adds some helper functions for handling Openstack.
// It also adds client-side library structures and
// functions used to connect to RabbitMQ and listen to notifications.

package libmq

import (
  "fmt"

  "github.com/streadway/amqp"
)

// RabbitMQClientConfig encapsulates the parameters required to create a
// RabbitMQ client.
type RabbitMQClientConfig struct {
  Exchange     string `toml:"exchange"`
  ExchangeType string `toml:"exchange_type"`
  Queue        string `toml:"queue"`
  BindingKey   string `toml:"binding_key"`
  ConsumerTag  string `toml:"consumer_tag"`
}

// RabbitMQClient represents that parameters and state of a RabbitMQ client.
type RabbitMQClient struct {
  uri      string
  cfg      *RabbitMQClientConfig
  conn     *amqp.Connection
  amqpChan *amqp.Channel
  mqCh     <-chan amqp.Delivery
  msgCh    chan []byte
  quitCh   chan struct{}
}

// RabbitMQFactory implements the MQClientFactory interface.
type RabbitMQFactory struct {
  clients []*RabbitMQClient
}

// NewRabbitMQFactory creates a RabbitMQFactory object.
func NewRabbitMQFactory() *RabbitMQFactory {
  return &RabbitMQFactory{}
}

// NewClient creates a new RabbitMQClient object with the given config.
func (f *RabbitMQFactory) NewClient(
  uri string, cfgInf interface{}) (MQClient, error) {

  cfg, ok := cfgInf.(*RabbitMQClientConfig)
  if !ok || cfg == nil {
    return nil, fmt.Errorf("invalid config input")
  }

  client := &RabbitMQClient{
    uri:    uri,
    cfg:    cfg,
    quitCh: make(chan struct{}),
  }
  f.clients = append(f.clients, client)
  return client, nil
}

// Init initializes the RabbitMQ client and returns the AMQP delivery channel on
// which messages/notifications will be received.
func (cl *RabbitMQClient) Init() (chan []byte, error) {
  var err error
  cl.conn, err = amqp.Dial(cl.uri)
  if err != nil {
    return nil, cl.verboseError("error connecting to RabbitMQ", err)
  }

  cl.amqpChan, err = cl.conn.Channel()
  if err != nil {
    return nil, cl.verboseError("error opening a channel to RabbitMQ", err)
  }

  q, err := cl.amqpChan.QueueDeclare(
    cl.cfg.Queue, // name
    false,        // durable
    false,        // delete when used
    false,        // exclusive
    false,        // noWait
    nil,          // arguments
  )
  if err != nil {
    return nil, cl.verboseError("error declaring queue", err)
  }

  err = cl.amqpChan.QueueBind(
    q.Name,            // queue name
    cl.cfg.BindingKey, // routing key
    cl.cfg.Exchange,   // exchange
    false,             // noWait
    nil,               // arguments
  )
  if err != nil {
    return nil, cl.verboseError("error binding to queue", err)
  }

  cl.mqCh, err = cl.amqpChan.Consume(
    q.Name,
    cl.cfg.ConsumerTag,
    true,  // autoAck
    false, // exclusive
    false, // noLocal
    false, // noWait
    nil,   // arguments
  )
  if err != nil {
    return nil, cl.verboseError("error registering a consumer", err)
  }

  cl.msgCh = make(chan []byte)

  go cl.processMsgs()

  return cl.msgCh, nil
}

// processMsgs reads messages from the message queue channel hook (mqCh) and
// puts the body of the amqp.Delivery onto the msgCh. It quits when Stop is
// called by monitoring the quit channel.
func (cl *RabbitMQClient) processMsgs() {

  for {
    select {
    case msg := <-cl.mqCh:
      // Avoid a deadlock in case of a broken channel.
      select {
      case cl.msgCh <- msg.Body:
      case <-cl.quitCh:
        cl.closeConn()
        return
      }
    case <-cl.quitCh:
      fmt.Printf("received quit, closing RabbitMQ client connection")
      cl.closeConn()
      return
    }
  }
}

// closeConn closes a client connection and  its amqpChann.
func (cl *RabbitMQClient) closeConn() {
  if cl.conn != nil {
    cl.conn.Close()
    cl.conn = nil
  }
  if cl.amqpChan != nil {
    cl.amqpChan.Close()
    cl.amqpChan = nil
  }
}

// Stop stops listening to the amqp channel.
func (cl *RabbitMQClient) Stop() error {
  if cl == nil {
    return fmt.Errorf("stop on an uninitalized client")
  }
  if cl.quitCh != nil {
    close(cl.quitCh)
  }

  return nil
}

// verboseError is a helper function that returns a more verbose version of the
// provided error according to the provided .
func (cl *RabbitMQClient) verboseError(errStr string, err error) error {
  if err == nil {
    return nil
  }
  return fmt.Errorf("%v :: uri: %v, cfg %v, err: %v", errStr, cl.uri, cl.cfg,
    err)
}
