/*
 * Copyright (C) 2019 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package pingpong

import (
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mysteriumnetwork/node/core/storage/boltdb"
	"github.com/mysteriumnetwork/node/eventbus"
	"github.com/mysteriumnetwork/node/identity"
	"github.com/stretchr/testify/assert"
)

func TestConsumerTotalStorage(t *testing.T) {
	dir, err := os.MkdirTemp("", "consumerTotalsStorageTest")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	bolt, err := boltdb.NewStorage(dir)
	assert.NoError(t, err)
	defer bolt.Close()

	consumerTotalsStorage := NewConsumerTotalsStorage(bolt, eventbus.New())

	channelAddress := identity.FromAddress("someAddress")
	hermesAddress := common.HexToAddress("someOtherAddress")
	var amount = big.NewInt(12)

	// check if errors are wrapped correctly
	_, err = consumerTotalsStorage.Get(1, channelAddress, hermesAddress)
	assert.Equal(t, ErrNotFound, err)

	// store and check that total is stored correctly
	err = consumerTotalsStorage.Store(1, channelAddress, hermesAddress, amount)
	assert.NoError(t, err)

	a, err := consumerTotalsStorage.Get(1, channelAddress, hermesAddress)
	assert.NoError(t, err)
	assert.Equal(t, amount, a)

	var newAmount = big.NewInt(123)
	// overwrite the amount, check if it is overwritten
	err = consumerTotalsStorage.Store(1, channelAddress, hermesAddress, newAmount)
	assert.NoError(t, err)

	a, err = consumerTotalsStorage.Get(1, channelAddress, hermesAddress)
	assert.NoError(t, err)
	assert.EqualValues(t, newAmount, a)

	someOtherChannel := identity.FromAddress("someOtherChannel")
	// store two amounts, check if both are gotten correctly
	err = consumerTotalsStorage.Store(1, someOtherChannel, hermesAddress, amount)
	assert.NoError(t, err)

	a, err = consumerTotalsStorage.Get(1, channelAddress, hermesAddress)
	assert.NoError(t, err)
	assert.EqualValues(t, newAmount, a)

	a, err = consumerTotalsStorage.Get(1, someOtherChannel, hermesAddress)
	assert.NoError(t, err)
	assert.EqualValues(t, amount, a)

	var addAmount = big.NewInt(10)
	err = consumerTotalsStorage.Add(1, someOtherChannel, hermesAddress, addAmount)
	assert.NoError(t, err)

	a, err = consumerTotalsStorage.Get(1, someOtherChannel, hermesAddress)
	assert.NoError(t, err)
	assert.EqualValues(t, new(big.Int).Add(amount, addAmount), a)
}
