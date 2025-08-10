package operations

type SuperflowOps struct {
	Client ClientWrapper
}

func (s *SuperflowOps) AddAction(flowid int, typ string, actionid int, source string) (interface{}, error) {
	params := map[string]interface{}{
		"flowid":   flowid,
		"type":     typ,
		"actionid": actionid,
		"source":   source,
	}
	return s.Client.Post("/superflow/operations/addAction", params)
}

func (s *SuperflowOps) AddFlow(flowParams map[string]interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"flowParams": flowParams,
	}
	return s.Client.Post("/superflow/operations/addFlow", params)
}

func (s *SuperflowOps) AddHost(hostParams map[string]interface{}, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"hostParams": hostParams,
		"force":      force,
	}
	return s.Client.Post("/superflow/operations/addHost", params)
}

func (s *SuperflowOps) Delete(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return s.Client.Post("/superflow/operations/delete", params)
}

func (s *SuperflowOps) ImportResource(name string, filename string, force bool, typ string) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
		"type":     typ,
	}
	return s.Client.Import("/superflow/operations/importResource", filename, params)
}

func (s *SuperflowOps) Load(template string) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
	}
	return s.Client.Post("/superflow/operations/load", params)
}

func (s *SuperflowOps) New(template *string) (interface{}, error) {
	params := map[string]interface{}{}
	if template != nil {
		params["template"] = *template
	}
	return s.Client.Post("/superflow/operations/new", params)
}

func (s *SuperflowOps) RemoveAction(id int) (interface{}, error) {
	params := map[string]interface{}{
		"id": id,
	}
	return s.Client.Post("/superflow/operations/removeAction", params)
}

func (s *SuperflowOps) RemoveFlow(id int) (interface{}, error) {
	params := map[string]interface{}{
		"id": id,
	}
	return s.Client.Post("/superflow/operations/removeFlow", params)
}

func (s *SuperflowOps) Save(name *string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"force": force,
	}
	if name != nil {
		params["name"] = *name
	}
	return s.Client.Post("/superflow/operations/save", params)
}

func (s *SuperflowOps) SaveAs(name string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":  name,
		"force": force,
	}
	return s.Client.Post("/superflow/operations/saveAs", params)
}

func (s *SuperflowOps) Search(searchString string, limit string, sort string, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return s.Client.Post("/superflow/operations/search", params)
}

func (s *SuperflowOps) GetActionChoices(id int) (interface{}, error) {
	params := map[string]interface{}{
		"id": id,
	}
	return s.Client.Post("/superflow/actions/operations/getActionChoices", params)
}

func (s *SuperflowOps) GetActionInfo(id int) (interface{}, error) {
	params := map[string]interface{}{
		"id": id,
	}
	return s.Client.Post("/superflow/actions/operations/getActionInfo", params)
}

func (s *SuperflowOps) GetCannedFlows() (interface{}, error) {
	return s.Client.Post("/superflow/flows/operations/getCannedFlows", map[string]interface{}{})
}

func (s *SuperflowOps) GetFlowChoices(id int, name string) (interface{}, error) {
	params := map[string]interface{}{
		"id":   id,
		"name": name,
	}
	return s.Client.Post("/superflow/flows/operations/getFlowChoices", params)
}