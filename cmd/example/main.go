package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"bps-client-go/pkg/client"
)

var (
	modelName        string
	networkName      string
	bpsSystem        string
	bpsUser          string
	bpsPass          string
	slotNumber       = 1
	portList         = []int{0, 1, 4, 5}
	modelComponentAct []string
)

func loginToBps() *client.BPS {
	bps := client.NewBPS(bpsSystem, bpsUser, bpsPass, true)
	if _, err := bps.Login(); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	return bps
}

func setActiveComponents(bps *client.BPS) error {
	fmt.Println("Ajustando estado dos componentes do Test Model...")

	compsRaw, err := bps.TestModel.ListComponents()
	if err != nil {
		return fmt.Errorf("falha ao listar componentes: %w", err)
	}

	comps, ok := compsRaw.([]interface{})
	if !ok {
		return fmt.Errorf("resposta inesperada ao listar componentes: %#v", compsRaw)
	}

	activeWanted := make(map[string]bool)
	for _, name := range modelComponentAct {
		activeWanted[strings.TrimSpace(name)] = true
	}

	changed := false
	for _, comp := range comps {
		c, ok := comp.(map[string]interface{})
		if !ok {
			continue
		}

		label := fmt.Sprintf("%v", c["label"])

		idStr, ok := c["id"].(string)
		if !ok {
			return fmt.Errorf("id do componente não é string: %#v", c["id"])
		}

		active := boolValue(c["active"])
		wantActive := activeWanted[label]

		if active != wantActive {
			fmt.Printf("Alterando componente '%s' para ativo=%v\n", label, wantActive)
			if _, err := bps.TestModel.SetComponentActive(idStr, wantActive); err != nil {
				return fmt.Errorf("falha ao atualizar componente %s: %w", label, err)
			}
			changed = true
		} else {
			fmt.Printf("Componente '%s' já está ativo=%v, sem alterações\n", label, active)
		}
	}

	if changed {
		fmt.Println("Salvando alterações dos componentes...")
		_, err := bps.TestModel.Save(modelName, true)
		if err != nil {
			return fmt.Errorf("falha ao salvar componentes: %w", err)
		}
	}

	return nil
}

func searchTestModel(bps *client.BPS) {
	fmt.Println("Searching for test model...")
	searchResult, err := bps.TestModel.Search(modelName, 5, "name", "ascending")
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	if searchResult == nil {
		fmt.Println("Test Model not found.")
		bps.Logout()
		os.Exit(1)
	}
	fmt.Printf("Search result: %+v\n", searchResult)

	fmt.Println("Loading test model...")
	if _, err := bps.TestModel.Load(modelName, true); err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
}

func searchAndLoadNetworkConfig(bps *client.BPS) {
	fmt.Println("Searching for network config...")
	nnList, err := bps.NetworkOps.Search(networkName, "SRVIP", "", "ascending", "name", 10, 0)
	if err != nil || nnList == nil {
		fmt.Println("Network Config not found.")
		bps.Logout()
		os.Exit(1)
	}
	fmt.Printf("Network search result: %+v\n", nnList)

	fmt.Println("Loading network config...")
	if _, err := bps.NetworkOps.Load(networkName); err != nil {
		log.Fatalf("Failed to load network config: %v", err)
	}
}

func reservePorts(bps *client.BPS) {
	fmt.Println("Reserving Ports")
	for _, port := range portList {
		res := []map[string]interface{}{
			{"slot": slotNumber, "port": port, "group": 2},
		}
		if _, err := bps.TopologyOps.Reserve(res, true); err != nil {
			log.Fatalf("Failed to reserve port %d: %v", port, err)
		}
	}
}

func unreservePorts(bps *client.BPS) {
	fmt.Println("Unreserving ports")
	for _, port := range portList {
		req := []map[string]interface{}{
			{"slot": slotNumber, "port": port},
		}
		if _, err := bps.TopologyOps.Unreserve(req); err != nil {
			log.Printf("Failed to unreserve port %d: %v", port, err)
		}
	}
}

