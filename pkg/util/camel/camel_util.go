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

package camel

import (
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util/log"
	"github.com/blang/semver"
)

// FindBestMatch --
func FindBestMatch(version string, catalogs []v1alpha1.CamelCatalog) (*RuntimeCatalog, error) {

	parseRange, err := semver.ParseRange(version)
	//
	// if the version is not a constraint, use exact match
	//
	if err != nil || parseRange == nil {
		if err != nil {
			log.Debug("Unable to parse constraint: %s, error:\n", version, err.Error())
		}
		if parseRange == nil {
			log.Debug("Unable to parse constraint: %s\n", version)
		}

		return FindExactMatch(version, catalogs)
	}

	return FindBestSemVerMatch(version, catalogs)
}

// FindExactMatch --
func FindExactMatch(version string, catalogs []v1alpha1.CamelCatalog) (*RuntimeCatalog, error) {
	for _, catalog := range catalogs {
		if catalog.Spec.Version == version {
			return NewRuntimeCatalog(catalog.Spec), nil
		}
	}

	return nil, nil
}

func reverse(versions []semver.Version) []semver.Version {
	newVersions := make([]semver.Version, 0, len(versions))
	for i := len(versions) - 1; i >= 0; i-- {
		newVersions = append(newVersions, versions[i])
	}
	return newVersions
}

// FindBestSemVerMatch --
func FindBestSemVerMatch(constraint string, catalogs []v1alpha1.CamelCatalog) (*RuntimeCatalog, error) {
	versions := make([]semver.Version, 0)

	for _, catalog := range catalogs {
		v, err := semver.Parse(catalog.Spec.Version)
		if err != nil {
			log.Debugf("Invalid semver version %s, skip it", catalog.Spec.Version)
			continue
		}

		versions = append(versions, v)
	}

	semver.Sort(versions)
	reversedVersions := reverse(versions)

	for _, v := range reversedVersions {
		ver := v
		constraintRange, err := semver.ParseRange(constraint)
		if err != nil {
			log.Debugf("Could not parse range \"%s\"", constraint)
			continue
		}

		if constraintRange(ver) {
			for _, catalog := range catalogs {
				if catalog.Spec.Version == ver.String() {
					return NewRuntimeCatalog(catalog.Spec), nil
				}
			}
		}
	}

	return nil, nil
}
