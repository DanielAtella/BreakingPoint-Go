package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bps-client-go/pkg/client"
)

var (
	modelName         string
	networkName       string
	bpsSystem         string
	bpsUser           string
	bpsPass           string
	slotNumber        = 1
	portList          = []int{0, 1, 4, 5}
	modelComponentAct []string
	runID             int
)

func loginToBps() *client.BPS {
	bps := client.NewBPS(bpsSystem, bpsUser, bpsPass, true)
	if _, err := bps.Login(); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	return bps
}

func setActiveComponents(bps *client.BPS) error {
	fmt.Println("Adjusting Test Model components state...")

	compsRaw, err := bps.TestModel.ListComponents()
	if err != nil {
		return fmt.Errorf("failed to list components: %w", err)
	}

	comps, ok := compsRaw.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected response while listing components: %#v", compsRaw)
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
			return fmt.Errorf("component ID is not a string: %#v", c["id"])
		}

		active := boolValue(c["active"])
		wantActive := activeWanted[label]

		if active != wantActive {
			stateStr := "inactive"
			if wantActive {
				stateStr = "active"
			}
			fmt.Printf("Component '%s' state changed to %s\n", label, stateStr)

			if _, err := bps.TestModel.SetComponentActive(idStr, wantActive); err != nil {
				return fmt.Errorf("failed to update component %s: %w", label, err)
			}
			changed = true
		}
	}

	if changed {
		fmt.Println("Saving component changes...")
		_, err := bps.TestModel.Save(modelName, true)
		if err != nil {
			return fmt.Errorf("failed to save components: %w", err)
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

func runTestAndPoll(bps *client.BPS, modelName string) error {
	runID, err := bps.RunTest(modelName, 2, false)
	if err != nil {
		return fmt.Errorf("run test error: %w", err)
	}

	fmt.Printf("Test running with runID: %v\n", runID)

	var lastData map[string]interface{}
	for {
		lastData, err = bps.PollTestProgress(runID)
		if err != nil {
			return fmt.Errorf("poll test progress error: %w", err)
		}

		if gone, ok := lastData["resource_gone"].(bool); ok && gone {
			fmt.Println("\nTest completed/cancelled; resource no longer available, ending polling.")
			break
		}

		progress := intValue(lastData["progress"])
		completed := boolValue(lastData["completed"])

		if completed || progress >= 100 {
			fmt.Println("Test completed.")
			break
		}

		time.Sleep(5 * time.Second)
	}

	runIDInt := toInt(runID)
	stats, err := bps.TestModel.RealTimeStats(runIDInt, "summary", -1, 1, "", []string{})
	if err != nil {
		return fmt.Errorf("failed to get realTimeStats: %w", err)
	}
	var finalProgress interface{}
	if statsMap, ok := stats.(map[string]interface{}); ok {
		finalProgress = statsMap["progress"]
	}
	fmt.Printf("Final progress: %v%%\n", finalProgress)

	report, err := bps.Reports.GetReportTable(runIDInt, "3.4")
	if err != nil {
		return fmt.Errorf("failed to get report table: %w", err)
	}
	fmt.Printf("Report section 3.4:\n%+v\n", report)
	return nil
}

func intValue(v interface{}) int {
	switch vv := v.(type) {
	case float64:
		return int(vv)
	case int:
		return vv
	default:
		return 0
	}
}

func boolValue(v interface{}) bool {
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func toInt(v interface{}) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		i, _ := strconv.Atoi(v)
		return i
	default:
		return 0
	}
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
		log.Fatalf("Failed to configure components: %v", err)
	}

	searchAndLoadNetworkConfig(bps)
	reservePorts(bps)

	if err := runTestAndPoll(bps, modelName); err != nil {
		log.Fatalf("Test run failed: %v", err)
	}

	unreservePorts(bps)
}
