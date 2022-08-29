#!/bin/sh

kubectl -n spire exec spire-server-0 -- \
/opt/spire/bin/spire-server entry create \
-ttl 7200 \
-spiffeID spiffe://example.org/ns/spire/sa/spire-agent \
-selector k8s_psat:cluster:nsm-cluster \
-selector k8s_psat:agent_ns:spire \
-selector k8s_psat:agent_sa:spire-agent \
-node
