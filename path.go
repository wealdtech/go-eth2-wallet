// Copyright 2020 Weald Technology Trading
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
	"strings"
)

// WalletAndAccountNames breaks an account in to wallet and account names.
func WalletAndAccountNames(account string) (string, string, error) {
	if len(account) == 0 {
		return "", "", errors.New("invalid account format")
	}
	index := strings.Index(account, "/")
	if index == -1 {
		// Just the wallet
		return account, "", nil
	}
	if index == 0 {
		return "", "", errors.New("invalid account format")
	}
	if index == len(account)-1 {
		// Trailing /
		return account[:index], "", nil
	}
	return account[:index], account[index+1:], nil
}
