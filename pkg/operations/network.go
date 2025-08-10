package operations

type NetworkOps struct {
	Client ClientWrapper
}

func (n *NetworkOps) ExportNetwork(name string, attachments bool, filepath string) error {
	params := map[string]interface{}{
		"name":        name,
		"attachments": attachments,
		"filepath":    filepath,
	}
	return n.Client.Export("/network/operations/exportNetwork", filepath, params)
}

func (n *NetworkOps) ImportNetwork(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
	}
	return n.Client.Import("/network/operations/importNetwork", filename, params)
}

func (n *NetworkOps) Search(searchString, userid, clazz, sortorder, sort string, limit, offset int) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"userid":       userid,
		"class":        clazz,
		"sortorder":    sortorder,
		"sort":         sort,
		"limit":        limit,
		"offset":       offset,
	}
	return n.Client.Post("/network/operations/search", params)
}

func (n *NetworkOps) Load(template string) (interface{}, error) {
    params := map[string]interface{}{
        "template": template,
    }
    return n.Client.Post("/network/operations/load", params)
}