package games

import "fmt"

const (
	baseRewardNoClue = 20
	rewardClueOne    = 14
	rewardClueTwo    = 9
	rewardClueThree  = 5
)

func RewardForClueCount(clueCount int) int {
	switch {
	case clueCount <= 0:
		return baseRewardNoClue
	case clueCount == 1:
		return rewardClueOne
	case clueCount == 2:
		return rewardClueTwo
	default:
		return rewardClueThree
	}
}

func RewardGuide() string {
	return fmt.Sprintf("Reward: tanpa clue +%d balance | clue 1 +%d | clue 2 +%d | clue 3 +%d | nyerah +0", baseRewardNoClue, rewardClueOne, rewardClueTwo, rewardClueThree)
}
