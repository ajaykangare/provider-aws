package certificateauthority

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"

	"github.com/crossplane/provider-aws/apis/certificatemanager/v1alpha1"
)

// Client defines the CertificateManager operations
type Client interface {
	CreateCertificateAuthorityRequest(*acmpca.CreateCertificateAuthorityInput) acmpca.CreateCertificateAuthorityRequest
	CreatePermissionRequest(*acmpca.CreatePermissionInput) acmpca.CreatePermissionRequest
	DeleteCertificateAuthorityRequest(*acmpca.DeleteCertificateAuthorityInput) acmpca.DeleteCertificateAuthorityRequest
	DeletePermissionRequest(*acmpca.DeletePermissionInput) acmpca.DeletePermissionRequest
	UpdateCertificateAuthorityRequest(*acmpca.UpdateCertificateAuthorityInput) acmpca.UpdateCertificateAuthorityRequest
	DescribeCertificateAuthorityRequest(*acmpca.DescribeCertificateAuthorityInput) acmpca.DescribeCertificateAuthorityRequest
	ListTagsRequest(*acmpca.ListTagsInput) acmpca.ListTagsRequest
	UntagCertificateAuthorityRequest(*acmpca.UntagCertificateAuthorityInput) acmpca.UntagCertificateAuthorityRequest
	TagCertificateAuthorityRequest(*acmpca.TagCertificateAuthorityInput) acmpca.TagCertificateAuthorityRequest
	// GetCertificateAuthorityCertificateRequest(*acmpca.GetCertificateAuthorityCertificateInput) acmpca.GetCertificateAuthorityCertificateRequest
	// ImportCertificateAuthorityCertificateRequest(*acmpca.ImportCertificateAuthorityCertificateInput) acmpca.ImportCertificateAuthorityCertificateRequest
	// GetCertificateAuthorityCsrRequest(*acmpca.GetCertificateAuthorityCsrInput) acmpca.GetCertificateAuthorityCsrRequest
	// IssueCertificateRequest(*acmpca.IssueCertificateInput) acmpca.IssueCertificateRequest
}

// NewClient returns a new client using AWS credentials as JSON encoded data.
func NewClient(conf *aws.Config) (Client, error) {
	return acmpca.New(*conf), nil
}

// GenerateCreateCertificateAuthorityInput from certificateAuthorityParameters
func GenerateCreateCertificateAuthorityInput(p *v1alpha1.CertificateAuthorityParameters) *acmpca.CreateCertificateAuthorityInput {
	m := &acmpca.CreateCertificateAuthorityInput{

		IdempotencyToken:                  p.IdempotencyToken,
		CertificateAuthorityConfiguration: GenerateCertificateAuthorityConfiguration(p),
		RevocationConfiguration:           GenerateRevocationConfiguration(p),
	}

	if strings.EqualFold(p.Type, "ROOT") {
		m.CertificateAuthorityType = acmpca.CertificateAuthorityTypeRoot
	} else if strings.EqualFold(p.Type, "SUBORDINATE") {
		m.CertificateAuthorityType = acmpca.CertificateAuthorityTypeSubordinate
	}

	m.Tags = make([]acmpca.Tag, len(p.Tags))
	for i, val := range p.Tags {
		m.Tags[i] = acmpca.Tag{
			Key:   aws.String(val.Key),
			Value: aws.String(val.Value),
		}
	}

	return m
}

// GenerateCertificateAuthorityConfiguration from certificateAuthorityParameters
func GenerateCertificateAuthorityConfiguration(p *v1alpha1.CertificateAuthorityParameters) *acmpca.CertificateAuthorityConfiguration { // nolint:gocyclo

	m := &acmpca.CertificateAuthorityConfiguration{
		Subject: &acmpca.ASN1Subject{
			CommonName:                 p.CommonName,
			Country:                    p.Country,
			DistinguishedNameQualifier: p.DistinguishedNameQualifier,
			GenerationQualifier:        p.GenerationQualifier,
			GivenName:                  p.GivenName,
			Initials:                   p.Initials,
			Locality:                   p.Locality,
			Organization:               p.Organization,
			OrganizationalUnit:         p.OrganizationalUnit,
			Pseudonym:                  p.Pseudonym,
			SerialNumber:               p.SerialNumber,
			State:                      p.State,
			Surname:                    p.Surname,
			Title:                      p.Title,
		},
	}

	switch p.SigningAlgorithm {
	case "SHA256WITHECDSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha256withecdsa
	case "SHA384WITHECDSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha384withecdsa
	case "SHA512WITHECDSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha512withecdsa
	case "SHA256WITHRSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha256withrsa
	case "SHA384WITHRSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha384withrsa
	case "SHA512WITHRSA":
		m.SigningAlgorithm = acmpca.SigningAlgorithmSha512withrsa
	}

	switch p.KeyAlgorithm {
	case "RSA_2048":
		m.KeyAlgorithm = acmpca.KeyAlgorithmRsa2048
	case "RSA_4096":
		m.KeyAlgorithm = acmpca.KeyAlgorithmRsa4096
	case "EC_prime256v1":
		m.KeyAlgorithm = acmpca.KeyAlgorithmEcPrime256v1
	case "EC_secp384r1":
		m.KeyAlgorithm = acmpca.KeyAlgorithmEcSecp384r1
	}

	return m

}

