/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package crypto

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type EthCrypto struct {
	homestead bool
	ks        *keystore.KeyStore
}

//签名验证
func (c EthCrypto) ValidateSignatureValues(v byte, r, s []byte) bool {
	return ethCrypto.ValidateSignatureValues(v, new(big.Int).SetBytes(r), new(big.Int).SetBytes(s), c.homestead)
}

//生成hash
func (c EthCrypto) GenerateHash(data ...[]byte) []byte {
	return ethCrypto.Keccak256(data...)
}

//签名回复到地址
func (c EthCrypto) SigToAddress(hashPre, sig []byte) ([]byte, error) {
	hash := c.GenerateHash([]byte("\x19Ethereum Signed Message:\n32"), hashPre)
	pubKey, err := ethCrypto.SigToPub(hash, sig)
	if nil != err {
		return nil, err
	} else {
		return ethCrypto.PubkeyToAddress(*pubKey).Bytes(), nil
	}
}

func (c EthCrypto) VRSToSig(v byte, r, s []byte) (sig []byte, err error) {
	sig = make([]byte, 65)
	vUint8 := uint8(v)
	if vUint8 >= 27 {
		vUint8 -= 27
	}
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = byte(vUint8)
	return sig, nil
}

func (c EthCrypto) Sign(hashPre []byte, signer accounts.Account) ([]byte, error) {
	hash := c.GenerateHash([]byte("\x19Ethereum Signed Message:\n32"), hashPre)
	return c.ks.SignHash(signer, hash)
}

func (c EthCrypto) SigToVRS(sig []byte) (v byte, r []byte, s []byte) {
	r = make([]byte, 32)
	s = make([]byte, 32)
	v = byte(uint8(sig[64]) + uint8(27))
	copy(r, sig[0:32])
	copy(s, sig[32:64])
	return v, r, s
}

func (c EthCrypto) UnlockAccount(acc accounts.Account, passphrase string) error {
	return c.ks.Unlock(acc, passphrase)
}
func (c EthCrypto) SignTx(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return c.ks.SignTx(a, tx, chainID)
}

func NewCrypto(homestead bool, ks *keystore.KeyStore) EthCrypto {
	return EthCrypto{homestead: homestead, ks: ks}
}
