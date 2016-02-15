// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"github.com/FactomProject/factomd/common/directoryBlock"
	"github.com/FactomProject/factomd/common/entryCreditBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/log"
	
	"github.com/FactomProject/factomd/util"
	"testing"
)

var _ = fmt.Print
var _ = log.Print
var _ = util.ReadConfig

// We need to create states with a testing database!  Need to add
// back the mapping database for testing!  Until then, don't run
// tests by default that require us to delete or manipulate the 
// database and rebuild Factom's past.
var runDBTests = true


func Test_DBState1(t *testing.T) {
	if !runDBTests {
		return
	}
	
	state := GetState()

	var prev interfaces.IDirectoryBlock	// First call gets a nil, rest the previous DirectoryBlock
	
	var i uint32
	for i = 0; i < 10; i++ {

		// p ends up with the DirectoryBlock or nil.  All's good.
		p, _ := prev.(*directoryBlock.DirectoryBlock)
		if i > 0 && p == nil {
			t.Error("Should not fail to have a previous lbock")
		}
		dblk := directoryBlock.NewDirectoryBlock(i,p)
		prev = dblk
		ablk := state.NewAdminBlock(i)
		eblk := entryCreditBlock.NewECBlock()
		fblk := state.GetFactoidState().GetCurrentBlock()
		state.GetFactoidState().ProcessEndOfBlock(state)
		
		state.DBStates.NewDBState(true, dblk, ablk, fblk, eblk)
		state.DBStates.Process()
		
		h := dblk.GetHeader().GetDBHeight()
		if i != h {
			t.Errorf("Height error.  Expecting %d and got %d", i, h)
		}
		if state.DBHeight != i {
			t.Errorf("DBHeight error.  Expecting %d and got %d", i, state.DBHeight)
		}
	}

	dblks := make([]interfaces.IDirectoryBlock, 0)

	for j := uint32(0); j < i; j++ {
		dblk, _ := state.DB.FetchDBlockByHeight(j)
		if dblk == nil {
			fmt.Println("last dblk found:", j)
			break
		}
		dblks = append(dblks, dblk)
	}

	/*
		 * ecblkHash := dblks[len(dblks)-1].DBEntries[1].KeyMR

			i := 0
			for i = 0; ecblkHash != nil; i++ {
				fmt.Printf(" %x\n",ecblkHash.Bytes())
				ecblk, _ := db.FetchECBlockByHash(ecblkHash)
				if ecblk == nil {
					break
				}
				ecblkHash = ecblk.Header.PrevHeaderHash
			}
			fmt.Println ("End found after",i,"ec blocks")
		}
	*/
}