// GenerateRevocationConfiguration from certificateAuthorityParameters
func GenerateRevocationConfiguration(p *v1alpha1.CertificateAuthorityParameters) *acmpca.RevocationConfiguration {

	m := &acmpca.RevocationConfiguration{
		CrlConfiguration: &acmpca.CrlConfiguration{
			CustomCname:      p.CustomCname,
			Enabled:          p.RevocationConfigurationEnabled,
			ExpirationInDays: p.ExpirationInDays,
			S3BucketName:     p.S3BucketName,
		},
	}

	return m
}

// GenerateCertificateAuthorityStatus from status
func GenerateCertificateAuthorityStatus(status string) acmpca.CertificateAuthorityStatus {

	var m acmpca.CertificateAuthorityStatus
	switch strings.ToUpper(status) {
	case "CREATING":
		m = acmpca.CertificateAuthorityStatusCreating
	case "PENDING_CERTIFICATE":
		m = acmpca.CertificateAuthorityStatusPendingCertificate
	case "ACTIVE":
		m = acmpca.CertificateAuthorityStatusActive
	case "DELETED":
		m = acmpca.CertificateAuthorityStatusDeleted
	case "DISABLED":
		m = acmpca.CertificateAuthorityStatusDisabled
	case "EXPIRED":
		m = acmpca.CertificateAuthorityStatusExpired
	case "FAILED":
		m = acmpca.CertificateAuthorityStatusFailed
	}
	return m
}

// GenerateUpdateCertificateAuthorityInput from CertificateAuthority
func GenerateUpdateCertificateAuthorityInput(cr *v1alpha1.CertificateAuthority) *acmpca.UpdateCertificateAuthorityInput {

	return &acmpca.UpdateCertificateAuthorityInput{
		CertificateAuthorityArn: aws.String(cr.Status.AtProvider.CertificateAuthorityArn),
		RevocationConfiguration: GenerateRevocationConfiguration(&cr.Spec.ForProvider),
		Status:                  GenerateCertificateAuthorityStatus(cr.Spec.ForProvider.Status),
	}
}

// LateInitializeCertificateAuthority fills the empty fields in *v1beta1.CertificateAuthorityParameters with
// the values seen in acmpca.CertificateAuthority.
func LateInitializeCertificateAuthority(in *v1alpha1.CertificateAuthorityParameters, certificateAuthority *acmpca.CertificateAuthority) { // nolint:gocyclo
	if certificateAuthority == nil {
		return
	}

	if in.Type == "" && string(certificateAuthority.Type) != "" {
		in.Type = string(certificateAuthority.Type)
	}

	if in.Status == "" && string(certificateAuthority.Status) != "" {
		in.Status = string(certificateAuthority.Status)
	}

	if aws.StringValue(in.SerialNumber) == "" && aws.StringValue(certificateAuthority.Serial) != "" {
		in.SerialNumber = certificateAuthority.Serial
	}

	if in.ExpirationInDays == nil && certificateAuthority.RevocationConfiguration.CrlConfiguration.ExpirationInDays != nil {
		in.ExpirationInDays = certificateAuthority.RevocationConfiguration.CrlConfiguration.ExpirationInDays
	}

}

// IsCertificateAuthorityUpToDate checks whether there is a change in any of the modifiable fields.
func IsCertificateAuthorityUpToDate(p *v1alpha1.CertificateAuthority, cd acmpca.CertificateAuthority, tags []acmpca.Tag) bool { // nolint:gocyclo

	if !strings.EqualFold(aws.StringValue(p.Spec.ForProvider.CustomCname), aws.StringValue(cd.RevocationConfiguration.CrlConfiguration.CustomCname)) {
		return false
	}

	if !strings.EqualFold(aws.StringValue(p.Spec.ForProvider.S3BucketName), aws.StringValue(cd.RevocationConfiguration.CrlConfiguration.S3BucketName)) {
		return false
	}

	if aws.BoolValue(p.Spec.ForProvider.RevocationConfigurationEnabled) != aws.BoolValue(cd.RevocationConfiguration.CrlConfiguration.Enabled) {
		return false
	}

	if aws.Int64Value(p.Spec.ForProvider.ExpirationInDays) != aws.Int64Value(cd.RevocationConfiguration.CrlConfiguration.ExpirationInDays) {
		return false
	}

	if len(p.Spec.ForProvider.Tags) != len(tags) {
		return false
	}

	pTags := make(map[string]string, len(p.Spec.ForProvider.Tags))
	for _, tag := range p.Spec.ForProvider.Tags {
		pTags[tag.Key] = tag.Value
	}
	for _, tag := range tags {
		val, ok := pTags[aws.StringValue(tag.Key)]
		if !ok || !strings.EqualFold(val, aws.StringValue(tag.Value)) {
			return false
		}
	}

	return p.Spec.ForProvider.CertificateRenewalPermissionAllow == p.Status.AtProvider.RenewalPermission
}

// GenerateCertificateAuthorityExternalStatus is used to produce CertificateAuthorityExternalStatus from acmpca.certificateAuthorityStatus and v1alpha1.CertificateAuthority
func GenerateCertificateAuthorityExternalStatus(certificateAuthority acmpca.CertificateAuthority, p *v1alpha1.CertificateAuthority) v1alpha1.CertificateAuthorityExternalStatus {
	return v1alpha1.CertificateAuthorityExternalStatus{
		CertificateAuthorityArn: aws.StringValue(certificateAuthority.Arn),
		RenewalPermission:       p.Spec.ForProvider.CertificateRenewalPermissionAllow,
	}
}

// IsErrorNotFound returns true if the error code indicates that the item was not found
func IsErrorNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == acmpca.ErrCodeResourceNotFoundException {
			return true
		}
	}

	return false
}
