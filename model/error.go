package model

type (
	ErrNotFound struct {
		Message string `json:"message"`
	}
)

func (e *ErrNotFound) Error() string {
	return e.Message
}