package operations

type StatisticsOps struct {
	Client ClientWrapper
}

func (s *StatisticsOps) GetStatsDefinitions() (interface{}, error) {
	return s.Client.Post("/statistics/operations/getStatsDefinitions", map[string]interface{}{})
}

func (s *StatisticsOps) GetStatisticsByType(statType string) (interface{}, error) {
	params := map[string]interface{}{
		"type": statType,
	}
	return s.Client.Post("/statistics/operations/getStatisticsByType", params)
}

func (s *StatisticsOps) GetStatisticValues(componentID string, statisticName string, runID int) (interface{}, error) {
	params := map[string]interface{}{
		"componentId":   componentID,
		"statisticName": statisticName,
		"runId":         runID,
	}
	return s.Client.Post("/statistics/operations/getStatisticValues", params)
}

func (s *StatisticsOps) Search(searchString string, limit string, sort string, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return s.Client.Post("/statistics/operations/search", params)
}