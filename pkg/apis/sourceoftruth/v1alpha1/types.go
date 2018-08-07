package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UwsgiAppList struct describing all Uwsgi Apps in the cluster
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type UwsgiAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []UwsgiApp `json:"items"`
}

// UwsgiApp top level struct describing the CR elements
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type UwsgiApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              UwsgiAppSpec   `json:"spec"`
	Status            UwsgiAppStatus `json:"status,omitempty"`
}

// UwsgiAppSpec Definition of Uwsgi Application Spec
type UwsgiAppSpec struct {
	Image              string                 `json:"image"`
	ConfigVersion      string                 `json:"conf_version"`
	ApplicationName    string                 `json:"application_name"`
	ApplicationVersion string                 `json:"application_version,omitempty"`
	Replicas           int32                  `json:"replicas"`
	Command            []string               `json:"command"`
	Ports              []corev1.ContainerPort `json:"ports"`
}

// UwsgiAppStatus Definition of Uwsgi Application Status
type UwsgiAppStatus struct {
	Nodes []string `json:"nodes"`
}
