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
package database

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsrds "github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/test"

	"github.com/crossplane/provider-aws/apis/database/v1beta1"
	awsv1alpha3 "github.com/crossplane/provider-aws/apis/v1alpha3"
	awsclients "github.com/crossplane/provider-aws/pkg/clients"
	"github.com/crossplane/provider-aws/pkg/clients/rds"
	"github.com/crossplane/provider-aws/pkg/clients/rds/fake"
)

const (
	providerName    = "aws-creds"
	secretNamespace = "crossplane-system"
	testRegion      = "us-east-1"

	connectionSecretName = "my-little-secret"
	secretKey            = "credentials"
	credData             = "confidential!"
)

var (
	masterUsername = "root"
	engineVersion  = "5.6"

	replaceMe = "replace-me!"
	errBoom   = errors.New("boom")
)

type args struct {
	rds  rds.Client
	kube client.Client
	cr   *v1beta1.RDSInstance
}

type rdsModifier func(*v1beta1.RDSInstance)

func withMasterUsername(s *string) rdsModifier {
	return func(r *v1beta1.RDSInstance) { r.Spec.ForProvider.MasterUsername = s }
}

func withConditions(c ...runtimev1alpha1.Condition) rdsModifier {
	return func(r *v1beta1.RDSInstance) { r.Status.ConditionedStatus.Conditions = c }
}

func withBindingPhase(p runtimev1alpha1.BindingPhase) rdsModifier {
	return func(r *v1beta1.RDSInstance) { r.Status.SetBindingPhase(p) }
}

func withEngineVersion(s *string) rdsModifier {
	return func(r *v1beta1.RDSInstance) { r.Spec.ForProvider.EngineVersion = s }
}

func withTags(tagMaps ...map[string]string) rdsModifier {
	var tagList []v1beta1.Tag
	for _, tagMap := range tagMaps {
		for k, v := range tagMap {
			tagList = append(tagList, v1beta1.Tag{Key: k, Value: v})
		}
	}
	return func(r *v1beta1.RDSInstance) { r.Spec.ForProvider.Tags = tagList }
}

func withDBInstanceStatus(s string) rdsModifier {
	return func(r *v1beta1.RDSInstance) { r.Status.AtProvider.DBInstanceStatus = s }
}

func instance(m ...rdsModifier) *v1beta1.RDSInstance {
	cr := &v1beta1.RDSInstance{
		Spec: v1beta1.RDSInstanceSpec{
			ResourceSpec: runtimev1alpha1.ResourceSpec{
				ProviderReference: &corev1.ObjectReference{Name: providerName},
			},
		},
	}
	for _, f := range m {
		f(cr)
	}
	return cr
}

var _ managed.ExternalClient = &external{}
var _ managed.ExternalConnecter = &connector{}

