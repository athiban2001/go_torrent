package bencoding

import (
	"errors"
	"strconv"
)

var errorInvalidFormat = errors.New("invalid bencoding format")
var errorInvalidInteger = errors.New("invalid integer")

type ParserData struct {
	Data   []byte
	Index  int
	Length int
}

func (p ParserData) isString() bool {
	return p.Index < p.Length && p.Data[p.Index] >= byte(48) /* Digit 0 */ && p.Data[p.Index] <= byte(57) // Digit 9
}

func (p ParserData) isInt() bool {
	return p.Index < p.Length && p.Data[p.Index] == byte(105) // i
}

func (p *ParserData) isList() bool {
	return p.Index < p.Length && p.Data[p.Index] == byte(108) // l
}

func (p *ParserData) isDict() bool {
	return p.Index < p.Length && p.Data[p.Index] == byte(100) // d
}

func (p *ParserData) isEnd() bool {
	return p.Index < p.Length && p.Data[p.Index] == byte(101) // e
}

// ReadString : Format <string-length>:<string>
func (p *ParserData) ReadString() (string, error) {
	if !p.isString() {
		return "", errorInvalidFormat
	}

	length := 0
	colonIndex := -1

	for p.Index < p.Length {
		if p.Data[p.Index] == byte(58) {
			colonIndex = p.Index
			break
		}

		digit, err := strconv.Atoi(string(p.Data[p.Index]))
		if err != nil {
			return "", errorInvalidInteger
		}

		length *= 10
		length += digit
		p.Index++
	}

	if colonIndex == -1 {
		return "", errorInvalidFormat
	}

	// offseting index for next token
	p.Index += length + 1

	return string(p.Data[colonIndex+1 : colonIndex+1+length]), nil
}

// ReadInt : Format "i-123e"
func (p *ParserData) ReadInt() (int64, error) {
	var (
		intValue   int64 = 0
		isNegative bool  = false
		i          int   = 0
	)

	if !p.isInt() {
		return 0, errorInvalidFormat
	}

	// incrementing to offset "i"
	p.Index++

	for !p.isEnd() {
		if p.Data[p.Index] == byte(45) { // "-"
			isNegative = true
		} else if p.Data[p.Index] == byte(48) /* digit 0 */ &&
			((!isNegative && i == 0 && p.Index+1 < p.Length && p.Data[p.Index+1] != byte(101) /* "e" */) || (isNegative && i == 1)) { // invalid formats: "i-0e","i003e","i-003e"
			return 0, errorInvalidInteger
		} else {
			digit, err := strconv.Atoi(string(p.Data[p.Index]))
			if err != nil {
				return 0, errorInvalidInteger
			}
			intValue *= 10
			intValue += int64(digit)
		}

		p.Index++
		i++
	}

	// incrementing to offset "e"
	p.Index++

	if isNegative {
		return -intValue, nil
	}
	return intValue, nil
}

// ReadDictionary : Format "d<bencoded-key-value-pair>...e"
func (p *ParserData) ReadDictionary() (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	if !p.isDict() {
		return nil, errorInvalidFormat
	}

	// incrementing to offset "d"
	p.Index++

	for p.Index < p.Length && !p.isEnd() {
		key, err := p.ReadString()
		if err != nil {
			return nil, err
		} else if p.isString() {
			val, err := p.ReadString()
			if err != nil {
				return nil, err
			}
			dict[key] = val
		} else if p.isInt() {
			val, err := p.ReadInt()
			if err != nil {
				return nil, err
			}
			dict[key] = val
		} else if p.isList() {
			val, err := p.ReadList()
			if err != nil {
				return nil, err
			}
			dict[key] = val
		} else if p.isDict() {
			val, err := p.ReadDictionary()
			if err != nil {
				return nil, err
			}
			dict[key] = val
		} else {
			return nil, errorInvalidFormat
		}
	}

	// incrementing to offset "e"
	p.Index++

	return dict, nil
}

// ReadList : Format "l<bencoded-items>...e"
func (p *ParserData) ReadList() ([]interface{}, error) {
	list := []interface{}{}
	if !p.isList() {
		return nil, errorInvalidFormat
	}

	// incrementing to offset "l"
	p.Index++

	for !p.isEnd() {
		if p.isString() {
			val, err := p.ReadString()
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		} else if p.isInt() {
			val, err := p.ReadInt()
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		} else if p.isList() {
			val, err := p.ReadList()
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		} else {
			val, err := p.ReadDictionary()
			if err != nil {
				return nil, err
			}
			list = append(list, val)
		}
	}

	// incrementing to offset "e"
	p.Index++

	return list, nil
}
