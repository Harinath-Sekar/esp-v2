// Copyright 2019 Google Cloud Platform Proxy Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"fmt"
	"testing"
	"time"

	"cloudesf.googlesource.com/gcpproxy/src/go/util"
	"cloudesf.googlesource.com/gcpproxy/tests/endpoints/echo/client"
	"cloudesf.googlesource.com/gcpproxy/tests/env"
	"cloudesf.googlesource.com/gcpproxy/tests/utils"
	"google.golang.org/genproto/googleapis/api/annotations"

	comp "cloudesf.googlesource.com/gcpproxy/tests/env/components"
	conf "google.golang.org/genproto/googleapis/api/serviceconfig"
)

var testDynamicRoutingArgs = []string{
	"--service=test-echo",
	"--version=test-config-id",
	"--backend_protocol=http1",
	"--rollout_strategy=fixed",
	"--enable_backend_routing",
}

func NewDynamicRoutingTestEnv(port uint16) *env.TestEnv {
	s := env.NewTestEnv(port, "echo", nil)
	s.EnableDynamicRoutingBackend()
	s.AppendHttpRules([]*annotations.HttpRule{
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetById",
			Pattern: &annotations.HttpRule_Get{
				Get: "/pet/{pet_id}/num/{number}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPet",
			Pattern: &annotations.HttpRule_Get{
				Get: "/searchpet",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchDogsWithSlash",
			Pattern: &annotations.HttpRule_Get{
				Get: "/searchdog",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRoot",
			Pattern: &annotations.HttpRule_Get{
				Get: "/searchroot",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRootWithSlash",
			Pattern: &annotations.HttpRule_Get{
				Get: "/searchrootwithslash",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListPets",
			Pattern: &annotations.HttpRule_Get{
				Get: "/pets/{category}/year/{no}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListShelves",
			Pattern: &annotations.HttpRule_Get{
				Get: "/shelves",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookInfoWithSnakeCase",
			Pattern: &annotations.HttpRule_Get{
				Get: "/shelves/{s_h_e_l_f}/books/info/{b_o_o_k}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookIdWithSnakeCase",
			Pattern: &annotations.HttpRule_Get{
				Get: "/shelves/{s_h_e_l_f.i_d}/books/id/{b_o_o_k.id}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPetWithServiceControlVerification",
			Pattern: &annotations.HttpRule_Post{
				Post: "/sc/searchpet",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetByIdWithServiceControlVerification",
			Pattern: &annotations.HttpRule_Post{
				Post: "/sc/pet/{pet_id}/num/{number}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.BearertokenConstantAddress",
			Pattern: &annotations.HttpRule_Get{
				Get: "/bearertoken/constant/{foo}",
			},
		},
		{
			Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.BearertokenAppendAddress",
			Pattern: &annotations.HttpRule_Get{
				Get: "/bearertoken/append",
			},
		},
	})
	s.AppendUsageRules(
		[]*conf.UsageRule{
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetById",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPet",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchDogsWithSlash",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRoot",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRootWithSlash",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListPets",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListShelves",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookInfoWithSnakeCase",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookIdWithSnakeCase",
				AllowUnregisteredCalls: true,
			},
			{
				Selector:               "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.BearertokenConstantAddress",
				AllowUnregisteredCalls: true,
			},
		})

	s.AppendBackendRules(
		[]*conf.BackendRule{
			{
				Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetById",
				Address:         "https://localhost:-1/dynamicrouting/getpetbyid",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/getpetbyid",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPet",
				Address:         "https://localhost:-1/dynamicrouting/searchpet",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/searchpet",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchDogsWithSlash",
				Address:         "https://localhost:-1/dynamicrouting/searchdogs/",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/searchpet",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRoot",
				Address:         "https://localhost:-1",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/searchroot",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.AppendToRootWithSlash",
				Address:         "https://localhost:-1/",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/searchrootwithslash",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListPets",
				Address:         "https://localhost:-1/dynamicrouting/listpet",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/listpet",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.ListShelves",
				Address:         "https://localhost:-1/dynamicrouting/shelves",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/shelves",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookInfoWithSnakeCase",
				Address:         "https://localhost:-1/dynamicrouting/bookinfo",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/bookinfo",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetBookIdWithSnakeCase",
				Address:         "https://localhost:-1/dynamicrouting/bookid",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting/bookid",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPetWithServiceControlVerification",
				Address:         "https://localhost:-1/dynamicrouting/",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetByIdWithServiceControlVerification",
				Address:         "https://localhost:-1/dynamicrouting/",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/dynamicrouting",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.BearertokenConstantAddress",
				Address:         "https://localhost:-1/bearertoken/constant",
				PathTranslation: conf.BackendRule_CONSTANT_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/bearertoken/constant",
				},
			},
			{
				Selector:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.BearertokenAppendAddress",
				Address:         "https://localhost:-1",
				PathTranslation: conf.BackendRule_APPEND_PATH_TO_ADDRESS,
				Authentication: &conf.BackendRule_JwtAudience{
					JwtAudience: "https://localhost/bearertoken/append",
				},
			},
		})

	return s
}

func TestDynamicRouting(t *testing.T) {
	s := NewDynamicRoutingTestEnv(comp.TestDynamicRouting)
	if err := s.Setup(testDynamicRoutingArgs); err != nil {
		t.Fatalf("fail to setup test env, %v", err)
	}
	defer s.TearDown()

	testData := []struct {
		desc          string
		path          string
		method        string
		message       string
		wantResp      string
		httpCallError error
	}{
		{
			desc:     "Succeed, no path translation (no re-routing needed)",
			path:     "/echo?key=api-key",
			method:   "POST",
			message:  "hello",
			wantResp: `{"message":"hello"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct",
			path:     "/pet/123/num/987",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/getpetbyid?pet_id=123&number=987"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct, original URL has query parameters, original query parameters should appear first and query parameters converted from path parameters appear later",
			path:     "/pet/31/num/565?lang=US&zone=us-west1",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/getpetbyid?lang=US&zone=us-west1&pet_id=31&number=565"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct, original URL has query parameters, original query parameters should appear first and query parameters converted from path parameters appear later",
			path:     "/pet/31/num/565?lang=US&zone=us-west1",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/getpetbyid?lang=US&zone=us-west1&pet_id=31&number=565"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation with snake case is correct",
			path:     "/shelves/123/books/info/987",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/bookinfo?SHELF=123&BOOK=987"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation with snake case is correct, supports {foo.bar} style path, if corresponding jsonName not found, origin snake case path is used.",
			path:     "/shelves/221/books/id/2019",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/bookid?SHELF.i_d=221&BOOK.id=2019"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct for cases that does not have path parameter",
			path:     "/shelves",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/shelves"}`,
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct for cases that does not have path parameter but has query parameter",
			path:     "/shelves?q=story",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/shelves?q=story"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, appends original URL to backend address (https://domain/base/path)",
			path:     "/searchpet",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/searchpet/searchpet"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation with query parameter is correct, appends original URL to backend address (https://domain/base/path)",
			path:     "/searchpet?timezone=PST&lang=US",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/searchpet/searchpet?timezone=PST&lang=US"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, appends original URL to backend address that ends with slash (https://domain/base/path/)",
			path:     "/searchdog",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/searchdogs/searchdog"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation with query parameter is correct, appends original URL to backend address that ends with slash (https://domain/base/path/)",
			path:     "/searchdog?timezone=UTC",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/searchdogs/searchdog?timezone=UTC"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, original URL has path parameters",
			path:     "/pets/cat/year/2018",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/listpet/pets/cat/year/2018"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, original URL has path parameters and query parameters",
			path:     "/pets/dog/year/2019?lang=US&zone=us-west1",
			method:   "GET",
			wantResp: `{"RequestURI":"/dynamicrouting/listpet/pets/dog/year/2019?lang=US&zone=us-west1"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, backend address is root path with slash (https://domain/)",
			path:     "/searchrootwithslash",
			method:   "GET",
			wantResp: `{"RequestURI":"/searchrootwithslash"}`,
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation with query parameter is correct, backend address is root path with slash (https://domain/)",
			path:     "/searchroot?zone=us-central1&lang=en",
			method:   "GET",
			wantResp: `{"RequestURI":"/searchroot?zone=us-central1&lang=en"}`,
		},
		{
			desc:          "Fail, there is not backend rule specified for this path",
			path:          "/searchdogs",
			method:        "GET",
			httpCallError: fmt.Errorf("http response status is not 200 OK: 404 Not Found"),
		},
	}
	for _, tc := range testData {
		url := fmt.Sprintf("http://localhost:%v%v", s.Ports().ListenerPort, tc.path)
		var gotResp []byte
		var err error
		if tc.method == "GET" {
			gotResp, err = client.DoGet(url)

		} else if tc.method == "POST" {
			gotResp, err = client.DoPost(url, tc.message)
		} else {
			t.Fatalf("unknown HTTP method (%v) to call", tc.method)
		}

		if tc.httpCallError == nil {
			if err != nil {
				t.Fatal(err)
			}
		} else {
			if tc.httpCallError.Error() != err.Error() {
				t.Errorf("expected Http call error: %v, got: %v", tc.httpCallError, err)
			}
			continue
		}
		gotRespStr := string(gotResp)
		if !utils.JsonEqual(gotRespStr, tc.wantResp) {
			t.Errorf("response want: %s, got: %s", tc.wantResp, gotRespStr)
		}
	}
}

func TestBackendAuth(t *testing.T) {
	s := NewDynamicRoutingTestEnv(comp.TestBackendAuth)
	s.OverrideMockMetadata(
		map[string]string{
			util.IdentityTokenSuffix + "?audience=https://localhost/bearertoken/constant&format=standard": "ya29.constant",
			util.IdentityTokenSuffix + "?audience=https://localhost/bearertoken/append&format=standard":   "ya29.append",
		})

	if err := s.Setup(testDynamicRoutingArgs); err != nil {
		t.Fatalf("fail to setup test env, %v", err)
	}
	defer s.TearDown()

	testData := []struct {
		desc     string
		method   string
		path     string
		message  string
		wantResp string
	}{
		{
			desc:     "Add Bearer token for CONSTANT_ADDRESS backend that requires JWT token",
			method:   "GET",
			path:     "/bearertoken/constant/42",
			wantResp: `{"Authorization": "Bearer ya29.constant", "RequestURI": "/bearertoken/constant?foo=42"}`,
		},
		{
			desc:     "Add Bearer token for APPEND_PATH_TO_ADDRESS backend that requires JWT token",
			method:   "GET",
			path:     "/bearertoken/append?key=api-key",
			wantResp: `{"Authorization": "Bearer ya29.append", "RequestURI": "/bearertoken/append?key=api-key"}`,
		},
		{
			desc:     "Do not reject backend that doesn't require JWT token",
			method:   "POST",
			path:     "/echo?key=api-key",
			message:  "hello",
			wantResp: `{"message":"hello"}`,
		},
	}

	for _, tc := range testData {
		url := fmt.Sprintf("http://localhost:%v%v", s.Ports().ListenerPort, tc.path)
		var resp []byte
		var err error
		switch tc.method {
		case "GET":
			resp, err = client.DoGet(url)
		case "POST":
			resp, err = client.DoPost(url, tc.message)
		default:
			t.Fatalf("Test Desc(%s): unsupported HTTP Method %s", tc.desc, tc.path)
		}

		if err != nil {
			t.Fatalf("Test Desc(%s): %v", tc.desc, err)
		}

		gotResp := string(resp)
		if !utils.JsonEqual(gotResp, tc.wantResp) {
			t.Errorf("Test Desc(%s): want: %s, got: %s", tc.desc, tc.wantResp, gotResp)
		}
	}
}

func TestServiceControlRequestForDynamicRouting(t *testing.T) {
	s := NewDynamicRoutingTestEnv(comp.TestServiceControlRequestInDynamicRouting)
	if err := s.Setup(testDynamicRoutingArgs); err != nil {
		t.Fatalf("fail to setup test env, %v", err)
	}
	defer s.TearDown()

	testData := []struct {
		desc           string
		path           string
		message        string
		wantResp       string
		wantScRequests []interface{}
	}{
		{
			desc:     "Succeed, no path translation (no re-routing needed)",
			path:     "/echo?key=api-key",
			message:  "hello",
			wantResp: `{"message":"hello"}`,
			wantScRequests: []interface{}{
				&utils.ExpectedCheck{
					Version:         utils.APIProxyVersion,
					ServiceName:     "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID: "test-config-id",
					ConsumerID:      "api_key:api-key",
					OperationName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
					CallerIp:        "127.0.0.1",
				},
				&utils.ExpectedReport{
					Version:           utils.APIProxyVersion,
					ServiceName:       "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID:   "test-config-id",
					URL:               "/echo?key=api-key",
					ApiKey:            "api-key",
					ApiMethod:         "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
					ProducerProjectID: "producer-project",
					ConsumerProjectID: "123456",
					FrontendProtocol:  "http",
					HttpMethod:        "POST",
					LogMessage:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo is called",
					RequestSize:       20,
					ResponseSize:      19,
					RequestBytes:      20,
					ResponseBytes:     19,
					ResponseCode:      200,
					Platform:          util.GCE,
					Location:          "test-zone",
				},
			},
		},
		{
			desc:     "Succeed, APPEND_PATH_TO_ADDRESS path translation is correct, service control check request and report request are correct",
			path:     "/sc/searchpet?key=api-key&timezone=EST",
			message:  "hello",
			wantResp: `{"RequestURI":"/dynamicrouting/sc/searchpet?key=api-key&timezone=EST"}`,
			wantScRequests: []interface{}{
				&utils.ExpectedCheck{
					Version:         utils.APIProxyVersion,
					ServiceName:     "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID: "test-config-id",
					ConsumerID:      "api_key:api-key",
					OperationName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPetWithServiceControlVerification",
					CallerIp:        "127.0.0.1",
				},
				&utils.ExpectedReport{
					Version:           utils.APIProxyVersion,
					ServiceName:       "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID:   "test-config-id",
					URL:               "/sc/searchpet?key=api-key&timezone=EST",
					ApiKey:            "api-key",
					ApiMethod:         "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPetWithServiceControlVerification",
					ProducerProjectID: "producer-project",
					ConsumerProjectID: "123456",
					FrontendProtocol:  "http",
					HttpMethod:        "POST",
					LogMessage:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.SearchPetWithServiceControlVerification is called",
					RequestSize:       20,
					ResponseSize:      71,
					RequestBytes:      20,
					ResponseBytes:     71,
					ResponseCode:      200,
					Platform:          util.GCE,
					Location:          "test-zone",
				},
			},
		},
		{
			desc:     "Succeed, CONSTANT_ADDRESS path translation is correct, service control check request and report request are correct",
			path:     "/sc/pet/0325/num/2019?key=api-key&lang=en",
			message:  "hello",
			wantResp: `{"RequestURI":"/dynamicrouting?key=api-key&lang=en&pet_id=0325&number=2019"}`,
			wantScRequests: []interface{}{
				&utils.ExpectedCheck{
					Version:         utils.APIProxyVersion,
					ServiceName:     "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID: "test-config-id",
					ConsumerID:      "api_key:api-key",
					OperationName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetByIdWithServiceControlVerification",
					CallerIp:        "127.0.0.1",
				},
				&utils.ExpectedReport{
					Version:           utils.APIProxyVersion,
					ServiceName:       "echo-api.endpoints.cloudesf-testing.cloud.goog",
					ServiceConfigID:   "test-config-id",
					URL:               "/sc/pet/0325/num/2019?key=api-key&lang=en",
					ApiKey:            "api-key",
					ApiMethod:         "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetByIdWithServiceControlVerification",
					ProducerProjectID: "producer-project",
					ConsumerProjectID: "123456",
					FrontendProtocol:  "http",
					HttpMethod:        "POST",
					LogMessage:        "1.echo_api_endpoints_cloudesf_testing_cloud_goog.dynamic_routing.GetPetByIdWithServiceControlVerification is called",
					RequestSize:       20,
					ResponseSize:      77,
					RequestBytes:      20,
					ResponseBytes:     77,
					ResponseCode:      200,
					Platform:          util.GCE,
					Location:          "test-zone",
				},
			},
		},
	}
	for _, tc := range testData {
		url := fmt.Sprintf("http://localhost:%v%v", s.Ports().ListenerPort, tc.path)
		var gotResp []byte
		var err error
		gotResp, err = client.DoPost(url, tc.message)

		if err != nil {
			t.Fatal(err)
		}

		gotRespStr := string(gotResp)

		if !utils.JsonEqual(gotRespStr, tc.wantResp) {
			t.Errorf("Test Desc(%s): response want: %s, got: %s", tc.desc, tc.wantResp, gotRespStr)
		}

		scRequests, err := s.ServiceControlServer.GetRequests(len(tc.wantScRequests), 2*time.Second)
		if err != nil {
			t.Fatalf("Test Desc(%s): GetRequests returns error: %v", tc.desc, err)
		}

		for i, wantScRequest := range tc.wantScRequests {
			reqBody := scRequests[i].ReqBody
			switch wantScRequest.(type) {
			case *utils.ExpectedCheck:
				if scRequests[i].ReqType != comp.CHECK_REQUEST {
					t.Errorf("Test Desc(%s): service control request %v: should be Check", tc.desc, i)
				}
				if err := utils.VerifyCheck(reqBody, wantScRequest.(*utils.ExpectedCheck)); err != nil {
					t.Error(err)
				}
			case *utils.ExpectedReport:
				if scRequests[i].ReqType != comp.REPORT_REQUEST {
					t.Errorf("Test Desc(%s): service control request %v: should be Report", tc.desc, i)
				}
				if err := utils.VerifyReport(reqBody, wantScRequest.(*utils.ExpectedReport)); err != nil {
					t.Error(err)
				}
			default:
				t.Fatalf("Test Desc(%s): unknown service control response type", tc.desc)
			}
		}
	}
}
