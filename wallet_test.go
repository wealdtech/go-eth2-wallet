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
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	wallet "github.com/wealdtech/go-eth2-wallet"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
)

func _byte(input string) []byte {
	res, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	return res
}

func TestCreateWallet(t *testing.T) {
	require.NoError(t, wallet.UseStore(scratch.New()))
	tests := []struct {
		name    string
		options []wallet.Option
		err     string
	}{
		{
			name: "Nil",
		},
		{
			name: "EncryptorNil",
			options: []wallet.Option{
				wallet.WithEncryptor(nil),
				wallet.WithStore(scratch.New()),
				wallet.WithPassphrase([]byte("secret")),
				wallet.WithType("hd"),
				wallet.WithSeed(_byte("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			},
			err: "no encryptor specified",
		},
		{
			name: "EncryptorNil",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(nil),
				wallet.WithPassphrase([]byte("secret")),
				wallet.WithType("hd"),
				wallet.WithSeed(_byte("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			},
			err: "no store specified",
		},
		{
			name: "SeedNil",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(scratch.New()),
				wallet.WithPassphrase([]byte("secret")),
				wallet.WithType("hd"),
				wallet.WithSeed(nil),
			},
			err: "no seed specified",
		},
		{
			name: "TypeUnknown",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(scratch.New()),
				wallet.WithPassphrase([]byte("secret")),
				wallet.WithType("unknown"),
				wallet.WithSeed(nil),
			},
			err: "unhandled wallet type \"unknown\"",
		},
		{
			name: "HDGood",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(scratch.New()),
				wallet.WithPassphrase([]byte("secret")),
				wallet.WithType("hd"),
				wallet.WithSeed(_byte("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
			},
		},
		{
			name: "DistributedGood",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(scratch.New()),
				wallet.WithType("distributed"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet, err := wallet.CreateWallet(test.name, test.options...)
			if test.err == "" {
				require.NoError(t, err)
				require.NotNil(t, wallet)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}

func TestOpenWallet(t *testing.T) {
	store := scratch.New()
	require.NoError(t, wallet.UseStore(store))
	_, err := wallet.CreateWallet("Good",
		wallet.WithEncryptor(keystorev4.New()),
		wallet.WithStore(store),
		wallet.WithPassphrase([]byte("secret")),
		wallet.WithType("hd"),
		wallet.WithSeed(_byte("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		options []wallet.Option
		err     string
	}{
		{
			name: "Nil",
			err:  "wallet not found",
		},
		{
			name: "EncryptorNil",
			options: []wallet.Option{
				wallet.WithEncryptor(nil),
				wallet.WithStore(scratch.New()),
			},
			err: "no encryptor specified",
		},
		{
			name: "EncryptorNil",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(nil),
			},
			err: "no store specified",
		},
		{
			name: "Missing",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(store),
			},
			err: "wallet not found",
		},
		{
			name: "Good",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(store),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet, err := wallet.OpenWallet(test.name, test.options...)
			if test.err == "" {
				require.NoError(t, err)
				require.NotNil(t, wallet)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}

func TestWallets(t *testing.T) {
	store := scratch.New()
	require.NoError(t, wallet.UseStore(store))
	_, err := wallet.CreateWallet("Good",
		wallet.WithEncryptor(keystorev4.New()),
		wallet.WithStore(store),
		wallet.WithPassphrase([]byte("secret")),
		wallet.WithType("hd"),
		wallet.WithSeed(_byte("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		options []wallet.Option
		result  bool
	}{
		{
			name:   "Nil",
			result: true,
		},
		{
			name: "EncryptorNil",
			options: []wallet.Option{
				wallet.WithEncryptor(nil),
				wallet.WithStore(scratch.New()),
			},
		},
		{
			name: "StoreNil",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(nil),
			},
		},
		{
			name: "StoreEmpty",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(scratch.New()),
			},
		},
		{
			name: "Good",
			options: []wallet.Option{
				wallet.WithEncryptor(keystorev4.New()),
				wallet.WithStore(store),
			},
			result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := wallet.Wallets(test.options...)
			// Sleep to allow goroutine to fetch the wallet(s).
			time.Sleep(time.Second)
			if test.result {
				select {
				case <-ch:
					// Wallet found.
				default:
					require.Fail(t, "missing expected wallet")
				}
			} else {
				select {
				case w := <-ch:
					if w != nil {
						require.Fail(t, fmt.Sprintf("unexpected wallet %s", w.Name()))
					}
				default:
					// Nothing found.
				}
			}
		})
	}
}

func TestImportWallet(t *testing.T) {
	require.NoError(t, wallet.UseStore(scratch.New()))
	tests := []struct {
		name       string
		data       []byte
		passphrase []byte
		err        string
	}{
		{
			name: "Nil",
			err:  "encrypted data must be at least 81 bytes",
		},
		{
			name: "PassphraseMissing",
			data: _byte("0x0153eaead344082ae20d4d8b7f4c6666c54ac46480c3bc1fc817491d1cbc5af1cbc5bacd1a8dca7c64446bb94483efaa4c1af99f0271bf1ebe49ef50052020a07b8e000c2686702aa2e6a65fae5d51c84a0054018ecd4b170b76ccfa287e62c513f8110cd63373a13d73ab93a0abd41aaa784757f68cd669cb16588b99171b83a062b2a0a8a4886e0149a726ac9a8364f16643c3ff025d47e4ec5d1e1445e51cb5790aa32c8d2577b42ab8b89f3f93f3db402caa54cc51ec054d6338f9599cd0a5d7021d153758396f9d48d4d4d209f57222745666441a7619f8659504a1c27db0e5f5cdbfbbb9fc07544601b1a09fdc9d50b7c489bbfadbc09e2458c4e674dda1969dd700cf23181de8e603a7df915607166cefb30774ef20333dcea601979ad4f1f6f7d57aa2cac9996b60d2b980018883fa831d869beec2b0532b54530b143f3cda6b8d7962f47e5f2795ac35b78fe8a2eeacd399cad0d2c3bc6c7a1b799b99bc5bf0a5b149ab87d55c98a0a0fe6df2ba0c6931e26e014931b3e2067b2013a3d21c73b7a67bfd965de720f27fab4101cb56ddab57dc06e7e9118b4a450ae4b4401e13cdc92716c803470c6e2ae7f5120cd545d9c8d39a3b0e87d9a0262c1493abb26419a8ae03bd0b564fc6b217ce0fc63dd6557727d96cc51d67099cf3caf4f379af7219a81257dec743392322b57cca10c9a6defe807f8e6ce8147a64706056ec2bf43b97e6f30db98093e90ad9ce613592fd0e257cff1da6b2aa3341d57c8e979fbde3bbe1d9e933fbac5eef7dcb879cb6090ab4ebf59c86ca1bbe808378f672235ffd79a312c4fc22f14beb66d71247ea2218de9997afdc2844569db8de3366841f7d1d6c199b0a9633a763003008e9869963d90504d226e4585f8fe8e28f1f8e163278e3caf08d04167ab667c62a089367e0f945e77eb44b6f7d51f09cff840b1b49adc34905032797d56581c24bf9ee034a085eb2f1ed1e0b74ea213d1a448f32227991b6b349649315fb41604a857044c9f2f7657a502200c449744c4eb682cbbdcf93be3bc936dda438979c2fd9d5b910502a10bb9028411523e2942a765495e4f6f49b9702f329643fa475eb3bc76370f086e5312de705ef86d6e35ba012beaff3c68bcf1d6aa6c192e476d91db0a3f68e3b189662bab2463ecdc961bbad7680ae5ff4172b0536d443fe4a5e609aabf270a690b7ea3863fc82cab1e2480983ef84994a0599d80f85c927918f38425a55c2bd5ea9"),
			err:  "invalid key",
		},
		{
			name:       "PassphraseIncorrect",
			data:       _byte("0x0153eaead344082ae20d4d8b7f4c6666c54ac46480c3bc1fc817491d1cbc5af1cbc5bacd1a8dca7c64446bb94483efaa4c1af99f0271bf1ebe49ef50052020a07b8e000c2686702aa2e6a65fae5d51c84a0054018ecd4b170b76ccfa287e62c513f8110cd63373a13d73ab93a0abd41aaa784757f68cd669cb16588b99171b83a062b2a0a8a4886e0149a726ac9a8364f16643c3ff025d47e4ec5d1e1445e51cb5790aa32c8d2577b42ab8b89f3f93f3db402caa54cc51ec054d6338f9599cd0a5d7021d153758396f9d48d4d4d209f57222745666441a7619f8659504a1c27db0e5f5cdbfbbb9fc07544601b1a09fdc9d50b7c489bbfadbc09e2458c4e674dda1969dd700cf23181de8e603a7df915607166cefb30774ef20333dcea601979ad4f1f6f7d57aa2cac9996b60d2b980018883fa831d869beec2b0532b54530b143f3cda6b8d7962f47e5f2795ac35b78fe8a2eeacd399cad0d2c3bc6c7a1b799b99bc5bf0a5b149ab87d55c98a0a0fe6df2ba0c6931e26e014931b3e2067b2013a3d21c73b7a67bfd965de720f27fab4101cb56ddab57dc06e7e9118b4a450ae4b4401e13cdc92716c803470c6e2ae7f5120cd545d9c8d39a3b0e87d9a0262c1493abb26419a8ae03bd0b564fc6b217ce0fc63dd6557727d96cc51d67099cf3caf4f379af7219a81257dec743392322b57cca10c9a6defe807f8e6ce8147a64706056ec2bf43b97e6f30db98093e90ad9ce613592fd0e257cff1da6b2aa3341d57c8e979fbde3bbe1d9e933fbac5eef7dcb879cb6090ab4ebf59c86ca1bbe808378f672235ffd79a312c4fc22f14beb66d71247ea2218de9997afdc2844569db8de3366841f7d1d6c199b0a9633a763003008e9869963d90504d226e4585f8fe8e28f1f8e163278e3caf08d04167ab667c62a089367e0f945e77eb44b6f7d51f09cff840b1b49adc34905032797d56581c24bf9ee034a085eb2f1ed1e0b74ea213d1a448f32227991b6b349649315fb41604a857044c9f2f7657a502200c449744c4eb682cbbdcf93be3bc936dda438979c2fd9d5b910502a10bb9028411523e2942a765495e4f6f49b9702f329643fa475eb3bc76370f086e5312de705ef86d6e35ba012beaff3c68bcf1d6aa6c192e476d91db0a3f68e3b189662bab2463ecdc961bbad7680ae5ff4172b0536d443fe4a5e609aabf270a690b7ea3863fc82cab1e2480983ef84994a0599d80f85c927918f38425a55c2bd5ea9"),
			passphrase: []byte("incorrect"),
			err:        "invalid key",
		},
		{
			name:       "NDGood",
			data:       _byte("0x0153eaead344082ae20d4d8b7f4c6666c54ac46480c3bc1fc817491d1cbc5af1cbc5bacd1a8dca7c64446bb94483efaa4c1af99f0271bf1ebe49ef50052020a07b8e000c2686702aa2e6a65fae5d51c84a0054018ecd4b170b76ccfa287e62c513f8110cd63373a13d73ab93a0abd41aaa784757f68cd669cb16588b99171b83a062b2a0a8a4886e0149a726ac9a8364f16643c3ff025d47e4ec5d1e1445e51cb5790aa32c8d2577b42ab8b89f3f93f3db402caa54cc51ec054d6338f9599cd0a5d7021d153758396f9d48d4d4d209f57222745666441a7619f8659504a1c27db0e5f5cdbfbbb9fc07544601b1a09fdc9d50b7c489bbfadbc09e2458c4e674dda1969dd700cf23181de8e603a7df915607166cefb30774ef20333dcea601979ad4f1f6f7d57aa2cac9996b60d2b980018883fa831d869beec2b0532b54530b143f3cda6b8d7962f47e5f2795ac35b78fe8a2eeacd399cad0d2c3bc6c7a1b799b99bc5bf0a5b149ab87d55c98a0a0fe6df2ba0c6931e26e014931b3e2067b2013a3d21c73b7a67bfd965de720f27fab4101cb56ddab57dc06e7e9118b4a450ae4b4401e13cdc92716c803470c6e2ae7f5120cd545d9c8d39a3b0e87d9a0262c1493abb26419a8ae03bd0b564fc6b217ce0fc63dd6557727d96cc51d67099cf3caf4f379af7219a81257dec743392322b57cca10c9a6defe807f8e6ce8147a64706056ec2bf43b97e6f30db98093e90ad9ce613592fd0e257cff1da6b2aa3341d57c8e979fbde3bbe1d9e933fbac5eef7dcb879cb6090ab4ebf59c86ca1bbe808378f672235ffd79a312c4fc22f14beb66d71247ea2218de9997afdc2844569db8de3366841f7d1d6c199b0a9633a763003008e9869963d90504d226e4585f8fe8e28f1f8e163278e3caf08d04167ab667c62a089367e0f945e77eb44b6f7d51f09cff840b1b49adc34905032797d56581c24bf9ee034a085eb2f1ed1e0b74ea213d1a448f32227991b6b349649315fb41604a857044c9f2f7657a502200c449744c4eb682cbbdcf93be3bc936dda438979c2fd9d5b910502a10bb9028411523e2942a765495e4f6f49b9702f329643fa475eb3bc76370f086e5312de705ef86d6e35ba012beaff3c68bcf1d6aa6c192e476d91db0a3f68e3b189662bab2463ecdc961bbad7680ae5ff4172b0536d443fe4a5e609aabf270a690b7ea3863fc82cab1e2480983ef84994a0599d80f85c927918f38425a55c2bd5ea9"),
			passphrase: []byte("pass"),
		},
		{
			name:       "HDGood",
			data:       _byte("0x018eb13b2e650eb6670d741bf855a3879b27daf5c0a21cbd8992ef8a8aeee60b917296ea853cf9f95a86537f6a5c010b639a7143992c05c2d7097bc5b05e57bfa61b27ba1e87add0c431f605cdd8d48ee6e01075ef4fbc46eb128e37fcf6d7ecd8f593c24af5db5d6728d80fbe47ee548c5ef93dc1b1d0dff820ab9ba2edb097c4a2f60ad3c991a6279a6ed196fbc81778e54bdd055befdc5d2b877b06f49dd3000c605f33d7297416dfd2e67dbc84b10cc7427f51fe2f76c04c0f9ca172e6f6d1f033fd0591ad76fa0122af9fc4b4eea99e5fc8df9ac28c55ea7edba3eab98c743e16130d90c00b5cc7e2fbb87d62c14cf779417476b39513d60c16257ca05213527dc79df35c38c901886c5692e5430923269b6ce35ed350ceb2d0bcc76a0414977042d8a7e62aa027c985af7c7212dbae5f2d175939d0b1b169268cfa2cd67e6e398100852f9918b85083b9a8a2343026f063097e459f2ccaf1c5fe67d7927c88b8c1a84c9d99408d1811b358cfb0064ccbaec5be4bbad2cc9f1500e33b9b49865c05a2e7b41af5f1f6c9cb3e37361ee9a8f01550886ab8c02192674963aefebb45729e49c7a183cac8e354c5b4a6a6c3f16489bcf12a9255d18e5b088c055a1f8ef2a2e67c86b97b555222ed616dfdb33fd42e846a0dd28c37fb49228d0789903b96e550dcf5ca40e2f9d1cdfda1de63e0a64da6755e64b9bca71579b46924e191bdda57ba169830337463d03c44c5bd7025b3f628c8abc7a9b3ef6a5c68145509d606a339be65be993aa443910eaf5554843046d775f5755d28181aaa0b8a0da7807ab4d9eec9a459ac8fda50fcce3d698a6333ec9cfb2f45713b64d0f964b2dcbe6a02fc35c4e186dca1b41e1bff97d9f340c869d88dabd486b983387109df0b741a8e89f97423e2bd6c16be4cd92be3f028400f41b8b02cdbdc6e412da1acb43ed8a63f451f56d23cfd578940db047f27ae5fbc44c73ea1082a2723e7a8990c567c8ec0588572f5be139a6128a35c3002de55ef454004b1055199a405d2de693195b6b45b0118fe621edd88bcf159ae05d101d4cbb1c493884049af5cba341edd2f4bae2e9b711bafccdbaa66bc421285bd1f5bfe44d331b69c567c9a558e42f6ed474f4d6d5c8404b81ff3be0b2c20da28839314edbb9a6ab4bf4c68c3586014b7310dbb206a0d6cc4b099da2b54595195219f3fcff4541301f312ff9d4872e9de3c5805b137363b123fbe8919007f828a1c45bb60bec7d05e67057784e9673e7e60218c66a97168613a662332546dd0b2979e787dcef3244cf491e7a616f9c08c45c308ede0480e08c2580a949ecc0eda1a3b76927ea25a894da4ba1205b03c6abb91e647d06933cb5206d5860adf4a857d33aa0524ea2347350dc54b2cbb5b26763d95a2664bfe8fdb2099857ac29b62bed0b1390817262423dd481330e46dbab632b873c2efa60c691094f46118a88ea365c5b0d395679eda134331a26ea741ce63d44f4ef3d6fba7696e3599b6411adce5b31fe80fb30705cb81579994ca5fad64049b046b1db585767244f3aeefdac5add2440deca2fb938c4954d2b7d51bf90ad9f4183e86b400335bc488b83aa3f30ed1ef6890150e67b090c583717b9990c1a757585ddaeafee11ad35a4a8f4c9700088e6ce9414c6fbd7ca4041d4252a278e1cd36720ea41bab0890244d1a504a01737a180a6c5832cba0b8a5034fc6efb05ab52d12a82344ae5d8a3d2cb123747f92f06cad8a6e24588daa85abe7a06a2379451daf715c3e8284c3fab2b5164daf314c460de285b1c3d69ca87fc86c16b0ba5a2b8288caf46badafa423c6d185f86755236870d08cd883b0748e1a313ee279c1781592a333720d57c7c3b316c551fb47b5c80fbd17d11d44ff4c07d9b293846a3ed0123eada3b439ba0f738a388e80d0e618e144371e5fc9764b259a14e3bd4fc33e3fa0178a209c2a7fc2e4378878004e19b39cad7ed6d2d9301e395f9f934803aa5eb8e14173bf8bec3c7ff4488fd9daf18f96e7ea40eda2c98502bfe8"),
			passphrase: []byte("pass"),
		},
		{
			name:       "DistributedGood",
			data:       _byte("0x0158d6cc7f70f3c7488cf042c8e5c3e4c50c0ecc14db29a04f47d6e8697511c671b7591e3909345f12ceb8d3d1e11022b16396f0d806f20dbc23071b304274035977c147639a4b808488a35be9dc6f6ea6c9e37b8588213e1af1c8b67388f33ba63183b2eda4218f97b3855bbc05b7d8c7c1ed7f82043c6b9641ca54590708315059b8b9cb64b025dabecce10c1f11eb27aa171c29ba00e55813dfa4e071dadc72a2f3dc833e22c2ff1ce945201dbce3a2fa7f722708043d2af42afd002149f0aeed5e54b774f1f773b8b1411f8ec54ff788c5eed5decf66a736cb"),
			passphrase: []byte("pass"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w, err := wallet.ImportWallet(test.data, test.passphrase)
			if test.err == "" {
				require.NoError(t, err)
				require.NotNil(t, w)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}
