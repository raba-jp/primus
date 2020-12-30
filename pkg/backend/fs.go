package backend

import "github.com/spf13/afero"

func NewFs() afero.Fs {
	return afero.NewOsFs()
}
