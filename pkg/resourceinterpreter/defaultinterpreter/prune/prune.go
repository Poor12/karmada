package prune

import (
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/karmada-io/karmada/pkg/util"
	"github.com/karmada-io/karmada/pkg/util/helper"
)

// RemoveIrrelevantField used to remove fields that generated by kube-apiserver and no need(or can't) propagate to
// member clusters.
func RemoveIrrelevantField(workload *unstructured.Unstructured, extraHooks ...func(*unstructured.Unstructured)) error {
	// populated by the kubernetes.
	unstructured.RemoveNestedField(workload.Object, "metadata", "creationTimestamp")

	// populated by the kubernetes.
	// The kubernetes will set this fields in case of graceful deletion. This field is read-only and can't propagate to
	// member clusters.
	unstructured.RemoveNestedField(workload.Object, "metadata", "deletionTimestamp")

	// populated by the kubernetes.
	// The kubernetes will set this fields in case of graceful deletion. This field is read-only and can't propagate to
	// member clusters.
	unstructured.RemoveNestedField(workload.Object, "metadata", "deletionGracePeriodSeconds")

	// populated by the kubernetes.
	unstructured.RemoveNestedField(workload.Object, "metadata", "generation")

	// This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field.
	// Remove this field to keep 'Work' clean and tidy.
	unstructured.RemoveNestedField(workload.Object, "metadata", "managedFields")

	// populated by the kubernetes.
	unstructured.RemoveNestedField(workload.Object, "metadata", "resourceVersion")

	// populated by the kubernetes and has been deprecated by kubernetes.
	unstructured.RemoveNestedField(workload.Object, "metadata", "selfLink")

	// populated by the kubernetes.
	unstructured.RemoveNestedField(workload.Object, "metadata", "uid")

	unstructured.RemoveNestedField(workload.Object, "metadata", "ownerReferences")

	unstructured.RemoveNestedField(workload.Object, "metadata", "finalizers")

	unstructured.RemoveNestedField(workload.Object, "status")

	if workload.GroupVersionKind() == corev1.SchemeGroupVersion.WithKind(util.ServiceKind) {
		// In the case spec.clusterIP is set to `None`, means user want a headless service,  then it shouldn't be removed.
		clusterIP, exist, _ := unstructured.NestedString(workload.Object, "spec", "clusterIP")
		if exist && clusterIP != corev1.ClusterIPNone {
			unstructured.RemoveNestedField(workload.Object, "spec", "clusterIP")
			unstructured.RemoveNestedField(workload.Object, "spec", "clusterIPs")
		}
	}

	if workload.GroupVersionKind() == batchv1.SchemeGroupVersion.WithKind(util.JobKind) {
		job := &batchv1.Job{}
		err := helper.ConvertToTypedObject(workload, job)
		if err != nil {
			return err
		}
		if job.Spec.ManualSelector == nil || !*job.Spec.ManualSelector {
			if err = removeGenerateSelectorOfJob(workload); err != nil {
				return err
			}
		}
	}

	if workload.GroupVersionKind() == corev1.SchemeGroupVersion.WithKind(util.ServiceAccountKind) {
		secrets, exist, _ := unstructured.NestedSlice(workload.Object, "secrets")
		// If 'secrets' exists in ServiceAccount, remove the automatic generation secrets(e.g. default-token-xxx)
		if exist && len(secrets) > 0 {
			tokenPrefix := fmt.Sprintf("%s-token-", workload.GetName())
			for idx := 0; idx < len(secrets); idx++ {
				if strings.HasPrefix(secrets[idx].(map[string]interface{})["name"].(string), tokenPrefix) {
					secrets = append(secrets[:idx], secrets[idx+1:]...)
				}
			}
			_ = unstructured.SetNestedSlice(workload.Object, secrets, "secrets")
		}
	}

	for i := range extraHooks {
		extraHooks[i](workload)
	}
	return nil
}

func removeGenerateSelectorOfJob(workload *unstructured.Unstructured) error {
	matchLabels, exist, err := unstructured.NestedStringMap(workload.Object, "spec", "selector", "matchLabels")
	if err != nil {
		return err
	}
	if exist {
		if util.GetLabelValue(matchLabels, "controller-uid") != "" {
			delete(matchLabels, "controller-uid")
		}
		err = unstructured.SetNestedStringMap(workload.Object, matchLabels, "spec", "selector", "matchLabels")
		if err != nil {
			return err
		}
	}

	templateLabels, exist, err := unstructured.NestedStringMap(workload.Object, "spec", "template", "metadata", "labels")
	if err != nil {
		return err
	}
	if exist {
		if util.GetLabelValue(templateLabels, "controller-uid") != "" {
			delete(templateLabels, "controller-uid")
		}

		if util.GetLabelValue(templateLabels, "job-name") != "" {
			delete(templateLabels, "job-name")
		}

		err = unstructured.SetNestedStringMap(workload.Object, templateLabels, "spec", "template", "metadata", "labels")
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveJobTTLSeconds removes the '.spec.ttlSecondsAfterFinished' from a Job.
// The reason for removing it is that the Job propagated by Karmada probably be automatically deleted
// from member clusters(by 'ttl-after-finished' controller in member clusters). That will cause a conflict if
// Karmada tries to re-create it. See https://github.com/karmada-io/karmada/issues/2197 for more details.
//
// It is recommended to enable the `ttl-after-finished` controller in the Karmada control plane.
// See https://karmada.io/docs/administrator/configuration/configure-controllers#ttl-after-finished for more details.
func RemoveJobTTLSeconds(workload *unstructured.Unstructured) {
	if workload.GroupVersionKind() == batchv1.SchemeGroupVersion.WithKind(util.JobKind) {
		unstructured.RemoveNestedField(workload.Object, "spec", "ttlSecondsAfterFinished")
	}
}
