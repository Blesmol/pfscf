package canvas

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

// Entry represents an entry in the 'canvas' section
type Entry struct {
	id         string
	isResolved bool

	X, Y, X2, Y2 *float64
	Parent       *string
}

// NewEntry returns a new canvas entry
func NewEntry() (e Entry) {
	e = Entry{}
	return
}

func contentValErr(entry *Entry, errIn error) (errOut error) {
	return fmt.Errorf("Error validating canvas '%v': %v", entry.id, errIn)
}

// ID returns the ID of this entry
func (e *Entry) ID() string {
	return e.id
}

func (e *Entry) isValid() (err error) {
	if err = utils.CheckFieldsAreSet(e, "X", "Y", "X2", "Y2"); err != nil {
		return contentValErr(e, err)
	}

	if err = utils.CheckFieldsAreInRange(e, 0.0, 100.0, "X", "Y", "X2", "Y2"); err != nil {
		return contentValErr(e, err)
	}

	if e.X == e.X2 {
		err = fmt.Errorf("X coordinates are equal, both %.2f", *e.X)
		return contentValErr(e, err)
	}
	if e.Y == e.Y2 {
		err = fmt.Errorf("Y coordinates are equal, both %.2f", *e.Y)
		return contentValErr(e, err)
	}

	return nil
}

func (e *Entry) resolve(s *Store, resolveChain ...string) (err error) {
	if e.isResolved {
		return nil
	}

	// check for cyclic dependency
	for idx, otherID := range resolveChain {
		if e.id == otherID {
			outputChain := append(resolveChain[idx:], otherID) // reduce to relevant part, include conflicting ID again
			return fmt.Errorf("Error resolving canvas '%v': Cyclic dependency, chain is %v", e.id, outputChain)
		}
	}

	if err = e.isValid(); err != nil {
		return err
	}

	// ensure that coordinates are sorted, X<=X2, Y<=Y2
	if *e.X > *e.X2 {
		e.X, e.X2 = e.X2, e.X
	}
	if *e.Y > *e.Y2 {
		e.Y, e.Y2 = e.Y2, e.Y
	}

	if e.Parent != nil { // No parent canvas? Nothing more to do!
		parent, exists := s.Get(*e.Parent)

		// check that the entry from which we inherit really exists
		if !exists {
			return fmt.Errorf("Canvas '%v': Cannot find parent canvas '%v'", e.id, e.Parent)
		}

		// ensure that parent is already resolved
		resolveChain = append(resolveChain, e.id)
		if err = parent.resolve(s, resolveChain...); err != nil {
			return err
		}

		// calculate absolute coordinates
		e.X = calcAbsCoord(*e.X, *parent.X, *parent.X2)
		e.X2 = calcAbsCoord(*e.X2, *parent.X, *parent.X2)
		e.Y = calcAbsCoord(*e.Y, *parent.Y, *parent.Y2)
		e.Y2 = calcAbsCoord(*e.Y2, *parent.Y, *parent.Y2)
	}

	e.isResolved = true

	return nil
}

func calcAbsCoord(input, pCoord1, pCoord2 float64) *float64 {
	utils.Assert(input >= 0.0 && input <= 100.0, "Expected canvas coordinates to be in range 0.0-100.0")
	utils.Assert(pCoord1 <= pCoord2, "Coordinates should already be sorted")

	pWidth := pCoord2 - pCoord1
	result := pCoord1 + (pWidth * input / 100.0)

	return &result
}

func (e *Entry) deepCopy() *Entry {
	utils.Assert(e.isResolved == false, "copies should not happen after entry was resolved")

	copy := NewEntry()
	copy.id = e.id
	copy.X = utils.CopyFloat(e.X)
	copy.Y = utils.CopyFloat(e.Y)
	copy.X2 = utils.CopyFloat(e.X2)
	copy.Y2 = utils.CopyFloat(e.Y2)
	copy.Parent = utils.CopyString(e.Parent)

	return &copy
}

func (e *Entry) inherit(other *Entry) {
	utils.AddMissingValues(e, *other)
}
