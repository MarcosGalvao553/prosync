package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexInt é um tipo que aceita int ou string no JSON e converte para int
type FlexInt int

// UnmarshalJSON implementa json.Unmarshaler para FlexInt
func (fi *FlexInt) UnmarshalJSON(data []byte) error {
	// Tenta fazer unmarshal como int primeiro
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fi = FlexInt(i)
		return nil
	}

	// Se falhar, tenta como string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FlexInt deve ser int ou string: %w", err)
	}

	// Converte string para int
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("erro ao converter string para int: %w", err)
	}

	*fi = FlexInt(i)
	return nil
}

// Int retorna o valor como int
func (fi FlexInt) Int() int {
	return int(fi)
}

// FlexString é um tipo que aceita string ou número no JSON e converte para string
type FlexString string

// UnmarshalJSON implementa json.Unmarshaler para FlexString
func (fs *FlexString) UnmarshalJSON(data []byte) error {
	// Tenta fazer unmarshal como string primeiro
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fs = FlexString(s)
		return nil
	}

	// Se falhar, tenta como número (int ou float)
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*fs = FlexString(strconv.Itoa(i))
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*fs = FlexString(fmt.Sprintf("%v", f))
		return nil
	}

	return fmt.Errorf("FlexString deve ser string ou número")
}

// String retorna o valor como string
func (fs FlexString) String() string {
	return string(fs)
}
