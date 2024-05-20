// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package abi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/theQRL/go-zond/accounts/abi"
	"github.com/theQRL/go-zond/common"
)

func Test_Abi_Decode(t *testing.T) {
	const definition = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isBar","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		panic(err)
	}
	out, err := abi.Pack("isBar", common.HexToAddress("01"))
	if err != nil {
		t.Error(err) // panic(err)
	}
	fmt.Printf("%x\n", out)

	i, _ := abi.Unpack("isBar", out[4:])
	fmt.Printf("%v\n", i)
}
