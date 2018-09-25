/*
Copyright 2018 The Kubernetes Authors.

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

package transformerconfig

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ReferencePathConfig contains the configuration of a field that references
// the name of another resource whose GroupVersionKind is specified in referencedGVK.
// e.g. pod.spec.template.volumes.configMap.name references the name of a configmap
// Its corresponding referencePathConfig will look like:
//
//	ReferencePathConfig{
//	referencedGVK: schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"},
//	pathConfigs: []PathConfig{
//		{
//			GroupVersionKind: &schema.GroupVersionKind{Version: "v1", Kind: "Pod"},
//			Path:             []string{"spec", "volumes", "configMap", "name"},
//		},
//	}
type ReferencePathConfig struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
	// PathConfig is the gvk that is referencing the referencedGVK object's name.
	PathConfigs []PathConfig `json:"pathConfigs,omitempty" yaml:"pathConfigs,omitempty"`
}

// Gvk returns GroupVersionKind of the reference pathConfig
func (p ReferencePathConfig) Gvk() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   p.Group,
		Version: p.Version,
		Kind:    p.Kind,
	}
}

func merge(configs []ReferencePathConfig, config ReferencePathConfig) []ReferencePathConfig {
	var result []ReferencePathConfig
	found := false
	for _, c := range configs {
		if c.Gvk() == config.Gvk() {
			c.PathConfigs = append(c.PathConfigs, config.PathConfigs...)
			found = true
		}
		result = append(result, c)
	}

	if !found {
		result = append(result, config)
	}
	return result
}

func mergeNameReferencePathConfigs(a, b []ReferencePathConfig) []ReferencePathConfig {
	for _, r := range b {
		a = merge(a, r)
	}
	return a
}