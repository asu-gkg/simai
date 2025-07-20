package mock_nccl

type NcclGroupType int

const (
	NcclGroupTypeTp = iota
	NcclGroupTypeDp
	Pp
	Ep
	DpEp
	None
)
