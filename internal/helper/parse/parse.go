package parse

import "github.com/google/uuid"

// Convert []string (UUID as string) to []uuid.UUID, skip empty strings.
func ParseUUIDStrings(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}
