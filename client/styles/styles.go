package styles

import (
	"strings"

	"github.com/fatih/color"
)

var msgColors = map[string]*color.Color{
	"ERR":  color.New(color.FgRed, color.Bold),
	"INFO": color.New(color.FgYellow),
	"MSG":  color.New(color.FgCyan),
}

func StyleMsg(rawMsg string) string {
	/*
	  Check if the text has a message type and if not return as is.
	  If it does have a message type check if it's one of the recognized ones (msgColors)
	  and style accordingly. Also have a check for special use cases like MSG.
	*/
	msgTypeStart := strings.Index(rawMsg, "[")
	msgTypeEnd := strings.Index(rawMsg, "]")
	if msgTypeStart == -1 || msgTypeEnd == -1 {
		return rawMsg
	}
	msgType := rawMsg[msgTypeStart+1 : msgTypeEnd]
	filteredMsg := strings.TrimSpace(rawMsg[msgTypeEnd+1:])
	_, ok := msgColors[msgType]
	if !ok {
		return rawMsg
	}
	// Unless it's a MSG where we only want part of the message to be coloured, colour the
	// message based on the type of msgColors.
	if msgType == "MSG" {
		nameStart := strings.Index(filteredMsg, "[")
		nameEnd := strings.Index(filteredMsg, "]")
		if nameStart == -1 || nameEnd == -1 {
			return filteredMsg
		}
		username := msgColors["MSG"].Sprint(filteredMsg[nameStart+1 : nameEnd])
		message := filteredMsg[nameEnd+1:]
		return "[" + username + "]" + message
	}
	return msgColors[msgType].Sprint(filteredMsg)
}
