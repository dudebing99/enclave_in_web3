package vsock

import (
	"bufio"
	"enclave_in_web3/utils"
	"errors"
	"github.com/mdlayher/vsock"
	"github.com/spf13/viper"
	"io"
	"time"
)

func Trigger() {
	cid := viper.GetUint32("gateway.cid")
	if cid == 0 {
		return
	}

	port := viper.GetUint32("gateway.port")
	if port == 0 {
		return
	}

	conn, err := vsock.Dial(cid, port, nil)
	if err != nil {
		return
	}

	defer conn.Close()
}

func Process(msgType utils.MessageType, msg []byte) (utils.MessageType, []byte, error) {
	cid := viper.GetUint32("gateway.cid")
	if cid == 0 {
		return utils.UnknownMessageType, []byte{}, errors.New("invalid cid")
	}

	port := viper.GetUint32("gateway.port")
	if port == 0 {
		return utils.UnknownMessageType, []byte{}, errors.New("invalid port")
	}

	conn, err := vsock.Dial(cid, port, nil)
	if err != nil {
		return utils.UnknownMessageType, []byte{}, err
	}
	defer conn.Close()

	// 设置超时时间
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	conn.SetWriteDeadline(time.Now().Add(time.Second * 30))

	// 写数据
	encoded, err := utils.Encode(msgType, msg)
	if err != nil {
		return utils.UnknownMessageType, []byte{}, err
	}

	_, err = conn.Write(encoded)
	if err != nil {
		return utils.UnknownMessageType, []byte{}, err
	}

	// 读数据
	reader := bufio.NewReader(conn)
	msgType, msg, err = utils.Decode(reader)
	if err == io.EOF {
		return utils.UnknownMessageType, []byte{}, err
	}
	if err != nil {
		return utils.UnknownMessageType, []byte{}, err
	}

	return msgType, msg, nil
}
