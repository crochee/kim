/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PolicySpec defines the desired state of Policy
type PolicySpec struct {
	// Description of the policy
	Desc string `json:"desc"`
	// JSON-formatted policy statement
	Statement string `json:"statement"`
	// List of rules defining the policy
	Rules []Rule `json:"rules"`
}

// Rule defines a single policy rule
type Rule struct {
	// Resource type the rule applies to
	Resource string `json:"resource"`
	// List of actions allowed by the rule
	Actions []string `json:"actions"`
	// Effect of the rule (Allow/Deny)
	Effect string `json:"effect"`
}

// PolicyStatus defines the observed state of Policy
type PolicyStatus struct {
	// Conditions represent the latest available observations of the policy's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Policy is the Schema for the policies API
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PolicySpec   `json:"spec,omitempty"`
	Status PolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PolicyList contains a list of Policy
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Policy{}, &PolicyList{})
}
