// Copyright 2018 The Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package nuki provides an simple implementation of the bridge api
// Versoin: v1.7 (30.03.2018)
package nuki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL = "http://%s:%s/%s"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// Nuki holds you access token and connection information
type Nuki struct {
	ip    string
	port  string
	token string
}

// NewNuki returns a new nuki entrypoint without a token
// to request a new token, the auth function can be called.
// this connects to the bridge and stores the token internally.
func NewNuki(ip, port string) *Nuki {
	return &Nuki{
		ip:   ip,
		port: port,
	}
}

// NewNukiWithToken returns a new nuki entrypoint
func NewNukiWithToken(ip, port, token string) *Nuki {
	return &Nuki{
		ip:    ip,
		port:  port,
		token: token,
	}
}

// Token returns a token provided by nuki
func (n *Nuki) Token() string {
	return n.token
}

func (n *Nuki) doRequest(method, path string, i interface{}, param url.Values) error {
	url := fmt.Sprintf(baseURL, n.ip, n.port, path)
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	r.URL.RawQuery = param.Encode()
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	if err := json.NewDecoder(res.Body).Decode(i); err != nil {
		return err
	}
	return ErrorFromStatus(res.StatusCode)
}

// Auth Enables the api (if not yet enabled) and returns the api token.
// If no api token has yet been set, a new (random) one is generated.
// When issuing this API-call the bridge turns on its LED for 30 seconds.
// The button of the bridge has to be pressed within this timeframe.
// Otherwise the bridge returns a negative success and no token.
func (n *Nuki) Auth() (*AuthReponse, error) {
	const path = "auth"
	auth := &AuthReponse{}
	if err := n.doRequest(http.MethodGet, path, &auth, url.Values{}); err != nil {
		return nil, err
	}
	n.token = auth.Token
	return auth, nil
}

// ConfigAuth enables or disables the authorization via /auth and the
// publication of the local IP and port to the discovery URL
// The api token configured via the Nuki app when enabling the API
// enable: Flag (0 or 1) indicating whether or not the authorization
// should be enabled
func (n *Nuki) ConfigAuth(enable bool) (*ConfigAuthResponse, error) {
	const path = "configAuth"
	config := &ConfigAuthResponse{}
	if err := n.doRequest(http.MethodGet, path, &config, nil); err != nil {
		return nil, err
	}
	return config, nil
}

// List returns all paired Smart Locks a valid token is required
func (n *Nuki) List() ([]ListResponse, error) {
	const path = "list"
	param := url.Values{}
	param.Set("token", n.token)
	var list []ListResponse
	if err := n.doRequest(http.MethodGet, path, &list, param); err != nil {
		return nil, err
	}
	return list, nil
}

// LockState returns the current lock state of a smart lock by give id
func (n *Nuki) LockState(nukiID int) (*LockStateResponse, error) {
	const path = "lockState"
	param := url.Values{}
	param.Set("nukiId", strconv.Itoa(nukiID))
	param.Set("token", n.token)
	lockState := &LockStateResponse{}
	if err := n.doRequest(http.MethodGet, path, &lockState, param); err != nil {
		return nil, err
	}
	return lockState, nil
}

// LockAction performs a lock operation on the given smart lock given by the id
// action: the desired lock action
// noWait: indicating whether or not to wait for the lock action to complete and
// return its result
func (n *Nuki) LockAction(nukiID int, action Action, noWait bool) (*LockActionResponse, error) {
	const path = "lockAction"
	param := url.Values{}
	param.Set("nukiId", strconv.Itoa(nukiID))
	param.Set("action", strconv.Itoa(int(action)))
	param.Set("noWait", strconv.FormatBool(noWait))
	param.Set("token", n.token)
	lockAction := &LockActionResponse{}
	if err := n.doRequest(http.MethodGet, path, &lockAction, param); err != nil {
		return nil, err
	}
	return lockAction, nil
}

// Unpair removes the pairing with a given smart lock
func (n *Nuki) Unpair(nukiID int) (*UnpairResponse, error) {
	const path = "unpair"
	param := url.Values{}
	param.Set("token", n.token)
	param.Set("nukiId", strconv.Itoa(nukiID))
	unpair := &UnpairResponse{}
	if err := n.doRequest(http.MethodGet, path, &unpair, param); err != nil {
		return nil, err
	}
	return unpair, nil
}

// Info returns all smart locks in range and some device information of the
// bridge itself
func (n *Nuki) Info() (*InfoResponse, error) {
	const path = "info"
	param := url.Values{}
	param.Set("token", n.token)
	info := &InfoResponse{}
	if err := n.doRequest(http.MethodGet, path, &info, param); err != nil {
		return nil, err
	}
	return info, nil
}

// CallbackAdd registers a new callback url
func (n *Nuki) CallbackAdd(nukiURL string) (*CallbackReponse, error) {
	const path = "callback/add"
	param := url.Values{}
	param.Set("token", n.token)
	param.Set("url", nukiURL)
	callback := &CallbackReponse{}
	if err := n.doRequest(http.MethodGet, path, &callback, param); err != nil {
		return nil, err
	}
	return callback, nil
}

// CallbackList returns a CallbackReponse with all registered url callbacks
func (n *Nuki) CallbackList() (*CallbackReponse, error) {
	const path = "callback/list"
	param := url.Values{}
	param.Set("token", n.token)
	callback := &CallbackReponse{}
	if err := n.doRequest(http.MethodGet, path, &callback, param); err != nil {
		return nil, err
	}
	return callback, nil
}

// CallbackRemove removes a previously added callback by ID
func (n *Nuki) CallbackRemove(callbackID int) (*CallbackReponse, error) {
	const path = "callback/remove"
	param := url.Values{}
	param.Set("token", n.token)
	param.Set("id", strconv.Itoa(callbackID))
	callback := &CallbackReponse{}
	if err := n.doRequest(http.MethodGet, path, &callback, param); err != nil {
		return nil, err
	}
	return callback, nil
}

// The following endpoints are available for maintenance purposes of the
// hardware bridge. Therefore they are not available on the software bridge.
// TODO

// Log returns a log of the Bridge
func (n *Nuki) Log() (*LogResponse, error) {
	panic("not implemented yet")
}

// ClearLog clears the log of the Bridge
func (n *Nuki) ClearLog() error {
	panic("not implemented yet")
}

// FWUpdate immediately checks for a new firmware update and installs it
func (n *Nuki) FWUpdate() error {
	panic("not implemented yet")
}

// Reboot s the bridge
func (n *Nuki) Reboot() error {
	panic("not implemented yet")
}

// FactoryReset performs a factory reset
func (n *Nuki) FactoryReset() error {
	panic("not implemented yet")
}
