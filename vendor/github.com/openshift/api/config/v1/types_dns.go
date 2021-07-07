package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 03397665 (update api)
// DNS holds cluster-wide information about DNS. The canonical name is `cluster`
type DNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec holds user settable values for configuration
	// +kubebuilder:validation:Required
<<<<<<< HEAD
=======
// DNS holds cluster-wide information about DNS.  The canonical name is `cluster`
// TODO this object is an example of a possible grouping and is subject to change or removal
=======
// DNS holds cluster-wide information about DNS. The canonical name is `cluster`
>>>>>>> e879a141 (alibabacloud machine-api provider)
type DNS struct {
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
	Spec DNSSpec `json:"spec"`
	// status holds observed values from the cluster. They may not be overridden.
	// +optional
	Status DNSStatus `json:"status"`
}

type DNSSpec struct {
	// baseDomain is the base domain of the cluster. All managed DNS records will
	// be sub-domains of this base.
	//
	// For example, given the base domain `openshift.example.com`, an API server
	// DNS record may be created for `cluster-api.openshift.example.com`.
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	//
	// Once set, this field cannot be changed.
	BaseDomain string `json:"baseDomain"`
	// publicZone is the location where all the DNS records that are publicly accessible to
	// the internet exist.
	//
	// If this field is nil, no public records should be created.
	//
	// Once set, this field cannot be changed.
	//
=======
=======
	//
	// Once set, this field cannot be changed.
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	//
	// Once set, this field cannot be changed.
>>>>>>> 03397665 (update api)
	BaseDomain string `json:"baseDomain"`
	// publicZone is the location where all the DNS records that are publicly accessible to
	// the internet exist.
	//
	// If this field is nil, no public records should be created.
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 79bfea2d (update vendor)
=======
	//
	// Once set, this field cannot be changed.
	//
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	//
	// Once set, this field cannot be changed.
	//
>>>>>>> 03397665 (update api)
	// +optional
	PublicZone *DNSZone `json:"publicZone,omitempty"`
	// privateZone is the location where all the DNS records that are only available internally
	// to the cluster exist.
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	//
	// If this field is nil, no private records should be created.
	//
	// Once set, this field cannot be changed.
	//
<<<<<<< HEAD
<<<<<<< HEAD
=======
	// If this field is nil, no private records should be created.
>>>>>>> 79bfea2d (update vendor)
=======
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
>>>>>>> 03397665 (update api)
	// +optional
	PrivateZone *DNSZone `json:"privateZone,omitempty"`
}

// DNSZone is used to define a DNS hosted zone.
// A zone can be identified by an ID or tags.
type DNSZone struct {
	// id is the identifier that can be used to find the DNS hosted zone.
	//
	// on AWS zone can be fetched using `ID` as id in [1]
	// on Azure zone can be fetched using `ID` as a pre-determined name in [2],
	// on GCP zone can be fetched using `ID` as a pre-determined name in [3].
	//
	// [1]: https://docs.aws.amazon.com/cli/latest/reference/route53/get-hosted-zone.html#options
	// [2]: https://docs.microsoft.com/en-us/cli/azure/network/dns/zone?view=azure-cli-latest#az-network-dns-zone-show
	// [3]: https://cloud.google.com/dns/docs/reference/v1/managedZones/get
	// +optional
	ID string `json:"id,omitempty"`

	// tags can be used to query the DNS hosted zone.
	//
	// on AWS, resourcegroupstaggingapi [1] can be used to fetch a zone using `Tags` as tag-filters,
	//
	// [1]: https://docs.aws.amazon.com/cli/latest/reference/resourcegroupstaggingapi/get-resources.html#options
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

type DNSStatus struct {
	// dnsSuffix (service-ca amongst others)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DNSList struct {
	metav1.TypeMeta `json:",inline"`
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	metav1.ListMeta `json:"metadata"`

	Items []DNS `json:"items"`
=======
	// Standard object's metadata.
	metav1.ListMeta `json:"metadata"`
	Items           []DNS `json:"items"`
>>>>>>> 79bfea2d (update vendor)
=======
	metav1.ListMeta `json:"metadata"`

	Items []DNS `json:"items"`
>>>>>>> e879a141 (alibabacloud machine-api provider)
=======
	metav1.ListMeta `json:"metadata"`

	Items []DNS `json:"items"`
>>>>>>> 03397665 (update api)
}
