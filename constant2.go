package main

// MessageType is the type of a message.
type MessageType string

const (
	MessageTypeEvent  MessageType = "EVENT"  // NIP-01
	MessageTypeReq    MessageType = "REQ"    // NIP-01
	MessageTypeClose  MessageType = "CLOSE"  // NIP-01
	MessageTypeNotice MessageType = "NOTICE" // NIP-01
	MessageTypeEOSE   MessageType = "EOSE"   // NIP-01
	MessageTypeOK     MessageType = "OK"     // NIP-20
)

const (
	EventKindSetMetadata             EventKind = 0    // NIP-01
	EventKindTextNote                EventKind = 1    // NIP-01
	EventKindRecommendServer         EventKind = 2    // NIP-01
	EventKindContacts                EventKind = 3    // NIP-02
	EventKindEncryptedDirectMessages EventKind = 4    // NIP-04
	EventKindEventDeletion           EventKind = 5    // NIP-09
	EventKindReposts                 EventKind = 6    // NIP-18
	EventKindReaction                EventKind = 7    // NIP-25
	EventKindBadgeAward              EventKind = 8    // NIP-58
	EventKindChannelCreation         EventKind = 40   // NIP-28
	EventKindChannelMetadata         EventKind = 41   // NIP-28
	EventKindChannelMessage          EventKind = 42   // NIP-28
	EventKindChannelHideMessage      EventKind = 43   // NIP-28
	EventKindChannelMuteUser         EventKind = 44   // NIP-28
	EventKindFileMetadata            EventKind = 1063 // NIP-94
	EventKindReporting               EventKind = 1984 // NIP-56
	EventKindZapRequest              EventKind = 9734 // NIP-57
	EventKindZap                     EventKind = 9735 // NIP-57
	EventKindJobResult               EventKind = 7000 // NIP-90
)

const (
	JobRequstBeforeEventContentKey string = "event:id:enckey:"
)

const (
	RESPONSE_SUCCESS_CODE   int    = 200
	RESPONSE_SUCCESS_STRING string = "SUCCESS"
)
