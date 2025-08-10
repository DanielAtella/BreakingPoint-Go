package operations

type AdministrationOps struct {
	Client ClientWrapper
}

func (a *AdministrationOps) ImportAtiLicense(filename, name string) (interface{}, error) {
	params := map[string]interface{}{
		"filename": filename,
		"name":     name,
	}
	return a.Client.Import("/administration/atiLicensing/operations/importAtiLicense", filename, params)
}

func (a *AdministrationOps) ConfigPurge(configPurge interface{}) (interface{}, error) {
	return a.Client.Post("/administration/operations/configPurge", map[string]interface{}{
		"configPurge": configPurge,
	})
}

func (a *AdministrationOps) ExportAllTests(filepath string) error {
	params := map[string]interface{}{
		"filepath": filepath,
	}
	return a.Client.Export("/administration/operations/exportAllTests", filepath, params)
}

func (a *AdministrationOps) ImportAllTests(name, filename string, force bool) (interface{}, error) {
	params := map[string]interface{}{
		"name":     name,
		"filename": filename,
		"force":    force,
	}
	return a.Client.Import("/administration/operations/importAllTests", filename, params)
}