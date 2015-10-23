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

package main

// OSNotificationData represents generic Openstack notification data. It is only
// used to extract the event_type.
type OSNotificationData struct {
  EventType string `json:"event_type"`
}

// EventType lists all types of event.
// Note: Absolute event index = Entity index * 100 + relative event index
type EventType int32

const (
  EVENT_INVALID EventType = 0
  // Compute
  EVENT_HOST_METRIC                    EventType = 101
  EVENT_FIRST                          EventType = 101
  EVENT_HOST_ACCOUNT_ID_SET            EventType = 102
  EVENT_HOST_ACCOUNT_ID_SET_ERROR      EventType = 103
  EVENT_HOST_ACCOUNT_ID_DELETE         EventType = 104
  EVENT_HOST_ACCOUNT_ID_DELETE_ERROR   EventType = 105
  EVENT_VM_METRIC                      EventType = 201
  EVENT_VM_MIN                         EventType = 201
  EVENT_VM_CREATE_START                EventType = 202
  EVENT_VM_CREATE_ERROR                EventType = 203
  EVENT_VM_CREATE_END                  EventType = 204
  EVENT_VM_DELETE_START                EventType = 205
  EVENT_VM_DELETE_END                  EventType = 206
  EVENT_VM_SHUTDOWN_START              EventType = 207
  EVENT_VM_SHUTDOWN_END                EventType = 208
  EVENT_VM_POWEROFF_START              EventType = 209
  EVENT_VM_POWEROFF_END                EventType = 210
  EVENT_VM_POWERON_START               EventType = 211
  EVENT_VM_POWERON_END                 EventType = 212
  EVENT_VM_REBUILD_START               EventType = 213
  EVENT_VM_REBUILD_END                 EventType = 214
  EVENT_VM_SNAPSHOT_START              EventType = 215
  EVENT_VM_SNAPSHOT_END                EventType = 216
  EVENT_VM_SUSPEND_START               EventType = 217
  EVENT_VM_SUSPEND_END                 EventType = 218
  EVENT_VM_RESUME_START                EventType = 219
  EVENT_VM_RESUME_END                  EventType = 220
  EVENT_VM_PAUSE_START                 EventType = 221
  EVENT_VM_PAUSE_END                   EventType = 222
  EVENT_VM_UNPAUSE_START               EventType = 223
  EVENT_VM_UNPAUSE_END                 EventType = 224
  EVENT_VM_RESIZE_START                EventType = 225
  EVENT_VM_RESIZE_END                  EventType = 226
  EVENT_VM_FINISH_RESIZE_START         EventType = 227
  EVENT_VM_FINISH_RESIZE_END           EventType = 228
  EVENT_VM_RESIZE_PREP_START           EventType = 229
  EVENT_VM_RESIZE_PREP_END             EventType = 230
  EVENT_VM_RESIZE_CONFIRM_START        EventType = 231
  EVENT_VM_RESIZE_CONFIRM_END          EventType = 232
  EVENT_VM_RESIZE_REVERT_START         EventType = 233
  EVENT_VM_RESIZE_REVERT_END           EventType = 234
  EVENT_VM_EXISTS                      EventType = 235
  EVENT_VM_UPDATE                      EventType = 236
  EVENT_VM_VOLUME_ATTACH               EventType = 237
  EVENT_VM_VOLUME_DETACH               EventType = 238
  EVENT_VM_SUSPEND                     EventType = 239
  EVENT_VM_RESUME                      EventType = 240
  EVENT_LIVE_MIGRATION_PRE_START       EventType = 241
  EVENT_LIVE_MIGRATION_PRE_END         EventType = 242
  EVENT_LIVE_MIGRATION_POST_DEST_START EventType = 243
  EVENT_LIVE_MIGRATION_POST_DEST_END   EventType = 244
  EVENT_LIVE_MIGRATION_POST_START      EventType = 245
  EVENT_LIVE_MIGRATION_POST_END        EventType = 246
  EVENT_VM_SYNC                        EventType = 300
  EVENT_VM_MAX                         EventType = 300
  EVENT_PROCESS_METRIC                 EventType = 301
  EVENT_SERVER_GROUP_CREATE            EventType = 401
  EVENT_SERVER_GROUP_DELETE            EventType = 402
  EVENT_SERVER_GROUP_UPDATE            EventType = 403
  EVENT_SERVER_GROUP_ADD_MEMBER        EventType = 404
  // Storage
  EVENT_VOLUME_CREATE_START EventType = 1101
  EVENT_VOLUME_CREATE_END   EventType = 1102
  EVENT_VOLUME_DELETE_START EventType = 1103
  EVENT_VOLUME_DELETE_END   EventType = 1104
  EVENT_VOLUME_USAGE        EventType = 1105
  EVENT_VOLUME_ATTACH_START EventType = 1106
  EVENT_VOLUME_ATTACH_END   EventType = 1107
  EVENT_VOLUME_DETACH_START EventType = 1108
  EVENT_VOLUME_DETACH_END   EventType = 1109
  // Network
  EVENT_NETWORK_CREATE_START    EventType = 2101
  EVENT_NETWORK_CREATE_END      EventType = 2102
  EVENT_NETWORK_DELETE_START    EventType = 2103
  EVENT_NETWORK_DELETE_END      EventType = 2104
  EVENT_SUBNET_CREATE_START     EventType = 2201
  EVENT_SUBNET_CREATE_END       EventType = 2202
  EVENT_SUBNET_DELETE_START     EventType = 2203
  EVENT_SUBNET_DELETE_END       EventType = 2204
  EVENT_PORT_CREATE_START       EventType = 2301
  EVENT_PORT_CREATE_END         EventType = 2302
  EVENT_PORT_DELETE_START       EventType = 2303
  EVENT_PORT_DELETE_END         EventType = 2304
  EVENT_ROUTER_CREATE_START     EventType = 2305
  EVENT_ROUTER_CREATE_END       EventType = 2306
  EVENT_ROUTER_UPDATE_START     EventType = 2307
  EVENT_ROUTER_UPDATE_END       EventType = 2308
  EVENT_ROUTER_DELETE_START     EventType = 2309
  EVENT_ROUTER_DELETE_END       EventType = 2310
  EVENT_ROUTER_INTERFACE_CREATE EventType = 2311
  EVENT_ROUTER_INTERFACE_DELETE EventType = 2312
  // Identity
  EVENT_GROUP_CREATED    EventType = 3101
  EVENT_GROUP_UPDATED    EventType = 3102
  EVENT_GROUP_DELETED    EventType = 3103
  EVENT_PROJECT_CREATED  EventType = 3201
  EVENT_PROJECT_UPDATED  EventType = 3202
  EVENT_PROJECT_DELETED  EventType = 3203
  EVENT_ROLE_CREATED     EventType = 3301
  EVENT_ROLE_UPDATED     EventType = 3302
  EVENT_ROLE_DELETED     EventType = 3303
  EVENT_USER_CREATED     EventType = 3401
  EVENT_USER_UPDATED     EventType = 3402
  EVENT_USER_DELETED     EventType = 3403
  EVENT_TRUST_CREATED    EventType = 3501
  EVENT_TRUST_DELETED    EventType = 3502
  EVENT_REGION_CREATED   EventType = 3601
  EVENT_REGION_UPDATED   EventType = 3602
  EVENT_REGION_DELETED   EventType = 3603
  EVENT_ENDPOINT_CREATED EventType = 3701
  EVENT_ENDPOINT_UPDATED EventType = 3702
  EVENT_ENDPOINT_DELETED EventType = 3703
  EVENT_SERVICE_CREATED  EventType = 3801
  EVENT_SERVICE_UPDATED  EventType = 3802
  EVENT_SERVICE_DELETED  EventType = 3803
  EVENT_POLICY_CREATED   EventType = 3901
  EVENT_POLICY_UPDATED   EventType = 3902
  EVENT_POLICY_DELETED   EventType = 3903
  // Image
  EVENT_IMAGE_CREATE   EventType = 4101
  EVENT_IMAGE_PREPARE  EventType = 4102
  EVENT_IMAGE_UPLOAD   EventType = 4103
  EVENT_IMAGE_ACTIVATE EventType = 4104
  EVENT_IMAGE_SEND     EventType = 4105
  EVENT_IMAGE_UPDATE   EventType = 4106
  EVENT_IMAGE_DELETE   EventType = 4107
  EVENT_OPENSTACK_LAST EventType = 4107
  // ENTITY_UNMONITORED
  EVENT_UNMONITORED EventType = 99999
  // Alias
  EVENT_LAST EventType = 99999
)

