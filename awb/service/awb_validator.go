package service

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/repository"
)

type AwbValidator struct {
	awbProvider *repository.OrderAwbsRepository
}

func NewAwbValidator(awbProvider *repository.OrderAwbsRepository) *AwbValidator {
	return &AwbValidator{awbProvider: awbProvider}
}

/**
 * IsValid
 * @desc - Search for the awb code -  if it exists in cache or oms then it is a valid one
 */
func (aw *AwbValidator) IsValid(code string) bool {
	list := aw.awbProvider.GetList()
	if len(list) == 0 {
		return false
	}

	if aw.awbProvider.Get(code) == nil {
		return false
	}

	return true
}
