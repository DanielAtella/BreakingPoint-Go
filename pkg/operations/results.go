package operations

type ResultsOps struct {
	Client ClientWrapper
}

func (r *ResultsOps) GetGroups(name string, dynamicEnums bool, includeOutputs bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":           name,
		"dynamicEnums":   dynamicEnums,
		"includeOutputs": includeOutputs,
	}
	return r.Client.Post("/results/operations/getGroups", params)
}

func (r *ResultsOps) GetHistoricalResultSize(runid int, componentid string, group string) (interface{}, error) {
	params := map[string]interface{}{
		"runid":       runid,
		"componentid": componentid,
		"group":       group,
	}
	return r.Client.Post("/results/operations/getHistoricalResultSize", params)
}

func (r *ResultsOps) GetHistoricalSeries(runid int, componentid string, dataindex int, group string) (interface{}, error) {
	params := map[string]interface{}{
		"runid":       runid,
		"componentid": componentid,
		"dataindex":   dataindex,
		"group":       group,
	}
	return r.Client.Post("/results/operations/getHistoricalSeries", params)
}