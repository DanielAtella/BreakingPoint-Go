package operations

type TestModelOps struct {
	Client ClientWrapper
}

func (t *TestModelOps) Load(template string, validate bool) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
		"validate": validate,
	}
	return t.Client.Post("/testmodel/operations/load", params)
}

func (t *TestModelOps) Run(modelname string, group int, allowMalware bool) (interface{}, error) {
	params := map[string]interface{}{
		"modelname":    modelname,
		"group":        group,
		"allowMalware": allowMalware,
	}
	return t.Client.Post("/testmodel/operations/run", params)
}

func (t *TestModelOps) ExportModel(name string, attachments bool, filepath string) error {
	params := map[string]interface{}{
		"name":        name,
		"attachments": attachments,
		"filepath":    filepath,
	}
	return t.Client.Export("/testmodel/operations/exportModel", filepath, params)
}

func (t *TestModelOps) ImportModel(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
	}
	return t.Client.Import("/testmodel/operations/importModel", filename, params)
}

func (t *TestModelOps) Search(searchString string, limit int, sort, sortorder string) (interface{}, error) {
	params := map[string]interface{}{
		"searchString": searchString,
		"limit":        limit,
		"sort":         sort,
		"sortorder":    sortorder,
	}
	return t.Client.Post("/testmodel/operations/search", params)
}

func (t *TestModelOps) Add(name, component, compType string, active bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":      name,
		"component": component,
		"type":      compType,
		"active":    active,
	}
	return t.Client.Post("/testmodel/operations/add", params)
}

func (t *TestModelOps) Save(name string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":  name,
		"force": force,
	}
	return t.Client.Post("/testmodel/operations/save", params)
}

func (t *TestModelOps) Clone(template, compType string, active bool, label string) (interface{}, error) {
	params := map[string]interface{}{
		"template": template,
		"type":     compType,
		"active":   active,
		"label":    label,
	}
	return t.Client.Post("/testmodel/operations/clone", params)
}

func (t *TestModelOps) Delete(name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}
	return t.Client.Post("/testmodel/operations/delete", params)
}

func (t *TestModelOps) RealTimeStats(runid int, rtsgroup string, numSeconds int, numDataPoints int, aggregate string, protocol []string) (interface{}, error) {
	params := map[string]interface{}{
		"runid":         runid,
		"rtsgroup":      rtsgroup,
		"numSeconds":    numSeconds,
		"numDataPoints": numDataPoints,
		"aggregate":     aggregate,
		"protocol":      protocol,
	}
	return t.Client.Post("/testmodel/operations/realTimeStats", params)
}

func (t *TestModelOps) Remove(id string) (interface{}, error) {
	params := map[string]interface{}{
		"id": id,
	}
	return t.Client.Post("/testmodel/operations/remove", params)
}

func (t *TestModelOps) Stop(runid int) (interface{}, error) {
	params := map[string]interface{}{
		"runid": runid,
	}
	return t.Client.Post("/testmodel/operations/stop", params)
}

func (t *TestModelOps) TestComponentDefinition(name string, dynamicEnums bool, includeOutputs bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":           name,
		"dynamicEnums":   dynamicEnums,
		"includeOutputs": includeOutputs,
	}
	return t.Client.Post("/testmodel/operations/testComponentDefinition", params)
}

func (t *TestModelOps) Validate(group string) (interface{}, error) {
	params := map[string]interface{}{
		"group": group,
	}
	return t.Client.Post("/testmodel/operations/validate", params)
}

func (t *TestModelOps) ListComponents() (interface{}, error) {
	return t.Client.Get("/testmodel/component", nil, nil)
}

func (t *TestModelOps) GetComponent(componentID string) (interface{}, error) {
	return t.Client.Get("/testmodel/component/"+componentID, nil, nil)
}

func (t *TestModelOps) SetComponentLabel(componentID, newLabel string) (interface{}, error) {
	params := map[string]interface{}{
		"label": newLabel,
	}
	err := t.Client.Patch("/testmodel/component/"+componentID, params)
	return nil, err
}

func (t *TestModelOps) SetComponentActive(componentID string, active bool) (interface{}, error) {
	params := map[string]interface{}{
		"active": active,
	}
	err := t.Client.Patch("/testmodel/component/"+componentID, params)
	return nil, err
}
