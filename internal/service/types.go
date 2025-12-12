package service

type SvcErr struct {
	Internal bool  `json:"internal"`
	Err      error `json:"msg"`
}

func (e *SvcErr) Error() string {
	return e.Err.Error()
}
