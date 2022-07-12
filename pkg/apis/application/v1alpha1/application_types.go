package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/types"
)

// ComponentType tracks the Component type of Application: objectTemplate, helmTemplate.
type ComponentType string

// Constants
const (
	// ObjectType Used to indicate that the component is packaged by the kubernetes object.
	ObjectType ComponentType = "objectTemplate"
	// HelmType Used to indicate that that the component is packaged by the Helm chart.
	HelmType ComponentType = "helmTemplate"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=all,shortName=app
// +kubebuilder:subresource:status

// Application is the Schema for the applications API.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApplicationList contains a list of Application.
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

// ApplicationSpec defines the specification for an Application.
type ApplicationSpec struct {
	// Components indicates the components of the Application.
	Components []Component `json:"components,omitempty"`

	// RevisionHistoryLimit limits the number of items kept in the application's revision history
	// which is mainly used for informational purposes as well as for rollbacks to previous versions.
	// Default is 10.
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`

	// DisplayInfo indicates the display info of the Application.
	DisplayInfo DisplayInfo `json:"displayInfo,omitempty"`
}

type Component struct {
	// Name indicates the name of Component.
	Name string `json:"name"`

	// Type indicates the type of Component, only support objectTemplate and helmTemplate now.
	Type ComponentType `json:"type,omitempty"`

	// ObjectTemplate indicates the Component source based on the kubernetes object.
	ObjectTemplate ObjectTemplate `json:"objectTemplate,omitempty"`

	// HelmTemplate indicates the Component source based on the helm chart.
	HelmTemplate HelmTemplate `json:"helmTemplate,omitempty"`

	// ComponentKinds indicates the CRDs which the component needs.
	ComponentKinds []metav1.GroupVersionKind `json:"componentKinds,omitempty"`

	// Resource indicates the resource info about component.
	Resource Resource `json:"resource,omitempty"`
}

type RepoType string

const (
	// HelmRepo indicates files of the Helm chart is from Helm repo.
	HelmRepo RepoType = "helm"

	// GitRepo indicates files of the Helm chart is from Git repo.
	GitRepo RepoType = "git"
)

type HelmTemplate struct {
	// RepoType indicates the repo type of the Helm chart.
	RepoType RepoType `json:"repoType,omitempty"`

	// Url indicates the url of the Helm chart.
	Url string `json:"url,omitempty"`

	// Chart indicates the name of the Helm chart.
	Chart string `json:"chart,omitempty"`

	// Version indicates the version of the Helm chart.
	Version string `json:"version,omitempty"`

	// Values indicates the values of the Helm chart.
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
}

type Resource struct {
	// Request indicates the resource request which the Helm chart needs.
	Request corev1.ResourceList `json:"request,omitempty"`
}

type ObjectTemplate struct {
	// ObjectReference indicates the objectReference of an component.
	ObjectReference corev1.ObjectReference `json:"objectReference,omitempty"`

	// Weight indicates the weight of the installation order.
	Weight *int32 `json:"weight,omitempty"`
}

type DisplayInfo struct {
	// ComponentGroupKinds is a list of Kinds for Application's components (e.g. Deployments, Pods, Services, CRDs). It
	// can be used in conjunction with the Application's Selector to list or watch the Applications components.
	ComponentGroupKinds []metav1.GroupKind `json:"componentKinds,omitempty"`

	// Descriptor regroups information and metadata about an application.
	Descriptor Descriptor `json:"descriptor,omitempty"`

	// Selector is a label query over kinds that created by the application. It must match the component objects' labels.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Info contains human readable key,value pairs for the Application.
	// +patchStrategy=merge
	// +patchMergeKey=name
	Info []InfoItem `json:"info,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// ApplicationStatus defines controller's the observed state of Application.
type ApplicationStatus struct {
	// ObservedGeneration is the most recent generation observed. It corresponds to the
	// Object's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// Conditions represents the latest state of the object.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,10,rep,name=conditions"`

	// History contains information about the application's revision history.
	History []RevisionHistory `json:"history,omitempty"`

	// Resources embeds a list of object statuses.
	// +optional
	ComponentList `json:",inline,omitempty"`

	// ComponentsReady: number of the components with a Ready Condition.
	// +optional
	ComponentsReady int64 `json:"componentsReady,omitempty"`

	// ComponentsTotal: total number of the components targeted by this Application.
	// +optional
	ComponentsTotal int64 `json:"componentsTotal,omitempty"`
}

type RevisionHistory struct {
	// Name is a name of the RevisionHistory.
	// For example, if the application name is "example" and revision is 1, then this field will be "example-v1".
	Name string `json:"name"`

	// Revision is an auto incrementing identifier of the RevisionHistory.
	Revision int64 `json:"revision"`

	// DeployedAt indicates the time the update of the application completed.
	DeployedAt metav1.Time `json:"deployedAt,omitempty"`

	// Components record the application source at that revision.
	Components []Component `json:"components,omitempty"`
}

// ComponentList is a generic status holder for the top level resource.
type ComponentList struct {
	// componentsStatus hold the status for all components.
	componentsStatus []ApplicationComponentStatus `json:"componentsStatus,omitempty"`
}

// ApplicationComponentStatus is a generic status holder for components.
type ApplicationComponentStatus struct {
	// Name of object
	Name string `json:"name,omitempty"`

	// Kind of object
	Kind string `json:"kind,omitempty"`

	// APIVersion of object
	APIVersion string `json:"apiVersion,omitempty"`

	// Applied indicates whether the component is applied
	Applied bool `json:"applied,omitempty"`

	// Healthy indicates whether the component is healthy
	Healthy bool `json:"healthy,omitempty"`
}

// Descriptor defines the Metadata and information about the Application.
type Descriptor struct {
	// Type is the type of the application (e.g. WordPress, MySQL, Cassandra).
	Type string `json:"type,omitempty"`

	// Version is an optional version indicator for the Application.
	Version string `json:"version,omitempty"`

	// Description is a brief string description of the Application.
	Description string `json:"description,omitempty"`

	// Icons is an optional list of icons for an application. Icon information includes the source, size,
	// and mime type.
	Icons []ImageSpec `json:"icons,omitempty"`

	// Maintainers is an optional list of maintainers of the application. The maintainers in this list maintain the
	// the source code, images, and package for the application.
	Maintainers []ContactData `json:"maintainers,omitempty"`

	// Owners is an optional list of the owners of the installed application. The owners of the application should be
	// contacted in the event of a planned or unplanned disruption affecting the application.
	Owners []ContactData `json:"owners,omitempty"`

	// Keywords is an optional list of key words associated with the application (e.g. MySQL, RDBMS, database).
	Keywords []string `json:"keywords,omitempty"`

	// Links are a list of descriptive URLs intended to be used to surface additional documentation, dashboards, etc.
	Links []Link `json:"links,omitempty"`

	// Notes contain a human readable snippets intended as a quick start for the users of the Application.
	// CommonMark markdown syntax may be used for rich text representation.
	Notes string `json:"notes,omitempty"`
}

// ImageSpec contains information about an image used as an icon.
type ImageSpec struct {
	// The source for image represented as either an absolute URL to the image or a Data URL containing
	// the image. Data URLs are defined in RFC 2397.
	Source string `json:"src"`

	// (optional) The size of the image in pixels (e.g., 25x25).
	Size string `json:"size,omitempty"`

	// (optional) The mine type of the image (e.g., "image/png").
	Type string `json:"type,omitempty"`
}

// ContactData contains information about an individual or organization.
type ContactData struct {
	// Name is the descriptive name.
	Name string `json:"name,omitempty"`

	// Url could typically be a website address.
	URL string `json:"url,omitempty"`

	// Email is the email address.
	Email string `json:"email,omitempty"`
}

// Link contains information about an URL to surface documentation, dashboards, etc.
type Link struct {
	// Description is human readable content explaining the purpose of the link.
	Description string `json:"description,omitempty"`

	// Url typically points at a website address.
	URL string `json:"url,omitempty"`
}

// InfoItem is a human readable key,value pair containing important information about how to access the Application.
type InfoItem struct {
	// Name is a human readable title for this piece of information.
	Name string `json:"name,omitempty"`

	// Type of the value for this InfoItem.
	Type InfoItemType `json:"type,omitempty"`

	// Value is human readable content.
	Value string `json:"value,omitempty"`

	// ValueFrom defines a reference to derive the value from another source.
	ValueFrom *InfoItemSource `json:"valueFrom,omitempty"`
}

// InfoItemType is a string that describes the value of InfoItem
type InfoItemType string

const (
	// ValueInfoItemType const string for value type
	ValueInfoItemType InfoItemType = "Value"
	// ReferenceInfoItemType const string for ref type
	ReferenceInfoItemType InfoItemType = "Reference"
)

// InfoItemSource represents a source for the value of an InfoItem.
type InfoItemSource struct {
	// Type of source.
	Type InfoItemSourceType `json:"type,omitempty"`

	// Selects a key of a Secret.
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty"`

	// Selects a key of a ConfigMap.
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

	// Select a Service.
	ServiceRef *ServiceSelector `json:"serviceRef,omitempty"`

	// Select an Ingress.
	IngressRef *IngressSelector `json:"ingressRef,omitempty"`
}

// InfoItemSourceType is a string
type InfoItemSourceType string

// Constants for info type
const (
	SecretKeyRefInfoItemSourceType    InfoItemSourceType = "SecretKeyRef"
	ConfigMapKeyRefInfoItemSourceType InfoItemSourceType = "ConfigMapKeyRef"
	ServiceRefInfoItemSourceType      InfoItemSourceType = "ServiceRef"
	IngressRefInfoItemSourceType      InfoItemSourceType = "IngressRef"
)

// ConfigMapKeySelector selects a key from a ConfigMap.
type ConfigMapKeySelector struct {
	// The ConfigMap to select from.
	corev1.ObjectReference `json:",inline"`
	// The key to select.
	Key string `json:"key,omitempty"`
}

// SecretKeySelector selects a key from a Secret.
type SecretKeySelector struct {
	// The Secret to select from.
	corev1.ObjectReference `json:",inline"`
	// The key to select.
	Key string `json:"key,omitempty"`
}

// ServiceSelector selects a Service.
type ServiceSelector struct {
	// The Service to select from.
	corev1.ObjectReference `json:",inline"`
	// The optional port to select.
	Port *int32 `json:"port,omitempty"`
	// The optional HTTP path.
	Path string `json:"path,omitempty"`
	// Protocol for the service
	Protocol string `json:"protocol,omitempty"`
}

// IngressSelector selects an Ingress.
type IngressSelector struct {
	// The Ingress to select from.
	corev1.ObjectReference `json:",inline"`
	// The optional host to select.
	Host string `json:"host,omitempty"`
	// The optional HTTP path.
	Path string `json:"path,omitempty"`
	// Protocol for the ingress
	Protocol string `json:"protocol,omitempty"`
}
