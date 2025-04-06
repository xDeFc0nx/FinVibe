package flag

import (
)

func Flag() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	// 	arg := os.Args[1]
	//
}