////////////////////////////////////////////////////////////////////////////////
//                            VM events data
////////////////////////////////////////////////////////////////////////////////

// ComputeInstanceDataImageMeta has the image reference.
type ComputeInstanceDataImageMeta struct {
  BaseImageRef string `json:"base_image_ref"`
}

// ComputeInstanceDataPayload is the payload that comes in ComputeInstanceData.
type ComputeInstanceDataPayload struct {
  VMUUID              string                       `json:"instance_id"`
  DisplayName         string                       `json:"display_name"`
  LaunchedAt          string                       `json:"launched_at"`
  CreatedAt           string                       `json:"created_at"`
  AvailabilityZone    string                       `json:"availability_zone"`
  HostName            string                       `json:"host"`
  InstanceType        string                       `json:"instance_type"`
  InstanceTypeID      int32                        `json:"instance_type_id"`
  NewInstanceType     string                       `json:"new_instance_type"`
  NewInstanceTypeID   int32                        `json:"new_instance_type_id"`
  ImageRefURL         string                       `json:"image_ref_url"`
  ImageMeta           ComputeInstanceDataImageMeta `json:"json_meta"`
  AuditPeriodBegining string                       `json:"audit_period_begining"`
  AuditPeriodEnding   string                       `json:"audit_period_ending"`
  OldState            string                       `json:"old_state"`
  State               string                       `json:"state"`
  StateDescription    string                       `json:"state_description"`
  VCPUs               int                          `json:"vcpus"`
  MemoryMB            int                          `json:"memory_mb"`
  DiskGB              int                          `json:"disk_gb"`
  FlavorID            string                       `json:"instance_flavor_id"`
  Details             string                       `json:"message"`
  DeletedAt           string                       `json:"deleted_at"`
}

