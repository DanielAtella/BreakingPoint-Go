package operations

type CaptureOps struct {
	Client ClientWrapper
}

func (c *CaptureOps) ImportCapture(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":  name,
		"force": force,
	}
	return c.Client.Import("/capture/operations/importCapture", filename, params)
}

func (c *CaptureOps) Search(searchString, limit, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return c.Client.Post("/capture/operations/search", params)
}
