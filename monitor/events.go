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

package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"
  "reflect"
  "strings"

  "github.com/golang/glog"

  "github.com/zerostackinc/libmq"
)

var (
  rabbitmqURL = flag.String("url", "", "url of the rabbitmq server")
)

var mqClientCfg = []libmq.RabbitMQClientConfig{
  libmq.RabbitMQClientConfig{
    Exchange:     "openstack",
    ExchangeType: "topic",
    Queue:        "notifications.info",
    BindingKey:   "#",
    ConsumerTag:  "librabbitmq",
  },
  libmq.RabbitMQClientConfig{
    Exchange:     "nova",
    ExchangeType: "topic",
    Queue:        "notifications.info",
    BindingKey:   "#",
    ConsumerTag:  "librabbitmq",
  },
}

func main() {
  flag.Parse()
  if *rabbitmqURL == "" {
    fmt.Printf("rabbit mq url param --url is required\n")
    os.Exit(1)
  }
  mqFactory := libmq.NewRabbitMQFactory()
  err := startEventController(mqFactory, mqClientCfg, *rabbitmqURL)
  if err != nil {
    fmt.Printf(err.Error())
  }
}

// startEventController sets up and starts the event controller to
// listen to events, messages and notifications on the event channels.
func startEventController(mqFactory libmq.MQClientFactory,
  mqConfig []libmq.RabbitMQClientConfig, rabbitmqURL string) error {

  // Sanity check.
  if mqConfig == nil || len(mqConfig) == 0 {
    return fmt.Errorf("no RabbitMQ client config objects specified")
  }

  _, mqMsgChans, err := setupEventController(mqFactory, mqConfig, rabbitmqURL)
  if err != nil {
    glog.Errorf("error setting up event controller :: %v", err)
    return fmt.Errorf("error seting up event controller :: %v", err)
  }

  // Start main event controller loop.
  glog.Info("starting event controller")
  eventController(mqMsgChans)
  return nil
}

// eventController is the main event controller loop that listens to
// events, messages and notifications on the event channels.
func eventController(mqMsgChans []chan []byte) {

  // Cases for RabbitMQ message channels
  numChans := len(mqMsgChans)

  cases := make([]reflect.SelectCase, numChans)

  for i, msgChan := range mqMsgChans {
    cases[i] = reflect.SelectCase{
      Dir:  reflect.SelectRecv,
      Chan: reflect.ValueOf(msgChan),
    }
  }

  glog.Infof("started monitoring %d message channels", numChans)

  for {
    // Select a case and extract received value.
    chosen, value, ok := reflect.Select(cases)

    // Handle failure/quit.
    if value.IsNil() || !ok {
      // Channel was closed - close all channels
      glog.Errorf("connection to RabbitMQ broken, chosen: %v, value: %v, "+
        "ok: %v", chosen, value, ok)
      return
    }

    // Sanity test for value.
    if !value.CanInterface() {
      glog.Errorf("unexpected type of value received on RabbitMQ message " +
        "channel")
      continue
    }

    // Extract event payload.
    rabbitMsg, ok := value.Interface().([]byte)
    if !ok {
      glog.Errorf("unexpected type of value (%v) on RabbitMQ message "+
        " channel :: %v", value.Type(), value)
      continue
    }

    // Parse event.
    err := parseNotification(rabbitMsg)
    if err != nil {
      glog.Errorf("error parsing RabbitMQ event string to compose "+
        "eventStr: %s, err: %s", rabbitMsg, err)
      continue
    }
  }
}

// setupEventController sets up clients to listen to events.
func setupEventController(mqFactory libmq.MQClientFactory,
  mqConfig []libmq.RabbitMQClientConfig, rabbitmqURI string) (
  []libmq.MQClient, []chan []byte, error) {

  var err error
  var mqMsgChans []chan []byte
  var mqClients []libmq.MQClient

  // Setup channels to communicate with RabbitMQ.
  glog.Infof("trying to connect to RabbitMQ at URI: %v", rabbitmqURI)

  for _, rabbitmqCfg := range mqConfig {
    client, msgChan, err := setupRabbitMQClient(mqFactory, rabbitmqURI,
      &rabbitmqCfg)
    if client == nil || msgChan == nil || err != nil {
      glog.Errorf("error in setup client: %v, msgChan: %+v :: %v",
        client, msgChan, err)

      if msgChan != nil {
        close(msgChan)
      }
      if client != nil {
        client.Stop()
      }
      return nil, nil, err
    }
    mqMsgChans = append(mqMsgChans, msgChan)
    mqClients = append(mqClients, client)
  }

  // Return event channels.
  if mqMsgChans == nil {
    return nil, nil, err
  }
  return mqClients, mqMsgChans, nil
}

