/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package epunittesting

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestGenerateTokenAPIWithValidInputAndEmailNotVerified(t *testing.T) {
	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
			}),
		},
	}

	unittesting.BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))

	resp, err := unittesting.SignupRequest("test@gmail.com", "testPass123", testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}
	assert.NoError(t, err)
	// cookieData := unittesting.ExtractInfoFromResponse(resp)
	assert.Equal(t, 200, resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var response map[string]interface{}
	_ = json.Unmarshal(data, &response)
	assert.Equal(t, "OK", response["status"])

	// emailpassword.CreateEmailVerificationToken()

	// user := response["user"]
	// userstr := fmt.Sprintf("%v", user)
	// fmt.Println(user.(map[string]string)["email"])
	// userDataArr := strings.Split(userstr, ":")
	// for _, v := range us {

	// }
	// defer unittesting.AfterEach()
	// defer func() {
	// 	testServer.Close()
	// }()
}
