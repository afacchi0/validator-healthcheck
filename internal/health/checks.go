package health

import "validator-healthcheck/internal/rpc"

func IsBonded(v *rpc.StakingValidatorResponse) bool {
	return v.Validator.Status == "BOND_STATUS_BONDED"
}

func IsJailed(v *rpc.StakingValidatorResponse) bool {
	return v.Validator.Jailed
}

func HasTokens(v *rpc.StakingValidatorResponse) bool {
	return v.Validator.Tokens != "0"
}
