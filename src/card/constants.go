package card

// Tile layout follows CardsList: three numbered suits, then seven honors.
const (
	SuitCount            = 3
	RanksPerSuit         = 9
	HonorKindCount       = 7
	SuitedKindCount      = SuitCount * RanksPerSuit
	TileKindCount        = SuitedKindCount + HonorKindCount
	CopiesPerTile        = 4
	TotalTiles           = TileKindCount * CopiesPerTile
	MinRank              = 1
	MaxRank              = RanksPerSuit
	SuitedTileCodeLength = 2
)

// Structural sizes count each kong as one meld plus an extra replacement tile.
const (
	PairSize             = 2
	MeldSize             = 3
	KongSize             = 4
	KongExtraTiles       = KongSize - MeldSize
	WinningMeldCount     = 4
	WinningGroupCount    = WinningMeldCount + 1
	WinningTileCount     = WinningMeldCount*MeldSize + PairSize
	InitialHandSize      = WinningTileCount - 1
	MaxSequenceStartRank = MaxRank - MeldSize + 1
)
