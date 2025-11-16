package core

import (
	"errors"

	data_structure "github.com/spaghetti-lover/multithread-redis/internal/data_structure/simple_set"
)

func cmdSADD(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0] // TODO: check key is used by other types or not
	set, exist := setStore[key]
	if !exist {
		set = data_structure.NewSimpleSet(key)
		setStore[key] = set
	}
	count := set.Add(args[1:]...)
	return Encode(count, false)
}

func cmdSREM(args []string) []byte {
	if len(args) < 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SADD' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		set = data_structure.NewSimpleSet(key)
		setStore[key] = set
	}
	count := set.Rem(args[1:]...)
	return Encode(count, false)
}

func cmdSMEMBERS(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SMEMBERS' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(make([]string, 0), false)
	}
	return Encode(set.Members(), false)
}

func cmdSISMEMBER(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SISMEMBER' command"), false)
	}
	key := args[0]
	set, exist := setStore[key]
	if !exist {
		return Encode(0, false)
	}
	return Encode(set.IsMember(args[1]), false)
}
