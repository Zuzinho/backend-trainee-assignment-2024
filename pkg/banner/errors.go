package banner

type IncorrectTypeFromDBError struct {
}

func (IncorrectTypeFromDBError) Error() string {
	return "wait []byte type, have incorrect type"
}

var (
	IncorrectTypeFromDBErr = IncorrectTypeFromDBError{}
)
