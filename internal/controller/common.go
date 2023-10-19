package controller

const (
	// AnnotationPauseReconcile when set it to 'true' means cluster should not be reconciled by operator until
	// the annotation set to false or removed
	AnnotationPauseReconcile = "paas.netease.com/pause-reconcile"
)
