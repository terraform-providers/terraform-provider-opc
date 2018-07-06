Oracle Cloud Infrastructure Load Balancer Classic Client
========================================================

Client implementation for the [REST API for Oracle Cloud Infrastructure Load Balancing Classic](https://docs.oracle.com/en/cloud/iaas/load-balancer-cloud/lbrap/op-vlbrs-compute_region-virtual_load_balancer_resource_id-originserverpools-post.html)

The `LBaaSClient` is the base client implementation for the Load Balancer Classic APIs, but is not intended to be use directly. Specialized clients are implemented for different LBaaS Service resources:

- *`LBaaSClient`* - base implementation
  - `SSLCertificateClient` - for **SSL Certificates**
  - `LoadBalancerClient` - for the main **Load Balancer** resource
  - *`LBaaSResourceClient`* - base client implementation for child resources of a Load Balancer instance:
    - `PolicyClient` - for **Policies**
    - `ListenerClient` - for **Listeners**
    - `OriginServerPoolClient` - for **Origin Server Pools**


Testing
-------

Setup the testing environment according as covered in [Running the SDK Integration Tests](../README.md#running-the-sdk-integration-tests)

To run all the Load Balancer Classic client acceptance tests

```
$ make testacc TEST="./lbaas"
```

To run a single test

```
$ make testacc TEST="./lbaas" TESTARG="-run=TestAccLoadBalancerLifeCycle"
```

### Settings the Test Region

The Load Balancer resources are created in a specific region. The default region used by the tests is `uscom-central-1`. To use a different region set the environment variable `OCP_TEST_LBAAS_REGION` e.g.

```
$ OPC_TEST_LBAAS_REGION=uscom-east-1  make testacc TEST="./lbaas"
```


### Speeding up testing during development

The full Lifecycle tests of the Listener, Origin Server Pool and Policy resources each create a separate parent Load Balance instance.  As the creation and destruction of the Load Balancer can take time a shortcut for development testing is provided:

To speed up testing of the Load Balancer child resources set the environment variable `OPC_TEST_USE_EXISTING_LB` to the ID of an existing Load Balancer instance in the format `<region>/<name>`, e.g.

```
$ OPC_TEST_USE_EXISTING_LB=uscom-central-1/lb1 make test acc TEST="./lbaas" TESTARG="-run=TestAccListenerLifeCycle"
```
