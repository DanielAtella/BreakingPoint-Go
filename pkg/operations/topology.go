package operations

type TopologyOps struct {
	Client ClientWrapper
}

func (t *TopologyOps) GetFanoutModes(cardId int) (interface{}, error) {
	params := map[string]interface{}{
		"cardId": cardId,
	}
	return t.Client.Post("/topology/operations/getFanoutModes", params)
}

func (t *TopologyOps) Reserve(reservation []map[string]interface{}, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"reservation": reservation,
		"force":       force,
	}
	return t.Client.Post("/topology/operations/reserve", params)
}

func (t *TopologyOps) ExportCapture(filepath string, args map[string]interface{}) error {
	params := map[string]interface{}{
		"filepath": filepath,
		"args":     args,
	}
	return t.Client.Export("/topology/operations/exportCapture", filepath, params)
}

func (t *TopologyOps) Unreserve(ports []map[string]interface{}) (interface{}, error) {
    params := map[string]interface{}{
        "unreservation": ports,
    }
    return t.Client.Post("/topology/operations/unreserve", params)
}