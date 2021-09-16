/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"

	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/util"
	"github.com/apache/camel-k/pkg/util/defaults"
)

// ToJSON marshal to json format
func ToJSON(value runtime.Object) ([]byte, error) {
	return json.Marshal(value)
}

// ToYAML marshal to yaml format
func ToYAML(value runtime.Object) ([]byte, error) {
	data, err := ToJSON(value)
	if err != nil {
		return nil, err
	}

	return util.JSONToYAML(data)
}

func MeteringLabels(integration string) map[string]string {
	var labels = map[string]string{
		v1.IntegrationLabel: integration,
		"com.company":       "Red_Hat",
		"rht.prod_name":     "Red_Hat_Integration",
		"rht.prod_ver":      defaults.RHIntegrationVersion,
		"rht.comp":          "Camel-K",
		"rht.comp_ver":      defaults.Version,
		"rht.subcomp":       integration,
		"rht.subcomp_t":     "application",
	}
	return labels
}
