package param

import "github.com/Blesmol/pfscf/pfscf/utils"

type commonFields struct {
	id      string
	group   string
	theRank int
}

func (e *commonFields) ID() string {
	return e.id
}

func (e *commonFields) setID(id string) {
	utils.Assert(!utils.IsSet(e.id), "Should only be called once per object")
	e.id = id
}

func (e *commonFields) Group() string {
	return e.group
}

func (e *commonFields) setGroup(groupID string) {
	e.group = groupID
}

func (e *commonFields) rank() int {
	return e.theRank
}

func (e *commonFields) setRank(rank int) {
	e.theRank = rank
}
