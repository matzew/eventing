#!/usr/bin/env bash

# shellcheck disable=SC1090
source "$(dirname "$0")/../vendor/knative.dev/hack/e2e-tests.sh"
source "$(dirname "$0")/e2e-common.sh"

set -Eeuox pipefail

export TEST_IMAGE_TEMPLATE="${IMAGE_FORMAT//\$\{component\}/knative-eventing-test-{{.Name}}}"

env

scale_up_workers || exit 1

failed=0

(( !failed )) && install_serverless || failed=1

(( !failed )) && install_knative_eventing || failed=1

(( !failed )) && install_tracing || failed=1
echo "**************************************"
echo "***  OS    RUN REKT TESTS          ***"
echo "**************************************"
(( !failed )) && run_e2e_rekt_tests || failed=1
echo "**************************************"
echo "*** OS UNINSTALL KNATIVE EVENTING  ***"
echo "**************************************"
(( !failed )) && uninstall_knative_eventing || failed=1
echo "**************************************"
echo "***  DUMP CLUSTER? FAILED=$failed  ***"
echo "**************************************"
(( failed )) && dump_cluster_state

(( failed )) && exit 1

success