func TestConnect(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      connectionSecretName,
			Namespace: secretNamespace,
		},
		Data: map[string][]byte{
			secretKey: []byte(credData),
		},
	}

	providerSA := func(saVal bool) awsv1alpha3.Provider {
		return awsv1alpha3.Provider{
			Spec: awsv1alpha3.ProviderSpec{
				Region:            testRegion,
				UseServiceAccount: &saVal,
				ProviderSpec: runtimev1alpha1.ProviderSpec{
					CredentialsSecretRef: &runtimev1alpha1.SecretKeySelector{
						SecretReference: runtimev1alpha1.SecretReference{
							Namespace: secretNamespace,
							Name:      connectionSecretName,
						},
						Key: secretKey,
					},
				},
			},
		}
	}
	type args struct {
		kube        client.Client
		newClientFn func(ctx context.Context, credentials []byte, region string, auth awsclients.AuthMethod) (rds.Client, error)
		cr          *v1beta1.RDSInstance
	}
	type want struct {
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				kube: &test.MockClient{
					MockGet: func(_ context.Context, key client.ObjectKey, obj runtime.Object) error {
						switch key {
						case client.ObjectKey{Name: providerName}:
							p := providerSA(false)
							p.DeepCopyInto(obj.(*awsv1alpha3.Provider))
							return nil
						case client.ObjectKey{Namespace: secretNamespace, Name: connectionSecretName}:
							secret.DeepCopyInto(obj.(*corev1.Secret))
							return nil
						}
						return errBoom
					},
				},
				newClientFn: func(_ context.Context, credentials []byte, region string, _ awsclients.AuthMethod) (i rds.Client, e error) {
					if diff := cmp.Diff(credData, string(credentials)); diff != "" {
						t.Errorf("r: -want, +got:\n%s", diff)
					}
					if diff := cmp.Diff(testRegion, region); diff != "" {
						t.Errorf("r: -want, +got:\n%s", diff)
					}
					return nil, nil
				},
				cr: instance(),
			},
		},
		"SuccessfulUseServiceAccount": {
			args: args{
				kube: &test.MockClient{
					MockGet: func(_ context.Context, key client.ObjectKey, obj runtime.Object) error {
						if key == (client.ObjectKey{Name: providerName}) {
							p := providerSA(true)
							p.DeepCopyInto(obj.(*awsv1alpha3.Provider))
							return nil
						}
						return errBoom
					},
				},
				newClientFn: func(_ context.Context, credentials []byte, region string, _ awsclients.AuthMethod) (i rds.Client, e error) {
					if diff := cmp.Diff("", string(credentials)); diff != "" {
						t.Errorf("r: -want, +got:\n%s", diff)
					}
					if diff := cmp.Diff(testRegion, region); diff != "" {
						t.Errorf("r: -want, +got:\n%s", diff)
					}
					return nil, nil
				},
				cr: instance(),
			},
		},
		"ProviderGetFailed": {
			args: args{
				kube: &test.MockClient{
					MockGet: func(_ context.Context, key client.ObjectKey, obj runtime.Object) error {
						return errBoom
					},
				},
				cr: instance(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetProvider),
			},
		},
		"SecretGetFailed": {
			args: args{
				kube: &test.MockClient{
					MockGet: func(_ context.Context, key client.ObjectKey, obj runtime.Object) error {
						switch key {
						case client.ObjectKey{Name: providerName}:
							p := providerSA(false)
							p.DeepCopyInto(obj.(*awsv1alpha3.Provider))
							return nil
						case client.ObjectKey{Namespace: secretNamespace, Name: connectionSecretName}:
							return errBoom
						default:
							return nil
						}
					},
				},
				cr: instance(),
			},
			want: want{
				err: errors.Wrap(errBoom, errGetProviderSecret),
			},
		},
		"SecretGetFailedNil": {
			args: args{
				kube: &test.MockClient{
					MockGet: func(_ context.Context, key client.ObjectKey, obj runtime.Object) error {
						switch key {
						case client.ObjectKey{Name: providerName}:
							p := providerSA(false)
							p.SetCredentialsSecretReference(nil)
							p.DeepCopyInto(obj.(*awsv1alpha3.Provider))
							return nil
						case client.ObjectKey{Namespace: secretNamespace, Name: connectionSecretName}:
							return errBoom
						default:
							return nil
						}
					},
				},
				cr: instance(),
			},
			want: want{
				err: errors.New(errGetProviderSecret),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := &connector{kube: tc.kube, newClientFn: tc.newClientFn}
			_, err := c.Connect(context.Background(), tc.args.cr)
			if diff := cmp.Diff(tc.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestObserve(t *testing.T) {
	type want struct {
		cr     *v1beta1.RDSInstance
		result managed.ExternalObservation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"SuccessfulAvailable": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{
									{
										DBInstanceStatus: aws.String(string(v1beta1.RDSInstanceStateAvailable)),
									},
								},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(
					withConditions(runtimev1alpha1.Available()),
					withBindingPhase(runtimev1alpha1.BindingPhaseUnbound),
					withDBInstanceStatus(string(v1beta1.RDSInstanceStateAvailable))),
				result: managed.ExternalObservation{
					ResourceExists:    true,
					ResourceUpToDate:  true,
					ConnectionDetails: rds.GetConnectionDetails(v1beta1.RDSInstance{}),
				},
			},
		},
		"DeletingState": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{
									{
										DBInstanceStatus: aws.String(string(v1beta1.RDSInstanceStateDeleting)),
									},
								},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(
					withConditions(runtimev1alpha1.Deleting()),
					withDBInstanceStatus(string(v1beta1.RDSInstanceStateDeleting))),
				result: managed.ExternalObservation{
					ResourceExists:    true,
					ResourceUpToDate:  true,
					ConnectionDetails: rds.GetConnectionDetails(v1beta1.RDSInstance{}),
				},
			},
		},
		"FailedState": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{
									{
										DBInstanceStatus: aws.String(string(v1beta1.RDSInstanceStateFailed)),
									},
								},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(
					withConditions(runtimev1alpha1.Unavailable()),
					withDBInstanceStatus(string(v1beta1.RDSInstanceStateFailed))),
				result: managed.ExternalObservation{
					ResourceExists:    true,
					ResourceUpToDate:  true,
					ConnectionDetails: rds.GetConnectionDetails(v1beta1.RDSInstance{}),
				},
			},
		},
		"FailedDescribeRequest": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr:  instance(),
				err: errors.Wrap(errBoom, errDescribeFailed),
			},
		},
		"NotFound": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errors.New(awsrds.ErrCodeDBInstanceNotFoundFault)},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(),
			},
		},
		"LateInitSuccess": {
			args: args{
				kube: &test.MockClient{
					MockUpdate: test.NewMockUpdateFn(nil),
				},
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{
									{
										EngineVersion:    aws.String(engineVersion),
										DBInstanceStatus: aws.String(string(v1beta1.RDSInstanceStateCreating)),
									},
								},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(
					withEngineVersion(&engineVersion),
					withDBInstanceStatus(string(v1beta1.RDSInstanceStateCreating)),
					withConditions(runtimev1alpha1.Creating()),
				),
				result: managed.ExternalObservation{
					ResourceExists:    true,
					ResourceUpToDate:  true,
					ConnectionDetails: rds.GetConnectionDetails(v1beta1.RDSInstance{}),
				},
			},
		},
		"LateInitFailedKubeUpdate": {
			args: args{
				kube: &test.MockClient{
					MockUpdate: test.NewMockUpdateFn(errBoom),
				},
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{
									{
										EngineVersion:    aws.String(engineVersion),
										DBInstanceStatus: aws.String(string(v1beta1.RDSInstanceStateCreating)),
									},
								},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(
					withEngineVersion(&engineVersion),
				),
				err: errors.Wrap(errBoom, errKubeUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.rds}
			o, err := e.Observe(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type want struct {
		cr     *v1beta1.RDSInstance
		result managed.ExternalCreation
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				rds: &fake.MockRDSClient{
					MockCreate: func(input *awsrds.CreateDBInstanceInput) awsrds.CreateDBInstanceRequest {
						return awsrds.CreateDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.CreateDBInstanceOutput{}},
						}
					},
				},
				cr: instance(withMasterUsername(&masterUsername)),
			},
			want: want{
				cr: instance(
					withMasterUsername(&masterUsername),
					withConditions(runtimev1alpha1.Creating())),
				result: managed.ExternalCreation{
					ConnectionDetails: managed.ConnectionDetails{
						runtimev1alpha1.ResourceCredentialsSecretUserKey:     []byte(masterUsername),
						runtimev1alpha1.ResourceCredentialsSecretPasswordKey: []byte(replaceMe),
					},
				},
			},
		},
		"SuccessfulNoNeedForCreate": {
			args: args{
				cr: instance(withDBInstanceStatus(v1beta1.RDSInstanceStateCreating)),
			},
			want: want{
				cr: instance(
					withDBInstanceStatus(v1beta1.RDSInstanceStateCreating),
					withConditions(runtimev1alpha1.Creating())),
			},
		},
		"SuccessfulNoUsername": {
			args: args{
				rds: &fake.MockRDSClient{
					MockCreate: func(input *awsrds.CreateDBInstanceInput) awsrds.CreateDBInstanceRequest {
						return awsrds.CreateDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.CreateDBInstanceOutput{}},
						}
					},
				},
				cr: instance(withMasterUsername(nil)),
			},
			want: want{
				cr: instance(
					withMasterUsername(nil),
					withConditions(runtimev1alpha1.Creating())),
				result: managed.ExternalCreation{
					ConnectionDetails: managed.ConnectionDetails{
						runtimev1alpha1.ResourceCredentialsSecretPasswordKey: []byte(replaceMe),
					},
				},
			},
		},
		"FailedRequest": {
			args: args{
				rds: &fake.MockRDSClient{
					MockCreate: func(input *awsrds.CreateDBInstanceInput) awsrds.CreateDBInstanceRequest {
						return awsrds.CreateDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr:  instance(withConditions(runtimev1alpha1.Creating())),
				err: errors.Wrap(errBoom, errCreateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.rds}
			o, err := e.Create(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if string(tc.want.result.ConnectionDetails[runtimev1alpha1.ResourceCredentialsSecretPasswordKey]) == replaceMe {
				tc.want.result.ConnectionDetails[runtimev1alpha1.ResourceCredentialsSecretPasswordKey] =
					o.ConnectionDetails[runtimev1alpha1.ResourceCredentialsSecretPasswordKey]
			}
			if diff := cmp.Diff(tc.want.result, o); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type want struct {
		cr     *v1beta1.RDSInstance
		result managed.ExternalUpdate
		err    error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				rds: &fake.MockRDSClient{
					MockModify: func(input *awsrds.ModifyDBInstanceInput) awsrds.ModifyDBInstanceRequest {
						return awsrds.ModifyDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.ModifyDBInstanceOutput{}},
						}
					},
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{{}},
							}},
						}
					},
					MockAddTags: func(input *awsrds.AddTagsToResourceInput) awsrds.AddTagsToResourceRequest {
						return awsrds.AddTagsToResourceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.AddTagsToResourceOutput{}},
						}
					},
				},
				cr: instance(withTags(map[string]string{"foo": "bar"})),
			},
			want: want{
				cr: instance(withTags(map[string]string{"foo": "bar"})),
			},
		},
		"AlreadyModifying": {
			args: args{
				cr: instance(withDBInstanceStatus(v1beta1.RDSInstanceStateModifying)),
			},
			want: want{
				cr: instance(withDBInstanceStatus(v1beta1.RDSInstanceStateModifying)),
			},
		},
		"FailedDescribe": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr:  instance(),
				err: errors.Wrap(errBoom, errDescribeFailed),
			},
		},
		"FailedModify": {
			args: args{
				rds: &fake.MockRDSClient{
					MockModify: func(input *awsrds.ModifyDBInstanceInput) awsrds.ModifyDBInstanceRequest {
						return awsrds.ModifyDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{{}},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr:  instance(),
				err: errors.Wrap(errBoom, errModifyFailed),
			},
		},
		"FailedAddTags": {
			args: args{
				rds: &fake.MockRDSClient{
					MockModify: func(input *awsrds.ModifyDBInstanceInput) awsrds.ModifyDBInstanceRequest {
						return awsrds.ModifyDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.ModifyDBInstanceOutput{}},
						}
					},
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{{}},
							}},
						}
					},
					MockAddTags: func(input *awsrds.AddTagsToResourceInput) awsrds.AddTagsToResourceRequest {
						return awsrds.AddTagsToResourceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
				},
				cr: instance(withTags(map[string]string{"foo": "bar"})),
			},
			want: want{
				cr:  instance(withTags(map[string]string{"foo": "bar"})),
				err: errors.Wrap(errBoom, errAddTagsFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.rds}
			u, err := e.Update(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.result, u); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type want struct {
		cr  *v1beta1.RDSInstance
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDelete: func(input *awsrds.DeleteDBInstanceInput) awsrds.DeleteDBInstanceRequest {
						return awsrds.DeleteDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DeleteDBInstanceOutput{}},
						}
					},
					MockModify: func(input *awsrds.ModifyDBInstanceInput) awsrds.ModifyDBInstanceRequest {
						return awsrds.ModifyDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.ModifyDBInstanceOutput{}},
						}
					},
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{{}},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(withConditions(runtimev1alpha1.Deleting())),
			},
		},
		"AlreadyDeleting": {
			args: args{
				cr: instance(withDBInstanceStatus(v1beta1.RDSInstanceStateDeleting)),
			},
			want: want{
				cr: instance(withDBInstanceStatus(v1beta1.RDSInstanceStateDeleting),
					withConditions(runtimev1alpha1.Deleting())),
			},
		},
		"AlreadyDeleted": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errors.New(awsrds.ErrCodeDBInstanceNotFoundFault)},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr: instance(withConditions(runtimev1alpha1.Deleting())),
			},
		},
		"Failed": {
			args: args{
				rds: &fake.MockRDSClient{
					MockDelete: func(input *awsrds.DeleteDBInstanceInput) awsrds.DeleteDBInstanceRequest {
						return awsrds.DeleteDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Error: errBoom},
						}
					},
					MockModify: func(input *awsrds.ModifyDBInstanceInput) awsrds.ModifyDBInstanceRequest {
						return awsrds.ModifyDBInstanceRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.ModifyDBInstanceOutput{}},
						}
					},
					MockDescribe: func(input *awsrds.DescribeDBInstancesInput) awsrds.DescribeDBInstancesRequest {
						return awsrds.DescribeDBInstancesRequest{
							Request: &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Data: &awsrds.DescribeDBInstancesOutput{
								DBInstances: []awsrds.DBInstance{{}},
							}},
						}
					},
				},
				cr: instance(),
			},
			want: want{
				cr:  instance(withConditions(runtimev1alpha1.Deleting())),
				err: errors.Wrap(errBoom, errDeleteFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &external{kube: tc.kube, client: tc.rds}
			err := e.Delete(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, test.EquateConditions()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	type want struct {
		cr  *v1beta1.RDSInstance
		err error
	}

	cases := map[string]struct {
		args
		want
	}{
		"Successful": {
			args: args{
				cr:   instance(withTags(map[string]string{"foo": "bar"})),
				kube: &test.MockClient{MockUpdate: test.NewMockUpdateFn(nil)},
			},
			want: want{
				cr: instance(withTags(resource.GetExternalTags(instance()), map[string]string{"foo": "bar"})),
			},
		},
		"UpdateFailed": {
			args: args{
				cr:   instance(),
				kube: &test.MockClient{MockUpdate: test.NewMockUpdateFn(errBoom)},
			},
			want: want{
				err: errors.Wrap(errBoom, errKubeUpdateFailed),
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := &tagger{kube: tc.kube}
			err := e.Initialize(context.Background(), tc.args.cr)

			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.cr, tc.args.cr, cmpopts.SortSlices(func(a, b v1beta1.Tag) bool { return a.Key > b.Key })); err == nil && diff != "" {
				t.Errorf("r: -want, +got:\n%s", diff)
			}
		})
	}
}
