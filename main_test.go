// Copyright © 2020 Weald Technology Trading
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

package wallet_test

import (
	"os"
	"testing"

	e2types "github.com/wealdtech/go-eth2-types/v2"
)

func TestMain(m *testing.M) {
	if err := e2types.InitBLS(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
