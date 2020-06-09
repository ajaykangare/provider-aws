package certificateauthoritypermission

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
)

// Client defines the CertificateManager operations
type Client interface {
	CreatePermissionRequest(*acmpca.CreatePermissionInput) acmpca.CreatePermissionRequest
	DeletePermissionRequest(*acmpca.DeletePermissionInput) acmpca.DeletePermissionRequest
	ListPermissionsRequest(*acmpca.ListPermissionsInput) acmpca.ListPermissionsRequest
}

// NewClient returns a new client using AWS credentials as JSON encoded data.
func NewClient(conf *aws.Config) (Client, error) {
	return acmpca.New(*conf), nil
}

// LateInitializeCertificateAuthority fills the empty fields in *v1beta1.CertificateAuthorityParameters with
// the values seen in acmpca.CertificateAuthority.
// func LateInitializeCertificateAuthority(in *v1alpha1.CertificateAuthorityParameters, certificateAuthority *acmpca.CertificateAuthority) { // nolint:gocyclo
// 	if certificateAuthority == nil {
// 		return
// 	}

// 	if string(in.Type) == "" && string(certificateAuthority.Type) != "" {
// 		in.Type = certificateAuthority.Type
// 	}

// 	if (string(in.Status) == "" || in.Status == acmpca.CertificateAuthorityStatusPendingCertificate) && string(certificateAuthority.Status) != "" {
// 		in.Status = certificateAuthority.Status
// 	}

// 	if aws.StringValue(in.SerialNumber) == "" && aws.StringValue(certificateAuthority.Serial) != "" {
// 		in.SerialNumber = certificateAuthority.Serial
// 	}

// 	if in.ExpirationInDays == nil && certificateAuthority.RevocationConfiguration.CrlConfiguration.ExpirationInDays != nil {
// 		in.ExpirationInDays = certificateAuthority.RevocationConfiguration.CrlConfiguration.ExpirationInDays
// 	}

// }

// IsErrorNotFound returns true if the error code indicates that the item was not found
func IsErrorNotFound(err error) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == acmpca.ErrCodeInvalidStateException {
			return true
		}
	}

	return false
}
