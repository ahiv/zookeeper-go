package protocol

import "fmt"

type Codec interface {
	Decode(packet []byte) (Message, error)
	Encode(message Message) ([]byte, error)
}

type Protocol struct {
	messageTypes []MessageType
}

// Decode decodes a packet
func (protocol *Protocol) Decode(packet []byte) (Message, error) {
	buffer := CreateBufferAroundSlice(packet)
	messageType, err := protocol.readMessageType(buffer)
	if err != nil {
		return nil, err
	}
	return messageType.Decode(buffer.ReadRemainingBytes())
}

func (protocol *Protocol) Encode(message Message) ([]byte, error) {
	return message.Type().Encode(message)
}

func (protocol *Protocol) readMessageType(buffer *Buffer) (MessageType, error) {
	typeId, err := readMessageTypeId(buffer)
	if err != nil {
		return nil, err
	}
	return protocol.findMessageTypeById(typeId)
}

func (protocol *Protocol) findMessageTypeById(typeId MessageTypeId) (MessageType, error) {
	if !protocol.isTypeIdInBounds(typeId) {
		return nil, createNoSuchMessageTypeError(typeId)
	}
	return protocol.messageTypes[typeId], nil
}

func createNoSuchMessageTypeError(typeId MessageTypeId) error {
	return fmt.Errorf("protocol does not know any message with id: %d", typeId)
}

func (protocol *Protocol) isTypeIdInBounds(typeId MessageTypeId) bool {
	return typeId >= 0 && int(typeId) < len(protocol.messageTypes)
}

func readMessageTypeId(buffer *Buffer) (MessageTypeId, error) {
	typeId, err := buffer.ReadInt32()
	return MessageTypeId(typeId), err
}

type Builder struct {
	messageTypes map[MessageTypeId] MessageType
	highestTypeId MessageTypeId
}

func CreateBuilder() *Builder {
	return &Builder{
		messageTypes:  map[MessageTypeId] MessageType{},
		highestTypeId: 0,
	}
}

func (builder *Builder) RegisterMessageType(messageType MessageType) *Builder {
	id := messageType.Id()
	builder.registerTypeId(id)
	builder.messageTypes[id] = messageType
	return builder
}

func (builder *Builder) registerTypeId(id MessageTypeId) {
	if id > builder.highestTypeId {
		builder.highestTypeId = id
	}
}

func (builder *Builder) CreateProtocol() *Protocol {
	flattenedTypes := builder.createFlattenedMessageTypeMap()
	return &Protocol{messageTypes:flattenedTypes}
}

func (builder *Builder) createFlattenedMessageTypeMap() []MessageType {
	highestIndex := int(builder.highestTypeId)
	return flattenMessageTypeMap(builder.messageTypes, highestIndex)
}

func flattenMessageTypeMap(
	input map[MessageTypeId] MessageType, highestIndex int) []MessageType {

	output := createInvalidMessageTypeSlice(highestIndex)
	for index, messageType := range input {
		output[index] = messageType
	}
	return output
}

func createInvalidMessageTypeSlice(length int) []MessageType {
	messages := make([]MessageType, length)
	for index := 0; index < length; index++ {
		typeId := MessageTypeId(index)
		messages[index] = createInvalidMessageType(typeId)
	}
	return messages
}