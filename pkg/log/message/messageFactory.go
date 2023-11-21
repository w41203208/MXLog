package message

type IMessageFactory interface {
	GenMessage() IMessage
}

type XMessageFactory struct {
}

func (xmf *XMessageFactory) GenMessage(msg string) XMessage {
	return XMessage{}
}
