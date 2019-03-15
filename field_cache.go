package validator

type (
	fieldCache struct {
		index     int
		isPrivate bool
		name      string
		tagValue  string
		tagChunk  *tagChunk
	}
)
