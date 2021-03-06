// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package auth

import (
	"testing"
	"time"

	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/stretchr/testify/assert"
)

var l = NewUserLock(2 * time.Second)

var adminServerLdapTestConfig = map[string]interface{}{
	common.ExtEndpoint:          "host01.com",
	common.AUTHMode:             "ldap_auth",
	common.DatabaseType:         "postgresql",
	common.PostGreSQLHOST:       "127.0.0.1",
	common.PostGreSQLPort:       5432,
	common.PostGreSQLUsername:   "postgres",
	common.PostGreSQLPassword:   "root123",
	common.PostGreSQLDatabase:   "registry",
	common.LDAPURL:              "ldap://127.0.0.1",
	common.LDAPSearchDN:         "cn=admin,dc=example,dc=com",
	common.LDAPSearchPwd:        "admin",
	common.LDAPBaseDN:           "dc=example,dc=com",
	common.LDAPUID:              "uid",
	common.LDAPFilter:           "",
	common.LDAPScope:            3,
	common.LDAPTimeout:          30,
	common.CfgExpiration:        5,
	common.AdminInitialPassword: "password",
}

func TestLock(t *testing.T) {
	t.Log("Locking john")
	l.Lock("john")
	if !l.IsLocked("john") {
		t.Errorf("John should be locked")
	}
	t.Log("Locking jack")
	l.Lock("jack")
	t.Log("Sleep for 2 seconds and check...")
	time.Sleep(2 * time.Second)
	if l.IsLocked("jack") {
		t.Errorf("After 2 seconds, jack shouldn't be locked")
	}
	if l.IsLocked("daniel") {
		t.Errorf("daniel has never been locked, he should not be locked")
	}
}

func TestDefaultAuthenticate(t *testing.T) {
	authHelper := DefaultAuthenticateHelper{}
	m := models.AuthModel{}
	user, err := authHelper.Authenticate(m)
	if user != nil || err == nil {
		t.Fatal("Default implementation should return nil")
	}
}

func TestDefaultOnBoardUser(t *testing.T) {
	user := &models.User{}
	authHelper := DefaultAuthenticateHelper{}
	err := authHelper.OnBoardUser(user)
	if err == nil {
		t.Fatal("Default implementation should return error")
	}
}

func TestDefaultMethods(t *testing.T) {
	authHelper := DefaultAuthenticateHelper{}
	_, err := authHelper.SearchUser("sample")
	if err == nil {
		t.Fatal("Default implementation should return error")
	}

	_, err = authHelper.SearchGroup("sample")
	if err == nil {
		t.Fatal("Default implementation should return error")
	}

	err = authHelper.OnBoardGroup(&models.UserGroup{}, "sample")
	if err == nil {
		t.Fatal("Default implementation should return error")
	}
}

func TestErrAuth(t *testing.T) {
	assert := assert.New(t)
	e := NewErrAuth("test")
	expectedStr := "Failed to authenticate user, due to error 'test'"
	assert.Equal(expectedStr, e.Error())
}
