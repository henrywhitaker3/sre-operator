/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TargetSpec struct {
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum:deployment;statefulset;daemonset
	Kind string `json:"kind"`

	//+kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	//+kubebuilder:validation:Required
	Name string `json:"name"`
}

// RolloutSpec defines the desired state of Rollout
type RolloutSpec struct {
	//+kubebuilder:validation:Required
	Triggers []Trigger `json:"triggers"`

	//+kubebuilder:validation:Required
	Target TargetSpec `json:"target"`

	//+kubebuilder:validation:Required
	//+kubebuilder:validation:Enum:pause;restart;resume;undo
	Action string `json:"action"`

	//+kubebuilder:validation:Optional
	// Throttle the action to 1 time per specified duration e.g. 5m, 1h
	Throttle string `json:"throttle,omitempty"`
}

// RolloutStatus defines the observed state of Rollout
type RolloutStatus struct {
	Registered bool   `json:"registered,omitempty"`
	Error      string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Rollout is the Schema for the rollouts API
type Rollout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RolloutSpec   `json:"spec,omitempty"`
	Status RolloutStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RolloutList contains a list of Rollout
type RolloutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rollout `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Rollout{}, &RolloutList{})
}
