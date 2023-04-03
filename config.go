// Copyright © 2023 Meroxa, Inc.
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

package ais

// Config contains shared config parameters, common to the source and
// destination. If you don't need shared parameters you can entirely remove this
// file.
type Config struct {
	// APIURL is the Spire API URL to use for accessing the Spire GraphQL API.
	APIURL string `json:"apiUrl" default:"https://api.spire.com/graphql"`

	// Token is the access token to use when accessing the Spire GraphQL API.
	Token string `json:"token" validate:"required"`

	// BatchSize is the quantity of vessels to retrieve per API call.
	BatchSize string `json:"batchSize"`
}
