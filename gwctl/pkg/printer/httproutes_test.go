/*
Copyright 2023 The Kubernetes Authors.

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

package printer

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	testingclock "k8s.io/utils/clock/testing"

	apisv1beta1 "sigs.k8s.io/gateway-api/apis/applyconfiguration/apis/v1beta1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/gwctl/pkg/common"
	"sigs.k8s.io/gateway-api/gwctl/pkg/resourcediscovery"
	"sigs.k8s.io/gateway-api/gwctl/pkg/utils"
)

func TestHTTPRoutesPrinter_PrintTable(t *testing.T) {
	fakeClock := testingclock.NewFakeClock(time.Now())
	objects := []runtime.Object{
		&gatewayv1.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo-gatewayclass-1",
			},
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "example.net/gateway-controller",
				Description:    common.PtrTo("random"),
			},
		},
		&gatewayv1.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "demo-gatewayclass-2",
			},
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "example.net/gateway-controller",
				Description:    common.PtrTo("random"),
			},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ns1",
			},
			Status: corev1.NamespaceStatus{
				Phase: corev1.NamespaceActive,
			},
		},
		&corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ns2",
			},
			Status: corev1.NamespaceStatus{
				Phase: corev1.NamespaceActive,
			},
		},
		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "demo-gateway-1",
				Namespace: "default",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "demo-gatewayclass-1",
			},
		},
		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "demo-gateway-2",
				Namespace: "ns2",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "demo-gatewayclass-2",
			},
		},
		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "demo-gateway-200",
				Namespace: "default",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "demo-gatewayclass-1",
			},
		},
		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "demo-gateway-345",
				Namespace: "ns1",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "demo-gatewayclass-2",
			},
		},
		&gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-httproute-1",
				Namespace: "default",
				CreationTimestamp: metav1.Time{
					Time: fakeClock.Now().Add(-24 * time.Hour),
				},
			},
			Spec: gatewayv1.HTTPRouteSpec{
				Hostnames: []gatewayv1.Hostname{"example.com", "example2.com", "example3.com"},
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Kind:      common.PtrTo(gatewayv1.Kind("Gateway")),
							Group:     common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
							Namespace: common.PtrTo(gatewayv1.Namespace("ns2")),
							Name:      "demo-gateway-2",
						},
					},
				},
			},
		},
		&gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "qmn-httproute-100",
				Namespace: "default",
				CreationTimestamp: metav1.Time{
					Time: fakeClock.Now().Add(-11 * time.Hour),
				},
			},
			Spec: gatewayv1.HTTPRouteSpec{
				Hostnames: []gatewayv1.Hostname{"example.com"},
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Kind:  common.PtrTo(gatewayv1.Kind("Gateway")),
							Group: common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
							Name:  "demo-gateway-1",
						},
						{
							Kind:  common.PtrTo(gatewayv1.Kind("Gateway")),
							Group: common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
							Name:  "demo-gateway-200",
						},
					},
				},
			},
		},
		&gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar-route-21",
				Namespace: "ns1",
				CreationTimestamp: metav1.Time{
					Time: fakeClock.Now().Add(-9 * time.Hour),
				},
			},
			Spec: gatewayv1.HTTPRouteSpec{
				Hostnames: []gatewayv1.Hostname{"foo.com", "bar.com", "example.com", "example2.com", "example3.com", "example4.com", "example5.com"},
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Kind:      common.PtrTo(gatewayv1.Kind("Gateway")),
							Group:     common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
							Namespace: common.PtrTo(gatewayv1.Namespace("default")),
							Name:      "demo-gateway-200",
						},
					},
				},
			},
		},
		&gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bax-httproute-18777",
				Namespace: "ns2",
				CreationTimestamp: metav1.Time{
					Time: fakeClock.Now().Add(-5 * time.Minute),
				},
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Kind:      common.PtrTo(gatewayv1.Kind("Gateway")),
							Group:     common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
							Namespace: common.PtrTo(gatewayv1.Namespace("ns1")),
							Name:      "demo-gateway-345",
						},
					},
				},
			},
		},
		&apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "timeoutpolicies.bar.com",
				Labels: map[string]string{
					gatewayv1alpha2.PolicyLabelKey: "direct",
				},
			},
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Scope:    apiextensionsv1.ClusterScoped,
				Group:    "bar.com",
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{{Name: "v1"}},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "timeoutpolicies",
					Kind:   "TimeoutPolicy",
				},
			},
		},
		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "bar.com/v1",
				"kind":       "TimeoutPolicy",
				"metadata": map[string]interface{}{
					"name": "timeout-policy-httproute",
				},
				"spec": map[string]interface{}{
					"condition": "path=/def",
					"seconds":   int64(60),
					"targetRef": map[string]interface{}{
						"group": "gateway.networking.k8s.io",
						"kind":  "HTTPRoute",
						"name":  "foo-httproute-1",
					},
				},
			},
		},
	}

	k8sClients := common.MustClientsForTest(t, objects...)
	policyManager := utils.MustPolicyManagerForTest(t, k8sClients)
	buff := &bytes.Buffer{}
	discoverer := resourcediscovery.Discoverer{
		K8sClients:    k8sClients,
		PolicyManager: policyManager,
	}
	resourceModel, err := discoverer.DiscoverResourcesForHTTPRoute(resourcediscovery.Filter{})
	if err != nil {
		t.Fatalf("Failed to construct resourceModel: %v", err)
	}

	hp := &HTTPRoutesPrinter{
		Writer: buff,
		Clock:  fakeClock,
	}

	hp.PrintTable(resourceModel, false)

	got := buff.String()
	want := `
NAMESPACE  NAME                 HOSTNAMES                          PARENT REFS  AGE
default    foo-httproute-1      example.com,example2.com + 1 more  1            24h
default    qmn-httproute-100    example.com                        2            11h
ns1        bar-route-21         foo.com,bar.com + 5 more           1            9h
ns2        bax-httproute-18777  None                               1            5m
`
	if diff := cmp.Diff(common.YamlString(want), common.YamlString(got), common.YamlStringTransformer); diff != "" {
		t.Errorf("Unexpected diff\ngot=\n%v\nwant=\n%v\ndiff (-want +got)=\n%v", got, want, diff)
	}

	buff.Reset()
	hp2 := &HTTPRoutesPrinter{
		Writer: buff,
		Clock:  fakeClock,
	}
	hp2.PrintTable(resourceModel, true)

	got2 := buff.String()
	want2 := `
NAMESPACE  NAME                 HOSTNAMES                          PARENT REFS  AGE  POLICIES
default    foo-httproute-1      example.com,example2.com + 1 more  1            24h  1
default    qmn-httproute-100    example.com                        2            11h  0
ns1        bar-route-21         foo.com,bar.com + 5 more           1            9h   0
ns2        bax-httproute-18777  None                               1            5m   0
`
	if diff := cmp.Diff(common.YamlString(want2), common.YamlString(got2), common.YamlStringTransformer); diff != "" {
		t.Errorf("Unexpected diff\ngot=\n%v\nwant=\n%v\ndiff (-want +got)=\n%v", got2, want2, diff)
	}
}

func TestHTTPRoutesPrinter_PrintDescribeView(t *testing.T) {
	fakeClock := testingclock.NewFakeClock(time.Now())
	objects := []runtime.Object{
		&gatewayv1.GatewayClass{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "GatewayClass",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo-gatewayclass",
			},
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "example.net/gateway-controller",
				Description:    common.PtrTo("random"),
			},
		},
		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "foo.com/v1",
				"kind":       "HealthCheckPolicy",
				"metadata": map[string]interface{}{
					"name": "health-check-gatewayclass",
				},
				"spec": map[string]interface{}{
					"override": map[string]interface{}{
						"key1": "value-parent-1",
						"key3": "value-parent-3",
						"key5": "value-parent-5",
					},
					"default": map[string]interface{}{
						"key2": "value-parent-2",
						"key4": "value-parent-4",
					},
					"targetRef": map[string]interface{}{
						"group": "gateway.networking.k8s.io",
						"kind":  "GatewayClass",
						"name":  "foo-gatewayclass",
					},
				},
			},
		},

		&gatewayv1.Gateway{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "Gateway",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-gateway",
				Namespace: "default",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "foo-gatewayclass",
			},
		},
		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "foo.com/v1",
				"kind":       "HealthCheckPolicy",
				"metadata": map[string]interface{}{
					"name":      "health-check-gateway",
					"namespace": "default",
				},
				"spec": map[string]interface{}{
					"override": map[string]interface{}{
						"key1": "value-child-1",
					},
					"default": map[string]interface{}{
						"key2": "value-child-2",
						"key5": "value-child-5",
					},
					"targetRef": map[string]interface{}{
						"group":     "gateway.networking.k8s.io",
						"kind":      "Gateway",
						"name":      "foo-gateway",
						"namespace": "default",
					},
				},
			},
		},

		&gatewayv1.HTTPRoute{
			TypeMeta: metav1.TypeMeta{
				APIVersion: gatewayv1.GroupVersion.String(),
				Kind:       "HTTPRoute",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo-httproute",
				Namespace: "default",
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{{
						Kind:  common.PtrTo(gatewayv1.Kind("Gateway")),
						Group: common.PtrTo(gatewayv1.Group("gateway.networking.k8s.io")),
						Name:  "foo-gateway",
					}},
				},
			},
		},
		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "bar.com/v1",
				"kind":       "TimeoutPolicy",
				"metadata": map[string]interface{}{
					"name":      "timeout-policy-httproute",
					"namespace": "default",
				},
				"spec": map[string]interface{}{
					"condition": "path=/def",
					"seconds":   int64(60),
					"targetRef": map[string]interface{}{
						"group": "gateway.networking.k8s.io",
						"kind":  "HTTPRoute",
						"name":  "foo-httproute",
					},
				},
			},
		},

		&apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "healthcheckpolicies.foo.com",
				Labels: map[string]string{
					gatewayv1alpha2.PolicyLabelKey: "inherited",
				},
			},
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Scope:    apiextensionsv1.ClusterScoped,
				Group:    "foo.com",
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{{Name: "v1"}},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "healthcheckpolicies",
					Kind:   "HealthCheckPolicy",
				},
			},
		},

		&apiextensionsv1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name: "timeoutpolicies.bar.com",
				Labels: map[string]string{
					gatewayv1alpha2.PolicyLabelKey: "direct",
				},
			},
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Scope:    apiextensionsv1.ClusterScoped,
				Group:    "bar.com",
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{{Name: "v1"}},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Plural: "timeoutpolicies",
					Kind:   "TimeoutPolicy",
				},
			},
		},
		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "bar.com/v1",
				"kind":       "TimeoutPolicy",
				"metadata": map[string]interface{}{
					"name": "timeout-policy-namespace",
				},
				"spec": map[string]interface{}{
					"condition": "path=/abc",
					"seconds":   int64(30),
					"targetRef": map[string]interface{}{
						"kind": "Namespace",
						"name": "default",
					},
				},
			},
		},
	}

	k8sClients := common.MustClientsForTest(t, objects...)
	policyManager := utils.MustPolicyManagerForTest(t, k8sClients)
	buff := &bytes.Buffer{}
	discoverer := resourcediscovery.Discoverer{
		K8sClients:    k8sClients,
		PolicyManager: policyManager,
	}
	resourceModel, err := discoverer.DiscoverResourcesForHTTPRoute(resourcediscovery.Filter{})
	if err != nil {
		t.Fatalf("Failed to construct resourceModel: %v", err)
	}

	hp := &HTTPRoutesPrinter{
		Writer:       buff,
		Clock:        fakeClock,
		EventFetcher: discoverer,
	}
	hp.PrintDescribeView(resourceModel)

	got := buff.String()
	want := `
Name: foo-httproute
Namespace: default
Label: null
Annotations: null
APIVersion: gateway.networking.k8s.io/v1
Kind: HTTPRoute
Metadata:
  creationTimestamp: null
  resourceVersion: "999"
Spec:
  parentRefs:
  - group: gateway.networking.k8s.io
    kind: Gateway
    name: foo-gateway
Status:
  parents: null
DirectlyAttachedPolicies:
  Type                   Name
  ----                   ----
  TimeoutPolicy.bar.com  default/timeout-policy-httproute
InheritedPolicies:
  Type                       Name                          Target Kind   Target Name
  ----                       ----                          -----------   -----------
  HealthCheckPolicy.foo.com  health-check-gatewayclass     GatewayClass  foo-gatewayclass
  HealthCheckPolicy.foo.com  default/health-check-gateway  Gateway       default/foo-gateway
EffectivePolicies:
  default/foo-gateway:
    HealthCheckPolicy.foo.com:
      key1: value-parent-1
      key2: value-child-2
      key3: value-parent-3
      key4: value-parent-4
      key5: value-parent-5
    TimeoutPolicy.bar.com:
      condition: path=/def
      seconds: 60
Events: <none>
`
	if diff := cmp.Diff(common.YamlString(want), common.YamlString(got), common.YamlStringTransformer); diff != "" {
		t.Errorf("Unexpected diff\ngot=\n%v\nwant=\n%v\ndiff (-want +got)=\n%v", got, want, diff)
	}
}

// TestHTTPRoutesPrinter_PrintJsonYaml tests the correctness of JSON/YAML output associated with -o json/yaml of `get` subcommand
func TestHTTPRoutesPrinter_PrintJsonYaml(t *testing.T) {
	fakeClock := testingclock.NewFakeClock(time.Now())
	creationTime := fakeClock.Now().Add(-24 * time.Hour).UTC() // UTC being necessary for consistently handling the time while marshaling/unmarshaling its JSON

	hrName, hrNamespace := "httproute-1", "default"
	hrConfig := apisv1beta1.HTTPRoute(hrName, hrNamespace)

	objects := []runtime.Object{
		&gatewayv1.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gatewayclass-1",
			},
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "example.net/gateway-controller",
				Description:    common.PtrTo("random"),
			},
		},

		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gateway-1",
				Namespace: "default",
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: "gatewayclass-1",
			},
		},
		&gatewayv1.HTTPRoute{
			TypeMeta: metav1.TypeMeta{
				APIVersion: *hrConfig.APIVersion,
				Kind:       *hrConfig.Kind,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      hrName,
				Namespace: hrNamespace,
				CreationTimestamp: metav1.Time{
					Time: creationTime,
				},
				Labels: map[string]string{"app": "foo", "env": "internal"},
			},
			Spec: gatewayv1.HTTPRouteSpec{
				Hostnames: []gatewayv1.Hostname{"example.com"},
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{
						{
							Name: "gateway-1",
						},
					},
				},
			},
		},
	}

	k8sClients := common.MustClientsForTest(t, objects...)
	policyManager := utils.MustPolicyManagerForTest(t, k8sClients)
	buff := &bytes.Buffer{}
	discoverer := resourcediscovery.Discoverer{
		K8sClients:    k8sClients,
		PolicyManager: policyManager,
	}

	resourceModel, err := discoverer.DiscoverResourcesForHTTPRoute(resourcediscovery.Filter{})
	if err != nil {
		t.Fatalf("Failed to discover resources: %v", err)
	}

	hp := &HTTPRoutesPrinter{
		Writer: buff,
		Clock:  fakeClock,
	}
	Print(hp, resourceModel, utils.OutputFormatJSON)

	gotJSON := common.JSONString(buff.String())
	wantJSON := common.JSONString(fmt.Sprintf(`
        {
          "apiVersion": "v1",
          "items": [
            {
              "apiVersion": "gateway.networking.k8s.io/v1beta1",
              "kind": "HTTPRoute",
              "metadata": {
                "creationTimestamp": "%s",
                "labels": {
                  "app": "foo",
                  "env": "internal"
                },
                "name": "httproute-1",
                "namespace": "default",
                "resourceVersion": "999"
              },
              "spec": {
                "hostnames": [
                  "example.com"
                ],
                "parentRefs": [
                  {
                    "name": "gateway-1"
                  }
                ]
              },
              "status": {
                "parents": null
              }
            }
          ],
          "kind": "List"
        }`, creationTime.Format(time.RFC3339)))
	diff, err := wantJSON.CmpDiff(gotJSON)
	if err != nil {
		t.Fatalf("Failed to compare the json diffs: %v", diff)
	}
	if diff != "" {
		t.Errorf("Unexpected diff\ngot=\n%v\nwant=\n%v\ndiff (-want +got)=\n%v", gotJSON, wantJSON, diff)
	}

	hp.Writer = &bytes.Buffer{}

	Print(hp, resourceModel, utils.OutputFormatYAML)

	gotYaml := common.YamlString(hp.Writer.(*bytes.Buffer).String())
	wantYaml := common.YamlString(fmt.Sprintf(`
apiVersion: v1
items:
- apiVersion: gateway.networking.k8s.io/v1beta1
  kind: HTTPRoute
  metadata:
    creationTimestamp: "%s"
    labels:
      app: foo
      env: internal
    name: httproute-1
    namespace: default
    resourceVersion: "999"
  spec:
    hostnames:
    - example.com
    parentRefs:
    - name: gateway-1
  status:
    parents: null
kind: List`, creationTime.Format(time.RFC3339)))
	if diff := cmp.Diff(wantYaml, gotYaml, common.YamlStringTransformer); diff != "" {
		t.Errorf("Unexpected diff\ngot=\n%v\nwant=\n%v\ndiff (-want +got)=\n%v", gotYaml, wantYaml, diff)
	}
}
