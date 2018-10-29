package uuid

import (
	"github.com/go-openapi/strfmt"
	satoriuuid "github.com/satori/go.uuid"
)

// New generates and returns new uuid
func New() strfmt.UUID {

	return strfmt.UUID(satoriuuid.NewV1().String())
}
