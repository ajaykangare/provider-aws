// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthority) DeepCopyInto(out *CertificateAuthority) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthority.
func (in *CertificateAuthority) DeepCopy() *CertificateAuthority {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthority)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CertificateAuthority) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityExternalStatus) DeepCopyInto(out *CertificateAuthorityExternalStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityExternalStatus.
func (in *CertificateAuthorityExternalStatus) DeepCopy() *CertificateAuthorityExternalStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityExternalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityList) DeepCopyInto(out *CertificateAuthorityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CertificateAuthority, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityList.
func (in *CertificateAuthorityList) DeepCopy() *CertificateAuthorityList {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CertificateAuthorityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityParameters) DeepCopyInto(out *CertificateAuthorityParameters) {
	*out = *in
	if in.IdempotencyToken != nil {
		in, out := &in.IdempotencyToken, &out.IdempotencyToken
		*out = new(string)
		**out = **in
	}
	if in.Organization != nil {
		in, out := &in.Organization, &out.Organization
		*out = new(string)
		**out = **in
	}
	if in.OrganizationalUnit != nil {
		in, out := &in.OrganizationalUnit, &out.OrganizationalUnit
		*out = new(string)
		**out = **in
	}
	if in.Country != nil {
		in, out := &in.Country, &out.Country
		*out = new(string)
		**out = **in
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
	if in.Locality != nil {
		in, out := &in.Locality, &out.Locality
		*out = new(string)
		**out = **in
	}
	if in.CommonName != nil {
		in, out := &in.CommonName, &out.CommonName
		*out = new(string)
		**out = **in
	}
	if in.RevocationConfigurationEnabled != nil {
		in, out := &in.RevocationConfigurationEnabled, &out.RevocationConfigurationEnabled
		*out = new(bool)
		**out = **in
	}
	if in.S3BucketName != nil {
		in, out := &in.S3BucketName, &out.S3BucketName
		*out = new(string)
		**out = **in
	}
	if in.CustomCname != nil {
		in, out := &in.CustomCname, &out.CustomCname
		*out = new(string)
		**out = **in
	}
	if in.ExpirationInDays != nil {
		in, out := &in.ExpirationInDays, &out.ExpirationInDays
		*out = new(int64)
		**out = **in
	}
	if in.PermanentDeletionTimeInDays != nil {
		in, out := &in.PermanentDeletionTimeInDays, &out.PermanentDeletionTimeInDays
		*out = new(int64)
		**out = **in
	}
	if in.DistinguishedNameQualifier != nil {
		in, out := &in.DistinguishedNameQualifier, &out.DistinguishedNameQualifier
		*out = new(string)
		**out = **in
	}
	if in.GenerationQualifier != nil {
		in, out := &in.GenerationQualifier, &out.GenerationQualifier
		*out = new(string)
		**out = **in
	}
	if in.Initials != nil {
		in, out := &in.Initials, &out.Initials
		*out = new(string)
		**out = **in
	}
	if in.GivenName != nil {
		in, out := &in.GivenName, &out.GivenName
		*out = new(string)
		**out = **in
	}
	if in.Pseudonym != nil {
		in, out := &in.Pseudonym, &out.Pseudonym
		*out = new(string)
		**out = **in
	}
	if in.SerialNumber != nil {
		in, out := &in.SerialNumber, &out.SerialNumber
		*out = new(string)
		**out = **in
	}
	if in.Surname != nil {
		in, out := &in.Surname, &out.Surname
		*out = new(string)
		**out = **in
	}
	if in.Title != nil {
		in, out := &in.Title, &out.Title
		*out = new(string)
		**out = **in
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]Tag, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityParameters.
func (in *CertificateAuthorityParameters) DeepCopy() *CertificateAuthorityParameters {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityPermission) DeepCopyInto(out *CertificateAuthorityPermission) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityPermission.
func (in *CertificateAuthorityPermission) DeepCopy() *CertificateAuthorityPermission {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityPermission)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CertificateAuthorityPermission) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityPermissionList) DeepCopyInto(out *CertificateAuthorityPermissionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CertificateAuthorityPermission, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityPermissionList.
func (in *CertificateAuthorityPermissionList) DeepCopy() *CertificateAuthorityPermissionList {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityPermissionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CertificateAuthorityPermissionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityPermissionParameters) DeepCopyInto(out *CertificateAuthorityPermissionParameters) {
	*out = *in
	if in.CertificateAuthorityArn != nil {
		in, out := &in.CertificateAuthorityArn, &out.CertificateAuthorityArn
		*out = new(string)
		**out = **in
	}
	if in.CertificateAuthorityArnRef != nil {
		in, out := &in.CertificateAuthorityArnRef, &out.CertificateAuthorityArnRef
		*out = new(corev1alpha1.Reference)
		**out = **in
	}
	if in.CertificateAuthorityArnSelector != nil {
		in, out := &in.CertificateAuthorityArnSelector, &out.CertificateAuthorityArnSelector
		*out = new(corev1alpha1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.Actions != nil {
		in, out := &in.Actions, &out.Actions
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Principal != nil {
		in, out := &in.Principal, &out.Principal
		*out = new(string)
		**out = **in
	}
	if in.SourceAccount != nil {
		in, out := &in.SourceAccount, &out.SourceAccount
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityPermissionParameters.
func (in *CertificateAuthorityPermissionParameters) DeepCopy() *CertificateAuthorityPermissionParameters {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityPermissionParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityPermissionSpec) DeepCopyInto(out *CertificateAuthorityPermissionSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityPermissionSpec.
func (in *CertificateAuthorityPermissionSpec) DeepCopy() *CertificateAuthorityPermissionSpec {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityPermissionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityPermissionStatus) DeepCopyInto(out *CertificateAuthorityPermissionStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityPermissionStatus.
func (in *CertificateAuthorityPermissionStatus) DeepCopy() *CertificateAuthorityPermissionStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityPermissionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthoritySpec) DeepCopyInto(out *CertificateAuthoritySpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthoritySpec.
func (in *CertificateAuthoritySpec) DeepCopy() *CertificateAuthoritySpec {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthoritySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateAuthorityStatus) DeepCopyInto(out *CertificateAuthorityStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	out.AtProvider = in.AtProvider
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateAuthorityStatus.
func (in *CertificateAuthorityStatus) DeepCopy() *CertificateAuthorityStatus {
	if in == nil {
		return nil
	}
	out := new(CertificateAuthorityStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tag) DeepCopyInto(out *Tag) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tag.
func (in *Tag) DeepCopy() *Tag {
	if in == nil {
		return nil
	}
	out := new(Tag)
	in.DeepCopyInto(out)
	return out
}
