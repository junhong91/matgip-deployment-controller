/*
Copyright 2022.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MatgipRedisServerSpec defines the desired state of MatgipRedisServer
type MatgipRedisServerSpec struct {
	// The port number of redis server to expose
	// +optional
	PortNumber *int32 `json:"portNumber,omitempty"`
}

// MatgipRedisServerStatus defines the observed state of MatgipRedisServer
type MatgipRedisServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MatgipRedisServer is the Schema for the matgipredisservers API
type MatgipRedisServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MatgipRedisServerSpec   `json:"spec,omitempty"`
	Status MatgipRedisServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MatgipRedisServerList contains a list of MatgipRedisServer
type MatgipRedisServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MatgipRedisServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MatgipRedisServer{}, &MatgipRedisServerList{})
}
