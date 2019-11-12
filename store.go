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

	filesystem "github.com/wealdtech/go-eth2-wallet-store-filesystem"
	types "github.com/wealdtech/go-eth2-wallet-types"
)

var store types.Store

func init() {
	// default store is filesystem
	store = filesystem.New()
}

// UseStore sets a store to use.
func UseStore(s types.Store) error {
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
