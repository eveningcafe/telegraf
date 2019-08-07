# Openstack plugin
This plugin provide metric about cloud build on openstack. Call openstack api to get metric of cloud control plane.
Here is it's flows:
* Call openstack keystone service to authenticate, get token identity and all openstack's endpoint
* Call openstack [keystone,nova,cinder,neutron] endpoint to get metric.

## Metric:

... About more than 50 metrics, you can get it when run unit test. I'm too lazy to list it here.
