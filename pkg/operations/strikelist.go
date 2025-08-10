package operations

type StrikeListOps struct {
	Client ClientWrapper
}

func (s *StrikeListOps) Add(strikes []map[string]interface{}, validate bool, toList *string) (interface{}, error) {
	params := map[string]interface{}{
		"strike":   strikes,
		"validate": validate,
	}
	if toList != nil {
		params["toList"] = *toList
	}
	return s.Client.Post("/strikeList/operations/add", params)
}

func (s *StrikeListOps) Delete(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return s.Client.Post("/strikeList/operations/delete", params)
}

func (s *StrikeListOps) ExportStrikeList(name, filepath string) error {
	params := map[string]interface{}{
		"name":     name,
		"filepath": filepath,
	}
	return s.Client.Export("/strikeList/operations/exportStrikeList", filepath, params)
}

func (s *StrikeListOps) ImportStrikeList(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
	}
	return s.Client.Import("/strikeList/operations/importStrikeList", filename, params)
}

func (s *StrikeListOps) Load(template string) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
	}
	return s.Client.Post("/strikeList/operations/load", params)
}

func (s *StrikeListOps) New(template *string) (interface{}, error) {
	params := map[string]interface{}{}
	if template != nil {
		params["template"] = *template
	}
	return s.Client.Post("/strikeList/operations/new", params)
}

func (s *StrikeListOps) Remove(strikes []map[string]interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"strike": strikes,
	}
	return s.Client.Post("/strikeList/operations/remove", params)
}

func (s *StrikeListOps) Save(name *string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"force": force,
	}
	if name != nil {
		params["name"] = *name
	}
	return s.Client.Post("/strikeList/operations/save", params)
}

func (s *StrikeListOps) SaveAs(name string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":  name,
		"force": force,
	}
	return s.Client.Post("/strikeList/operations/saveAs", params)
}

func (s *StrikeListOps) Search(searchString string, limit int, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return s.Client.Post("/strikeList/operations/search", params)
}