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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	// Paused specified whether cluster service continue to serve
	//
	// +optional
	Paused bool `json:"paused,omitempty"`

	// Storage specifies storage options
	//
	// +required
	Storage *StorageOption `json:"storage,omitempty"`

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

	// AvailableZones specified the available zone pod can be scheduled
	//
	// +optional
	AvailableZones []string `json:"availableZones,omitempty"`
}

// FastDFSStatus defines the observed state of FastDFS
type FastDFSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Information when was the last time the cr was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// Replicas is the number of replicas created in the cluster
	Replicas int32 `json:"replicas"`

	// CurrentStatefulSetReplicas is the number of statefulset replicas
	CurrentStatefulSetReplicas int32 `json:"currentStatefulSetReplicas"`

	// ObserverReplicas is the number of observer replicas created in the cluster
	ObserverReplicas int32 `json:"observerReplicas"`

	// ParticipantReplicas is the number of participant replicas created in the cluster
	ParticipantReplicas int32 `json:"participantReplicas"`

	// ReadyReplicas is the number of ready replicas in the cluster that are ready
	ReadyReplicas int32 `json:"readyReplicas"`
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

	// tracker & storage image
	//
	// +optional
	Images []Image `json:"images,omitempty"`
}

type Image struct {
	// container image name
	//
	// +optional
	Name string `json:"name,omitempty"`

	//container image version
	//
	// +optional
	Version string `json:"version,omitempty"`
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

/**
 * GetStatefulSetName is the name of the cluster statefulset
 *
 * @return string
 */
func (cluster *FastDFS) GetStatefulSetName() string {
	return fmt.Sprintf(StatefulsetName, cluster.Name)
}

/**
 * GetStatefulSetNamespacedName is the namespaced and name of the statefulset
 *
 * @return types.NamespacedName
 */
func (cluster *FastDFS) GetStatefulSetNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: cluster.Namespace, Name: cluster.GetStatefulSetName()}
}

/**
 * ResourceMatchingLabels is the labels that will be tagged
 * on all resources created by FastDFSCluster
 *
 * @return map[string]string
 */
func (cluster *FastDFS) ResourceMatchingLabels() map[string]string {
	return map[string]string{
		"clusterId": string(cluster.UID),
	}
}

/**
 * ResourceLabels is the labels that will be tagged
 * on all resources created by FastDFSCluster
 *
 * @return map[string]string
 */
func (cluster *FastDFS) ResourceLabels() map[string]string {
	labels := map[string]string{}
	for k, v := range cluster.Spec.Labels {
		labels[k] = v
	}

	// separate resources from different cluster
	for k, v := range cluster.ResourceMatchingLabels() {
		labels[k] = v
	}
	return labels
}

/**
 * HeadlessServiceName is the name of the headless service for the FastDFS cluster
 * - headless service name: <cluster-name>-headless
 * statefulset need a unique service tag to expose net
 * @return string
 */
func (cluster *FastDFS) GetHeadlessServiceName() string {
	return fmt.Sprintf(HeadlessServiceName, cluster.Name)
}

func (cluster *FastDFS) GetConfigMapName() string {
	return fmt.Sprintf(ConfigMapName, cluster.Name)
}

func (cluster *FastDFS) GetConfigMapNamespacedName() types.NamespacedName {
	return types.NamespacedName{Namespace: cluster.Namespace, Name: cluster.GetConfigMapName()}
}

func (cluster *FastDFS) IgnoreSchedulePolicy() bool {
	if cluster.Annotations != nil && cluster.Annotations[ScheduleTypeAnnotation] == ScheduleTypeAnnotationValueIgnore {
		return true
	}
	return false
}

type VolumeReclaimPolicy string

const (
	VolumeReclaimPolicyRetain VolumeReclaimPolicy = "Retain"
	VolumeReclaimPolicyDelete VolumeReclaimPolicy = "Delete"
)

type StorageOption struct {
	// DiskSize specifies the storage size of pod
	// unit Gi
	//
	// +required
	// +kubebuilder:validation:Minimum=1
	DiskSize int32 `json:"diskSize,omitempty"`

	// Unit specifies the unit of DiskSize, default Gi
	//
	// +optional
	// +kubebuilder:validation:Enum=Mi;Gi
	Unit string `json:"unit,omitempty"`

	// StorageClass specifies storageclass used by pvc
	//
	// +optional
	StorageClass *string `json:"storageClass,omitempty"`

	// VolumeReclaimPolicy is a zookeeper operator configuration. If it's set to Delete,
	// the corresponding PVCs will be deleted by the operator when zookeeper cluster is deleted.
	// The default value is Retain.
	//
	// +optional
	// +kubebuilder:validation:Enum="Delete";"Retain"
	VolumeReclaimPolicy VolumeReclaimPolicy `json:"reclaimPolicy,omitempty"`
}

func (cluster *FastDFS) NextReplicas() *int32 {
	nextReplicas := cluster.Status.CurrentStatefulSetReplicas

	if cluster.Status.CurrentStatefulSetReplicas < *cluster.Spec.Replicas {
		// when scale up, make sure previous replica be ready
		if cluster.Status.ReadyReplicas == cluster.Status.CurrentStatefulSetReplicas {
			nextReplicas = cluster.Status.ReadyReplicas + 1
		}
	} else if cluster.Status.CurrentStatefulSetReplicas > *cluster.Spec.Replicas {
		// when scale down, scale immediately
		nextReplicas = *cluster.Spec.Replicas
	}

	return &nextReplicas
}

func init() {
	SchemeBuilder.Register(&FastDFS{}, &FastDFSList{})
}
