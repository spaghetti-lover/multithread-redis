package core

func cmdHELP() []byte {
	helpCommands := []string{
		"--------------------------------",
		"PING [message] - Ping the server",
		"GET key - Get the value of a key",
		"SET key value - Set the value of a key",
		"EXISTS key - Check if a key exists",
		"TTL key - Get the time to live for a key",
		"SADD key member [member ...] - Add members to a set",
		"SREM key member [member ...] - Remove members from a set",
		"HELP - Show this help message",
		"CLEAR - Clear the terminal screen",
		"--------------------------------",
	}
	return Encode(helpCommands, false)
}
