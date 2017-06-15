package main

import (
	"fmt"
	"strings"
)

type CurrencyPair struct {
	Base   string
	Target string
}

func NewCurrencyPair(s string) (CurrencyPair, error) {
	curr := strings.Split(s, ",")
	if len(curr) != 2 {
		return CurrencyPair{}, MalformedCurrencyPairError{MalformedPair: s}
	}
	return CurrencyPair{Base: curr[0], Target: curr[1]}, nil
}

func (c CurrencyPair) String() string {
	return fmt.Sprintf("%v->%v", c.Base, c.Target)
}
