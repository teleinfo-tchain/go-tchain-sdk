/********************************************************************************
   This file is part of go-bif.
   go-bif is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-bif is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-bif.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package providers

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/tchain/go-tchain-sdk/providers/util"
)

type HTTPProvider struct {
	address string
	timeout int32
	secure  bool
	client  *http.Client
}

func NewHTTPProvider(address string, timeout int32, secure bool) *HTTPProvider {
	return NewHTTPProviderWithClient(address, timeout, secure, &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	})
}

func NewHTTPProviderWithClient(address string, timeout int32, secure bool, client *http.Client) *HTTPProvider {
	provider := new(HTTPProvider)
	provider.address = address
	provider.timeout = timeout
	provider.secure = secure
	provider.client = client

	return provider
}

func (provider HTTPProvider) SendRequest(v interface{}, method string, params interface{}) error {

	bodyString := util.JSONRPCObject{Version: util.Version, Method: method, Params: params, ID: rand.Intn(100)}

	prefix := "http://"
	if provider.secure {
		prefix = "https://"
	}

	body := strings.NewReader(bodyString.AsJsonString())
	req, err := http.NewRequest("POST", prefix+provider.address, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := provider.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var bodyBytes []byte

	if resp.StatusCode == 200 {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}

	return json.Unmarshal(bodyBytes, v)

}

func (provider HTTPProvider) Close() error { return nil }
