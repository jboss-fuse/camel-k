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

package trait

import (
	"strconv"
	"strings"

	"github.com/apache/camel-k/pkg/util/kubernetes"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/metadata"
	"github.com/apache/camel-k/pkg/util/envvar"
	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	knativeServingClassAnnotation    = "autoscaling.knative.dev/class"
	knativeServingMetricAnnotation   = "autoscaling.knative.dev/metric"
	knativeServingTargetAnnotation   = "autoscaling.knative.dev/target"
	knativeServingMinScaleAnnotation = "autoscaling.knative.dev/minScale"
	knativeServingMaxScaleAnnotation = "autoscaling.knative.dev/maxScale"
)

type knativeServiceTrait struct {
	BaseTrait `property:",squash"`
	Class     string `property:"autoscaling-class"`
	Metric    string `property:"autoscaling-metric"`
	Target    *int   `property:"autoscaling-target"`
	MinScale  *int   `property:"min-scale"`
	MaxScale  *int   `property:"max-scale"`
	Auto      *bool  `property:"auto"`
	deployer  deployerTrait
}

func newKnativeServiceTrait() *knativeServiceTrait {
	return &knativeServiceTrait{
		BaseTrait: newBaseTrait("knative-service"),
	}
}

func (t *knativeServiceTrait) Configure(e *Environment) (bool, error) {
	if t.Enabled != nil && !*t.Enabled {
		return false, nil
	}

	if !e.InPhase(v1alpha1.IntegrationKitPhaseReady, v1alpha1.IntegrationPhaseDeploying) {
		return false, nil
	}

	strategy, err := e.DetermineControllerStrategy(t.ctx, t.client)
	if err != nil {
		return false, err
	}
	if strategy != ControllerStrategyKnativeService {
		return false, nil
	}

	deployment := e.Resources.GetDeployment(func(d *appsv1.Deployment) bool {
		if name, ok := d.ObjectMeta.Labels["camel.apache.org/integration"]; ok {
			return name == e.Integration.Name
		}
		return false
	})
	if deployment != nil {
		// A controller is already present for the integration
		return false, nil
	}

	if t.Auto == nil || *t.Auto {
		// Check the right value for minScale, as not all services are allowed to scale down to 0
		if t.MinScale == nil {
			sources, err := kubernetes.ResolveIntegrationSources(t.ctx, t.client, e.Integration, e.Resources)
			if err != nil {
				return false, err
			}

			meta := metadata.ExtractAll(e.CamelCatalog, sources)
			if !meta.RequiresHTTPService || !meta.PassiveEndpoints {
				single := 1
				t.MinScale = &single
			}
		}
	}

	dt := e.Catalog.GetTrait("deployer")
	if dt != nil {
		t.deployer = *dt.(*deployerTrait)
	}

	return true, nil
}

func (t *knativeServiceTrait) Apply(e *Environment) error {
	svc := t.getServiceFor(e)
	maps := e.ComputeConfigMaps()

	e.Resources.Add(svc)
	e.Resources.AddAll(maps)

	return nil
}

func (t *knativeServiceTrait) getServiceFor(e *Environment) *serving.Service {
	labels := map[string]string{
		"camel.apache.org/integration": e.Integration.Name,
	}

	annotations := make(map[string]string)

	// Copy annotations from the integration resource
	if e.Integration.Annotations != nil {
		for k, v := range FilterTransferableAnnotations(e.Integration.Annotations) {
			annotations[k] = v
		}
	}

	// Resolve registry host names when used
	annotations["alpha.image.policy.openshift.io/resolve-names"] = "*"

	//
	// Set Knative Scaling behavior
	//
	if t.Class != "" {
		annotations[knativeServingClassAnnotation] = t.Class
	}
	if t.Metric != "" {
		annotations[knativeServingMetricAnnotation] = t.Metric
	}
	if t.Target != nil {
		annotations[knativeServingTargetAnnotation] = strconv.Itoa(*t.Target)
	}
	if t.MinScale != nil && *t.MinScale > 0 {
		annotations[knativeServingMinScaleAnnotation] = strconv.Itoa(*t.MinScale)
	}
	if t.MaxScale != nil && *t.MaxScale > 0 {
		annotations[knativeServingMaxScaleAnnotation] = strconv.Itoa(*t.MaxScale)
	}

	svc := serving.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: serving.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        e.Integration.Name,
			Namespace:   e.Integration.Namespace,
			Labels:      labels,
			Annotations: e.Integration.Annotations,
		},
		Spec: serving.ServiceSpec{
			RunLatest: &serving.RunLatestType{
				Configuration: serving.ConfigurationSpec{
					RevisionTemplate: serving.RevisionTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels:      labels,
							Annotations: annotations,
						},
						Spec: serving.RevisionSpec{
							ServiceAccountName: e.Integration.Spec.ServiceAccountName,
							Container: corev1.Container{
								Image: e.Integration.Status.Image,
								Env:   make([]corev1.EnvVar, 0),
							},
						},
					},
				},
			},
		},
	}

	paths := e.ComputeSourcesURI()
	environment := &svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.Env

	// combine Environment of integration with kit, integration
	for key, value := range e.CollectConfigurationPairs("env") {
		envvar.SetVal(environment, key, value)
	}

	// set env vars needed by the runtime
	envvar.SetVal(environment, "JAVA_MAIN_CLASS", "org.apache.camel.k.jvm.Application")

	// add a dummy env var to trigger deployment if everything but the code
	// has been changed
	envvar.SetVal(environment, "CAMEL_K_DIGEST", e.Integration.Status.Digest)

	envvar.SetVal(environment, "CAMEL_K_ROUTES", strings.Join(paths, ","))
	envvar.SetVal(environment, "CAMEL_K_CONF", "/etc/camel/conf/application.properties")
	envvar.SetVal(environment, "CAMEL_K_CONF_D", "/etc/camel/conf.d")

	// add env vars from traits
	for _, envVar := range t.getAllowedEnvVars(e) {
		envvar.SetVar(&svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.Env, envVar)
	}

	e.ConfigureVolumesAndMounts(
		&svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Volumes,
		&svc.Spec.RunLatest.Configuration.RevisionTemplate.Spec.Container.VolumeMounts,
	)

	return &svc
}

func (t *knativeServiceTrait) getAllowedEnvVars(e *Environment) []corev1.EnvVar {
	res := make([]corev1.EnvVar, 0, len(e.EnvVars))
	for _, env := range e.EnvVars {
		if env.ValueFrom == nil {
			// Standard env vars are supported
			res = append(res, env)
		} else if env.ValueFrom.FieldRef != nil && env.ValueFrom.FieldRef.FieldPath == "metadata.namespace" {
			// Namespace is known to the operator
			res = append(res, corev1.EnvVar{
				Name:  env.Name,
				Value: e.Integration.Namespace,
			})
		} else if env.ValueFrom.FieldRef != nil {
			t.L.Infof("Environment variable %s uses fieldRef and cannot be set on a Knative service", env.Name)
		} else if env.ValueFrom.ResourceFieldRef != nil {
			t.L.Infof("Environment variable %s uses resourceFieldRef and cannot be set on a Knative service", env.Name)
		} else {
			// Other downward APIs should be supported
			res = append(res, env)
		}
	}
	return res
}
