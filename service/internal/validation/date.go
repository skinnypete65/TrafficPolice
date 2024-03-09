package validation

import (
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

func IsDateOnly(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.DateOnly, fl.Field().String())
	log.Println(fl.Field().String())
	return err == nil
}
