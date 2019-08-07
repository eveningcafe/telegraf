package blockstorage


func ServiceListResponseBody() string{
	return `
{"services": [{"binary": "cinder-scheduler", "host": "controller", "zone": "nova", "status": "enabled", "state": "up", "updated_at": "2019-07-23T10:42:01.000000", "disabled_reason": null}, {"binary": "cinder-volume", "host": "controller@ceph", "zone": "nova", "status": "enabled", "state": "up", "updated_at": "2019-07-23T10:42:02.000000", "disabled_reason": null, "replication_status": "disabled", "active_backend_id": null, "frozen": false}]}
`
}

func StoragePoolListResponseBody() string{
	return `
{"pools": [{"name": "controller@ceph#RBD", "capabilities": {"vendor_name": "Open Source", "driver_version": "1.2.0", "storage_protocol": "ceph", "total_capacity_gb": 125.03, "free_capacity_gb": 125.03, "reserved_percentage": 0, "multiattach": true, "thin_provisioning_support": true, "max_over_subscription_ratio": "20.0", "location_info": "ceph:/etc/ceph/ceph.conf:2a0e1996-8aa5-495c-8bdd-02439f539d47:volumes:volumes", "backend_state": "up", "volume_backend_name": "RBD", "replication_enabled": false, "provisioned_capacity_gb": 0, "allocated_capacity_gb": 0, "filter_function": null, "goodness_function": null, "timestamp": "2019-07-22T11:56:28.684921"}}]}
`
}
func QuotasListResponseBody() string{
	return `
{"quota_set": {"volumes": {"limit": 10, "in_use": 0, "reserved": 0}, "per_volume_gigabytes": {"limit": -1, "in_use": 0, "reserved": 0}, "snapshots": {"limit": 10, "in_use": 0, "reserved": 0}, "gigabytes": {"limit": 1000, "in_use": 0, "reserved": 0}, "backups": {"limit": 10, "in_use": 0, "reserved": 0}, "backup_gigabytes": {"limit": 1000, "in_use": 0, "reserved": 0}, "groups": {"limit": 10, "in_use": 0, "reserved": 0}, "id": "33b03d1e28404ce68c8cf8c91506465b"}}
`
}
