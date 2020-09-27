package handlers

import (
	"time"

	"github.com/k0kubun/pp"
)

const installTimeout = 5 * time.Minute

func init() {
	pp.ColoringEnabled = false
}
