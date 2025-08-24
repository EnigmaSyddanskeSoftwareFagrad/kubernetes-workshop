Kubernetes requires a CSI driver to persist storage. One of such drivers is Longhorn local-path-provisioner.
It uses the host filesystem as a storage backend, meaning that it mounts itself on your host machine and stores data directly on to the filesystem.

First install the longhorn local-path-provisioner using the following command:

kubectl

```
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.26/deploy/local-path-storage.yaml
```

verify it works by running the following command:
```
kubectl get storageclass
```

```
kubectl get pods -n local-path-storage
```

If you are using a talos cluster, then your namespace enforces unprivileged containers for security reasons.
This enforcement can be disabled by overwriting the label on the namespace

```
kubectl label namespace local-path-storage pod-security.kubernetes.io/enforce=privileged --overwrite
```

If this does not work, try overwriting the label again on the namespace of your pods.
