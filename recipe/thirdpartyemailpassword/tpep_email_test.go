/*
 * Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package thirdpartyemailpassword

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultBackwardCompatibilityPasswordResetForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	EmailPasswordSignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Equal(t, emailpassword.PasswordResetDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestDefaultBackwardCompatibilityPasswordResetForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	ThirdPartySignInUp("custom", "user-id", tpepmodels.EmailStruct{ID: "test@example.com", IsVerified: true})
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestDefaultBackwardCompatibilityPasswordResetForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)
}

func TestBackwardCompatibilityResetPasswordForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				passwordResetLink = passwordResetURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	EmailPasswordSignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.True(t, customCalled)
}

func TestBackwardCompatibilityResetPasswordForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				passwordResetLink = passwordResetURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	ThirdPartySignInUp("custom", "user-id", tpepmodels.EmailStruct{ID: "test@example.com", IsVerified: true})
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestBackwardCompatibilityResetPasswordForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		ResetPasswordUsingTokenFeature: &epmodels.TypeInputResetPasswordUsingTokenFeature{
			CreateAndSendCustomEmail: func(user epmodels.User, passwordResetURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				passwordResetLink = passwordResetURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestCustomOverrideResetPasswordForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	EmailPasswordSignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.True(t, customCalled)
}

func TestCustomOverrideResetPasswordForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	ThirdPartySignInUp("custom", "user-id", tpepmodels.EmailStruct{ID: "test@example.com", IsVerified: true})
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestCustomOverrideResetPasswordForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	passwordResetLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.PasswordReset != nil {
						customCalled = true
						email = input.PasswordReset.User.Email
						passwordResetLink = input.PasswordReset.PasswordResetLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, customCalled)
}

func TestSMTPOverridePasswordResetForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPContent, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.SMTPContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	EmailPasswordSignUp("test@example.com", "1234abcd")
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, passwordResetLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPOverridePasswordResetForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPContent, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.SMTPContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	ThirdPartySignInUp("custom", "user-id", tpepmodels.EmailStruct{ID: "test@example.com", IsVerified: true})
	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, getContentCalled)
	assert.False(t, sendRawEmailCalled)
}

func TestSMTPOverridePasswordResetForNonExistantUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	passwordResetLink := ""

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPContent, error) {
				if input.PasswordReset != nil {
					email = input.PasswordReset.User.Email
					passwordResetLink = input.PasswordReset.PasswordResetLink
					getContentCalled = true
				}
				return emaildelivery.SMTPContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.PasswordResetTokenRequest("test@example.com", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler not called
	assert.Empty(t, email)
	assert.Empty(t, passwordResetLink)
	assert.False(t, getContentCalled)
	assert.False(t, sendRawEmailCalled)
}

func TestDefaultBackwardCompatibilityEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	testServer := supertokensInitForTest(t, session.Init(nil), Init(nil))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Equal(t, emailverification.EmailVerificationDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)
}

func TestDefaultBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	tpepConfig := &tpepmodels.TypeInput{
		Providers: []tpmodels.TypeProvider{
			customProviderForEmailVerification,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()

	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, emailverification.EmailVerificationEmailSentForTest)
	assert.Equal(t, emailverification.EmailVerificationDataForTest.User.Email, "test@example.com")
	assert.NotEmpty(t, emailverification.EmailVerificationDataForTest.EmailVerifyURLWithToken)
}

func TestBackwardCompatibilityEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailVerificationFeature: &tpepmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user tpepmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				emailVerifyLink = emailVerificationURLWithToken
				customCalled = true
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

func TestBackwardCompatibilityEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""
	var thirdparty *struct {
		ID     string `json:"id"`
		UserID string `json:"userId"`
	}

	tpepConfig := &tpepmodels.TypeInput{
		EmailVerificationFeature: &tpepmodels.TypeInputEmailVerificationFeature{
			CreateAndSendCustomEmail: func(user tpepmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
				email = user.Email
				emailVerifyLink = emailVerificationURLWithToken
				thirdparty = user.ThirdParty
				customCalled = true
			},
		},
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.NotNil(t, thirdparty)
	assert.True(t, customCalled)
}

func TestCustomOverrideEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.EmailVerification != nil {
						customCalled = true
						email = input.EmailVerification.User.Email
						emailVerifyLink = input.EmailVerification.EmailVerifyLink
					}
					return nil
				}
				return originalImplementation
			},
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)
	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

func TestCustomOverrideEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customCalled := false
	email := ""
	emailVerifyLink := ""

	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Override: func(originalImplementation emaildelivery.EmailDeliveryInterface) emaildelivery.EmailDeliveryInterface {
				*originalImplementation.SendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
					if input.EmailVerification != nil {
						customCalled = true
						email = input.EmailVerification.User.Email
						emailVerifyLink = input.EmailVerification.EmailVerifyLink
					}
					return nil
				}
				return originalImplementation
			},
		},
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	// Custom handler called
	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.True(t, customCalled)
}

func TestSMTPOverrideEmailVerifyForEmailPasswordUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	emailVerifyLink := ""

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPContent, error) {
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
					getContentCalled = true
				}
				return emaildelivery.SMTPContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	resp, err := unittesting.SignupRequest("test@example.com", "1234abcd", testServer.URL)
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}

func TestSMTPOverrideEmailVerifyForThirdpartyUser(t *testing.T) {
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	getContentCalled := false
	sendRawEmailCalled := false
	email := ""
	emailVerifyLink := ""

	smtpService := MakeSMTPService(emaildelivery.SMTPServiceConfig{
		Settings: emaildelivery.SMTPSettings{
			Host: "",
			From: emaildelivery.SMTPFrom{
				Name:  "Test User",
				Email: "",
			},
			Port:     123,
			Password: "",
		},
		Override: func(originalImplementation emaildelivery.SMTPServiceInterface) emaildelivery.SMTPServiceInterface {
			(*originalImplementation.GetContent) = func(input emaildelivery.EmailType, userContext supertokens.UserContext) (emaildelivery.SMTPContent, error) {
				if input.EmailVerification != nil {
					email = input.EmailVerification.User.Email
					emailVerifyLink = input.EmailVerification.EmailVerifyLink
					getContentCalled = true
				}
				return emaildelivery.SMTPContent{}, nil
			}

			(*originalImplementation.SendRawEmail) = func(input emaildelivery.SMTPContent, userContext supertokens.UserContext) error {
				sendRawEmailCalled = true
				return nil
			}

			return originalImplementation
		},
	})
	tpepConfig := &tpepmodels.TypeInput{
		EmailDelivery: &emaildelivery.TypeInput{
			Service: smtpService,
		},
		Providers: []tpmodels.TypeProvider{customProviderForEmailVerification},
	}
	testServer := supertokensInitForTest(t, session.Init(nil), Init(tpepConfig))
	defer testServer.Close()

	signinupPostData := PostDataForCustomProvider{
		ThirdPartyId: "custom",
		AuthCodeResponse: map[string]string{
			"access_token": "saodiasjodai",
		},
		RedirectUri: "http://127.0.0.1/callback",
	}

	postBody, err := json.Marshal(signinupPostData)
	resp, err := http.Post(testServer.URL+"/auth/signinup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)

	cookies := resp.Cookies()
	resp, err = unittesting.EmailVerificationTokenRequest(cookies, testServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Default handler not called
	assert.False(t, emailpassword.PasswordResetEmailSentForTest)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.User.Email)
	assert.Empty(t, emailpassword.PasswordResetDataForTest.PasswordResetURLWithToken)

	assert.Equal(t, email, "test@example.com")
	assert.NotEmpty(t, emailVerifyLink)
	assert.Equal(t, getContentCalled, true)
	assert.Equal(t, sendRawEmailCalled, true)
}
