package main

import (
	v "PoC/SearchVersion/versioningService"
	"fmt"
)

func main() {
	vs := v.NewVersioningService()

	// initialize version 1 data into DB, this can be treat as the inital search items data we currently have
	vs.InitData()
	vs.GetItemsByVersion(1)

	v2basev1 := []*v.Modification{
		{
			Action:       v.EDIT,
			TargetItemID: 1,
			NewName:      "a2",
		},
		{
			Action:       v.EDIT,
			TargetItemID: 4,
			NewName:      "d2",
		},
		{
			Action:  v.ADD,
			NewName: "e",
		},
	}
	vs.CreateSnapshot(1, v2basev1)

	v3basev1 := []*v.Modification{

		{
			Action:       v.DELETE,
			TargetItemID: 4,
			NewName:      "d",
		},
		{
			Action:  v.ADD,
			NewName: "f",
		},
	}
	vs.CreateSnapshot(1, v3basev1)

	v4basev2 := []*v.Modification{
		{
			Action:       v.DELETE,
			TargetItemID: 3,
		},
		{
			Action:       v.EDIT,
			TargetItemID: 7,
			NewName:      "e2",
		},
	}
	vs.CreateSnapshot(2, v4basev2)

	v5basev3 := []*v.Modification{
		{
			Action:       v.EDIT,
			TargetItemID: 2,
			NewName:      "b2",
		},
		{
			Action:       v.EDIT,
			TargetItemID: 3,
			NewName:      "c3",
		},
		{
			Action:       v.EDIT,
			TargetItemID: 8,
			NewName:      "f2",
		},
		{
			Action:  v.ADD,
			NewName: "h",
		},
	}
	vs.CreateSnapshot(3, v5basev3)

	v6basev2 := []*v.Modification{
		{
			Action:       v.DELETE,
			TargetItemID: 5,
		},
		{
			Action:       v.EDIT,
			TargetItemID: 3,
			NewName:      "c2",
		},
		{
			Action:       v.EDIT,
			TargetItemID: 7,
			NewName:      "e3",
		},
		{
			Action:  v.ADD,
			NewName: "g",
		},
	}
	vs.CreateSnapshot(2, v6basev2)

	v7basev6 := []*v.Modification{
		{
			Action:       v.EDIT,
			TargetItemID: 2,
			NewName:      "b10",
		},
		{
			Action:       v.DELETE,
			TargetItemID: 16,
		},
		{
			Action:  v.ADD,
			NewName: "i",
		},
	}
	vs.CreateSnapshot(6, v7basev6)

	v8basev6 := []*v.Modification{
		{
			Action:       v.EDIT,
			TargetItemID: 2,
			NewName:      "b11",
		},
		{
			Action:       v.DELETE,
			TargetItemID: 15,
		},
	}
	vs.CreateSnapshot(6, v8basev6)

	fmt.Println("\n", "-------- Get SearchItems by Version --------", "\n")
	vs.GetItemsByVersion(1)
	vs.GetItemsByVersion(2)
	vs.GetItemsByVersion(3)
	vs.GetItemsByVersion(4)
	vs.GetItemsByVersion(5)
	vs.GetItemsByVersion(6)
	vs.GetItemsByVersion(7)
	vs.GetItemsByVersion(8)

	fmt.Println("\n", "-------- List DB data --------", "\n")
	vs.ListDBs()
}
