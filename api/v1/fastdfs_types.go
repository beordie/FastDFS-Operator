/*
Copyright 2023.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FastDFSSpec defines the desired state of FastDFS
type FastDFSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Version specifies expect FastDFS image tag, except 3.6.3
	//
	// +required
	Version string `json:"version"`

	// server auth is used to authenticate servers
	//
	// +optional
	ServerAuth ServerAuth `json:"serverAuth,omitempty"`

	// Replicas is the expected size of the FastDFS cluster.
	// The operator will eventually make the size of the running cluster
	// equal to the expected size.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// ParticipantReplicas is the expected size of the FastDFS participants
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	ParticipantReplicas *int32 `json:"participantReplicas,omitempty"`

	// ClientPort is the port of FastDFS client service
	// Cannot be updated.
	//
	// +optional
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=0
	ClientPort int32 `json:"clientPort,omitempty"`

	// MetricsPort is the http port of FastDFS prometheus metrics, default 8080
	// Cannot be updated.
	//
	// +optional
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:validation:Minimum=0
	MetricsPort int32 `json:"metricsPort,omitempty"`

	// Paused specified whether cluster service continue to serve
	//
	// +optional
	Paused bool `json:"paused,omitempty"`

	// Labels specifies the labels that will be tagged
	// on all resources created by FastDFSCluster
	//
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Pod specifies deploy policy of pod
	//
	// +optional
	Pod *PodOption `json:"pod,omitempty"`

	// NodeSelector specified the pod affinity of nodes
	//
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Affinity describes node selector requirements that must be met
	//
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Tolerations enable pod to schedule on node that have taints on,
	// normally combine it with Affinity when deploy a schedule policy
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// FastDFSStatus defines the observed state of FastDFS
type FastDFSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FastDFS is the Schema for the fastdfs API
type FastDFS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FastDFSSpec   `json:"spec,omitempty"`
	Status FastDFSStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FastDFSList contains a list of FastDFS
type FastDFSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FastDFS `json:"items"`
}

type PodOption struct {
	// Annotations specifies the annotations to attach to pods the operator creates
	// for the FastDFS cluster.
	//
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// ImagePullRepository specifies the repository of FastDFS images
	//
	// +optional
	ImagePullRepository string `json:"imagePullRepository,omitempty"`

	// ImagePullPolicy describes a policy for if/when to pull a FastDFS image
	// default IfNotPresent
	//
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace
	// to use for pulling any of the images used by this PodSpec.
	//
	// +optional
	ImagePullSecrets string `json:"imagePullSecrets,omitempty"`

	// Resources specifies the resource needed per pod
	//
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type ServerAuth struct {
	// disabled auth default false
	//
	// +optional
	Token string `json:"token,omitempty"`

	// token TTL (time to live), seconds
	// default value is 600
	//
	// +optional
	TTL int32 `json:"ttl,omitempty"`

	// secret key to generate anti-steal token
	// this parameter must be set when http.anti_steal.check_token set to true
	// the length of the secret key should not exceed 128 bytes
	//
	// +optional
	Secret string `json:"secret,omitempty"`

	// return the content of the file when check token fail
	// default value is empty (no file sepecified)
	//
	// +optional
	FailFallBack string `json:"failFallBack,omitempty"`
}

func init() {
	SchemeBuilder.Register(&FastDFS{}, &FastDFSList{})
}
