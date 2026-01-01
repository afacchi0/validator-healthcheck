package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"validator-healthcheck/internal/health"
	"validator-healthcheck/internal/rpc"
)

func main() {
	rpcURL := flag.String(
		"rpc",
		"https://rpc.cosmos.directory/cosmoshub",
		"Tendermint RPC endpoint",
	)

	validator := flag.String(
		"validator",
		"",
		"Validator operator address (cosmosvaloper...)",
	)

	flag.Parse()

	if *validator == "" {
		fmt.Fprintln(os.Stderr, "--validator is required")
		os.Exit(2)
	}

	// Tendermint RPC (node status)
	client := rpc.New(*rpcURL)

	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Cosmos REST (validator state)
	cosmos := rpc.NewCosmos("https://rest.cosmos.directory/cosmoshub")

	v, err := cosmos.Validator(*validator)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Build report
	report := health.Report{
		Validator:  *validator,
		Jailed:     health.IsJailed(v),
		Bonded:     health.IsBonded(v),
		HasTokens:  health.HasTokens(v),
		CatchingUp: status.Result.SyncInfo.CatchingUp,
		Issues:     []string{},
	}

	if report.Jailed {
		report.Issues = append(report.Issues, "validator is jailed")
	}
	if !report.Bonded {
		report.Issues = append(report.Issues, "validator not bonded")
	}
	if !report.HasTokens {
		report.Issues = append(report.Issues, "validator has zero tokens")
	}
	if report.CatchingUp {
		report.Issues = append(report.Issues, "node is catching up")
	}

	report.Healthy = len(report.Issues) == 0

	// Emit JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if err := enc.Encode(report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Exit codes
	if report.Healthy {
		os.Exit(0)
	}

	os.Exit(1)
}
