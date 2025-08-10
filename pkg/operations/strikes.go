package operations

type StrikesOps struct {
	Client ClientWrapper
}

func (s *StrikesOps) Search(searchString string, limit int, sort string, sortorder string, offset int) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
		"offset":       offset,
	}
	return s.Client.Post("/strikes/operations/search", params)
}