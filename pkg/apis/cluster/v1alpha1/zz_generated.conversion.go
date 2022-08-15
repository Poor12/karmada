//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	url "net/url"
	unsafe "unsafe"

	cluster "github.com/karmada-io/karmada/pkg/apis/cluster"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*APIEnablement)(nil), (*cluster.APIEnablement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_APIEnablement_To_cluster_APIEnablement(a.(*APIEnablement), b.(*cluster.APIEnablement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.APIEnablement)(nil), (*APIEnablement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_APIEnablement_To_v1alpha1_APIEnablement(a.(*cluster.APIEnablement), b.(*APIEnablement), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*APIResource)(nil), (*cluster.APIResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_APIResource_To_cluster_APIResource(a.(*APIResource), b.(*cluster.APIResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.APIResource)(nil), (*APIResource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_APIResource_To_v1alpha1_APIResource(a.(*cluster.APIResource), b.(*APIResource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Cluster)(nil), (*cluster.Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Cluster_To_cluster_Cluster(a.(*Cluster), b.(*cluster.Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.Cluster)(nil), (*Cluster)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_Cluster_To_v1alpha1_Cluster(a.(*cluster.Cluster), b.(*Cluster), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterList)(nil), (*cluster.ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterList_To_cluster_ClusterList(a.(*ClusterList), b.(*cluster.ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterList)(nil), (*ClusterList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterList_To_v1alpha1_ClusterList(a.(*cluster.ClusterList), b.(*ClusterList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterProxyOptions)(nil), (*cluster.ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(a.(*ClusterProxyOptions), b.(*cluster.ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterProxyOptions)(nil), (*ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterProxyOptions_To_v1alpha1_ClusterProxyOptions(a.(*cluster.ClusterProxyOptions), b.(*ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterSpec)(nil), (*cluster.ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec(a.(*ClusterSpec), b.(*cluster.ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterSpec)(nil), (*ClusterSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterSpec_To_v1alpha1_ClusterSpec(a.(*cluster.ClusterSpec), b.(*ClusterSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ClusterStatus)(nil), (*cluster.ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus(a.(*ClusterStatus), b.(*cluster.ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ClusterStatus)(nil), (*ClusterStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus(a.(*cluster.ClusterStatus), b.(*ClusterStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LocalSecretReference)(nil), (*cluster.LocalSecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LocalSecretReference_To_cluster_LocalSecretReference(a.(*LocalSecretReference), b.(*cluster.LocalSecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.LocalSecretReference)(nil), (*LocalSecretReference)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_LocalSecretReference_To_v1alpha1_LocalSecretReference(a.(*cluster.LocalSecretReference), b.(*LocalSecretReference), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NodeSummary)(nil), (*cluster.NodeSummary)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NodeSummary_To_cluster_NodeSummary(a.(*NodeSummary), b.(*cluster.NodeSummary), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.NodeSummary)(nil), (*NodeSummary)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_NodeSummary_To_v1alpha1_NodeSummary(a.(*cluster.NodeSummary), b.(*NodeSummary), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceSummary)(nil), (*cluster.ResourceSummary)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary(a.(*ResourceSummary), b.(*cluster.ResourceSummary), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*cluster.ResourceSummary)(nil), (*ResourceSummary)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_cluster_ResourceSummary_To_v1alpha1_ResourceSummary(a.(*cluster.ResourceSummary), b.(*ResourceSummary), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*url.Values)(nil), (*ClusterProxyOptions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_url_Values_To_v1alpha1_ClusterProxyOptions(a.(*url.Values), b.(*ClusterProxyOptions), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_APIEnablement_To_cluster_APIEnablement(in *APIEnablement, out *cluster.APIEnablement, s conversion.Scope) error {
	out.GroupVersion = in.GroupVersion
	out.Resources = *(*[]cluster.APIResource)(unsafe.Pointer(&in.Resources))
	return nil
}

// Convert_v1alpha1_APIEnablement_To_cluster_APIEnablement is an autogenerated conversion function.
func Convert_v1alpha1_APIEnablement_To_cluster_APIEnablement(in *APIEnablement, out *cluster.APIEnablement, s conversion.Scope) error {
	return autoConvert_v1alpha1_APIEnablement_To_cluster_APIEnablement(in, out, s)
}

func autoConvert_cluster_APIEnablement_To_v1alpha1_APIEnablement(in *cluster.APIEnablement, out *APIEnablement, s conversion.Scope) error {
	out.GroupVersion = in.GroupVersion
	out.Resources = *(*[]APIResource)(unsafe.Pointer(&in.Resources))
	return nil
}

// Convert_cluster_APIEnablement_To_v1alpha1_APIEnablement is an autogenerated conversion function.
func Convert_cluster_APIEnablement_To_v1alpha1_APIEnablement(in *cluster.APIEnablement, out *APIEnablement, s conversion.Scope) error {
	return autoConvert_cluster_APIEnablement_To_v1alpha1_APIEnablement(in, out, s)
}

func autoConvert_v1alpha1_APIResource_To_cluster_APIResource(in *APIResource, out *cluster.APIResource, s conversion.Scope) error {
	out.Name = in.Name
	out.Kind = in.Kind
	return nil
}

// Convert_v1alpha1_APIResource_To_cluster_APIResource is an autogenerated conversion function.
func Convert_v1alpha1_APIResource_To_cluster_APIResource(in *APIResource, out *cluster.APIResource, s conversion.Scope) error {
	return autoConvert_v1alpha1_APIResource_To_cluster_APIResource(in, out, s)
}

func autoConvert_cluster_APIResource_To_v1alpha1_APIResource(in *cluster.APIResource, out *APIResource, s conversion.Scope) error {
	out.Name = in.Name
	out.Kind = in.Kind
	return nil
}

// Convert_cluster_APIResource_To_v1alpha1_APIResource is an autogenerated conversion function.
func Convert_cluster_APIResource_To_v1alpha1_APIResource(in *cluster.APIResource, out *APIResource, s conversion.Scope) error {
	return autoConvert_cluster_APIResource_To_v1alpha1_APIResource(in, out, s)
}

func autoConvert_v1alpha1_Cluster_To_cluster_Cluster(in *Cluster, out *cluster.Cluster, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Cluster_To_cluster_Cluster is an autogenerated conversion function.
func Convert_v1alpha1_Cluster_To_cluster_Cluster(in *Cluster, out *cluster.Cluster, s conversion.Scope) error {
	return autoConvert_v1alpha1_Cluster_To_cluster_Cluster(in, out, s)
}

func autoConvert_cluster_Cluster_To_v1alpha1_Cluster(in *cluster.Cluster, out *Cluster, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_cluster_ClusterSpec_To_v1alpha1_ClusterSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_cluster_Cluster_To_v1alpha1_Cluster is an autogenerated conversion function.
func Convert_cluster_Cluster_To_v1alpha1_Cluster(in *cluster.Cluster, out *Cluster, s conversion.Scope) error {
	return autoConvert_cluster_Cluster_To_v1alpha1_Cluster(in, out, s)
}

func autoConvert_v1alpha1_ClusterList_To_cluster_ClusterList(in *ClusterList, out *cluster.ClusterList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]cluster.Cluster, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_Cluster_To_cluster_Cluster(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_ClusterList_To_cluster_ClusterList is an autogenerated conversion function.
func Convert_v1alpha1_ClusterList_To_cluster_ClusterList(in *ClusterList, out *cluster.ClusterList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterList_To_cluster_ClusterList(in, out, s)
}

func autoConvert_cluster_ClusterList_To_v1alpha1_ClusterList(in *cluster.ClusterList, out *ClusterList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cluster, len(*in))
		for i := range *in {
			if err := Convert_cluster_Cluster_To_v1alpha1_Cluster(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_cluster_ClusterList_To_v1alpha1_ClusterList is an autogenerated conversion function.
func Convert_cluster_ClusterList_To_v1alpha1_ClusterList(in *cluster.ClusterList, out *ClusterList, s conversion.Scope) error {
	return autoConvert_cluster_ClusterList_To_v1alpha1_ClusterList(in, out, s)
}

func autoConvert_v1alpha1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in *ClusterProxyOptions, out *cluster.ClusterProxyOptions, s conversion.Scope) error {
	out.Path = in.Path
	return nil
}

// Convert_v1alpha1_ClusterProxyOptions_To_cluster_ClusterProxyOptions is an autogenerated conversion function.
func Convert_v1alpha1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in *ClusterProxyOptions, out *cluster.ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterProxyOptions_To_cluster_ClusterProxyOptions(in, out, s)
}

func autoConvert_cluster_ClusterProxyOptions_To_v1alpha1_ClusterProxyOptions(in *cluster.ClusterProxyOptions, out *ClusterProxyOptions, s conversion.Scope) error {
	out.Path = in.Path
	return nil
}

// Convert_cluster_ClusterProxyOptions_To_v1alpha1_ClusterProxyOptions is an autogenerated conversion function.
func Convert_cluster_ClusterProxyOptions_To_v1alpha1_ClusterProxyOptions(in *cluster.ClusterProxyOptions, out *ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_cluster_ClusterProxyOptions_To_v1alpha1_ClusterProxyOptions(in, out, s)
}

func autoConvert_url_Values_To_v1alpha1_ClusterProxyOptions(in *url.Values, out *ClusterProxyOptions, s conversion.Scope) error {
	// WARNING: Field TypeMeta does not have json tag, skipping.

	if values, ok := map[string][]string(*in)["path"]; ok && len(values) > 0 {
		if err := runtime.Convert_Slice_string_To_string(&values, &out.Path, s); err != nil {
			return err
		}
	} else {
		out.Path = ""
	}
	return nil
}

// Convert_url_Values_To_v1alpha1_ClusterProxyOptions is an autogenerated conversion function.
func Convert_url_Values_To_v1alpha1_ClusterProxyOptions(in *url.Values, out *ClusterProxyOptions, s conversion.Scope) error {
	return autoConvert_url_Values_To_v1alpha1_ClusterProxyOptions(in, out, s)
}

func autoConvert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec(in *ClusterSpec, out *cluster.ClusterSpec, s conversion.Scope) error {
	out.ID = in.ID
	out.SyncMode = cluster.ClusterSyncMode(in.SyncMode)
	out.APIEndpoint = in.APIEndpoint
	out.SecretRef = (*cluster.LocalSecretReference)(unsafe.Pointer(in.SecretRef))
	out.ImpersonatorSecretRef = (*cluster.LocalSecretReference)(unsafe.Pointer(in.ImpersonatorSecretRef))
	out.InsecureSkipTLSVerification = in.InsecureSkipTLSVerification
	out.ProxyURL = in.ProxyURL
	out.ProxyHeader = *(*map[string]string)(unsafe.Pointer(&in.ProxyHeader))
	out.Provider = in.Provider
	out.Region = in.Region
	out.Zone = in.Zone
	out.Taints = *(*[]v1.Taint)(unsafe.Pointer(&in.Taints))
	return nil
}

// Convert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec is an autogenerated conversion function.
func Convert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec(in *ClusterSpec, out *cluster.ClusterSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterSpec_To_cluster_ClusterSpec(in, out, s)
}

func autoConvert_cluster_ClusterSpec_To_v1alpha1_ClusterSpec(in *cluster.ClusterSpec, out *ClusterSpec, s conversion.Scope) error {
	out.ID = in.ID
	out.SyncMode = ClusterSyncMode(in.SyncMode)
	out.APIEndpoint = in.APIEndpoint
	out.SecretRef = (*LocalSecretReference)(unsafe.Pointer(in.SecretRef))
	out.ImpersonatorSecretRef = (*LocalSecretReference)(unsafe.Pointer(in.ImpersonatorSecretRef))
	out.InsecureSkipTLSVerification = in.InsecureSkipTLSVerification
	out.ProxyURL = in.ProxyURL
	out.ProxyHeader = *(*map[string]string)(unsafe.Pointer(&in.ProxyHeader))
	out.Provider = in.Provider
	out.Region = in.Region
	out.Zone = in.Zone
	out.Taints = *(*[]v1.Taint)(unsafe.Pointer(&in.Taints))
	// WARNING: in.ResourceModels requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus(in *ClusterStatus, out *cluster.ClusterStatus, s conversion.Scope) error {
	out.KubernetesVersion = in.KubernetesVersion
	out.APIEnablements = *(*[]cluster.APIEnablement)(unsafe.Pointer(&in.APIEnablements))
	out.Conditions = *(*[]metav1.Condition)(unsafe.Pointer(&in.Conditions))
	out.NodeSummary = (*cluster.NodeSummary)(unsafe.Pointer(in.NodeSummary))
	if in.ResourceSummary != nil {
		in, out := &in.ResourceSummary, &out.ResourceSummary
		*out = new(cluster.ResourceSummary)
		if err := Convert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ResourceSummary = nil
	}
	return nil
}

// Convert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus is an autogenerated conversion function.
func Convert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus(in *ClusterStatus, out *cluster.ClusterStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_ClusterStatus_To_cluster_ClusterStatus(in, out, s)
}

func autoConvert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus(in *cluster.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	out.KubernetesVersion = in.KubernetesVersion
	out.APIEnablements = *(*[]APIEnablement)(unsafe.Pointer(&in.APIEnablements))
	out.Conditions = *(*[]metav1.Condition)(unsafe.Pointer(&in.Conditions))
	out.NodeSummary = (*NodeSummary)(unsafe.Pointer(in.NodeSummary))
	if in.ResourceSummary != nil {
		in, out := &in.ResourceSummary, &out.ResourceSummary
		*out = new(ResourceSummary)
		if err := Convert_cluster_ResourceSummary_To_v1alpha1_ResourceSummary(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.ResourceSummary = nil
	}
	return nil
}

// Convert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus is an autogenerated conversion function.
func Convert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus(in *cluster.ClusterStatus, out *ClusterStatus, s conversion.Scope) error {
	return autoConvert_cluster_ClusterStatus_To_v1alpha1_ClusterStatus(in, out, s)
}

func autoConvert_v1alpha1_LocalSecretReference_To_cluster_LocalSecretReference(in *LocalSecretReference, out *cluster.LocalSecretReference, s conversion.Scope) error {
	out.Namespace = in.Namespace
	out.Name = in.Name
	return nil
}

// Convert_v1alpha1_LocalSecretReference_To_cluster_LocalSecretReference is an autogenerated conversion function.
func Convert_v1alpha1_LocalSecretReference_To_cluster_LocalSecretReference(in *LocalSecretReference, out *cluster.LocalSecretReference, s conversion.Scope) error {
	return autoConvert_v1alpha1_LocalSecretReference_To_cluster_LocalSecretReference(in, out, s)
}

func autoConvert_cluster_LocalSecretReference_To_v1alpha1_LocalSecretReference(in *cluster.LocalSecretReference, out *LocalSecretReference, s conversion.Scope) error {
	out.Namespace = in.Namespace
	out.Name = in.Name
	return nil
}

// Convert_cluster_LocalSecretReference_To_v1alpha1_LocalSecretReference is an autogenerated conversion function.
func Convert_cluster_LocalSecretReference_To_v1alpha1_LocalSecretReference(in *cluster.LocalSecretReference, out *LocalSecretReference, s conversion.Scope) error {
	return autoConvert_cluster_LocalSecretReference_To_v1alpha1_LocalSecretReference(in, out, s)
}

func autoConvert_v1alpha1_NodeSummary_To_cluster_NodeSummary(in *NodeSummary, out *cluster.NodeSummary, s conversion.Scope) error {
	out.TotalNum = in.TotalNum
	out.ReadyNum = in.ReadyNum
	return nil
}

// Convert_v1alpha1_NodeSummary_To_cluster_NodeSummary is an autogenerated conversion function.
func Convert_v1alpha1_NodeSummary_To_cluster_NodeSummary(in *NodeSummary, out *cluster.NodeSummary, s conversion.Scope) error {
	return autoConvert_v1alpha1_NodeSummary_To_cluster_NodeSummary(in, out, s)
}

func autoConvert_cluster_NodeSummary_To_v1alpha1_NodeSummary(in *cluster.NodeSummary, out *NodeSummary, s conversion.Scope) error {
	out.TotalNum = in.TotalNum
	out.ReadyNum = in.ReadyNum
	return nil
}

// Convert_cluster_NodeSummary_To_v1alpha1_NodeSummary is an autogenerated conversion function.
func Convert_cluster_NodeSummary_To_v1alpha1_NodeSummary(in *cluster.NodeSummary, out *NodeSummary, s conversion.Scope) error {
	return autoConvert_cluster_NodeSummary_To_v1alpha1_NodeSummary(in, out, s)
}

func autoConvert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary(in *ResourceSummary, out *cluster.ResourceSummary, s conversion.Scope) error {
	out.Allocatable = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocatable))
	out.Allocating = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocating))
	out.Allocated = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocated))
	return nil
}

// Convert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary is an autogenerated conversion function.
func Convert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary(in *ResourceSummary, out *cluster.ResourceSummary, s conversion.Scope) error {
	return autoConvert_v1alpha1_ResourceSummary_To_cluster_ResourceSummary(in, out, s)
}

func autoConvert_cluster_ResourceSummary_To_v1alpha1_ResourceSummary(in *cluster.ResourceSummary, out *ResourceSummary, s conversion.Scope) error {
	out.Allocatable = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocatable))
	out.Allocating = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocating))
	out.Allocated = *(*v1.ResourceList)(unsafe.Pointer(&in.Allocated))
	// WARNING: in.AllocatableModeling requires manual conversion: does not exist in peer-type
	return nil
}