// ComputeInstanceData identifies fields of interest from these Openstack
// notifications:
// compute.instance.create.{start,error,end} and
// compute.instance.{delete,shutdown,power_off,power_on}.{start,end}
// compute.instance.{rebuild,snapshot,suspend,resume,pause,unpause}.{start,end}
// compute.instance.{resize,finish_resize}.{start,end}
// compute.instance.resize.{prep,confirm,revert}.{start,end}
// compute.instance.{exists,update}
type ComputeInstanceData struct {
  EventType   string                     `json:"event_type"`
  Timestamp   string                     `json:"timestamp"`
  Priority    string                     `json:"priority"`
  MessageID   string                     `json:"message_id"`
  ProjectID   string                     `json:"_context_project_id"`
  ProjectName string                     `json:"_context_project_name"`
  UserID      string                     `json:"_context_user_id"`
  UserName    string                     `json:"_context_user_name"`
  Token       string                     `json:"_context_auth_token"`
  Payload     ComputeInstanceDataPayload `json:"payload"`
}

// ServerGroupData identifies fields of interest from these Openstack
// notifications:
// servergroup.{create,delete,update,addmember}
type ServerGroupData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Payload     struct {
    ServerGroupID string `json:"server_group_id"`
    Name          string `json:"name"`
  }
}