// setupRabbitMQClient sets up one client to RabbitMQ
func setupRabbitMQClient(mqFactory libmq.MQClientFactory,
  rabbitmqURI string, rabbitmqCfg *libmq.RabbitMQClientConfig) (
  libmq.MQClient, chan []byte, error) {

  // Create new RabbitMQ client.
  client, err := mqFactory.NewClient(rabbitmqURI, rabbitmqCfg)
  if client == nil || err != nil {
    return client, nil, err
  }

  // Initialize new RabbitMQ client.
  msgChan, err := client.Init()
  return client, msgChan, err
}

// parseNotification is a function that parses the provided event.
func parseNotification(notification []byte) error {

  var data OSNotificationData
  err := json.Unmarshal(notification, &data)
  if err != nil || data.EventType == "" {
    glog.Infof("unhandled notification :: %s, err: %s", notification, err)
    return nil
  }

  glog.Infof("received notification of type: %s", data.EventType)

  switch strings.TrimSpace(data.EventType) {
  case "compute.instance.create.start":
    return handleVMEvent(notification, EVENT_VM_CREATE_START)

  case "compute.instance.create.error":
    return handleVMEvent(notification, EVENT_VM_CREATE_ERROR)

  case "compute.instance.create.end":
    return handleVMEvent(notification, EVENT_VM_CREATE_END)

  case "compute.instance.delete.start":
    return handleVMEvent(notification, EVENT_VM_DELETE_START)

  case "compute.instance.delete.end":
    return handleVMEvent(notification, EVENT_VM_DELETE_END)

  case "compute.instance.shutdown.start":
    return handleVMEvent(notification, EVENT_VM_SHUTDOWN_START)

  case "compute.instance.shutdown.end":
    return handleVMEvent(notification, EVENT_VM_SHUTDOWN_END)

  case "compute.instance.power_off.start":
    return handleVMEvent(notification, EVENT_VM_POWEROFF_START)

  case "compute.instance.power_off.end":
    return handleVMEvent(notification, EVENT_VM_POWEROFF_END)

  case "compute.instance.power_on.start":
    return handleVMEvent(notification, EVENT_VM_POWERON_START)

  case "compute.instance.power_on.end":
    return handleVMEvent(notification, EVENT_VM_POWERON_END)

  case "compute.instance.rebuild.start":
    return handleVMEvent(notification, EVENT_VM_REBUILD_START)

  case "compute.instance.rebuild.end":
    return handleVMEvent(notification, EVENT_VM_REBUILD_END)

  case "compute.instance.snapshot.start":
    return handleVMEvent(notification, EVENT_VM_SNAPSHOT_START)

  case "compute.instance.snapshot.end":
    return handleVMEvent(notification, EVENT_VM_SNAPSHOT_END)

  case "compute.instance.suspend.start":
    return handleVMEvent(notification, EVENT_VM_SUSPEND_START)

  case "compute.instance.suspend.end":
    return handleVMEvent(notification, EVENT_VM_SUSPEND_END)

  case "compute.instance.resume.start":
    return handleVMEvent(notification, EVENT_VM_RESUME_START)

  case "compute.instance.resume.end":
    return handleVMEvent(notification, EVENT_VM_RESUME_END)

  case "compute.instance.pause.start":
    return handleVMEvent(notification, EVENT_VM_PAUSE_START)

  case "compute.instance.pause.end":
    return handleVMEvent(notification, EVENT_VM_PAUSE_END)

  case "compute.instance.unpause.start":
    return handleVMEvent(notification, EVENT_VM_UNPAUSE_START)

  case "compute.instance.unpause.end":
    return handleVMEvent(notification, EVENT_VM_UNPAUSE_END)

  case "compute.instance.resize.start":
    return handleVMEvent(notification, EVENT_VM_RESIZE_START)

  case "compute.instance.resize.end":
    return handleVMEvent(notification, EVENT_VM_RESIZE_END)

  case "compute.instance.finish_resize.start":
    return handleVMEvent(notification, EVENT_VM_FINISH_RESIZE_START)

  case "compute.instance.finish_resize.end":
    return handleVMEvent(notification, EVENT_VM_FINISH_RESIZE_END)

  case "compute.instance.resize.prep.start":
    return handleVMEvent(notification, EVENT_VM_RESIZE_PREP_START)

  case "compute.instance.resize.prep.end":
    return handleVMEvent(notification, EVENT_VM_RESIZE_PREP_END)

  case "compute.instance.resize.confirm.start":
    return handleVMEvent(notification, EVENT_VM_RESIZE_CONFIRM_START)

  case "compute.instance.resize.confirm.end":
    return handleVMEvent(notification, EVENT_VM_RESIZE_CONFIRM_END)

  case "compute.instance.resize.revert.start":
    return handleVMEvent(notification, EVENT_VM_RESIZE_REVERT_START)

  case "compute.instance.resize.revert.end":
    return handleVMEvent(notification, EVENT_VM_RESIZE_REVERT_END)

  case "compute.instance.exists":
    return handleVMEvent(notification, EVENT_VM_EXISTS)

  case "compute.instance.update":
    return handleVMEvent(notification, EVENT_VM_UPDATE)

  case "compute.instance.volume.attach":
    return handleVMEvent(notification, EVENT_VM_VOLUME_ATTACH)

  case "compute.instance.volume.detach":
    return handleVMEvent(notification, EVENT_VM_VOLUME_DETACH)

  // https://wiki.openstack.org/wiki/SystemUsageData describes {suspend,
  // resume}.{start,end} as events generated for vm suspend and resume. But we
  // only get suspend and resume for the activity. This could be a bug, and we
  // handle it below.
  case "compute.instance.suspend":
    return handleVMEvent(notification, EVENT_VM_SUSPEND)

  case "compute.instance.resume":
    return handleVMEvent(notification, EVENT_VM_RESUME)

  case "compute.instance.live_migration.pre.start":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_PRE_START)

  case "compute.instance.live_migration.pre.end":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_PRE_END)

  case "compute.instance.live_migration.post.dest.start":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_POST_DEST_START)

  case "compute.instance.live_migration.post.dest.end":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_POST_DEST_END)

  case "compute.instance.live_migration._post.start":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_POST_START)

  case "compute.instance.live_migration._post.end":
    return handleVMEvent(notification, EVENT_LIVE_MIGRATION_POST_END)

  case "servergroup.create":
    return handleServerGroupEvent(notification, EVENT_SERVER_GROUP_CREATE)

  case "servergroup.delete":
    return handleServerGroupEvent(notification, EVENT_SERVER_GROUP_DELETE)

  case "servergroup.update":
    return handleServerGroupEvent(notification, EVENT_SERVER_GROUP_UPDATE)

  case "servergroup.addmember":
    return handleServerGroupEvent(notification, EVENT_SERVER_GROUP_ADD_MEMBER)

  case "volume.create.start":
    return handleVolumeEvent(notification, EVENT_VOLUME_CREATE_START)

  case "volume.create.end":
    return handleVolumeEvent(notification, EVENT_VOLUME_CREATE_END)

  case "volume.delete.start":
    return handleVolumeEvent(notification, EVENT_VOLUME_DELETE_START)

  case "volume.delete.end":
    return handleVolumeEvent(notification, EVENT_VOLUME_DELETE_END)

  case "volume.usage":
    return handleVolumeEvent(notification, EVENT_VOLUME_USAGE)

  case "volume.attach.start":
    return handleVolumeEvent(notification, EVENT_VOLUME_ATTACH_START)

  case "volume.attach.end":
    return handleVolumeEvent(notification, EVENT_VOLUME_ATTACH_END)

  case "volume.detach.start":
    return handleVolumeEvent(notification, EVENT_VOLUME_DETACH_START)

  case "volume.detach.end":
    return handleVolumeEvent(notification, EVENT_VOLUME_DETACH_END)

  case "network.create.start":
    return handleNetEvent(notification, EVENT_NETWORK_CREATE_START)

  case "network.create.end":
    return handleNetEvent(notification, EVENT_NETWORK_CREATE_END)

  case "network.delete.start":
    return handleNetEvent(notification, EVENT_NETWORK_DELETE_START)

  case "network.delete.end":
    return handleNetEvent(notification, EVENT_NETWORK_DELETE_END)

  case "subnet.create.start":
    return handleSubnetEvent(notification, EVENT_SUBNET_CREATE_START)

  case "subnet.create.end":
    return handleSubnetEvent(notification, EVENT_SUBNET_CREATE_END)

  case "subnet.delete.start":
    return handleSubnetEvent(notification, EVENT_SUBNET_DELETE_START)

  case "subnet.delete.end":
    return handleSubnetEvent(notification, EVENT_SUBNET_DELETE_END)

  case "port.create.start":
    return handlePortEvent(notification, EVENT_PORT_CREATE_START)

  case "port.create.end":
    return handlePortEvent(notification, EVENT_PORT_CREATE_END)

  case "port.delete.start":
    return handlePortEvent(notification, EVENT_PORT_DELETE_START)

  case "port.delete.end":
    return handlePortEvent(notification, EVENT_PORT_DELETE_END)

  case "router.create.start":
    return handleRouterEvent(notification, EVENT_ROUTER_CREATE_START)

  case "router.create.end":
    return handleRouterEvent(notification, EVENT_ROUTER_CREATE_END)

  case "router.update.start":
    return handleRouterEvent(notification, EVENT_ROUTER_UPDATE_START)

  case "router.update.end":
    return handleRouterEvent(notification, EVENT_ROUTER_UPDATE_END)

  case "router.delete.start":
    return handleRouterEvent(notification, EVENT_ROUTER_DELETE_START)

  case "router.delete.end":
    return handleRouterEvent(notification, EVENT_ROUTER_DELETE_END)

  case "router.interface.create":
    return handleRouterEvent(notification, EVENT_ROUTER_INTERFACE_CREATE)

  case "router.interface.delete":
    return handleRouterEvent(notification, EVENT_ROUTER_INTERFACE_DELETE)

  case "identity.group.created":
    return handleGroupEvent(notification, EVENT_GROUP_CREATED)

  case "identity.group.updated":
    return handleGroupEvent(notification, EVENT_GROUP_UPDATED)

  case "identity.group.deleted":
    return handleGroupEvent(notification, EVENT_GROUP_DELETED)

  case "identity.project.created":
    return handleProjectEvent(notification, EVENT_PROJECT_CREATED)

  case "identity.project.updated":
    return handleProjectEvent(notification, EVENT_PROJECT_UPDATED)

  case "identity.project.deleted":
    return handleProjectEvent(notification, EVENT_PROJECT_DELETED)

  case "identity.role.created":
    return handleRoleEvent(notification, EVENT_ROLE_CREATED)

  case "identity.role.updated":
    return handleRoleEvent(notification, EVENT_ROLE_UPDATED)

  case "identity.role.deleted":
    return handleRoleEvent(notification, EVENT_ROLE_DELETED)

  case "identity.user.created":
    return handleUserEvent(notification, EVENT_USER_CREATED)

  case "identity.user.updated":
    return handleUserEvent(notification, EVENT_USER_UPDATED)

  case "identity.user.deleted":
    return handleUserEvent(notification, EVENT_USER_DELETED)

  case "identity.trust.created":
    return handleTrustEvent(notification, EVENT_TRUST_CREATED)

  // No case for "identity.trust.updated" because "identity.trust" is an
  // immutable resource.

  case "identity.trust.deleted":
    return handleTrustEvent(notification, EVENT_TRUST_DELETED)

  case "identity.region.created":
    return handleRegionEvent(notification, EVENT_REGION_CREATED)

  case "identity.region.updated":
    return handleRegionEvent(notification, EVENT_REGION_UPDATED)

  case "identity.region.deleted":
    return handleRegionEvent(notification, EVENT_REGION_DELETED)

  case "identity.endpoint.created":
    return handleEndpointEvent(notification, EVENT_ENDPOINT_CREATED)

  case "identity.endpoint.updated":
    return handleEndpointEvent(notification, EVENT_ENDPOINT_UPDATED)

  case "identity.endpoint.deleted":
    return handleEndpointEvent(notification, EVENT_ENDPOINT_DELETED)

  case "identity.service.created":
    return handleServiceEvent(notification, EVENT_SERVICE_CREATED)

  case "identity.service.updated":
    return handleServiceEvent(notification, EVENT_SERVICE_UPDATED)

  case "identity.service.deleted":
    return handleServiceEvent(notification, EVENT_SERVICE_DELETED)

  case "identity.policy.created":
    return handlePolicyEvent(notification, EVENT_POLICY_CREATED)

  case "identity.policy.updated":
    return handlePolicyEvent(notification, EVENT_POLICY_UPDATED)

  case "identity.policy.deleted":
    return handlePolicyEvent(notification, EVENT_POLICY_DELETED)

  case "image.create":
    return handleImageEvent(notification, EVENT_IMAGE_CREATE)

  case "image.prepare":
    return handleImageEvent(notification, EVENT_IMAGE_PREPARE)

  case "image.upload":
    return handleImageEvent(notification, EVENT_IMAGE_UPLOAD)

  case "image.activate":
    return handleImageEvent(notification, EVENT_IMAGE_ACTIVATE)

  case "image.send":
    // image.send notification payload is different from other image events.
    return handleImageSendEvent(notification, EVENT_IMAGE_SEND)

  case "image.update":
    return handleImageEvent(notification, EVENT_IMAGE_UPDATE)

  case "image.delete":
    return handleImageEvent(notification, EVENT_IMAGE_DELETE)

  default:
    glog.Infof("unhandled eventType: %s", data.EventType)
    return handleUnmonitoredEvent(notification, EVENT_UNMONITORED)
  }
}

