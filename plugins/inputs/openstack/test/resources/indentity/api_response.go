package indentity

func CreateTokenResponseBody(keystoneEndpoint string, novaEndpoint string, cinderEndpoint string,neutronEndpoint string ) string {
	return `
{
  "token": {
    "methods": [
      "password"
    ],
    "user": {
      "domain": {
        "id": "default",
        "name": "Default"
      },
      "id": "ceeefa7882b34847ae15763607ba0d69",
      "name": "admin",
      "password_expires_at": null
    },
    "audit_ids": [
      "HqWauhzPQuSG5H1o52QUUA"
    ],
    "expires_at": "2019-07-22T12:57:46.000000Z",
    "issued_at": "2019-07-22T11:57:46.000000Z",
    "project": {
      "domain": {
        "id": "default",
        "name": "Default"
      },
      "id": "7e985781250646e781010e3a31364590",
      "name": "admin"
    },
    "is_domain": false,
    "roles": [
      {
        "id": "f63ae2da0a83427dbbdbf749fe5860c6",
        "name": "admin"
      },
      {
        "id": "fdca3f78d8bd4c2ebdde2ae418ff4480",
        "name": "reader"
      },
      {
        "id": "46cede6266354be5b2905ee08708c4ab",
        "name": "member"
      }
    ],
    "catalog": [
      {
        "endpoints": [
          {
            "id": "023d2233bb024ae3badb1669dd02ef44",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "`+neutronEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "40022bf32ca840f08162c3403fd44fcd",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "`+neutronEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "757f5b70a5374c8f82be07860ba9e47d",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "`+neutronEndpoint+`",
            "region": "RegionOne"
          }
        ],
        "id": "41a38a41279e42be9e4912e81af0de9a",
        "type": "network",
        "name": "neutron"
      },
      {
        "endpoints": [
          {
            "id": "19725c659e3046f196df93694d3222c6",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "4a5e037a50b6497e99fa766612062cae",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "de3577cb88804e8397b220fc23d9e220",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          }
        ],
        "id": "620c9583cd00478ba46ababa2f48ac07",
        "type": "volumev2",
        "name": "cinderv2"
      },
      {
        "endpoints": [
          {
            "id": "07b8ec57044f4e80a89038331b644e17",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "`+novaEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "8a96a73bfe9343a4a511fa72061673d7",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "`+novaEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "c971e29563e94ca186e41691f1a5b43d",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "`+novaEndpoint+`",
            "region": "RegionOne"
          }
        ],
        "id": "68fc9e41339a45e68ea424e148d86e43",
        "type": "compute",
        "name": "nova"
      },
      {
        "endpoints": [
          {
            "id": "74158342913c456297eccf651e72a01d",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "`+keystoneEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "8816cd6fcdae47d59da697cf7b98f47d",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "`+keystoneEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "b867ee12aa35496aa31c98b6487f9e66",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "`+keystoneEndpoint+`",
            "region": "RegionOne"
          }
        ],
        "id": "af54a98106b643c593482ea0929fb929",
        "type": "identity",
        "name": "keystone"
      },
      {
        "endpoints": [
          {
            "id": "7c114538231243d59a9cbe6db76db08e",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "88fdd64a32b14cc3bcfb7dfcc854372f",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          },
          {
            "id": "9cf45c536fba4f0ab441e5cc17365ecb",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "`+cinderEndpoint+`",
            "region": "RegionOne"
          }
        ],
        "id": "f7e43ef00bbb4b48a71e958ad63b8fb2",
        "type": "volumev3",
        "name": "cinderv3"
      }
    ]
  }
}
`
}


func ServiceListResponseBody()  string{
	return `
{"services":[{"description":"OpenStack Networking","name":"neutron","id":"41a38a41279e42be9e4912e81af0de9a","type":"network","enabled":true,"links":{"self":"https://controller:5000/v3/services/41a38a41279e42be9e4912e81af0de9a"}},{"description":"OpenStack Block Storage","name":"cinderv2","id":"620c9583cd00478ba46ababa2f48ac07","type":"volumev2","enabled":true,"links":{"self":"https://controller:5000/v3/services/620c9583cd00478ba46ababa2f48ac07"}},{"description":"OpenStack Compute","name":"nova","id":"68fc9e41339a45e68ea424e148d86e43","type":"compute","enabled":true,"links":{"self":"https://controller:5000/v3/services/68fc9e41339a45e68ea424e148d86e43"}},{"name":"keystone","id":"af54a98106b643c593482ea0929fb929","type":"identity","enabled":true,"links":{"self":"https://controller:5000/v3/services/af54a98106b643c593482ea0929fb929"}},{"description":"Placement API","name":"placement","id":"baa1bcb28bba4a0495e8eb53c04757c2","type":"placement","enabled":true,"links":{"self":"https://controller:5000/v3/services/baa1bcb28bba4a0495e8eb53c04757c2"}},{"description":"OpenStack Image","name":"glance","id":"d77bc52b401849d8a629e075275dc7c2","type":"image","enabled":true,"links":{"self":"https://controller:5000/v3/services/d77bc52b401849d8a629e075275dc7c2"}},{"description":"OpenStack Block Storage","name":"cinderv3","id":"f7e43ef00bbb4b48a71e958ad63b8fb2","type":"volumev3","enabled":true,"links":{"self":"https://controller:5000/v3/services/f7e43ef00bbb4b48a71e958ad63b8fb2"}}],"links":{"next":null,"self":"https://controller:5000/v3/services","previous":null}}
`
}

func ProjectListResponseBody() string{
	return `
{
  "projects": [
    {
      "id": "33b03d1e28404ce68c8cf8c91506465b",
      "name": "demo",
      "domain_id": "default",
      "description": "Demo Project",
      "enabled": true,
      "parent_id": "default",
      "is_domain": false,
      "tags": [],
      "links": {
        "self": "https://controller:5000/v3/projects/33b03d1e28404ce68c8cf8c91506465b"
      }
    },
    {
      "id": "8794dc92419b4c65b43654aa39225aba",
      "name": "7e985781250646e781010e3a31364590-5c941b4b-a98c-4ebd-8f08-6a121dc",
      "domain_id": "400621da3fd64d85a935d159eeb17ce2",
      "description": "Heat stack user project",
      "enabled": true,
      "parent_id": "400621da3fd64d85a935d159eeb17ce2",
      "is_domain": false,
      "tags": [],
      "links": {
        "self": "https://controller:5000/v3/projects/8794dc92419b4c65b43654aa39225aba"
      }
    }
  ],
  "links": {
    "next": null,
    "self": "https://controller:5000/v3/projects",
    "previous": null
  }
}`
}

func UserListResponseBody() string{
	return `
{"links":{"next":null,"self":"https://controller:5000/v3/users","previous":null},"users":[{"id":"38ea3c7432fe45a09698b734775563d0","name":"demo","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/38ea3c7432fe45a09698b734775563d0"}},{"id":"621b8d51954846dd9001917a39711010","name":"nova","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/621b8d51954846dd9001917a39711010"}},{"id":"994c050126c94052a4baaabbe7eaa491","name":"neutron","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/994c050126c94052a4baaabbe7eaa491"}},{"id":"ac59933ac2ab4336b9fe3d97409d3847","name":"glance","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/ac59933ac2ab4336b9fe3d97409d3847"}},{"id":"b0b9979970bf43c5a552f43448e09447","name":"placement","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/b0b9979970bf43c5a552f43448e09447"}},{"id":"ceeefa7882b34847ae15763607ba0d69","name":"admin","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/ceeefa7882b34847ae15763607ba0d69"}},{"id":"de63cd9c75ac4373a9ff75a188413056","name":"cinder","domain_id":"default","enabled":true,"password_expires_at":null,"options":{},"links":{"self":"https://controller:5000/v3/users/de63cd9c75ac4373a9ff75a188413056"}}]}
`
}

func GroupListResponseBody()   string{
	return `
{"groups": [{"id": "b67020b4f2eb4165a1834f7fc51e0369", "name": "a", "domain_id": "default", "description": "", "links": {"self": "https://controller:5000/v3/groups/b67020b4f2eb4165a1834f7fc51e0369"}}], "links": {"next": null, "self": "https://controller:5000/v3/groups", "previous": null}}
`
}