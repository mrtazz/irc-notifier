package main

import (
	"encoding/json"
	"testing"
)

func TestReadLogLine(t *testing.T) {
	const testLogLine = `
  {"@source":"","@tags":["channelnotification","_grokparsefailure"],"@fields":{"irctime":["19:32:32"],"ircsender":["irccat"],"ircmessage":["foobla"]},"@timestamp":"2014-03-31T19:32:32.215Z","@source_host":"","@source_path":"","@message":"[19:32:32] <irccat> foobla","@type":"znclog"}
    `

	var notification Notification
	err := json.Unmarshal([]byte(testLogLine), &notification)
	if err != nil {
		t.Errorf("Parsing is broken with %s", err.Error())
		t.FailNow()
	}

	if notification.Fields.Time[0] != "19:32:32" {
		t.Errorf("wrong time read, expected 19:32:32 and got %s", notification.Fields.Time[0])
	}
	if notification.Fields.Sender[0] != "irccat" {
		t.Errorf("wrong time read, expected irccat and got %s",
			notification.Fields.Sender[0])
	}
	if notification.Fields.Message[0] != "foobla" {
		t.Errorf("wrong time read, expected foobla and got %s",
			notification.Fields.Message[0])
	}

}
