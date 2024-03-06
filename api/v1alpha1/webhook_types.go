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

// WebhookSpec defines the desired state of Webhook
type WebhookSpec struct {
	// The path of the webhook. Optional, if you don't specify one
	// it will be generated for you. If you don't set an ID, then this will
	// change each time the operator is restarted.
	ID string `json:"id,omitempty"`

	// The secret to validate incoming webhooks with. Not implemented yet.
	Secret string `json:"secret,omitempty"`
}

// WebhookStatus defines the observed state of Webhook
type WebhookStatus struct {
	// Whether the Webhhok configuration is valid and has been processed.
	Valid bool `json:"enabled"`
	// The ID id of the webhook
	ID string `json:"id"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Webhook is the Schema for the webhooks API
type Webhook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebhookSpec   `json:"spec,omitempty"`
	Status WebhookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WebhookList contains a list of Webhook
type WebhookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webhook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Webhook{}, &WebhookList{})
}
