#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

check_env_var() {
  if [ -z "${2}" ]; then
    echo "Error: ${1} env var not defined"
    exit 1
  fi
}

location=$(dirname $0)
rootdir=$location/../

cd $rootdir

bundle_dir="${rootdir}bundle"
dockerfile="${bundle_dir}/Dockerfile"

if [ ! -f "${dockerfile}" ]; then
  echo "Error: Cannot find bundle dockerfile."
	exit 1
fi

set +e

check_env_var "CSV_VERSION" ${CSV_VERSION}
check_env_var "OPERATOR_VERSION" ${OPERATOR_VERSION}

MINOR_VERSION=${OPERATOR_VERSION%\.[0-9]}

#
# Use the first line of the insert text to check if it has
# already been inserted
#
precomment="# Tells the pipeline that this is a bundle image and should be delivered via an index image"
if ! grep -q "${precomment}" "${dockerfile}"; then
  sed '/^FROM\ scratch/r'<(cat <<EOF

${precomment}
LABEL com.redhat.delivery.operator.bundle=true

# Tells the pipeline which versions of OpenShift the operator supports.
# This is used to control which index images should include this operator.
LABEL com.redhat.openshift.versions="v4.6"

# The rest of these labels are copies of the same content in annotations.yaml and are needed by OLM
# Note the package name and channels which are very important!
EOF
) -i -- "${dockerfile}"
fi

#
# Use the first line of the insert text to check if it has
# already been inserted
#
postcomment="# Log the camel-k operator version we built this with"
if ! grep -q "${postcomment}" "${dockerfile}"; then
  sed '/^COPY\ tests\/scorecard.*/r'<(cat <<EOF

${postcomment}
LABEL com.redhat.fuse.camel-k.operatorversion=camelk-${MINOR_VERSION}-openshift-rhel-8
LABEL com.redhat.fuse.camel-k.csvversion=${CSV_VERSION}

# This last block are standard Red Hat container labels
LABEL %%
        com.redhat.component="red-hat-camel-k-bundle-container" %%
        version="${OPERATOR_VERSION}" %%
        name="integration/camel-k-rhel8-operator-bundle" %%
        License="ASL 2.0" %%
        io.k8s.display-name="red-hat-camel-k bundle" %%
        io.k8s.description="bundle containing manifests for red-hat-camel-k" %%
        summary="bundle containing manifests for red-hat-camel-k" %%
        maintainer="Thomas Cunningham <tcunning@redhat.com>"
EOF
) -i -- "${dockerfile}"

  # Replace the '%%' above to be the correct backslashes.
  # Cannot add them in the template as they convert the text to a single line
  sed -i 's/%%/\\/g' "${dockerfile}"
fi

set -e
