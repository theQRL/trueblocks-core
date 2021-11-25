/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
package blockRange

import (
	"fmt"

	"encoding/json"
	"testing"
)

func TestPointToPointTypeBlock(t *testing.T) {
	point := &Point{Block: 100}
	result := getPointType(point)

	if result != BlockRangeBlockNumber {
		t.Error("Bad point type returned")
	}
}

func TestPointToPointTypeSpecial(t *testing.T) {
	point := &Point{Special: "london"}
	result := getPointType(point)

	if result != BlockRangeSpecial {
		t.Error("Bad point type returned")
	}
}

func TestPointToPointTypeDate(t *testing.T) {
	point := &Point{Date: "2021-10-03"}
	result := getPointType(point)

	if result != BlockRangeDate {
		t.Error("Bad point type returned")
	}
}

func TestModifierToModifierTypeStep(t *testing.T) {
	modifier := &Modifier{Step: 15}
	result := getModifierType(modifier)

	if result != BlockRangeStep {
		t.Error("Bad modifier type returned")
	}
}

func TestModifierToModifierTypePeriod(t *testing.T) {
	modifier := &Modifier{Period: "daily"}
	result := getModifierType(modifier)

	if result != BlockRangePeriod {
		t.Error("Bad modifier type returned")
	}
}

func TestNewBlocks(t *testing.T) {
	blockRange, err := New("10-1000:10")
	if err != nil {
		t.Error(err)
	}

	if blockRange.StartType != BlockRangeBlockNumber {
		t.Error("StartType is not block number")
	}

	if blockRange.Start.Block != 10 {
		t.Errorf("Wrong start")
	}

	if blockRange.EndType != BlockRangeBlockNumber {
		t.Error("EndType is not block number")
	}

	if blockRange.End.Block != 1000 {
		t.Error("Wrong end")
	}

	if blockRange.ModifierType != BlockRangeStep {
		t.Error("ModifierType is not step")
	}

	if blockRange.Modifier.Step != 10 {
		t.Error("Wrong modifier")
	}
}

func TestNewSpecial(t *testing.T) {
	blockRange, err := New("london:weekly")

	if err != nil {
		t.Error(err)
	}

	if blockRange.StartType != BlockRangeSpecial {
		t.Error("StartType is not special")
	}

	if blockRange.Start.Special != "london" {
		t.Errorf("Wrong start")
	}

	if blockRange.EndType != BlockRangeNotDefined {
		t.Error("EndType is not notdefined")
	}

	if blockRange.ModifierType != BlockRangePeriod {
		t.Error("ModifierType is not period")
	}

	if blockRange.Modifier.Period != "weekly" {
		t.Error("Wrong modifier")
	}
}

func TestHandleParserErrors(t *testing.T) {
	_, modifierErr := New("10-100:biweekly")

	if me, ok := modifierErr.(*WrongModifierError); ok {
		if me.Token != "biweekly" {
			t.Errorf("Wrong token: %s", me.Token)
		}
	} else {
		t.Error("Returned error is not WrongModifier")
		t.Error(modifierErr)
	}
}

func TestBlockRange_UnmarshalJSON(t *testing.T) {
	type SomeRecord struct {
		Blocks BlockRange `json:"blocks"`
	}

	var record SomeRecord
	source := []byte(`{"blocks":"000000000-10567003"}`)

	err := json.Unmarshal(source, &record)
	if err != nil {
		t.Error(err)
	}

	if record.Blocks.StartType != BlockRangeBlockNumber {
		t.Errorf("Wrong StartType %d", record.Blocks.StartType)
	}

	if record.Blocks.EndType != BlockRangeBlockNumber {
		t.Errorf("Wrong EndType %d", record.Blocks.EndType)
	}

	if record.Blocks.Start.Block != uint(0) {
		t.Error("Wrong start value")
	}

	if record.Blocks.End.Block != uint(10567003) {
		t.Errorf("Wrong end value %d", record.Blocks.End.Block)
	}
}

func TestToString(t *testing.T) {
	br, err := New("1234")
	if err != nil {
		t.Errorf("Could not parse block")
	}
	expected := "{\"StartType\":0,\"Start\":{\"Block\":1234,\"Date\":\"\",\"Special\":\"\"},\"EndType\":5,\"End\":{\"Block\":0,\"Date\":\"\",\"Special\":\"\"},\"ModifierType\":5,\"Modifier\":{\"Step\":0,\"Period\":\"\"}}\n"
	got := fmt.Sprintf("%s\n", br.MarshalJSON())
	if got != expected {
		t.Errorf("String printer for blockRange not equal to expected")
	}
}
