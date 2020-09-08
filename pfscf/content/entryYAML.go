package content

import (
	"fmt"

	"github.com/Blesmol/pfscf/pfscf/utils"
)

type entryYAML struct {
	e Entry
}

func (s *entryYAML) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
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
	case typeTextCell:
		s.e = newTextCell()
		err = unmarshal(s.e)
	case typeRectangle:
		s.e = newRectangle()
		err = unmarshal(s.e)
	case typeTrigger:
		s.e = newTrigger()
		err = unmarshal(s.e)
	default:
		err = fmt.Errorf("Unknown type: '%v'", ety.Type)
	}
	if err != nil {
		return err
	}

	return nil
}
