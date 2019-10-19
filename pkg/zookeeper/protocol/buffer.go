package protocol

import "unsafe"

type Buffer struct {
	bytes []byte
	index int
}

func CreateBufferAroundSlice(bytes []byte) *Buffer {
	return &Buffer{
		bytes: bytes,
		index: 0,
	}
}

const int32Size = unsafe.Sizeof(int32(0))

func (buffer* Buffer) ReadRemainingBytes() []byte {
	endIndex := len(buffer.bytes)
	remaining := buffer.bytes[buffer.index:endIndex]
	buffer.index = endIndex
	return remaining
}

func (buffer* Buffer) ReadInt32() (int32, error) {
	return 0, nil
}

// TODO: Return error on panic or do explicit array bounds checks
func (buffer *Buffer) ReadBytes(length int) ([]byte, error) {
	endIndex := buffer.index + length
	bytes := buffer.bytes[buffer.index:endIndex]
	buffer.index = endIndex
	return bytes, nil
}
