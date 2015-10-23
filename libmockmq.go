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

package libmq

import (
  "fmt"
)

// MockMQClient represents that parameters and state of a RabbitMQ client.
type MockMQClient struct {
  uri    string
  quitCh chan struct{}
  msgCh  chan []byte
}

// MockMQClientFactory generates MockMQClient and keeps track of them.
type MockMQClientFactory struct {
  Clients []*MockMQClient
}

// NewMockMQClientFactory creates a MockMQClientFactory object.
func NewMockMQClientFactory() *MockMQClientFactory {
  return &MockMQClientFactory{}
}

// NewClient creates a new MockMQClient object.
func (f *MockMQClientFactory) NewClient(
  uri string, cfgInf interface{}) (MQClient, error) {

  if cfgInf == nil {
    return nil, fmt.Errorf("invalid config input")
  }

  client := &MockMQClient{
    uri:    uri,
    quitCh: make(chan struct{}),
  }

  f.Clients = append(f.Clients, client)

  return client, nil
}

// Init initializes the mock client and returns the channel on
// which messages/notifications will be received.
func (cl *MockMQClient) Init() (chan []byte, error) {
  cl.msgCh = make(chan []byte)
  return cl.msgCh, nil
}

// Stop stops listening to the mq.
func (cl *MockMQClient) Stop() error {
  if cl == nil {
    return fmt.Errorf("Stop on an uninitialized client")
  }
  return nil
}

// SendMessage is a blocking function that sends the given message on the
// message channel.
func (cl *MockMQClient) SendMessage(msg []byte) error {
  if cl == nil || cl.msgCh == nil {
    return fmt.Errorf("invalid client in sending message")
  }
  cl.msgCh <- msg
  return nil
}
