package v1

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
=======
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
>>>>>>> 79bfea2d (update vendor)
=======
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
>>>>>>> 03397665 (update api)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
// Console holds cluster-wide configuration for the web console, including the
// logout URL, and reports the public URL of the console. The canonical name is
// `cluster`.
type Console struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
<<<<<<< HEAD
=======
// Console holds cluster-wide information about Console.  The canonical name is `cluster`
=======
// Console holds cluster-wide configuration for the web console, including the
// logout URL, and reports the public URL of the console. The canonical name is
// `cluster`.
>>>>>>> e879a141 (alibabacloud machine-api provider)
type Console struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
<<<<<<< HEAD
>>>>>>> 79bfea2d (update vendor)
=======
	// +kubebuilder:validation:Required
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	// +required
	Spec ConsoleSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ConsoleStatus `json:"status"`
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// ConsoleSpec is the specification of the desired behavior of the Console.
=======
>>>>>>> 79bfea2d (update vendor)
=======
// ConsoleSpec is the specification of the desired behavior of the Console.
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// ConsoleSpec is the specification of the desired behavior of the Console.
>>>>>>> 03397665 (update api)
type ConsoleSpec struct {
	// +optional
	Authentication ConsoleAuthentication `json:"authentication"`
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// ConsoleStatus defines the observed status of the Console.
=======
>>>>>>> 79bfea2d (update vendor)
=======
// ConsoleStatus defines the observed status of the Console.
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// ConsoleStatus defines the observed status of the Console.
>>>>>>> 03397665 (update api)
type ConsoleStatus struct {
	// The URL for the console. This will be derived from the host for the route that
	// is created for the console.
	ConsoleURL string `json:"consoleURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ConsoleList struct {
	metav1.TypeMeta `json:",inline"`
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	metav1.ListMeta `json:"metadata"`

	Items []Console `json:"items"`
}

// ConsoleAuthentication defines a list of optional configuration for console authentication.
=======
	// Standard object's metadata.
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	metav1.ListMeta `json:"metadata"`

	Items []Console `json:"items"`
}

<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 79bfea2d (update vendor)
=======
// ConsoleAuthentication defines a list of optional configuration for console authentication.
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
// ConsoleAuthentication defines a list of optional configuration for console authentication.
>>>>>>> 03397665 (update api)
type ConsoleAuthentication struct {
	// An optional, absolute URL to redirect web browsers to after logging out of
	// the console. If not specified, it will redirect to the default login page.
	// This is required when using an identity provider that supports single
	// sign-on (SSO) such as:
	// - OpenID (Keycloak, Azure)
	// - RequestHeader (GSSAPI, SSPI, SAML)
	// - OAuth (GitHub, GitLab, Google)
	// Logging out of the console will destroy the user's token. The logoutRedirect
	// provides the user the option to perform single logout (SLO) through the identity
	// provider to destroy their single sign-on session.
	// +optional
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	// +kubebuilder:validation:Pattern=`^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$`
=======
	// +kubebuilder:validation:Pattern=^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$
>>>>>>> 79bfea2d (update vendor)
=======
	// +kubebuilder:validation:Pattern=`^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$`
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	// +kubebuilder:validation:Pattern=`^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$`
>>>>>>> 03397665 (update api)
	LogoutRedirect string `json:"logoutRedirect,omitempty"`
}