////////////////////////////////////////////////////////////////////////////////
//                         Storage events data
////////////////////////////////////////////////////////////////////////////////

// VolumeData identifies fields of interest from these Openstack
// notifications:
// volume.{create,delete}.{start,end}
// volume.usage
type VolumeData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Token       string `json:"_context_auth_token"`
  Payload     struct {
    VolumeID            string `json:"volume_id"`
    VolumeType          string `json:"volume_type"`
    DisplayName         string `json:"display_name"`
    LaunchedAt          string `json:"launched_at"`
    CreatedAt           string `json:"created_at"`
    AvailabilityZone    string `json:"availability_zone"`
    Status              string `json:"status"`
    SizeGB              int    `json:"size"`
    SnapshotID          string `json:"snapshot_id"`
    AuditPeriodBegining string `json:"audit_period_begining"`
    AuditPeriodEnding   string `json:"audit_period_ending"`
  }
}

////////////////////////////////////////////////////////////////////////////////
//                         Network events data
////////////////////////////////////////////////////////////////////////////////

// NetworkData contains interesting data from the Openstack
// network.{create,delete}.{start,end} notifications.
type NetworkData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Payload     struct {
    Network struct {
      ID                     string `json:"id"`
      RouterExternal         bool   `json:"router:external"`
      Name                   string `json:"name"`
      AdminStateUp           bool   `json:"admin_state_up"`
      ProviderNetworkType    string `json:"provider:network_type"`
      Shared                 bool   `json:"shared"`
      ProviderSegmentationID int64  `json:"provider:segmentation_id"`
    }
  }
}

// SubnetData contains interesting data from the Openstack
// subnet.{create,delete}.{start,end} notifications.
type SubnetData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Payload     struct {
    Subnet struct {
      Name       string `json:"name"`
      NetworkID  string `json:"network_id"`
      EnableDHCP bool   `json:"enable_dhcp"`
      IPVersion  int64  `json:"ip_version"`
      CIDR       string `json:"cidr"`
    }
  }
}

// PortData contains interesting data from the Openstack
// port.{create,delete}.{start,end} notifications.
type PortData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Payload     struct {
    Port struct {
      Status          string `json:"status"`
      BindingHostID   string `json:"binding:host_id"`
      DeviceOwner     string `json:"device_owner"`
      ID              string `json:"id"`
      DeviceID        string `json:"device_id"`
      Name            string `json:"name"`
      AdminStateUp    bool   `json:"admin_state_up"`
      NetworkID       string `json:"network_id"`
      BindingVnicType string `json:"binding:vnic_type"`
      BindingVifType  string `json:"binding:vif_type"`
      MacAddress      string `json:"mac_address"`
    }
  }
}

// ExternalFixedIPs contains interesting data about fixed ips of an
// external gateway in a router.
type ExternalFixedIPs struct {
  SubnetID  string `json:"subnet_id"`
  IPAddress string `json:"ip_address"`
}

// ExternalGWInfo contains interesting data about external gateway of a router.
type ExternalGWInfo struct {
  NetID       string             `json:"network_id"`
  EnableSnat  bool               `json:"enable_snat"`
  ExtFixedIPs []ExternalFixedIPs `json:"external_fixed_ips"`
}

// RouterInfo contains interesting data about router.
type RouterInfo struct {
  Status       string         `json:"status"`
  ExtGWInfo    ExternalGWInfo `json:"external_gateway_info"`
  Name         string         `json:"name"`
  AdminStateUp bool           `json:"admin_state_up"`
  TenantID     string         `json:"tenant_id"`
  Distributed  bool           `json:"distributed"`
  Routes       []string       `json:"routes"`
  HA           bool           `json:"ha"`
  ID           string         `json:"id"`
}

// RouterInterfaceInfo contains interesting data about router's interface.
type RouterInterfaceInfo struct {
  SubnetID string `json:"subnet_id"`
  TenantID string `json:"tenant_id"`
  PortID   string `json:"port_id"`
  ID       string `json:"id"`
}

