package utils

import (
	"github.com/jinzhu/copier"
)

func CopyDeep(target interface{}, source interface{}) error {
	if err := copier.CopyWithOption(target, source, copier.Option{
		DeepCopy:   true,
		Converters: []copier.TypeConverter{},
	}); err != nil {
		return err
	}
	return nil
}
