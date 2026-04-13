package err_mgr

// Application-level error codes.
const (
	ErrOK            = 0
	ErrInvalidParam  = 1001
	ErrStorageRead   = 2001
	ErrStorageWrite  = 2002
	ErrStorageDelete = 2003
	ErrAIForward     = 3001
	ErrInternal      = 9999
)
