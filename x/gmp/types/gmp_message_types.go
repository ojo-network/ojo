package types

const (
	// TypeUnrecognized means coin type is unrecognized
	TypeUnrecognized = iota
	// TypeGeneralMessage is a pure message
	TypeGeneralMessage
	// TypeGeneralMessageWithToken is a general message with token
	TypeGeneralMessageWithToken
	// TypeSendToken is a direct token transfer
	TypeSendToken
)
