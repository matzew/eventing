#!/usr/bin/env bash

source $(dirname $0)/resolve.sh

release=$1

artifacts_dir="openshift/release/artifacts/"
rm -rf $artifacts_dir
mkdir -p $artifacts_dir

eventing_core="${artifacts_dir}knative-eventing-core.yaml"
eventing_imc="${artifacts_dir}knative-eventing-imc.yaml"
eventing_mt_broker="${artifacts_dir}knative-eventing-mt-broker.yaml"

if [ "$release" == "ci" ]; then
    image_prefix="registry.ci.openshift.org/openshift/knative-nightly:knative-eventing-"
    tag=""
else
    image_prefix="registry.ci.openshift.org/openshift/knative-${release}:knative-eventing-"
    tag=""
fi

# Eventing core
resolve_resources config/ $eventing_core $image_prefix $tag
resolve_resources config/sugar/ eventing-core-resolved.yaml $image_prefix $tag
cat eventing-core-resolved.yaml >> $eventing_core
rm eventing-core-resolved.yaml

# InMemoryChannel folders
## The root folder
resolve_resources config/channels/in-memory-channel/ imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml
## The configmaps folder
resolve_resources config/channels/in-memory-channel/configmaps imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml
## The deployments folder
resolve_resources config/channels/in-memory-channel/deployments imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml
## The resources folder
resolve_resources config/channels/in-memory-channel/resources imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml
## The roles folder
resolve_resources config/channels/in-memory-channel/roles imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml
## The webhooks folder
resolve_resources config/channels/in-memory-channel/webhooks imc-channel-resolved.yaml $image_prefix $tag
cat imc-channel-resolved.yaml >> $eventing_imc
rm imc-channel-resolved.yaml

# MT Broker
resolve_resources config/brokers/mt-channel-broker/ $eventing_mt_broker $image_prefix $tag
