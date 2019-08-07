package networking

import "github.com/influxdata/telegraf/plugins/inputs/openstack/test/resources"

func AgentsListResponseBody() string {
	return `
{"agents": [{"id": "61e842df-885c-4122-b071-970ec6764df6", "agent_type": "Open vSwitch agent", "binary": "neutron-openvswitch-agent", "topic": "N/A", "host": "compute02", "admin_state_up": true, "created_at": "2019-06-18 10:26:58", "started_at": "2019-07-22 08:20:26", "heartbeat_timestamp": "2019-07-23 10:55:38", "description": null, "resources_synced": null, "availability_zone": null, "alive": true, "configurations": {"arp_responder_enabled": false, "bridge_mappings": {"provider": "br-provider"}, "datapath_type": "system", "devices": 0, "enable_distributed_routing": false, "extensions": [], "in_distributed_mode": false, "integration_bridge": "br-int", "l2_population": false, "log_agent_heartbeats": false, "ovs_capabilities": {"datapath_types": ["netdev", "system"], "iface_types": ["erspan", "geneve", "gre", "internal", "ip6erspan", "ip6gre", "lisp", "patch", "stt", "system", "tap", "vxlan"]}, "ovs_hybrid_plug": true, "resource_provider_bandwidths": {}, "resource_provider_inventory_defaults": {"allocation_ratio": 1.0, "min_unit": 1, "step_size": 1, "reserved": 0}, "tunnel_types": [], "tunneling_ip": null, "vhostuser_socket_dir": "/var/run/openvswitch"}}, {"id": "7bfbd6d0-86ac-4efc-b1c8-2b438b7c0a8e", "agent_type": "Metadata agent", "binary": "neutron-metadata-agent", "topic": "N/A", "host": "controller", "admin_state_up": true, "created_at": "2019-06-18 07:20:41", "started_at": "2019-07-22 08:20:02", "heartbeat_timestamp": "2019-07-23 10:55:38", "description": null, "resources_synced": null, "availability_zone": null, "alive": true, "configurations": {"log_agent_heartbeats": false, "metadata_proxy_socket": "/var/lib/neutron/metadata_proxy", "nova_metadata_host": "127.0.0.1", "nova_metadata_port": 8775}}, {"id": "985b0c52-4f47-48af-801a-025dce3b9427", "agent_type": "Open vSwitch agent", "binary": "neutron-openvswitch-agent", "topic": "N/A", "host": "compute01", "admin_state_up": true, "created_at": "2019-06-18 08:00:51", "started_at": "2019-07-22 08:20:30", "heartbeat_timestamp": "2019-07-23 10:55:43", "description": null, "resources_synced": null, "availability_zone": null, "alive": true, "configurations": {"arp_responder_enabled": false, "bridge_mappings": {"provider": "br-provider"}, "datapath_type": "system", "devices": 0, "enable_distributed_routing": false, "extensions": [], "in_distributed_mode": false, "integration_bridge": "br-int", "l2_population": false, "log_agent_heartbeats": false, "ovs_capabilities": {"datapath_types": ["netdev", "system"], "iface_types": ["erspan", "geneve", "gre", "internal", "ip6erspan", "ip6gre", "lisp", "patch", "stt", "system", "tap", "vxlan"]}, "ovs_hybrid_plug": true, "resource_provider_bandwidths": {}, "resource_provider_inventory_defaults": {"allocation_ratio": 1.0, "min_unit": 1, "step_size": 1, "reserved": 0}, "tunnel_types": [], "tunneling_ip": null, "vhostuser_socket_dir": "/var/run/openvswitch"}}, {"id": "b76e8703-1a67-4423-8995-97745203fcac", "agent_type": "Open vSwitch agent", "binary": "neutron-openvswitch-agent", "topic": "N/A", "host": "controller", "admin_state_up": true, "created_at": "2019-06-18 07:20:37", "started_at": "2019-07-22 08:19:26", "heartbeat_timestamp": "2019-07-23 10:55:32", "description": null, "resources_synced": null, "availability_zone": null, "alive": true, "configurations": {"arp_responder_enabled": false, "bridge_mappings": {"provider": "br-provider"}, "datapath_type": "system", "devices": 1, "enable_distributed_routing": false, "extensions": [], "in_distributed_mode": false, "integration_bridge": "br-int", "l2_population": false, "log_agent_heartbeats": false, "ovs_capabilities": {"datapath_types": ["netdev", "system"], "iface_types": ["erspan", "geneve", "gre", "internal", "ip6erspan", "ip6gre", "lisp", "patch", "stt", "system", "tap", "vxlan"]}, "ovs_hybrid_plug": true, "resource_provider_bandwidths": {}, "resource_provider_inventory_defaults": {"allocation_ratio": 1.0, "min_unit": 1, "step_size": 1, "reserved": 0}, "tunnel_types": [], "tunneling_ip": null, "vhostuser_socket_dir": "/var/run/openvswitch"}}, {"id": "fdd27f14-a71a-417e-a31e-a761e98c8aec", "agent_type": "DHCP agent", "binary": "neutron-dhcp-agent", "topic": "dhcp_agent", "host": "controller", "admin_state_up": true, "created_at": "2019-06-18 07:20:40", "started_at": "2019-07-22 08:20:05", "heartbeat_timestamp": "2019-07-23 10:55:31", "description": null, "resources_synced": null, "availability_zone": "nova", "alive": true, "configurations": {"dhcp_driver": "neutron.agent.linux.dhcp.Dnsmasq", "dhcp_lease_duration": 86400, "log_agent_heartbeats": false, "networks": 1, "ports": 1, "subnets": 1}}]}
`
}

