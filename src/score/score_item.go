package score

type Item string
type Items []Item

const (
	Waiting        Item = "waiting"
	ExposedKong    Item = "exposed_kong"
	ConcealedKong  Item = "concealed_kong"
	Win            Item = "win"
	OwnDraw        Item = "own_draw"
	Lose           Item = "lose"
	Single         Item = "single"
	MissOneKind    Item = "miss_one_kind"
	MissTwoKind    Item = "miss_two_kind"
	Counter        Item = "counter"
	SevenPairs     Item = "seven_pairs"
	ThirteenOne    Item = "thirteen_one"
	ContinuousLine Item = "continuous_line"
)

func (i *Items) AddI(item Item) {
	*i = append(*i, item)
}

func (i *Items) indexI(item Item) int {
	for ix, it := range *i {
		if it == item {
			return ix
		}
	}
	return -1
}

func (i *Items) ContainsI(item Item) bool {
	return i.indexI(item) >= 0
}

func (i *Items) DelI(item Item) {
	ix := i.indexI(item)
	*i = append((*i)[:ix], (*i)[ix+1:]...)
}
