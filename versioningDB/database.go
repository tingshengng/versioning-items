package versioningdb

type VersioningModel struct {
	ID       int
	ItemID   int
	ItemName string
	Created  int
	Expired  int
}

type VersioningDB struct {
	Rows          []*VersioningModel
	LatestVersion int
}

func NewVersioningDB() *VersioningDB {
	return &VersioningDB{}
}

func (m *VersioningDB) InsertRow(row *VersioningModel) int {
	if row.Created > m.LatestVersion {
		m.LatestVersion = row.Created
	}
	row.ID = len(m.Rows) + 1
	m.Rows = append(m.Rows, row)
	return row.ID
}

func (m *VersioningDB) InsertRows(rows []*VersioningModel) {
	for _, row := range rows {
		row.ID = len(m.Rows) + 1
		if row.Created > m.LatestVersion {
			m.LatestVersion = row.Created
		}
		m.Rows = append(m.Rows, row)
	}
}

func (m *VersioningDB) GetRowsByVersion(version int) []*VersioningModel {
	res := []*VersioningModel{}
	for _, row := range m.Rows {
		if row.Created <= version && (row.Expired > version || row.Expired == 0) {
			res = append(res, row)
		}
	}
	return res
}

func (m *VersioningDB) GetOutdatedItemsByVersion(version int) map[int]*VersioningModel {
	versionItems := m.GetRowsByVersion(version)

	res := map[int]*VersioningModel{}
	for _, item := range versionItems {
		if item.Expired != 0 {
			res[item.ItemID] = item
		}
	}
	return res
}

func (m *VersioningDB) GetRowByItemIDAndVersion(itemID, version int) *VersioningModel {
	versionItems := m.GetRowsByVersion(version)
	for _, item := range versionItems {
		if item.ItemID == itemID {
			return item
		}
	}
	return nil
}

func (m *VersioningDB) UpdateRowByVersioningModel(row *VersioningModel) {
	for i, r := range m.Rows {
		if r.ID == row.ID {
			m.Rows[i] = row
		}
	}
}

func (m *VersioningDB) ExpireItemsFromVersion(baseVersion, currentVersion int) {
	for _, row := range m.Rows {
		if row.Created > baseVersion && row.Expired == 0 {
			row.Expired = currentVersion
		}
	}
}

func (m *VersioningDB) ListItems() []*VersioningModel {
	return m.Rows
}
