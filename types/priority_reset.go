package types

import "math"

var (
	// PriorityResetHeightInterval determines the interval at which the validator priority is reset.
	// If set to 100, the priority is reset every 100 blocks (e.g., at heights 100, 200, 300, ...).
	// When the chain reaches a height that is a multiple of this value, the validator priority is reset.
	PriorityResetHeightInterval int64 = math.MaxInt64
	// PriorityResetRound defines the maximum number of rounds a block can remain unprocessed at a specific height.
	// If the chain is stuck at a height for more than this number of rounds (e.g., if PriorityResetRound is 100),
	// the validator priority will be reset to ensure progress.
	PriorityResetRound int32 = math.MaxInt32
)
