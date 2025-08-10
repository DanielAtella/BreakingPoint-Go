package operations

type ReportsOps struct {
	Client ClientWrapper
}

func (r *ReportsOps) Delete(runid int) (interface{}, error) {
	params := map[string]interface{}{
		"runid": runid,
	}
	return r.Client.Post("/reports/operations/delete", params)
}

func (r *ReportsOps) ExportReport(filepath string, runid int, reportType string, sectionIds string, dataType string) error {
	params := map[string]interface{}{
		"filepath":   filepath,
		"runid":      runid,
		"reportType": reportType,
		"sectionIds": sectionIds,
		"dataType":   dataType,
	}
	return r.Client.Export("/reports/operations/exportReport", filepath, params)
}

func (r *ReportsOps) GetReportContents(runid int, getTableOfContents bool) (interface{}, error) {
	params := map[string]interface{}{
		"runid":              runid,
		"getTableOfContents": getTableOfContents,
	}
	return r.Client.Post("/reports/operations/getReportContents", params)
}

func (r *ReportsOps) GetReportTable(runid int, sectionId string) (interface{}, error) {
	params := map[string]interface{}{
		"runid":     runid,
		"sectionId": sectionId,
	}
	return r.Client.Post("/reports/operations/getReportTable", params)
}

func (r *ReportsOps) Search(searchString, limit, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return r.Client.Post("/reports/operations/search", params)
}