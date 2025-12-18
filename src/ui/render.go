package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/card"
	l "github.com/txdywmwddhxs/mahjong-engine-cli/src/language"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
)

func RenderCards(c *card.Cards) string {
	if c == nil || len(*c) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(*c))
	for _, i := range *c {
		parts = append(parts, string(i))
	}
	return "[" + strings.Join(parts, ", ") + "]"
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

var scoreTrans = map[score.Item]func() string{
	score.Waiting:        l.Waiting,
	score.ExposedKong:    l.ExposedKong,
	score.ConcealedKong:  l.ConcealedKong,
	score.Win:            l.Win,
	score.OwnDraw:        l.OwnDraw,
	score.Lose:           l.Lose,
	score.Single:         l.Single,
	score.MissTwoKind:    l.MissTwoKind,
	score.MissOneKind:    l.MissOneKind,
	score.Counter:        l.Counter,
	score.SevenPairs:     l.SevenPairs,
	score.ThirteenOne:    l.ThirteenOne,
	score.ContinuousLine: l.ContinuousLine,
}

func RenderScoreDetail(detail map[score.Item]int) string {
	if len(detail) == 0 {
		return ""
	}

	exists := make(map[score.Item]bool, len(detail))
	for k := range detail {
		exists[k] = true
	}

	parts := make([]string, 0, len(detail))
	seen := make(map[score.Item]bool, len(detail))

	for _, item := range sortBase {
		v, ok := detail[item]
		if !ok {
			continue
		}
		seen[item] = true
		nameFn, ok := scoreTrans[item]
		name := string(item)
		if ok && nameFn != nil {
			name = nameFn()
		}
		parts = append(parts, fmt.Sprintf("%s: %d", name, v))
	}

	// Append any extra items in a stable order.
	extra := make([]score.Item, 0)
	for item := range exists {
		if seen[item] {
			continue
		}
		extra = append(extra, item)
	}
	sort.Slice(extra, func(i, j int) bool { return string(extra[i]) < string(extra[j]) })
	for _, item := range extra {
		nameFn, ok := scoreTrans[item]
		name := string(item)
		if ok && nameFn != nil {
			name = nameFn()
		}
		parts = append(parts, fmt.Sprintf("%s: %d", name, detail[item]))
	}

	return strings.Join(parts, ", ")
}
