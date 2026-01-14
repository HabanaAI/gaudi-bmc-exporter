package rasmonitoringapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertData(t *testing.T) {

	data, err := convertData(OpcodeData{
		Data: []byte{163, 29},
	}, ReverseString)

	assert.Nil(t, err)

	assert.Equal(t, "1da3", data)

	data, err = convertData(OpcodeData{
		Data: []byte{67, 121, 235, 141},
	}, Int)
	assert.Nil(t, err)

	assert.Equal(t, "2381019459", data)

	data, err = convertData(OpcodeData{
		Data: []byte{255, 254, 63, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}, BinArray)
	assert.Nil(t, err)
	assert.Equal(t, "111111110111111111111100", data)

	data, err = convertData(OpcodeData{
		Data: []byte{16, 32, 72, 76, 50, 48, 56, 48},
	}, AsciiString)
	assert.Nil(t, err)
	assert.Equal(t, "HL2080", data)

	data, err = convertData(OpcodeData{
		Data: []byte{29, 163},
	}, String)
	assert.Nil(t, err)
	assert.Equal(t, "1da3", data)
}
