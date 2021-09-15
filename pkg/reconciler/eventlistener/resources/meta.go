/*
Copyright 2021 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	"context"
	"regexp"

	"github.com/tektoncd/triggers/pkg/apis/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"

	"github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
)

// ObjectMeta generates the object meta that should be used by all
// resources generated by the EventListener reconciler
func ObjectMeta(el *v1beta1.EventListener, filteredElLabels, staticResourceLabels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace:       el.Namespace,
		Name:            el.Status.Configuration.GeneratedResourceName,
		OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(el)},
		Labels:          kmeta.UnionMaps(filteredElLabels, GenerateLabels(el.Name, staticResourceLabels)),
		Annotations:     el.Annotations,
	}
}

// GenerateLabels generates the labels to be used on all generated resources.
func GenerateLabels(eventListenerName string, staticResourceLabels map[string]string) map[string]string {
	resourceLabels := kmeta.CopyMap(staticResourceLabels)
	resourceLabels["eventlistener"] = eventListenerName
	return resourceLabels
}

// FilterLabels filters label based on regex pattern defined in
// feature-flag `labels-exclusion-pattern`
func FilterLabels(ctx context.Context, labels map[string]string) map[string]string {
	cfg := config.FromContextOrDefaults(ctx)

	if len(labels) == 0 || cfg.FeatureFlags.LabelsExclusionPattern == "" {
		return labels
	}

	filteredLabels := make(map[string]string)
	r := regexp.MustCompile(cfg.FeatureFlags.LabelsExclusionPattern)

	for key, value := range labels {
		if !r.MatchString(key) {
			filteredLabels[key] = value
		}
	}

	return filteredLabels
}
