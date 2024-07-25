package errlog

import "errors"

// ErrBankEmployee404 indicates that the bank employee can't be found in the system
var ErrBankEmployee404 = errors.New("bank employee doesn't exist")
