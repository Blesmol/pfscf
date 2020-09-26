package template

type hierarchieStore []struct {
	ct  *Chronicle
	sub *hierarchieStore
}

func newHierarchieStore(s *Store, parentID string) (hs *hierarchieStore) {
	topLevelIDs := s.getTemplatesInheritingFrom(parentID) // lexically sorted list
	lhs := make(hierarchieStore, len(topLevelIDs))

	for idx, id := range topLevelIDs {
		lhs[idx].ct, _ = s.Get(id)
		lhs[idx].sub = newHierarchieStore(s, id)
	}

	return &lhs
}

func (hs *hierarchieStore) flatten() (result []*Chronicle) {
	result = make([]*Chronicle, 0)
	for _, e := range *hs {
		result = append(result, e.ct)
		if e.sub != nil {
			subList := e.sub.flatten()
			result = append(result, subList...)
		}
	}

	return result
}
