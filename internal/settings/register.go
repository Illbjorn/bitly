package settings

import (
	"encoding/json"
	"fmt"
)

var settings = &_Settings{}

type _Settings struct {
	ints    []*Setting[int]
	bools   []*Setting[bool]
	strings []*Setting[string]
}

func (self *_Settings) MarshalJSON() ([]byte, error) {
	x := map[string]any{}
	for _, i := range self.ints {
		x[i.Name()] = i.Get()
	}
	for _, b := range self.bools {
		x[b.Name()] = b.Get()
	}
	for _, s := range self.strings {
		x[s.Name()] = s.Get()
	}
	return json.Marshal(x)
}

func (self *_Settings) UnmarshalJSON(data []byte) error {
	x := map[string]any{}

	err := json.Unmarshal(data, &x)
	if err != nil {
		return err
	}

next:
	for k, v := range x {
		for _, i := range self.ints {
			if i.Name() == k {
				i.Set(v.(int))
				continue next
			}
		}
		for _, s := range self.strings {
			if s.Name() == k {
				s.Set(v.(string))
				continue next
			}
		}
		for _, b := range self.bools {
			if b.Name() == k {
				b.Set(v.(bool))
				continue next
			}
		}
		return fmt.Errorf("found unexpected setting [%s]", k)
	}

	return nil
}

func RegisterInt(s *Setting[int]) *Setting[int] {
	if s == nil {
		return nil
	}
	settings.ints = append(settings.ints, s)
	return s
}

func RegisterBool(s *Setting[bool]) *Setting[bool] {
	if s == nil {
		return nil
	}
	settings.bools = append(settings.bools, s)
	return s
}

func RegisterString(s *Setting[string]) *Setting[string] {
	if s == nil {
		return nil
	}
	settings.strings = append(settings.strings, s)
	return s
}
