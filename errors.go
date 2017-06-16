package main

import (
	"fmt"
)

type MalformedCurrencyPairError struct {
	MalformedPair string
}

type PairNotFoundError struct {
	C       CurrencyPair
	Message string
}

func (e MalformedCurrencyPairError) Error() string {
	return fmt.Sprintf("%T : %v, the correct format is \"Base,Target\"", e, e.MalformedPair)
}

func (e PairNotFoundError) Error() string {
	return fmt.Sprintf("%v : %v and/or %v is not a valid currency", e.Message, e.C.Base, e.C.Target)
}
