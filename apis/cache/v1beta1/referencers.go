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

package v1beta1

// // ResolveReferences of this ReplicationGroup
// func (mg *ReplicationGroup) ResolveReferences(ctx context.Context, c client.Reader) error {
// 	r := reference.NewAPIResolver(c, mg)

// 	// Resolve spec.forProvider.securityGroupIDs
// 	mrsp, err := r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
// 		CurrentValues: mg.Spec.ForProvider.SecurityGroupIDs,
// 		References:    mg.Spec.ForProvider.SecurityGroupIDRefs,
// 		Selector:      mg.Spec.ForProvider.SecurityGroupIDSelector,
// 		To:            reference.To{Managed: &v1beta1.SecurityGroup{}, List: &v1beta1.SecurityGroupList{}},
// 		Extract:       reference.ExternalName(),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	mg.Spec.ForProvider.SecurityGroupIDs = mrsp.ResolvedValues
// 	mg.Spec.ForProvider.SecurityGroupIDRefs = mrsp.ResolvedReferences

// 	// Resolve spec.forProvider.cacheSecurityGroupNames
// 	mrsp, err = r.ResolveMultiple(ctx, reference.MultiResolutionRequest{
// 		CurrentValues: mg.Spec.ForProvider.CacheSecurityGroupNames,
// 		References:    mg.Spec.ForProvider.CacheSecurityGroupNameRefs,
// 		Selector:      mg.Spec.ForProvider.CacheSecurityGroupNameSelector,
// 		To:            reference.To{Managed: &v1beta1.SecurityGroup{}, List: &v1beta1.SecurityGroupList{}},
// 		Extract:       v1beta1.SecurityGroupName(),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	mg.Spec.ForProvider.CacheSecurityGroupNames = mrsp.ResolvedValues
// 	mg.Spec.ForProvider.CacheSecurityGroupNameRefs = mrsp.ResolvedReferences

// 	return nil
// }
