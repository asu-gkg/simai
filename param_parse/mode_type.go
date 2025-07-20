package param_parse

type ModeType int

const (
	ModeTypeNONE = iota
	ModeTypeANALYTICAL
	ModeTypeNS3
	ModeTypeTECCL
	ModeTypeHTSIM
)
