# validator-healthcheck

`validator-healthcheck` is a lightweight CLI tool for assessing the **health of a Cosmos SDK validator** using live on-chain data.

It combines **Tendermint RPC** and **Cosmos REST** endpoints to produce a normalized, machine-readable health report suitable for automation, monitoring, and scripting.

---

## What this tool does

Given:
- a Tendermint RPC endpoint
- a validator operator address (`cosmosvaloper...`)

`validator-healthcheck` determines whether a validator is **healthy** based on objective chain signals.

The tool outputs a structured JSON report and exits with deterministic exit codes.

---

## Supported checks

The following signals are evaluated:

- Validator is **not jailed**
- Validator is **bonded**
- Validator has **non-zero tokens**
- Node is **not catching up**

All checks are derived from:
- Tendermint RPC (`/status`)
- Cosmos REST (`/cosmos/staking/v1beta1/validators/{valoper}`)

---

## What this tool is (and isn’t)

**validator-healthcheck is:**
- a validator-aware health probe
- protocol-accurate for Cosmos SDK chains
- suitable for scripts, cron jobs, CI, and monitoring

**validator-healthcheck is not:**
- a full validator monitoring stack
- a Prometheus exporter (yet)
- a multi-chain validator abstraction

---

## Installation

Clone the repository and build:

```bash
go build -o validator-healthcheck ./cmd/health
```

## Usage

Syntax

```bash
./validator-healthcheck  --rpc <rpc-url> --validator <validator address>
```

Example

```bash
./validator-healthcheck   --rpc https://rpc.cosmos.directory/cosmoshub   --validator cosmosvaloper1hkqejlyrj9h8knms9hwvrn
p9xhqvjqlseh06vq
{
  "validator": "cosmosvaloper1hkqejlyrj9h8knms9hwvrnp9xhqvjqlseh06vq",
  "healthy": true,
  "issues": [],
  "jailed": false,
  "bonded": true,
  "has_tokens": true,
  "catching_up": false
}
```
