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

	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	types "github.com/wealdtech/go-eth2-wallet-types"
)

var encryptor types.Encryptor

func init() {
	encryptor = keystorev4.New()
}

// UseEncryptor sets an encryptor to use.
func UseEncryptor(e types.Encryptor) error {
	if e == nil {
		return errors.New("no encryptor supplied")
	}
	encryptor = e
	return nil
}

// GetEncryptor returns the name of the current encryptor.
func GetEncryptor() string {
	return encryptor.Name()
}
