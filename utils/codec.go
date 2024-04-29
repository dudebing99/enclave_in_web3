package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// 消息: 4 字节长度 + 4 字节类型 + 变长内容
const (
	FixedMessageType uint32 = 4
	FixedMessageLen  uint32 = 4
)

// Encode 将消息编码
func Encode(messageType MessageType, message []byte) ([]byte, error) {
	// 计算消息的长度
	length := uint32(len(message))
	pkg := new(bytes.Buffer)

	// 写入消息长度
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	// 写入消息类型
	err = binary.Write(pkg, binary.LittleEndian, messageType)
	if err != nil {
		return nil, err
	}

	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (MessageType /* messageType */, []byte /* message */, error) {
	// 读取消息长度
	lengthByte, _ := reader.Peek(int(FixedMessageLen)) // 读取前4个字节
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length uint32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return UnknownMessageType, []byte{}, err
	}

	//log.Println("message len: ", length)

	// 读取消息类型
	messageTypeByte, _ := reader.Peek(int(FixedMessageLen + FixedMessageType)) // 读取前8个字节
	messageTypeBuff := bytes.NewBuffer(messageTypeByte[FixedMessageLen:])
	var messageType MessageType
	err = binary.Read(messageTypeBuff, binary.LittleEndian, &messageType)
	if err != nil {
		return UnknownMessageType, []byte{}, err
	}

	//log.Println("message type: ", messageType)

	// Buffered返回缓冲中现有的可读取的字节数。
	if int64(reader.Buffered()) < int64(length+FixedMessageLen+FixedMessageType) {
		return UnknownMessageType, []byte{}, err
	}

	// 读取真正的消息数据
	pack := make([]byte, int(length+FixedMessageLen+FixedMessageType))
	_, err = reader.Read(pack)
	if err != nil {
		return UnknownMessageType, []byte{}, err
	}

	return messageType, pack[FixedMessageLen+FixedMessageType:], nil
}
