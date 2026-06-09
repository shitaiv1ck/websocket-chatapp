package domains

import "encoding/json"

type Nullable[T any] struct {
	Value *T
	Set   bool
}

func (n *Nullable[T]) UnmarshalJSON(bytes []byte) error {
	n.Set = true

	if string(bytes) == "null" {
		n.Value = nil

		return nil
	}

	var value T
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}
