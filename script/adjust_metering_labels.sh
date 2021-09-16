#!/bin/sh

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

location=$(dirname $0)
rootdir=$location/../

files="config/manager/operator-deployment.yaml helm/camel-k/templates/operator.yaml"

# the RH_INTEGRATION_PRODUCT_VERSION
rhi_version=$(grep ^RH_INTEGRATION_PRODUCT_VERSION script/Makefile|cut -d' ' -f3)

# the Camel K version
ck_version=$(grep ^'VERSION ?=' script/Makefile|cut -d' ' -f3)

echo "Checking Red Hat Integration metering labels"
for f in ${files}; do
    rhi_v=$(grep rht.prod_ver ${f} |awk '{print $2}')
    if [[ ${rhi_v} != ${rhi_version} ]]; then
        echo "Adjust ${f} to Red Hat Integration version ${rhi_version}"
        sed -i "/rht.prod_ver/ s/:.*/: ${rhi_version}/g" ${f}
    fi

    ck_v=$(grep rht.comp_ver ${f} |awk '{print $2}')
    if [[ ${ck_v} != ${ck_version} ]]; then
        echo "Adjust ${f} to Camel K version ${ck_version}"
        sed -i "/rht.comp_ver/ s/:.*/: ${ck_version}/g" ${f}
    fi
done


