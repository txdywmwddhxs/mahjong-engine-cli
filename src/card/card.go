package card

import (
	"errors"
	"sort"
)

type Card string
type Cards []Card

func (c *Cards) Init() Cards {
	return *new(Cards)
}

func (c *Cards) ContainsC(card Card) bool {
	i := c.index(card)
	return i >= 0
}

func (c *Cards) index(card Card) int {
	for i, e := range *c {
		if card == e {
			return i
		}
	}
	return -1
}

func (c *Cards) CountC(card Card) int {
	i := 0
	for _, e := range *c {
		if card == e {
			i++
		}
	}
	return i
}

func (c *Cards) DelT(card Card, t int) {
	for i := 0; i < t; i++ {
		c.delC(card)
	}
}

func (c *Cards) DelE(card Card) error {
	if !c.ContainsC(card) {
		return errors.New("error")
	}
	c.delC(card)
	return nil
}

func (c *Cards) DelM(es ...Card) {
	for _, e := range es {
		c.delC(e)
	}
}

func (c *Cards) DelL() Card {
	l := len(*c)
	r := (*c)[l-1]
	*c = (*c)[:l-1]
	return r
}

func (c *Cards) delC(card Card) {
	i := c.index(card)
	*c = append((*c)[:i], (*c)[i+1:]...)
}

func (c *Cards) NewC(card ...Card) {
	*c = append(*c, card...)
}

func (c *Cards) Sort() {
	major := Cards{}
	minor := Cards{}
	for _, i := range *c {
		if len(i) == 2 {
			major = append(major, i)
		} else {
			minor = append(minor, i)
		}
	}
	sort.Sort(&major)
	sort.Sort(&minor)
	*c = append(major, minor...)
}

func (c *Cards) Beauty() string {
	r := "["
	for i, v := range *c {
		if i == 0 {
			r += string(v)
		} else {
			r = r + ", " + string(v)
		}
	}
	r += "]"
	return r
}

func (c *Cards) Less(i, j int) bool {
	return (*c)[i] < (*c)[j]
}

func (c *Cards) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *Cards) Len() int {
	return len(*c)
}
