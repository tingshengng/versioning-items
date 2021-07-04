package versioningservice

import (
	searchitem "PoC/SearchVersion/searchitemDB"
	versioning "PoC/SearchVersion/versioningDB"
	"fmt"
	"sort"
)

const (
	ADD    = "ADD"
	EDIT   = "EDIT"
	DELETE = "DELETE"
)

type Modification struct {
	Action       string
	TargetItemID int
	NewName      string
	OldItemID    int
}

// SearchItemSorter implements the sort interface so that we could sort the SearchItems based on the ItemName
type SearchItemSorter []*versioning.VersioningModel

func (s SearchItemSorter) Len() int           { return len(s) }
func (s SearchItemSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SearchItemSorter) Less(i, j int) bool { return s[i].ItemName < s[j].ItemName }

func NewVersioningService() *VersioningService {
	return &VersioningService{
		versionDAO:    versioning.NewVersioningDB(),
		searchitemDAO: searchitem.NewSearchItemDB(),
	}
}

type VersioningService struct {
	versionDAO    *versioning.VersioningDB
	searchitemDAO *searchitem.SearchItemDB
}

func (v *VersioningService) CreateSnapshot(baseVersion int, modifications []*Modification) {
	currentVersion := v.versionDAO.LatestVersion + 1

	// If the changes are not based on latest version, we need to temporary "revert" back to the baseVersion state
	// We need to check outdated items from the base version and store temporaily.
	// We also need to expire all the changes made after the base version.
	outdatedItems := map[int]*versioning.VersioningModel{}
	if baseVersion < v.versionDAO.LatestVersion {
		outdatedItems = v.versionDAO.GetOutdatedItemsByVersion(baseVersion)
		v.versionDAO.ExpireItemsFromVersion(baseVersion, currentVersion)
	}

	for _, mod := range modifications {
		switch mod.Action {
		case ADD:
			v.AddModification(currentVersion, mod)
		case EDIT:
			v.EditModification(currentVersion, baseVersion, mod)
		case DELETE:
			v.DeleteModification(currentVersion, baseVersion, mod)
		}
		delete(outdatedItems, mod.TargetItemID)
	}

	for _, outdatedItem := range outdatedItems {
		newVersioningItem := &versioning.VersioningModel{
			ItemID:   outdatedItem.ItemID,
			ItemName: outdatedItem.ItemName,
			Created:  currentVersion,
		}
		v.versionDAO.InsertRow(newVersioningItem)
	}
}

func (v *VersioningService) GetItemsByVersion(version int) {
	versionSearchItems := v.versionDAO.GetRowsByVersion(version)
	sort.Sort(SearchItemSorter(versionSearchItems))
	fmt.Println("Search items in Version: ", version)
	for _, item := range versionSearchItems {
		fmt.Printf("%d:\t%s\n", item.ItemID, item.ItemName)
	}
}

func (v *VersioningService) AddModification(curVersion int, mod *Modification) {
	newItem := &searchitem.SearchItemModel{
		Name: mod.NewName,
	}
	newItemID := v.searchitemDAO.InsertRow(newItem)
	newVersioningItem := &versioning.VersioningModel{
		ItemID:   newItemID,
		ItemName: newItem.Name,
		Created:  curVersion,
	}
	v.versionDAO.InsertRow(newVersioningItem)
}

func (v *VersioningService) EditModification(curVersion, baseVersion int, mod *Modification) {
	newItem := &searchitem.SearchItemModel{
		Name: mod.NewName,
	}
	newItemID := v.searchitemDAO.InsertRow(newItem)
	newVersioningItem := &versioning.VersioningModel{
		ItemID:   newItemID,
		ItemName: newItem.Name,
		Created:  curVersion,
	}
	v.versionDAO.InsertRow(newVersioningItem)

	oldVersioningItem := v.versionDAO.GetRowByItemIDAndVersion(mod.TargetItemID, baseVersion)
	if oldVersioningItem.Expired == 0 {
		oldVersioningItem.Expired = curVersion
		v.versionDAO.UpdateRowByVersioningModel(oldVersioningItem)
	}
}

func (v *VersioningService) DeleteModification(curVersion, baseVersion int, mod *Modification) {
	oldVersioningItem := v.versionDAO.GetRowByItemIDAndVersion(mod.TargetItemID, baseVersion)
	if oldVersioningItem.Expired == 0 {
		oldVersioningItem.Expired = curVersion
		v.versionDAO.UpdateRowByVersioningModel(oldVersioningItem)
	}
}

func (v *VersioningService) ListDBs() {
	fmt.Println("SearchItemDB", "\n")
	searchItems := v.searchitemDAO.ListItems()
	for _, s := range searchItems {
		fmt.Printf("%+v\n", s)
	}
	fmt.Println("\n", "VersioningDB", "\n")
	for _, ver := range v.versionDAO.ListItems() {
		fmt.Printf("%+v\n", ver)
	}
}

// InitData seeds data into the DB, it is equivalent to seeding all the search items we currently have into DB
func (v *VersioningService) InitData() {
	searchItems := []*searchitem.SearchItemModel{
		{
			Name: "a",
		},
		{
			Name: "b",
		},
		{
			Name: "c",
		},
		{
			Name: "d",
		},
	}
	v.searchitemDAO.InsertRows(searchItems)
	versioningData := []*versioning.VersioningModel{
		{
			ItemID:   1,
			ItemName: "a",
			Created:  1,
			Expired:  0,
		},
		{
			ItemID:   2,
			ItemName: "b",
			Created:  1,
			Expired:  0,
		},
		{
			ItemID:   3,
			ItemName: "c",
			Created:  1,
			Expired:  0,
		},
		{
			ItemID:   4,
			ItemName: "d",
			Created:  1,
			Expired:  0,
		},
	}
	v.versionDAO.InsertRows(versioningData)
}
