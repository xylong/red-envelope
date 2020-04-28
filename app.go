package red_envelope

import (
	"red-envelope/infra"
	"red-envelope/infra/base"
)

func init() {
	infra.Register(&base.PropsStart{})
	infra.Register(&base.DbxDatabaseStarter{})
}
