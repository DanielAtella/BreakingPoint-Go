package operations

type RemoteOps struct {
	Client ClientWrapper
}

func (r *RemoteOps) ConnectChassis(address string, remote string) (interface{}, error) {
	params := map[string]interface{}{
		"address": address,
		"remote":  remote,
	}
	return r.Client.Post("/remote/operations/connectChassis", params)
}

func (r *RemoteOps) DisconnectChassis(address string, port *int) (interface{}, error) {
	params := map[string]interface{}{
		"address": address,
	}
	if port != nil {
		params["port"] = *port
	}
	return r.Client.Post("/remote/operations/disconnectChassis", params)
}