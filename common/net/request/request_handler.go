package request

import (
	logger "gitlab.altex.ro/plug/go_logger"
	"io"
)

type NormalizerInterface interface {
	Normalize(reader io.Reader, toNormalize interface{}) error
}

type ValidatorInterface interface {
	Validate(interface{}) (bool, error)
}

type Handler struct {
	normalizer NormalizerInterface
	validator  ValidatorInterface
	logger     logger.Logger
	err        error
	dest       interface{}
}

func NewRequestHandler(normalizer NormalizerInterface, validator ValidatorInterface, logger logger.Logger) *Handler {
	return &Handler{
		normalizer: normalizer,
		validator:  validator,
		logger:     logger,
		err:        nil,
		dest:       nil,
	}
}

func (h *Handler) Normalize(r io.Reader, dest interface{}) *Handler {
	h.err = h.normalizer.Normalize(r, dest)
	h.dest = dest
	return h
}

func (h *Handler) Validate() *Handler {
	if h.err != nil {
		return h
	}

	_, h.err = h.validator.Validate(h.dest)

	return h
}

func (h *Handler) Result() *RequestError {
	h.dest = nil
	if h.err != nil {
		return NewNormalizerErrorHandler(h.logger).HandleError(h.err)
	}

	return nil
}
