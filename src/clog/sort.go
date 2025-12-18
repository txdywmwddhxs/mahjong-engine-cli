package clog

import "github.com/txdywmwddhxs/mahjong-engine-cli/src/score"

type itemScore struct {
	Item  score.Item
	Score int
}

var sortBase = []score.Item{
	score.Lose,
	score.Win,
	score.Waiting,
	score.Single,
	score.ConcealedKong,
	score.ExposedKong,
	score.Counter,
	score.MissOneKind,
	score.MissTwoKind,
	score.SevenPairs,
	score.ContinuousLine,
	score.ThirteenOne,
	score.OwnDraw,
}

func sortMapByArray(m map[score.Item]int) []itemScore {
	var sortedPairs []itemScore

	// 创建映射以快速查找
	exists := make(map[score.Item]bool)
	for _, key := range sortBase {
		if value, found := m[key]; found {
			sortedPairs = append(sortedPairs, itemScore{key, value})
			exists[key] = true
		}
	}

	// 处理 sortBase 中的元素，只取存在于 map 的元素
	for _, key := range sortBase {
		if !exists[key] {
			continue // 如果在 map 中不存在该元素，则跳过
		}
	}

	return sortedPairs
}
