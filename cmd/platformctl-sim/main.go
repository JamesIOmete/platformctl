package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/JamesIOmete/platformctl/internal/sim"
	"github.com/JamesIOmete/platformctl/internal/storage"
)

func main() {
	// Plugin args start at [1] because [0] is the executable name
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		printUsage()
		return
	}

	switch args[0] {
	case "ls":
		runSimList()
	case "run":
		// parse --scenario
		scenario := "default-warehouse"
		if len(args) > 1 && len(args[1]) > 11 && args[1][:11] == "--scenario=" {
			scenario = args[1][11:]
		}
		runSimRun(scenario)
	case "logs":
		if len(args) < 2 {
			fmt.Println("Error: Missing job ID. Usage: platformctl sim logs <id>")
			os.Exit(1)
		}
		runSimLogs(args[1])
	default:
		fmt.Printf("Error: Unknown sim subcommand '%s'. Try: platformctl sim help\n", args[0])
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: platformctl sim [ls | run --scenario=<name> | logs <id>]")
}

func runSimList() {
	s, err := storage.Load()
	if err != nil {
		fmt.Printf("Error: Failed to load state: %v\n", err)
		os.Exit(1)
	}

	if len(s.Simulations) == 0 {
		fmt.Println("No simulations found.")
		return
	}

	fmt.Println("JOB ID      \tSCENARIO            \tSTATUS    \tCREATED")
	fmt.Println("------------\t--------------------\t----------\t--------------------")
	for _, job := range s.Simulations {
		fmt.Printf("%-12s\t%-20s\t%-10s\t%s\n", job.ID, job.Scenario, job.Status, job.CreatedAt)
	}
}

func runSimRun(scenario string) {
	s, err := storage.Load()
	if err != nil {
		fmt.Printf("Error: Failed to load state: %v\n", err)
		os.Exit(1)
	}

	jobID := fmt.Sprintf("job-%d", rand.Intn(9999))
	newJob := sim.Job{
		ID:        jobID,
		Scenario:  scenario,
		Status:    "Running", // Start as running
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	s.Simulations = append(s.Simulations, newJob)
	if err := storage.Save(s); err != nil {
		fmt.Printf("Error: Failed to save job: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Simulation submitted successfully.\n")
	fmt.Printf("Job ID: %s\n", jobID)
	fmt.Printf("Scenario: %s\n", scenario)
	fmt.Println("Run 'platformctl sim ls' to view status.")
}

func runSimLogs(id string) {
	s, err := storage.Load()
	if err != nil {
		fmt.Printf("Error: Failed to load state: %v\n", err)
		os.Exit(1)
	}
	
	found := false
	for _, job := range s.Simulations {
		if job.ID == id {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("Error: Job %s not found.\n", id)
		os.Exit(1)
	}

	fmt.Printf("Streaming logs for simulation %s...\n", id)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[SIM] Loading physics engine... 100%")
	fmt.Println("[SIM] Spawning agents...")
	fmt.Println("[SIM] Agent 001 spawned at (0,0,0)")
	fmt.Println("[SIM] Scenario started.")
}
