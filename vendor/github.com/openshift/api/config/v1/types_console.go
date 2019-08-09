package v1

<<<<<<< HEAD
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
=======
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
>>>>>>> 79bfea2d (update vendor)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

<<<<<<< HEAD
// Console holds cluster-wide configuration for the web console, including the
// logout URL, and reports the public URL of the console. The canonical name is
// `cluster`.
type Console struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
=======
// Console holds cluster-wide information about Console.  The canonical name is `cluster`
type Console struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
>>>>>>> 79bfea2d (update vendor)
	// +required
	Spec ConsoleSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status ConsoleStatus `json:"status"`
}

<<<<<<< HEAD
// ConsoleSpec is the specification of the desired behavior of the Console.
=======
>>>>>>> 79bfea2d (update vendor)
type ConsoleSpec struct {
	// +optional
	Authentication ConsoleAuthentication `json:"authentication"`
}

<<<<<<< HEAD
// ConsoleStatus defines the observed status of the Console.
=======
>>>>>>> 79bfea2d (update vendor)
type ConsoleStatus struct {
	// The URL for the console. This will be derived from the host for the route that
	// is created for the console.
	ConsoleURL string `json:"consoleURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ConsoleList struct {
	metav1.TypeMeta `json:",inline"`
<<<<<<< HEAD
	metav1.ListMeta `json:"metadata"`

	Items []Console `json:"items"`
}

// ConsoleAuthentication defines a list of optional configuration for console authentication.
=======
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata"`
	Items           []Console `json:"items"`
}

>>>>>>> 79bfea2d (update vendor)
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
	// +kubebuilder:validation:Pattern=`^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$`
=======
	// +kubebuilder:validation:Pattern=^$|^((https):\/\/?)[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/?))$
>>>>>>> 79bfea2d (update vendor)
	LogoutRedirect string `json:"logoutRedirect,omitempty"`
}
