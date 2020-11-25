# Activate SCTP on cluster

edit load-sctp-module.yaml label must match master in case of compact cluster installation (master and worker on same host), worker otherwise.

```bash
oc create -f load-sctp-module.yaml
```

Wait for configuration to apply (nodes will reboot):

```bash
watch -n 10 oc get nodes
```
