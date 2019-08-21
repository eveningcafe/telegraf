package placement

func ListResourcesProviderResponseBody() string {
	return `
{
  "resource_providers": [
    {
      "uuid": "6750f406-55bd-4c96-85d3-d5ca7f58fd4a",
      "name": "compute05",
      "generation": 6,
      "links": [
        {
          "rel": "self",
          "href": "/resource_providers/6750f406-55bd-4c96-85d3-d5ca7f58fd4a"
        },
        {
          "rel": "inventories",
          "href": "/resource_providers/6750f406-55bd-4c96-85d3-d5ca7f58fd4a/inventories"
        },
        {
          "rel": "usages",
          "href": "/resource_providers/6750f406-55bd-4c96-85d3-d5ca7f58fd4a/usages"
        }
      ]
    }
  ]
}
`
}

func GetResourcesInventoriesResponseBody() string {
	return `
{
  "resource_provider_generation": 6,
  "inventories": {
    "VCPU": {
      "total": 32,
      "reserved": 0,
      "min_unit": 1,
      "max_unit": 32,
      "step_size": 1,
      "allocation_ratio": 16
    },
    "MEMORY_MB": {
      "total": 31785,
      "reserved": 512,
      "min_unit": 1,
      "max_unit": 31785,
      "step_size": 1,
      "allocation_ratio": 1.5
    },
    "DISK_GB": {
      "total": 22354,
      "reserved": 0,
      "min_unit": 1,
      "max_unit": 22354,
      "step_size": 1,
      "allocation_ratio": 1
    }
  }
}
`
}

func GetResourcesUsagesResponseBody()  string{
	return `
{"resource_provider_generation": 6, "usages": {"VCPU": 16, "MEMORY_MB": 16384, "DISK_GB": 0}}
`
}