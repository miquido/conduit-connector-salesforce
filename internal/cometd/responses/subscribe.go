// Copyright © 2022 Meroxa, Inc. and Miquido
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

package responses

import "fmt"

// SubscribeResponse represents subscription response.
// See: https://docs.cometd.org/current7/reference/#_subscribe_response
type SubscribeResponse struct {
	Channel      string      `json:"channel"`
	Successful   bool        `json:"successful"`
	Subscription interface{} `json:"subscription"`
	Error        string      `json:"error,omitempty"`
	Advice       *advice     `json:"advice,omitempty"`
	Ext          *ext        `json:"ext,omitempty"`
	ClientID     string      `json:"clientId,omitempty"`
	ID           string      `json:"id,omitempty"`
}

func (s SubscribeResponse) GetSubscriptions() []string {
	switch subscription := s.Subscription.(type) {
	case string:
		return []string{subscription}

	case []string:
		return subscription
	}

	panic(fmt.Errorf("unexpected subscriptions data: %#v", s.Subscription))
}
