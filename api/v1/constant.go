package v1

const (
	StatefulsetName      = "%s-statefulset"
	ConfigMapName        = "%s-configmap"
	HeadlessServiceName  = "%s-headless-service"
	StorageValueUnit     = "%d%s"
	ConfigVolumeName     = "zookeeper-config"
	StorageContainerName = "storage"
	TrackerContainerName = "tracker"
	PvcName              = "fastdfs-storage-data"
	DataDir              = "/data"
)

const (
	ScheduleTypeAnnotation = "schedule.type"
	TopologyKey            = "failure-domain.beta.kubernetes.io/zone"
)

const (
	ScheduleTypeAnnotationValueIgnore = "ignore"

	DefaultStoragePort = 22122
	DefaultTrackerPort = 23000
	DefaultDHTPort     = 11411
)
