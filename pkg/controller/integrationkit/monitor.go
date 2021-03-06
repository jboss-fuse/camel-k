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

package integrationkit

import (
	"context"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util/digest"
)

// NewMonitorAction creates a new monitoring handling action for the kit
func NewMonitorAction() Action {
	return &monitorAction{}
}

type monitorAction struct {
	baseAction
}

func (action *monitorAction) Name() string {
	return "monitor"
}

func (action *monitorAction) CanHandle(kit *v1alpha1.IntegrationKit) bool {
	return kit.Status.Phase == v1alpha1.IntegrationKitPhaseReady || kit.Status.Phase == v1alpha1.IntegrationKitPhaseError
}

func (action *monitorAction) Handle(ctx context.Context, kit *v1alpha1.IntegrationKit) error {
	hash, err := digest.ComputeForIntegrationKit(kit)
	if err != nil {
		return err
	}
	if hash != kit.Status.Digest {
		action.L.Info("IntegrationKit needs a rebuild")

		target := kit.DeepCopy()
		target.Status.Digest = hash
		target.Status.Phase = v1alpha1.IntegrationKitPhaseBuildSubmitted

		action.L.Info("IntegrationKit state transition", "phase", target.Status.Phase)

		return action.client.Status().Update(ctx, target)
	}

	return nil
}
