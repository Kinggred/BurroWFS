package utils

import (
	"burrowfs/core/types"
	"io"
)

func NewCombinedResponses(file *types.FileDTO, data io.ReadCloser, address string) *types.CombinedResponses {
	return &types.CombinedResponses{
		File:    file,
		Data:    data,
		Address: address,
	}
}
