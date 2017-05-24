# kube-aws-labeller

Add node labels to Kubernetes workers based on EC2 properties and tags, to allow deployments to be targeted to subsets of worker instances.

Currently, it adds a label "spot-instance", which is either true or false depending on the lifecycle of the EC2 instance.

The best way to deploy the kube-aws-labeller is with a DaemonSet across all instances:

A Docker image based on this Git repo is available at https://hub.docker.com/r/grrywlsn/kube-aws-labeller

```
kind: DaemonSet
apiVersion: extensions/v1beta1
metadata:
  name: worker-labeller
  namespace: kube-system
  labels:
    k8s-app: worker-labeller
spec:
  selector:
    matchLabels:
      k8s-app: worker-labeller
  template:
    metadata:
      labels:
        k8s-app: worker-labeller
    spec:
      containers:
        - name: node-labeller
          image: grrywlsn/kube-aws-labeller:latest
          imagePullPolicy: Always
          env:
          - name: AWS_REGION
            value: "eu-west-1"
```
