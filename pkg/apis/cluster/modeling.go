package cluster

import (
	corev1 "k8s.io/api/core/v1"
)

func InitSummary(rsList []corev1.ResourceName) *ResourceSummary {
	return nil
}

func (rs *ResourceSummary) AddToResourceSummary() {

}

func (rs *ResourceSummary) DeleteFromResourceSummary() {

}

func (rs *ResourceSummary) UpdateInResourceSummary() {

}
