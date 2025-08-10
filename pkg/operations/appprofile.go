package operations

type AppProfileOps struct {
	Client ClientWrapper
}

func (a *AppProfileOps) Add(add []map[string]interface{}) (interface{}, error) {
	return a.Client.Post("/appProfile/operations/add", map[string]interface{}{"add": add})
}

func (a *AppProfileOps) Delete(name string) (interface{}, error) {
	return a.Client.Post("/appProfile/operations/delete", map[string]interface{}{"name": name})
}

func (a *AppProfileOps) ExportAppProfile(name string, attachments bool, filepath string) error {
	params := map[string]interface{}{
		"name":        name,
		"attachments": attachments,
		"filepath":    filepath,
	}
	return a.Client.Export("/appProfile/operations/exportAppProfile", filepath, params)
}

func (a *AppProfileOps) ImportAppProfile(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
	}
	return a.Client.Import("/appProfile/operations/importAppProfile", filename, params)
}

func (a *AppProfileOps) Search(searchString, limit, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return a.Client.Post("/appProfile/operations/search", params)
}