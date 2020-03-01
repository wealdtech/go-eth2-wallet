// Copyright © 2019 Weald Technology Trading
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

package wallet

import (
	"errors"
	"fmt"

	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
	s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

var store wtypes.Store

func init() {
	// default store is filesystem
	store = filesystem.New()
}

// SetStore sets a store to use given its name and optional passphrase.
// This does not allow access to all advanced features of stores.  To access these create the stores yourself and set them with
// `UseStore()`.
func SetStore(name string, passphrase []byte) error {
	var store wtypes.Store
	var err error
	switch name {
	case "s3":
		store, err = s3.New(s3.WithPassphrase(passphrase))
	case "filesystem":
		store = filesystem.New(filesystem.WithPassphrase(passphrase))
	default:
		err = fmt.Errorf("unknown wallet store %q", name)
	}
	if err != nil {
		return err
	}

	return UseStore(store)
}

// UseStore sets a store to use.
func UseStore(s wtypes.Store) error {
	if s == nil {
		return errors.New("no store supplied")
	}
	store = s
	return nil
}

// GetStore returns the name of the current store.
func GetStore() string {
	return store.Name()
}
