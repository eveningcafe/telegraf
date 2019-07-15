package scheduler

// StoragePool represents an individual StoragePool retrieved from the
// schedulerstats API.
type StoragePool struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
}

// Capabilities represents the information of an individual StoragePool.
type Capabilities struct {
	// The following fields should be present in all storage drivers.
	DriverVersion     string  `json:"driver_version"`
	FreeCapacityGB    float64 `json:"-"`
	StorageProtocol   string  `json:"storage_protocol"`
	TotalCapacityGB   float64 `json:"-"`
	VendorName        string  `json:"vendor_name"`
	VolumeBackendName string  `json:"volume_backend_name"`

	// The following fields are optional and may have empty values depending
	// on the storage driver in use.
	ReservedPercentage       int64   `json:"reserved_percentage"`
	LocationInfo             string  `json:"location_info"`
	QoSSupport               bool    `json:"QoS_support"`
	ProvisionedCapacityGB    float64 `json:"provisioned_capacity_gb"`
	MaxOverSubscriptionRatio string  `json:"max_over_subscription_ratio"`
	ThinProvisioningSupport  bool    `json:"thin_provisioning_support"`
	ThickProvisioningSupport bool    `json:"thick_provisioning_support"`
	TotalVolumes             int64   `json:"total_volumes"`
	FilterFunction           string  `json:"filter_function"`
	GoodnessFuction          string  `json:"goodness_function"`
	Multiattach              bool    `json:"multiattach"`
	SparseCopyVolume         bool    `json:"sparse_copy_volume"`
}