package resources

const GetTokenResponseBody  = `
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
            "url": "http://controller:9696",
            "region": "RegionOne"
          },
          {
            "id": "40022bf32ca840f08162c3403fd44fcd",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:9696",
            "region": "RegionOne"
          },
          {
            "id": "757f5b70a5374c8f82be07860ba9e47d",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "http://controller:9696",
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
            "url": "http://controller:8776/v2/7e985781250646e781010e3a31364590",
            "region": "RegionOne"
          },
          {
            "id": "4a5e037a50b6497e99fa766612062cae",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:8776/v2/7e985781250646e781010e3a31364590",
            "region": "RegionOne"
          },
          {
            "id": "de3577cb88804e8397b220fc23d9e220",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "http://controller:8776/v2/7e985781250646e781010e3a31364590",
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
            "url": "http://controller:8774/v2.1",
            "region": "RegionOne"
          },
          {
            "id": "8a96a73bfe9343a4a511fa72061673d7",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "http://controller:8774/v2.1",
            "region": "RegionOne"
          },
          {
            "id": "c971e29563e94ca186e41691f1a5b43d",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:8774/v2.1",
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
            "url": "https://controller:5000/v3/",
            "region": "RegionOne"
          },
          {
            "id": "8816cd6fcdae47d59da697cf7b98f47d",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "https://controller:5000/v3/",
            "region": "RegionOne"
          },
          {
            "id": "b867ee12aa35496aa31c98b6487f9e66",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "https://controller:5000/v3/",
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
            "id": "7280b3bb3659438098afdaf090b1e940",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "http://controller:8778",
            "region": "RegionOne"
          },
          {
            "id": "ae977e2672ed48ffaa0681aecb93bada",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:8778",
            "region": "RegionOne"
          },
          {
            "id": "bd087b7c26f74312beeebe607d7fe362",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "http://controller:8778",
            "region": "RegionOne"
          }
        ],
        "id": "baa1bcb28bba4a0495e8eb53c04757c2",
        "type": "placement",
        "name": "placement"
      },
      {
        "endpoints": [
          {
            "id": "245ecff116f441448909dd9558edf5ab",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "http://controller:9292",
            "region": "RegionOne"
          },
          {
            "id": "406128b38213463d8296be79b545e892",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "http://controller:9292",
            "region": "RegionOne"
          },
          {
            "id": "95f8826a45924be5b6992bc81e586c2e",
            "interface": "internal",
            "region_id": "NorthVN",
            "url": "http://controller:9292",
            "region": "NorthVN"
          },
          {
            "id": "ab1597824d6a49bb82c42dc588d76fd2",
            "interface": "public",
            "region_id": "NorthVN",
            "url": "http://controller:9292",
            "region": "NorthVN"
          },
          {
            "id": "eb7bd4ea843e418590f85008fe3cda9f",
            "interface": "admin",
            "region_id": "NorthVN",
            "url": "http://controller:9292",
            "region": "NorthVN"
          },
          {
            "id": "f5dfcdc3be654ac6bf7b182b1219b21d",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:9292",
            "region": "RegionOne"
          }
        ],
        "id": "d77bc52b401849d8a629e075275dc7c2",
        "type": "image",
        "name": "glance"
      },
      {
        "endpoints": [
          {
            "id": "7c114538231243d59a9cbe6db76db08e",
            "interface": "public",
            "region_id": "RegionOne",
            "url": "http://controller:8776/v3/7e985781250646e781010e3a31364590",
            "region": "RegionOne"
          },
          {
            "id": "88fdd64a32b14cc3bcfb7dfcc854372f",
            "interface": "admin",
            "region_id": "RegionOne",
            "url": "http://controller:8776/v3/7e985781250646e781010e3a31364590",
            "region": "RegionOne"
          },
          {
            "id": "9cf45c536fba4f0ab441e5cc17365ecb",
            "interface": "internal",
            "region_id": "RegionOne",
            "url": "http://controller:8776/v3/7e985781250646e781010e3a31364590",
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
