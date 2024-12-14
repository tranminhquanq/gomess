package storage

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/tranminhquanq/gomess/internal/config"
)

type Connection struct {
	*pop.Connection
}

func Dial(config *config.GlobalConfiguration) (*Connection, error) {
	return nil, nil
}
