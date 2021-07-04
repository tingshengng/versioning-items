package searchitemdb

type SearchItemModel struct {
	ID   int
	Name string
}

type SearchItemDB struct {
	Rows []*SearchItemModel
}

func NewSearchItemDB() *SearchItemDB {
	return &SearchItemDB{}
}

func (s *SearchItemDB) InsertRow(item *SearchItemModel) int {
	item.ID = len(s.Rows) + 1
	s.Rows = append(s.Rows, item)
	return item.ID
}

func (s *SearchItemDB) InsertRows(items []*SearchItemModel) {
	for _, item := range items {
		item.ID = len(s.Rows) + 1
		s.Rows = append(s.Rows, item)
	}
}

func (s *SearchItemDB) GetByID(id int) *SearchItemModel {
	for _, row := range s.Rows {
		if row.ID == id {
			return row
		}
	}
	return nil
}

func (s *SearchItemDB) ListItems() []*SearchItemModel {
	return s.Rows
}
