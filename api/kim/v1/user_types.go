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

type Claim struct {
	Email               *string `json:"email,omitempty"`
	EmailVerified       *bool   `json:"emailVerified,omitempty"`
	GivenName           *string `json:"givenName,omitempty"`
	FamilyName          *string `json:"familyName,omitempty"`
	MiddleName          *string `json:"middleName,omitempty"`
	NickName            *string `json:"nickName,omitempty"`
	PreferredUsername   *string `json:"preferredUsername,omitempty"`
	Profile             *string `json:"profile,omitempty"`
	Picture             *string `json:"picture,omitempty"`
	Website             *string `json:"website,omitempty"`
	Gender              *string `json:"gender,omitempty"`
	Birthdate           *string `json:"birthdate,omitempty"`
	Zoneinfo            *string `json:"zoneinfo,omitempty"`
	Locale              *string `json:"locale,omitempty"`
	PhoneNumber         *string `json:"phoneNumber,omitempty"`
	PhoneNumberVerified *bool   `json:"phoneNumberVerified,omitempty"`
	Address             *string `json:"address,omitempty"`
}

// UserSpec defines the desired state of User
type UserSpec struct {
	Desc       string `json:"desc"`
	SecretName string `json:"secretName"`
	Claim      `json:",inline"`
}

// UserStatus defines the observed state of User.
type UserStatus struct{}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// User is the Schema for the users API
type User struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of User
	// +required
	Spec UserSpec `json:"spec"`

	// status defines the observed state of User
	// +optional
	Status UserStatus `json:"status,omitempty,omitzero"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// UserList contains a list of User
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
