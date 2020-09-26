package param

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

type entryYAML struct {
	e Entry
}

var (
	rankCounter int
)

func (ey *entryYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	// determine type of entry
	type entryTypeYAML struct{ Type string }
	var ety entryTypeYAML
	err = unmarshal(&ety)
	if err != nil {
		return err
	}
	if !utils.IsSet(ety.Type) {
		return fmt.Errorf("Missing or empty 'type' field")
	}

	// read concrete object based on type information
	switch ety.Type {
	case typeText:
		var e textEntry
		err = unmarshal(&e)
		ey.e = &e
	case typeSocietyID:
		var e societyidEntry
		err = unmarshal(&e)
		ey.e = &e
	case typeChoice:
		var e choiceEntry
		err = unmarshal(&e)
		ey.e = &e
	default:
		err = fmt.Errorf("Unknown parameter type: '%v'", ety.Type)
	}
	if err != nil {
		return err
	}

	rankCounter++
	ey.e.setRank(rankCounter)

	return nil
}
