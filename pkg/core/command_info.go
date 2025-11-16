package core

import (
	"bytes"
	"errors"
	"fmt"
)

func cmdINFO(args []string) []byte {
	if len(args) == 0 {
		return Encode("All sections. I will implement later. You can try `INFO keyspace` command\n", false)
	}
	if len(args) > 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'INFO' command"), false)
	}
	switch args[0] {
	case "keyspace":
		var info []byte
		buf := bytes.NewBuffer(info)
		buf.WriteString("# Keyspace\r\n")
		fmt.Fprintf(buf, "db0:keys=%d,expires=%d,avg_ttl=%d\r\n", len(dictStore.GetDictStore()), dictStore.ExpiringKeysCount(), dictStore.TLL_Avg())
		return Encode(buf.String(), false)
	default:
		return Encode(errors.New("(error) ERR unknown INFO section"), false)
	}
}