// handleVMEvent parses VM related notifications from Openstack.
func handleVMEvent(notification []byte, evType EventType) error {

  var data ComputeInstanceData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleServerGroupEvent parses ServerGroup related notifications from
// Openstack.
func handleServerGroupEvent(notification []byte, evType EventType) error {

  var data ServerGroupData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleVolumeEvent parses Volume related notifications from Openstack.
func handleVolumeEvent(notification []byte, evType EventType) error {

  var data VolumeData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleNetEvent parses Network related notifications from Openstack.
func handleNetEvent(notification []byte, evType EventType) error {

  var data NetworkData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleSubnetEvent parses Subnet related notifications from Openstack.
func handleSubnetEvent(notification []byte, evType EventType) error {

  var data SubnetData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handlePortEvent parses Port related notifications from Openstack.
func handlePortEvent(notification []byte, evType EventType) error {

  var data PortData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleRouterEvent parses Router related notifications from Openstack.
func handleRouterEvent(notification []byte, evType EventType) error {

  var data RouterData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleGroupEvent parses Group related Keystone notifications from Openstack.
func handleGroupEvent(notification []byte, evType EventType) error {

  var data IdentityGroupData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleProjectEvent parses Project related notifications from Openstack.
func handleProjectEvent(notification []byte, evType EventType) error {

  var data IdentityProjectData
  if err := json.Unmarshal(notification, &data); err != nil {
    return err
  }

  return nil
}

// handleRoleEvent parses Role related Keystone notifications from Openstack.
func handleRoleEvent(notification []byte, evType EventType) error {

  var data IdentityRoleData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleUserEvent parses User related Keystone notifications from Openstack.
func handleUserEvent(notification []byte, evType EventType) error {

  var data IdentityUserData
  if err := json.Unmarshal(notification, &data); err != nil {
    return err
  }
  return nil
}

// handleTrustEvent parses Trust related Keystone notifications from Openstack.
func handleTrustEvent(notification []byte, evType EventType) error {

  var data IdentityTrustData
  if err := json.Unmarshal(notification, &data); err != nil {
    return err
  }
  return nil
}

// handleRegionEvent parses Region related Keystone notifications from
// Openstack.
func handleRegionEvent(notification []byte, evType EventType) error {

  var data IdentityRegionData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleEndpointEvent parses Endpoint related Keystone notifications from
// Openstack.
func handleEndpointEvent(notification []byte, evType EventType) error {

  var data IdentityEndpointData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleServiceEvent parses Service related Keystone notifications from
// Openstack.
func handleServiceEvent(notification []byte, evType EventType) error {

  var data IdentityServiceData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handlePolicyEvent parses Policy related Keystone notifications from
// Openstack.
func handlePolicyEvent(notification []byte, evType EventType) error {

  var data IdentityPolicyData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleImageEvent parses Image related Glance notifications from Openstack.
func handleImageEvent(notification []byte, evType EventType) error {

  var data ImageData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  // http://docs.openstack.org/developer/glance/notifications.html
  // Seems like a bug in glance notification that fails to populate the owner
  // field of some of the glance events.

  return nil
}

// handleImageSendEvent parses image.send Glance notification from Openstack.
func handleImageSendEvent(notification []byte, evType EventType) error {

  var data ImageSendData
  err := json.Unmarshal(notification, &data)
  if err != nil {
    return err
  }
  return nil
}

// handleUnmonitoredEvent creates a bolt with notification as its data.
func handleUnmonitoredEvent(notification []byte, evType EventType) error {

  return nil
}

/*
  We could just call parsing like this but it will not enable any custom
	handling of the each event type separately:

    return parseNotification(notification, EVENT_POLICY_DELETED, &IdentityPolicyData{})

type eventWithMessageID interface {
  GetMessageID() string
}

// handleNotification parses notifications using the interface
func handleNotification(notification []byte, evType EventType, data interface{}) error) {

  err := json.Unmarshal(notification, &data)
  if err != nil {
    return nil, err
  }
	return nil
}
*/
