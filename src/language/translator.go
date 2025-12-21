package language

import (
	"fmt"
	"strings"

	"github.com/txdywmwddhxs/mahjong-engine-cli/src/score"
	"github.com/txdywmwddhxs/mahjong-engine-cli/src/utils"
)

// Translator resolves Keys to localized strings with fallback.
type Translator struct {
	lang         utils.Lang
	fallbackLang utils.Lang
}

func New(lang utils.Lang) Translator {
	if strings.TrimSpace(string(lang)) == "" {
		lang = utils.Chinese
	}
	return Translator{
		lang:         lang,
		fallbackLang: utils.Chinese,
	}
}

func (t Translator) Lang() utils.Lang { return t.lang }

func (t Translator) T(key Key, args ...any) string {
	tmpl := lookup(t.lang, key)
	if tmpl == "" {
		tmpl = lookup(t.fallbackLang, key)
	}
	if tmpl == "" {
		tmpl = string(key)
	}
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}

func lookup(lang utils.Lang, key Key) string {
	cat, ok := catalogs[lang]
	if !ok || cat == nil {
		return ""
	}
	return cat[key]
}

// ScoreItemName returns localized name for a score item.
func (t Translator) ScoreItemName(item score.Item) string {
	switch item {
	case score.Waiting:
		return t.T(KeyScoreWaiting)
	case score.ExposedKong:
		return t.T(KeyScoreExposedKong)
	case score.ConcealedKong:
		return t.T(KeyScoreConcealedKong)
	case score.Win:
		return t.T(KeyScoreWin)
	case score.OwnDraw:
		return t.T(KeyScoreOwnDraw)
	case score.Lose:
		return t.T(KeyScoreLose)
	case score.Single:
		return t.T(KeyScoreSingle)
	case score.MissTwoKind:
		return t.T(KeyScoreMissTwoKind)
	case score.MissOneKind:
		return t.T(KeyScoreMissOneKind)
	case score.Counter:
		return t.T(KeyScoreCounter, utils.CountOfWinAward)
	case score.SevenPairs:
		return t.T(KeyScoreSevenPairs)
	case score.ThirteenOne:
		return t.T(KeyScoreThirteenOne)
	case score.ContinuousLine:
		return t.T(KeyScoreContinuousLine)
	default:
		return string(item)
	}
}