// RouterData contains interesting data from the Openstack
// router.{create,update,delete}.{start,end} notifications.
type RouterData struct {
  EventType   string `json:"event_type"`
  Timestamp   string `json:"timestamp"`
  Priority    string `json:"priority"`
  MessageID   string `json:"message_id"`
  ProjectID   string `json:"_context_project_id"`
  ProjectName string `json:"_context_project_name"`
  UserID      string `json:"_context_user_id"`
  UserName    string `json:"_context_user_name"`
  Payload     struct {
    Router          RouterInfo          `json:"router"`
    RouterID        string              `json:"router_id"`
    RouterInterface RouterInterfaceInfo `json:"router_interface"`
  }
}

////////////////////////////////////////////////////////////////////////////////
//                          Identity events data
////////////////////////////////////////////////////////////////////////////////

// IdentityGroupData identifies fields of interest from Openstack
// identity.group.{created,updated,deleted} notifications.
type IdentityGroupData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    GroupID string `json:"resource_info"`
  }
}

// IdentityProjectData identifies fields of interest from Openstack
// identity.project.{created,updated,deleted} notifications.
type IdentityProjectData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    ProjectID string `json:"resource_info"`
  }
}

// IdentityRoleData identifies fields of interest from Openstack
// identity.role.{created,updated,deleted} notifications.
type IdentityRoleData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    RoleID string `json:"resource_info"`
  }
}

// IdentityUserData identifies fields of interest from Openstack
// identity.user.{created,updated,deleted} notifications.
type IdentityUserData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    UserID string `json:"resource_info"`
  }
}

// IdentityTrustData identifies fields of interest from Openstack
// identity.trust.{created,deleted} notifications.
type IdentityTrustData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    TrustID string `json:"resource_info"`
  }
}

// IdentityRegionData identifies fields of interest from Openstack
// identity.region.{created,updated,deleted} notifications.
type IdentityRegionData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    RegionID string `json:"resource_info"`
  }
}

// IdentityEndpointData identifies fields of interest from Openstack
// identity.endpoint.{created,updated,deleted} notifications.
type IdentityEndpointData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    EndpointID string `json:"resource_info"`
  }
}

// IdentityServiceData identifies fields of interest from Openstack
// identity.service.{created,updated,deleted} notifications.
type IdentityServiceData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    ServiceID string `json:"resource_info"`
  }
}

// IdentityPolicyData identifies fields of interest from Openstack
// identity.policy.{created,updated,deleted} notifications.
type IdentityPolicyData struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
  Payload   struct {
    PolicyID string `json:"resource_info"`
  }
}

////////////////////////////////////////////////////////////////////////////////
//                           Image events data
////////////////////////////////////////////////////////////////////////////////

// ImageAttributes identifies attributes of interest from Openstack
// image.{create,prepare,upload,activate,send,update,delete} notifications.
type ImageAttributes struct {
  EventType string `json:"event_type"`
  Timestamp string `json:"timestamp"`
  Priority  string `json:"priority"`
  MessageID string `json:"message_id"`
}

// ImageData identifies fields of interest from Openstack
// image.{create,prepare,upload,activate,update,delete} notifications.
type ImageData struct {
  ImageAttributes
  Payload struct {
    ImageID    string `json:"id"`
    Name       string `json:"name"`
    DiskFormat string `json:"disk_format"`
    Owner      string `json:"owner"`
  }
}

// ImageSendData identifies fields of interest from Openstack
// image.send notifications.
type ImageSendData struct {
  ImageAttributes
  Payload struct {
    ReceiverTenantID string `json:"receiver_tenant_id"`
    ReceiverUserID   string `json:"receiver_user_id"`
    DestinationIP    string `json:"destination_ip"`
    ImageID          string `json:"image_id"`
    Owner            string `json:"owner_id"`
  }
}
