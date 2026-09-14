package destination

const (
	CreateDispositionIfNeeded CreateDisposition = "CREATE_IF_NEEDED"
	CreateDispositionNever    CreateDisposition = "CREATE_NEVER"
)

const (
	WriteDispositionTruncate WriteDisposition = "WRITE_TRUNCATE"
	WriteDispositionEmpty    WriteDisposition = "WRITE_EMPTY"
	WriteDispositionAppend   WriteDisposition = "WRITE_APPEND"
)

type (
	CreateDisposition string
	WriteDisposition  string
)
