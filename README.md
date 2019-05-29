# Networkmachinery-operators

NetworkMachinery is currently a PoC  meant for providing a catalog of network troubleshooting operators ( currently in-tree).
 
Below are some of Network machinery Operators and CRDs that are either implemented or planned to be implemented:

- [x] Network Connectivity Operator, with CRDs, 
    - NetworkConnectivityTest
- [x] Network Traffic Shaper, with CRDs:
    - NetworkTrafficShaper 
- [x] Network Monitoring, with CRDs:
    - NetworkMonitor 
    - NetworkNotification
- [ ] Network Controller, with CRDs: 
    - NetworkControl 
- [ ] Network Module Validator
- [ ] CNIBenchmark
- [ ] CNIPerformance
- [ ] Network Scalability Tester (IPVS and IPTables)

## How to use

- You can use the `Make` to package and build the container images via `make package && make tag && make push`.
- To deploy to a cluster, have a look at an example helm chart under the kubernetes director.

```bash
kubernetes
└── networkconnectivity
    ├── Chart.yaml
    ├── templates
    │   ├── deployment.yaml
    │   ├── rbac.yaml
    │   └── serviceaccount.yaml
    └── values.yaml
```

- To test and build locally, run `make <controller-name>` for example: 

```bash
  make start-network-monitor
  make start-networkconnectivity-test
  make start-network-control-controller
  make start-network-traffic-shaper
```