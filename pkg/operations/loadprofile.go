package operations

type LoadProfileOps struct {
	Client ClientWrapper
}

func (l *LoadProfileOps) CreateNew(loadProfile string) (interface{}, error) {
	params := map[string]interface{}{
		"loadProfile": loadProfile,
	}
	return l.Client.Post("/loadprofile/operations/createNew", params)
}

func (l *LoadProfileOps) Delete(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return l.Client.Post("/loadprofile/operations/delete", params)
}

func (l *LoadProfileOps) Load(template string) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
	}
	return l.Client.Post("/loadprofile/operations/load", params)
}

func (l *LoadProfileOps) Save() (interface{}, error) {
	params := map[string]interface{}{}
	return l.Client.Post("/loadprofile/operations/save", params)
}

func (l *LoadProfileOps) SaveAs(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return l.Client.Post("/loadprofile/operations/saveAs", params)
}

func (l *LoadProfileOps) Search(searchString, limit, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return l.Client.Post("/loadprofile/operations/search", params)
}

func (l *LoadProfileOps) SearchDynamic(searchString string, limit string, sort string, sortorder string, offset string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
		"offset":       offset,
	}
	return l.Client.Post("/loadprofile/operations/searchDynamic", params)
}