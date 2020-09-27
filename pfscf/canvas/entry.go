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
func (entry *Entry) ID() string {
	return entry.id
}

func (entry *Entry) isValid() (err error) {
	if err = utils.CheckFieldsAreSet(entry, "X", "Y", "X2", "Y2"); err != nil {
		return contentValErr(entry, err)
	}

	if err = utils.CheckFieldsAreInRange(entry, 0.0, 100.0, "X", "Y", "X2", "Y2"); err != nil {
		return contentValErr(entry, err)
	}

	if entry.X == entry.X2 {
		err = fmt.Errorf("X coordinates are equal, both %.2f", *entry.X)
		return contentValErr(entry, err)
	}
	if entry.Y == entry.Y2 {
		err = fmt.Errorf("Y coordinates are equal, both %.2f", *entry.Y)
		return contentValErr(entry, err)
	}

	return nil
}

func (entry *Entry) resolve(s *Store, resolveChain ...string) (err error) {
	if entry.isResolved {
		return nil
	}

	// check for cyclic dependency
	for idx, otherID := range resolveChain {
		if entry.id == otherID {
			outputChain := append(resolveChain[idx:], otherID) // reduce to relevant part, include conflicting ID again
			return fmt.Errorf("Error resolving canvas '%v': Cyclic dependency, chain is %v", entry.id, outputChain)
		}
	}

	if err = entry.isValid(); err != nil {
		return err
	}

	// ensure that coordinates are sorted, X<=X2, Y<=Y2
	if *entry.X > *entry.X2 {
		entry.X, entry.X2 = entry.X2, entry.X
	}
	if *entry.Y > *entry.Y2 {
		entry.Y, entry.Y2 = entry.Y2, entry.Y
	}

	if entry.Parent != nil { // No parent canvas? Nothing more to do!
		parent, exists := s.Get(*entry.Parent)

		// check that the entry from which we inherit really exists
		if !exists {
			return fmt.Errorf("Canvas '%v': Cannot find parent canvas '%v'", entry.id, entry.Parent)
		}

		// ensure that parent is already resolved
		resolveChain = append(resolveChain, entry.id)
		if err = parent.resolve(s, resolveChain...); err != nil {
			return err
		}

		// calculate absolute coordinates
		entry.X = calcAbsCoord(*entry.X, *parent.X, *parent.X2)
		entry.X2 = calcAbsCoord(*entry.X2, *parent.X, *parent.X2)
		entry.Y = calcAbsCoord(*entry.Y, *parent.Y, *parent.Y2)
		entry.Y2 = calcAbsCoord(*entry.Y2, *parent.Y, *parent.Y2)
	}

	entry.isResolved = true

	return nil
}

func calcAbsCoord(input, pCoord1, pCoord2 float64) *float64 {
	utils.Assert(input >= 0.0 && input <= 100.0, "Expected canvas coordinates to be in range 0.0-100.0")
	utils.Assert(pCoord1 <= pCoord2, "Coordinates should already be sorted")

	pWidth := pCoord2 - pCoord1
	result := pCoord1 + (pWidth * input / 100.0)

	return &result
}

func (entry *Entry) deepCopy() *Entry {
	utils.Assert(entry.isResolved == false, "copies should not happen after entry was resolved")

	copy := NewEntry()
	copy.id = entry.id
	copy.X = utils.CopyFloat(entry.X)
	copy.Y = utils.CopyFloat(entry.Y)
	copy.X2 = utils.CopyFloat(entry.X2)
	copy.Y2 = utils.CopyFloat(entry.Y2)
	copy.Parent = utils.CopyString(entry.Parent)

	return &copy
}

func (entry *Entry) inherit(other *Entry) {
	utils.AddMissingValues(entry, *other)
}
