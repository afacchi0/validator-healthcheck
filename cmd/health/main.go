package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"validator-healthcheck/internal/health"
	"validator-healthcheck/internal/rpc"
)

func main() {
	rpcURL := flag.String(
		"rpc",
		"https://rpc.cosmos.directory/cosmoshub",
		"Tendermint RPC endpoint",
	)

	restURL := flag.String(
		"rest",
		"https://rest.cosmos.directory/cosmoshub",
		"Cosmos REST API endpoint",
	)

	timeout := flag.Duration(
		"timeout",
		5*time.Second,
		"HTTP timeout for RPC and REST requests",
	)

	validator := flag.String(
		"validator",
		"",
		"Validator operator address (cosmosvaloper...)",
	)

	outputJSON := flag.Bool(
		"json",
		false,
		"Emit compact JSON output",
	)

	outputPretty := flag.Bool(
		"pretty",
		false,
		"Emit pretty-printed JSON output",
	)

	flag.Parse()

	if *validator == "" {
		fmt.Fprintln(os.Stderr, "--validator is required")
		os.Exit(2)
	}

	// Tendermint RPC (node status)
	client := rpc.NewWithTimeout(*rpcURL, *timeout)

	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Cosmos REST (validator state)
	cosmos := rpc.NewCosmosWithTimeout(*restURL, *timeout)

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
	if *outputJSON && *outputPretty {
		fmt.Fprintln(os.Stderr, "--json and --pretty are mutually exclusive")
		os.Exit(2)
	}

	enc := json.NewEncoder(os.Stdout)
	pretty := *outputPretty
	if !pretty && !*outputJSON {
		pretty = true
	}
	if pretty {
		enc.SetIndent("", "  ")
	}

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
