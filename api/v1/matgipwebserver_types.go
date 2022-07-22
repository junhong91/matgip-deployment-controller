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

// MatgipWebServerSpec defines the desired state of MatgipWebServer
type MatgipWebServerSpec struct {
	// The port number of matgip web server to expose in k8s environment
	// +optional
	PortNumber *int32 `json:"portNumber,omitempty"`

	//+kubebuilder:validation:MinLength=0

	// The database service name to connect in k8s enviroment
	DatabaseName string `json:"databaseName"`

	// The ingress name to expose to enable remote access
	// +optional
	IngressHostName string `json:"ingressHostName,omitempty"`

	//+kubebuilder:validation:MinLength=0

	// The newsClientId is the client id for requesting to new API server
	NewsClientId string `json:"newsClientId"`

	//+kubebuilder:validation:MinLength=0

	// The newsClientSecret is the client secret for requesting to news API server
	NewsClientSecret string `json:"newsClientSecret"`

	//+kubebuilder:validation:MinLength=0

	// The authTokenSecret is the API access token to communicate with backend server
	AuthTokenSecret string `json:"authTokenSecret"`
}

// MatgipWebServerStatus defines the observed state of MatgipWebServer
type MatgipWebServerStatus struct {
	// Reconciled defines whether the host has been successfully reconciled
	// at least onece. If further changes are made they will be ignored by the
	// reconciler.
	Reconciled bool `json:"reconciled"`

	// Defines whether the resource has been provisioned on the target system.
	InSync bool `json:"inSync"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MatgipWebServer is the Schema for the matgipwebservers API
type MatgipWebServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MatgipWebServerSpec   `json:"spec,omitempty"`
	Status MatgipWebServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MatgipWebServerList contains a list of MatgipWebServer
type MatgipWebServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MatgipWebServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MatgipWebServer{}, &MatgipWebServerList{})
}
