package protocol

import "errors"

type Message interface {
	Accept(visitor MessageVisitor)
	Type() MessageType
}

type MessageTypeId int

type MessageType interface {
	Id() MessageTypeId
	Encode(message interface{}) ([]byte, error)
	Decode(encoded []byte) (Message, error)
}

type MessageVisitor interface {}

var ErrInvalidMessageType = errors.New("the message type is invalid")

type invalidMessageType struct {
	id MessageTypeId
}

func createInvalidMessageType(id MessageTypeId) MessageType {
	return invalidMessageType{id: id}
}

func (invalidType invalidMessageType) Id() MessageTypeId {
	return invalidType.id
}

func (invalidMessageType) Encode(interface{}) ([]byte, error) {
	return nil, ErrInvalidMessageType
}

func (invalidMessageType) Decode([]byte) (Message, error) {
	return nil, ErrInvalidMessageType
}
