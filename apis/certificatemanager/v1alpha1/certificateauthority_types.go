/*
Copyright 2019 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Tag represents user-provided metadata that can be associated
type Tag struct {

	// The key name that can be used to look up or retrieve the associated value.
	Key string `json:"key"`

	// The value associated with this tag.
	// +optional
	Value string `json:"value,omitempty"`
}

// CertificateAuthoritySpec defines the desired state of CertificateAuthority
type CertificateAuthoritySpec struct {
	runtimev1alpha1.ResourceSpec `json:",inline"`
	ForProvider                  CertificateAuthorityParameters `json:"forProvider"`
}

// CertificateAuthorityExternalStatus keeps the state of external resource
type CertificateAuthorityExternalStatus struct {
	// String that contains the ARN of the issued certificate Authority
	CertificateAuthorityArn string `json:"certificateArn"`
	RenewalPermission       bool   `json:"renewalPermission"`
}

// An CertificateAuthorityStatus represents the observed state of an CertificateAuthority manager.
type CertificateAuthorityStatus struct {
	runtimev1alpha1.ResourceStatus `json:",inline"`
	AtProvider                     CertificateAuthorityExternalStatus `json:"atProvider"`
}

// CertificateAuthorityParameters defines the desired state of an AWS CertificateAuthority.
type CertificateAuthorityParameters struct {
	// Type of the certificate authority
	Type string `json:"type"`

	// Status of the certificate authority
	// +optional
	Status string `json:"status"`

	// Token to distinguish between calls to RequestCertificate.
	// +optional
	IdempotencyToken *string `json:"idempotencyToken,omitempty"`

	// Organization legal name
	Organization *string `json:"organization"`

	// Organization's subdivision or unit
	OrganizationalUnit *string `json:"organizationalUnit"`

	// Two-digit code that specifies the country
	Country *string `json:"country"`

	// State in which the subject of the certificate is located
	State *string `json:"state"`

	// The locality such as a city or town
	Locality *string `json:"locality"`

	// FQDN associated with the certificate subject
	CommonName *string `json:"commonName"`

	// Type of the public key algorithm
	KeyAlgorithm string `json:"keyAlgorithm"`

	// Algorithm that private CA uses to sign certificate requests
	SigningAlgorithm string `json:"signingAlgorithm"`

	// Boolean value that specifies certificate revocation
	RevocationConfigurationEnabled *bool `json:"revocationConfigurationEnabled"`

	// Name of the S3 bucket that contains the CRL
	S3BucketName *string `json:"s3BucketName"`

	// Alias for the CRL distribution point
	// +optional
	CustomCname *string `json:"customCname,omitempty"`

	// Number of days until a certificate expires
	ExpirationInDays *int64 `json:"expirationInDays,omitempty"`

	// The number of days to make a CA restorable after it has been deleted
	// +optional
	PermanentDeletionTimeInDays *int64 `json:"permanentDeletionTimeInDays,omitempty"`

	// The CertificateRenewalPermissionAllow decides Permissions for ACM renewals
	CertificateRenewalPermissionAllow bool `json:"certificateRenewalPermissionAllow"`

	// Disambiguating information for the certificate subject.
	// +optional
	DistinguishedNameQualifier *string `json:"distinguishedNameQualifier,omitempty"`

	// Typically a qualifier appended to the name of an individual
	// +optional
	GenerationQualifier *string `json:"generationQualifier,omitempty"`

	// Concatenation of first letter of the GivenName, Middle name and SurName.
	// +optional
	Initials *string `json:"initials,omitempty"`

	// First name
	// +optional
	GivenName *string `json:"givenName,omitempty"`

	// Shortened version of a longer GivenName
	// +optional
	Pseudonym *string `json:"pseudonym,omitempty"`

	// The certificate serial number.
	// +optional
	SerialNumber *string `json:"serialNumber,omitempty"`

	// Surname
	// +optional
	Surname *string `json:"surname,omitempty"`

	// Title
	// +optional
	Title *string `json:"title,omitempty"`

	// One or more resource tags to associate with the certificateAuthority.
	Tags []Tag `json:"tags,omitempty"`
}

// +kubebuilder:object:root=true

// CertificateAuthority is a managed resource that represents an AWS CertificateAuthority Manager.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type CertificateAuthority struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateAuthoritySpec   `json:"spec,omitempty"`
	Status CertificateAuthorityStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CertificateAuthorityList contains a list of CertificateAuthority
type CertificateAuthorityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertificateAuthority `json:"items"`
}
