package memory

import (
    "errors"
)

var (
    NOT_FOUND_ADDRESS   =   errors.New("Not found address.")
    OUT_OF_SIZE         =   errors.New("Out of block size.")
)