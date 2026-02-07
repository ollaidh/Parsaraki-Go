package main

type MessageProcessor interface {
	ProcessMsg(string) string
}