func FloatingIpsListResponseBody() string {
	return `
{
  "floatingips": [
    {
      "router_id": null,
      "description": "for test",
      "dns_domain": "my-domain.org.",
      "dns_name": "myfip2",
      "created_at": "2016-12-21T11:55:50Z",
      "updated_at": "2016-12-21T11:55:53Z",
      "revision_number": 2,
      "project_id": "` + resources.ProjectId + `",
      "tenant_id": "` + resources.ProjectId + `",
      "floating_network_id": "376da547-b977-4cfe-9cba-275c80debf57",
      "fixed_ip_address": null,
      "floating_ip_address": "172.24.4.227",
      "port_id": null,
      "id": "61cea855-49cb-4846-997d-801b70c71bdd",
      "status": "DOWN",
      "port_details": null,
      "tags": [
        "tag1,tag2"
      ],
      "port_forwardings": []
    }
  ]
}
`
}
func IpAvailabilityListResponseBody() string {
	return `
{"network_ip_availabilities": [{"network_id": "e30aada8-68a8-4f58-ab3b-c8c4d36a35fd", "network_name": "provider", "tenant_id": "7e985781250646e781010e3a31364590", "project_id": "7e985781250646e781010e3a31364590", "subnet_ip_availability": [{"subnet_id": "fe6c1225-e190-4f5b-89fb-0738e48d9d44", "ip_version": 4, "cidr": "192.168.33.0/24", "subnet_name": "provider-sub", "used_ips": 1, "total_ips": 52}], "used_ips": 1, "total_ips": 52}]}
`
}
func NetworkListResponseBody() string {
	return `
{"networks":[{"id":"e30aada8-68a8-4f58-ab3b-c8c4d36a35fd","name":"provider","tenant_id":"7e985781250646e781010e3a31364590","admin_state_up":true,"mtu":1500,"status":"ACTIVE","subnets":["fe6c1225-e190-4f5b-89fb-0738e48d9d44"],"shared":true,"availability_zone_hints":[],"availability_zones":["nova"],"ipv4_address_scope":null,"ipv6_address_scope":null,"router:external":true,"description":"","port_security_enabled":true,"tags":[],"created_at":"2019-06-18T11:21:28Z","updated_at":"2019-07-10T00:27:32Z","revision_number":4,"project_id":"7e985781250646e781010e3a31364590","provider:network_type":"flat","provider:physical_network":"provider","provider:segmentation_id":null}]}
`
}
func QuotasListResponseBody() string {
	return `
{"quota": {"network": {"limit": 100, "used": 1, "reserved": 0}, "subnet": {"limit": 100, "used": 1, "reserved": 0}, "subnetpool": {"limit": -1, "used": 0, "reserved": 0}, "port": {"limit": 500, "used": 1, "reserved": 0}, "rbac_policy": {"limit": 10, "used": 2, "reserved": 0}, "security_group": {"limit": 10, "used": 1, "reserved": 0}, "security_group_rule": {"limit": 100, "used": 6, "reserved": 0}}}
`
}
