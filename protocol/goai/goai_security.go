// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/internal/json"
)

type SecurityScheme struct {
	Type             string      `json:"type,omitempty"             yaml:"type,omitempty"`
	Description      string      `json:"description,omitempty"      yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty"             yaml:"name,omitempty"`
	In               string      `json:"in,omitempty"               yaml:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"           yaml:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"     yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"            yaml:"flows,omitempty"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

type SecuritySchemes map[string]SecuritySchemeRef

type SecuritySchemeRef struct {
	Ref   string
	Value *SecurityScheme
}

type SecurityRequirements []SecurityRequirement

type SecurityRequirement map[string][]string

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"          yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"          yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"         yaml:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"       yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"                     yaml:"scopes"`
}

func (r SecuritySchemeRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}
