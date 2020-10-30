---
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      networkservicemesh.io/app: "icmp-responder-nsc"
  replicas: 4
  template:
    metadata:
      labels:
        networkservicemesh.io/app: "icmp-responder-nsc"
    spec:
      serviceAccount: nsmgr-acc
      containers:
        - name: alpine-img
          image: alpine:latest
          securityContext:
            privileged: true
          command: ['tail', '-f', '/dev/null']
metadata:
  name: icmp-responder-nsc
  namespace: {{ .Release.Namespace }}
  annotations:
    ns.networkservicemesh.io: icmp-responder?app=icmp
