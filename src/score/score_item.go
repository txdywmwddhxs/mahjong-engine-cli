package score

type Item string
type Items []Item

const (
	Waiting         Item = "waiting"
	ExposedKong     Item = "exposed_kong"
	ConcealedKong   Item = "concealed_kong"
	Win             Item = "win"
	OwnDraw         Item = "own_draw"
	Lose            Item = "lose"
	Single          Item = "single"
	MissOneKind     Item = "miss_one_kind"
	MissTwoKind     Item = "miss_two_kind"
	Counter         Item = "counter"
	SevenPairs      Item = "seven_pairs"
	ThirteenOrphans Item = "thirteen_one" // Keep the recorded item ID for replay compatibility.
	ContinuousLine  Item = "continuous_line"
)

func (items *Items) Add(item Item) {
	*items = append(*items, item)
}

func (items *Items) index(item Item) int {
	for itemIndex, candidate := range *items {
		if candidate == item {
			return itemIndex
		}
	}
	return -1
}

func (items *Items) Contains(item Item) bool {
	return items.index(item) >= 0
}

func (items *Items) Remove(item Item) {
	itemIndex := items.index(item)
	*items = append((*items)[:itemIndex], (*items)[itemIndex+1:]...)
}
