package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestCreateAccessTokenPayloadWithSessionClaims(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							trueClaim, _ := TrueClaim()
							accessTokenPayload, err := trueClaim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(res, userID, accessTokenPayload, sessionData, userContext)
						}
						*originalImplementation.CreateNewSession = nCreateNewSession
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer
	accessTokenPayload := map[string]interface{}{}

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(rw, "rope", accessTokenPayload, map[string]interface{}{})
		assert.NoError(t, err)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	accessTokenPayload = sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 1, len(accessTokenPayload))
	assert.NotNil(t, accessTokenPayload["st-true"])
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])
	assert.Greater(t, accessTokenPayload["st-true"].(map[string]interface{})["t"], float64(time.Now().UnixNano()/1000000-1000))
}

func TestNotCreateAccessTokenPayloadWithNilClaim(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							nilClaim, _ := NilClaim()
							accessTokenPayload, err := nilClaim.Build(userID, accessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							return oCreateNewSession(res, userID, accessTokenPayload, sessionData, userContext)
						}
						*originalImplementation.CreateNewSession = nCreateNewSession
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer
	accessTokenPayload := map[string]interface{}{}

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(rw, "rope", accessTokenPayload, map[string]interface{}{})
		assert.NoError(t, err)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	accessTokenPayload = sessionContainer.GetAccessTokenPayload()
	assert.Equal(t, 0, len(accessTokenPayload))
}

func TestMergeClaimsAndPassedAccessTokenPayload(t *testing.T) {
	payloadParam := map[string]interface{}{
		"initial": true,
	}
	custom2 := map[string]interface{}{
		"nilProp": nil,
		"inner":   "asdf",
	}
	customClaims := map[string]interface{}{
		"user-custom":  "asdf",
		"user-custom2": custom2,
		"user-custom3": nil,
	}

	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				Override: &sessmodels.OverrideStruct{
					Functions: func(originalImplementation sessmodels.RecipeInterface) sessmodels.RecipeInterface {
						oCreateNewSession := *originalImplementation.CreateNewSession
						nCreateNewSession := func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
							nAccessTokenPayload := map[string]interface{}{}
							for k, v := range accessTokenPayload {
								nAccessTokenPayload[k] = v
							}
							trueClaim, _ := TrueClaim()
							nAccessTokenPayload, err := trueClaim.Build(userID, nAccessTokenPayload, userContext)
							if err != nil {
								return nil, err
							}
							for k, v := range customClaims {
								nAccessTokenPayload[k] = v
							}
							return oCreateNewSession(res, userID, nAccessTokenPayload, sessionData, userContext)
						}
						*originalImplementation.CreateNewSession = nCreateNewSession
						return originalImplementation
					},
				},
			}),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		t.Error(err.Error())
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		t.Error(err.Error())
	}

	includesNullInPayload := unittesting.MaxVersion(cdiVersion, "2.14") != "2.14"

	mux := http.NewServeMux()
	var sessionContainer sessmodels.SessionContainer

	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		var err error
		sessionContainer, err = CreateNewSession(rw, "rope", payloadParam, map[string]interface{}{})
		assert.NoError(t, err)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/create", nil)
	assert.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// The passed object should be unchanged
	assert.Equal(t, 1, len(payloadParam))

	accessTokenPayload := sessionContainer.GetAccessTokenPayload()
	if includesNullInPayload {
		assert.Equal(t, 5, len(accessTokenPayload))
	} else {
		assert.Equal(t, 4, len(accessTokenPayload))
	}

	// We have the prop from the payload param
	assert.Equal(t, true, accessTokenPayload["initial"])

	// We have the boolean claim
	assert.NotNil(t, accessTokenPayload["st-true"])
	assert.Equal(t, true, accessTokenPayload["st-true"].(map[string]interface{})["v"])
	assert.Greater(t, accessTokenPayload["st-true"].(map[string]interface{})["t"], float64(time.Now().UnixNano()/1000000-1000))

	// We have the custom claim
	// The resulting payload is different from the input: it doesn't container nil values
	assert.Equal(t, "asdf", accessTokenPayload["user-custom"])
	if includesNullInPayload {
		assert.Equal(t, custom2, accessTokenPayload["user-custom2"])
	} else {
		assert.Equal(t, map[string]interface{}{
			"inner": "asdf",
		}, accessTokenPayload["user-custom2"])
	}
}