func pollTestProgress(bps *client.BPS, runID interface{}) (map[string]interface{}, error) {
	fmt.Println("Polling test progress...")
	runIDStr := fmt.Sprintf("%v", runID)
	path := fmt.Sprintf("/topology/runningTest/TEST-%s", runIDStr)

	var lastData map[string]interface{}
	for {
		resp, err := bps.Get(path, nil, nil)
		if err != nil {
			return lastData, fmt.Errorf("error getting runningTest info: %w", err)
		}

		if resp == nil {
			fmt.Println("Test completed and resource no longer available, ending polling.")
			break
		}

		data, ok := resp.(map[string]interface{})
		if !ok {
			return lastData, fmt.Errorf("unexpected response type: %#v", resp)
		}
		lastData = data

		progress := intValue(data["progress"])
		initProgress := intValue(data["initProgress"])
		phase := strValue(data["phase"])
		state := strValue(data["state"])
		completed := boolValue(data["completed"])

		fmt.Printf("Phase: %s, State: %s, Progress: %d%%, InitProgress: %d%%, Completed: %v\n",
			phase, state, progress, initProgress, completed)

		if completed || progress >= 100 {
			fmt.Println("Test completed.")
			break
		}

		time.Sleep(5 * time.Second)
	}
	return lastData, nil
}

func intValue(v interface{}) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case nil:
		return 0
	default:
		return 0
	}
}

func strValue(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func boolValue(v interface{}) bool {
	switch t := v.(type) {
	case bool:
		return t
	case nil:
		return false
	default:
		return false
	}
}

func runTestAndPoll(bps *client.BPS, modelName string) error {
	result, err := bps.TestModel.Run(modelName, 2, false)
	if err != nil {
		return fmt.Errorf("run test error: %w", err)
	}

	runID, ok := result.(map[string]interface{})["runid"]
	if !ok {
		return fmt.Errorf("runid not found in run result")
	}

	fmt.Printf("Test running with runid: %v\n", runID)
	_, err = pollTestProgress(bps, runID)
	if err != nil {
		return err
	}

	stats, err := bps.TestModel.RealTimeStats(
		int(runID.(float64)), 
		"summary",            
		-1,                   
		1,                    
		"",                   
		[]string{},           
	)

	if err != nil {
		return fmt.Errorf("failed to get realTimeStats: %w", err)
	}

	var finalProgress interface{}
	if statsMap, ok := stats.(map[string]interface{}); ok {
		finalProgress = statsMap["progress"]
	}

	fmt.Printf("Final progress: %v%%\n", finalProgress)

	report, err := bps.Reports.GetReportTable(int(runID.(float64)), "3.4")
	if err != nil {
		return fmt.Errorf("failed to get report table: %w", err)
	}
	fmt.Printf("Report section 3.4:\n%+v\n", report)

	return nil
}

func main() {
	modelName = strings.TrimSpace(os.Getenv("MODEL_NAME"))
	networkName = strings.TrimSpace(os.Getenv("NETWORK_NAME"))
	bpsSystem = strings.TrimSpace(os.Getenv("BPS_SYSTEM"))
	bpsUser = strings.TrimSpace(os.Getenv("BPS_USER"))
	bpsPass = strings.TrimSpace(os.Getenv("BPS_PASS"))

	compActEnv := strings.TrimSpace(os.Getenv("COMPONENT_ACT"))
	if compActEnv != "" {
		modelComponentAct = strings.Split(compActEnv, ",")
	}

	if modelName == "" || networkName == "" || bpsSystem == "" || bpsUser == "" || bpsPass == "" || len(modelComponentAct) == 0 {
		log.Fatal("Please set MODEL_NAME, NETWORK_NAME, BPS_SYSTEM, BPS_USER, BPS_PASS and COMPONENT_ACT environment variables")
	}

	bps := loginToBps()
	defer func() {
		fmt.Println("Logging out from BPS session")
		if err := bps.Logout(); err != nil {
			log.Printf("Logout error: %v", err)
		}
	}()

	searchTestModel(bps)

	if err := setActiveComponents(bps); err != nil {
		log.Fatalf("Falha ao configurar componentes: %v", err)
	}

	searchAndLoadNetworkConfig(bps)
	reservePorts(bps)

	if err := runTestAndPoll(bps, modelName); err != nil {
		log.Fatalf("Test run failed: %v", err)
	}

	unreservePorts(bps)
}