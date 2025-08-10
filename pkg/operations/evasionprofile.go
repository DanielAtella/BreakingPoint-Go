package operations

type EvasionProfileOps struct {
	Client ClientWrapper
}

func (e *EvasionProfileOps) GetStrikeOptions() (interface{}, error) {
	return e.Client.Post("/evasionProfile/StrikeOptions/operations/getStrikeOptions", map[string]interface{}{})
}

func (e *EvasionProfileOps) Delete(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return e.Client.Post("/evasionProfile/operations/delete", params)
}

func (e *EvasionProfileOps) Load(template string) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
	}
	return e.Client.Post("/evasionProfile/operations/load", params)
}

func (e *EvasionProfileOps) New(template *string) (interface{}, error) {
	params := map[string]interface{}{}
	if template != nil {
		params["template"] = *template
	}
	return e.Client.Post("/evasionProfile/operations/new", params)
}

func (e *EvasionProfileOps) Save(name *string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"force": force,
	}
	if name != nil {
		params["name"] = *name
	}
	return e.Client.Post("/evasionProfile/operations/save", params)
}

func (e *EvasionProfileOps) SaveAs(name string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":  name,
		"force": force,
	}
	return e.Client.Post("/evasionProfile/operations/saveAs", params)
}

func (e *EvasionProfileOps) Search(searchString string, limit string, sort string, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return e.Client.Post("/evasionProfile/operations/search", params)
}