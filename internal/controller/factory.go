package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"
	"fmt"

	"fastdfs_operator/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *FastDFSReconciler) makeStatefulSet(cluster *v1.FastDFS) *appsv1.StatefulSet {
	nn := cluster.GetStatefulSetNamespacedName()
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
	}
}

func (r *FastDFSReconciler) makePVCStorageSize(cluster *v1.FastDFS) resource.Quantity {
	return resource.MustParse(fmt.Sprintf(v1.StorageValueUnit, cluster.Spec.Storage.DiskSize, cluster.Spec.Storage.Unit))
}

func (r *FastDFSReconciler) mutateStatefulSet(cluster *v1.FastDFS, sts *appsv1.StatefulSet) error {
	if sts.ObjectMeta.CreationTimestamp.IsZero() {
		// sts resource was not created yet, or happened any error
		sts.ObjectMeta.Labels = cluster.ResourceLabels()
		sts.Spec.Selector = &metav1.LabelSelector{MatchLabels: cluster.ResourceMatchingLabels()}
		sts.Spec.ServiceName = cluster.GetHeadlessServiceName()
		sts.Spec.PodManagementPolicy = appsv1.ParallelPodManagement

		if sts.Spec.VolumeClaimTemplates == nil || len(sts.Spec.VolumeClaimTemplates) == 0 {
			sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{{}}
		}
		pvc := &sts.Spec.VolumeClaimTemplates[0]
		pvc.Name = v1.PvcName
		pvc.Namespace = cluster.GetNamespace()
		pvc.Labels = cluster.ResourceLabels()
		pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		pvc.Spec.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: r.makePVCStorageSize(cluster),
			},
		}
		pvc.Spec.StorageClassName = cluster.Spec.Storage.StorageClass
	}
	sts.Spec.Replicas = cluster.NextReplicas()
	sts.Spec.UpdateStrategy = appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType}
	// Template
	sts.Spec.Template.Labels = cluster.ResourceLabels()
	annotations, err := r.makePodAnnotations(cluster)
	if err != nil {
		return err
	}
	sts.Spec.Template.Annotations = annotations
	sts.Spec.Template.Spec.ImagePullSecrets = utils.GetReferencesFromStringSlice(cluster.Spec.Pod.ImagePullSecrets)
	if sts.Spec.Template.Spec.Affinity == nil {
		sts.Spec.Template.Spec.Affinity = r.makePodAffinity(cluster)
	}

	sts.Spec.Template.Spec.Tolerations = cluster.Spec.Tolerations
	sts.Spec.Template.Spec.NodeSelector = cluster.Spec.NodeSelector

	// Template.Spec.Volumes
	if sts.Spec.Template.Spec.Volumes == nil || len(sts.Spec.Template.Spec.Volumes) == 0 {
		sts.Spec.Template.Spec.Volumes = []corev1.Volume{{}}
	}
	volume := &sts.Spec.Template.Spec.Volumes[0]
	volume.Name = v1.PvcName

	if volume.VolumeSource.ConfigMap == nil {
		volume.VolumeSource.ConfigMap = &corev1.ConfigMapVolumeSource{}
	}
	volume.VolumeSource.ConfigMap.LocalObjectReference = corev1.LocalObjectReference{Name: cluster.GetConfigMapName()}

	sts.Spec.Template.Spec.Containers = r.makePodImage(cluster)
	return controllerutil.SetControllerReference(cluster, sts, r.Scheme)
}

func (r *FastDFSReconciler) makePodAnnotations(cluster *v1.FastDFS) (map[string]string, error) {
	annotations := map[string]string{}
	if cluster.Spec.Pod.Annotations != nil {
		for k, v := range cluster.Spec.Pod.Annotations {
			annotations[k] = v
		}
	}

	cm := &corev1.ConfigMap{}
	if err := r.Client.Get(context.TODO(), cluster.GetConfigMapNamespacedName(), cm); err != nil {
		return nil, err
	}

	return annotations, nil
}

func (r *FastDFSReconciler) makePodAffinity(cluster *v1.FastDFS) *corev1.Affinity {
	if cluster.IgnoreSchedulePolicy() {
		return nil
	}

	affinity := &corev1.Affinity{}
	if cluster.Spec.Affinity != nil {
		affinity = cluster.Spec.Affinity.DeepCopy()
	}

	makePodAntiAffinity(cluster, affinity)
	makePodNodeAffinity(cluster, affinity)
	return affinity
}

func makePodAntiAffinity(cluster *v1.FastDFS, affinity *corev1.Affinity) {
	if affinity.PodAntiAffinity != nil {
		return
	}

	affinity.PodAntiAffinity = &corev1.PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
			{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: cluster.ResourceLabels(),
				},
				TopologyKey: corev1.LabelHostname,
			},
		},
	}
}

func makePodNodeAffinity(cluster *v1.FastDFS, affinity *corev1.Affinity) {
	if len(cluster.Spec.AvailableZones) == 0 {
		return
	}

	if affinity.NodeAffinity == nil {
		affinity.NodeAffinity = &corev1.NodeAffinity{}
	}
	if affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{}
	}

	terms := []corev1.NodeSelectorTerm{}
	if affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms != nil {
		terms = affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms
	} else {
		terms = append(terms, corev1.NodeSelectorTerm{})
	}
	for i := range terms {
		terms[i].MatchExpressions = append(terms[i].MatchExpressions, corev1.NodeSelectorRequirement{
			Key:      v1.TopologyKey,
			Operator: corev1.NodeSelectorOpIn,
			Values:   cluster.Spec.AvailableZones,
		})
	}
	affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms = terms
}

func (r *FastDFSReconciler) makePodImage(cluster *v1.FastDFS) []corev1.Container {
	imagePullRepository := cluster.Spec.Pod.ImagePullRepository

	pod := cluster.Spec.Pod
	containers := []corev1.Container{}
	container := corev1.Container{}
	container.Name = pod.Image.Name
	container.ImagePullPolicy = pod.ImagePullPolicy
	container.Image = imagePullRepository + "/" + pod.Image.Name + ":" + pod.Image.Version
	container.Resources = pod.Resources
	container.Ports = makePodPorts(pod.Image.Name)
	container.LivenessProbe = &corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(8888),
			},
		},
		InitialDelaySeconds: 20,
		PeriodSeconds:       5,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      30,
	}
	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      v1.PvcName,
			MountPath: v1.DataDir,
		},
	}
	containers = append(containers, container)
	return containers
}

func makePodPorts(name string) []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Name:          name,
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 8888,
		},
	}
}

func makeConfigmap(cluster *v1.FastDFS) *corev1.ConfigMap {
	nn := cluster.GetConfigMapNamespacedName()
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
	}
}

func (r *FastDFSReconciler) mutateConfigmap(cluster *v1.FastDFS, cm *corev1.ConfigMap) error {
	cm.Labels = cluster.ResourceLabels()
	var cd ConfigMap = make(map[string]string)

	cm.Data = cd
	return controllerutil.SetControllerReference(cluster, cm, r.Scheme)
}
