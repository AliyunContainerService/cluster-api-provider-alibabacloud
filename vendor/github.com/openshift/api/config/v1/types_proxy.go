package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Proxy holds cluster-wide information on how to configure default proxies for the cluster. The canonical name is `cluster`
type Proxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
<<<<<<< HEAD
<<<<<<< HEAD

	// Spec holds user-settable values for the proxy configuration
	// +kubebuilder:validation:Required
	// +required
	Spec ProxySpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ProxyStatus `json:"status"`
}

// ProxySpec contains cluster proxy creation configuration.
=======
=======

>>>>>>> e879a141 (alibabacloud machine-api provider)
	// Spec holds user-settable values for the proxy configuration
	// +kubebuilder:validation:Required
	// +required
	Spec ProxySpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ProxyStatus `json:"status"`
}

<<<<<<< HEAD
>>>>>>> 79bfea2d (update vendor)
=======
// ProxySpec contains cluster proxy creation configuration.
>>>>>>> e879a141 (alibabacloud machine-api provider)
type ProxySpec struct {
	// httpProxy is the URL of the proxy for HTTP requests.  Empty means unset and will not result in an env var.
	// +optional
	HTTPProxy string `json:"httpProxy,omitempty"`

	// httpsProxy is the URL of the proxy for HTTPS requests.  Empty means unset and will not result in an env var.
	// +optional
	HTTPSProxy string `json:"httpsProxy,omitempty"`

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	// noProxy is a comma-separated list of hostnames and/or CIDRs for which the proxy should not be used.
	// Empty means unset and will not result in an env var.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`

	// readinessEndpoints is a list of endpoints used to verify readiness of the proxy.
	// +optional
	ReadinessEndpoints []string `json:"readinessEndpoints,omitempty"`

	// trustedCA is a reference to a ConfigMap containing a CA certificate bundle.
	// The trustedCA field should only be consumed by a proxy validator. The
	// validator is responsible for reading the certificate bundle from the required
	// key "ca-bundle.crt", merging it with the system default trust bundle,
	// and writing the merged trust bundle to a ConfigMap named "trusted-ca-bundle"
	// in the "openshift-config-managed" namespace. Clients that expect to make
	// proxy connections must use the trusted-ca-bundle for all HTTPS requests to
	// the proxy, and may use the trusted-ca-bundle for non-proxy HTTPS requests as
	// well.
	//
	// The namespace for the ConfigMap referenced by trustedCA is
	// "openshift-config". Here is an example ConfigMap (in yaml):
	//
	// apiVersion: v1
	// kind: ConfigMap
	// metadata:
	//  name: user-ca-bundle
	//  namespace: openshift-config
	//  data:
	//    ca-bundle.crt: |
	//      -----BEGIN CERTIFICATE-----
	//      Custom CA certificate bundle.
	//      -----END CERTIFICATE-----
	//
	// +optional
	TrustedCA ConfigMapNameReference `json:"trustedCA,omitempty"`
}

// ProxyStatus shows current known state of the cluster proxy.
type ProxyStatus struct {
	// httpProxy is the URL of the proxy for HTTP requests.
	// +optional
	HTTPProxy string `json:"httpProxy,omitempty"`

	// httpsProxy is the URL of the proxy for HTTPS requests.
	// +optional
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// noProxy is a comma-separated list of hostnames and/or CIDRs for which the proxy should not be used.
<<<<<<< HEAD
=======
	// noProxy is the list of domains for which the proxy should not be used.  Empty means unset and will not result in an env var.
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ProxyList struct {
	metav1.TypeMeta `json:",inline"`
<<<<<<< HEAD
<<<<<<< HEAD
	metav1.ListMeta `json:"metadata"`

	Items []Proxy `json:"items"`
=======
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata"`
	Items           []Proxy `json:"items"`
>>>>>>> 79bfea2d (update vendor)
=======
	metav1.ListMeta `json:"metadata"`

	Items []Proxy `json:"items"`
>>>>>>> e879a141 (alibabacloud machine-api provider)
}
