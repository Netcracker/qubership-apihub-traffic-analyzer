// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orderedmap

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/openapi/yaml"
)

type Pair struct {
	key   string
	value interface{}
}

func (kv *Pair) Key() string {
	return kv.key
}

func (kv *Pair) Value() interface{} {
	return kv.value
}

type ByPair struct {
	Pairs    []*Pair
	LessFunc func(a *Pair, j *Pair) bool
}

func (a ByPair) Len() int           { return len(a.Pairs) }
func (a ByPair) Swap(i, j int)      { a.Pairs[i], a.Pairs[j] = a.Pairs[j], a.Pairs[i] }
func (a ByPair) Less(i, j int) bool { return a.LessFunc(a.Pairs[i], a.Pairs[j]) }

type OrderedMap struct {
	keys       []string
	values     map[string]interface{}
	escapeHTML bool
}

func New() *OrderedMap {
	o := OrderedMap{}
	o.keys = []string{}
	o.values = map[string]interface{}{}
	o.escapeHTML = true
	return &o
}

func (o *OrderedMap) SetEscapeHTML(on bool) {
	o.escapeHTML = on
}

func (o *OrderedMap) Get(key string) (interface{}, bool) {
	val, exists := o.values[key]
	return val, exists
}

func (o *OrderedMap) Set(key string, value interface{}) {
	_, exists := o.values[key]
	if !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

func (o *OrderedMap) Delete(key string) {
	// check key is in use
	_, ok := o.values[key]
	if !ok {
		return
	}
	// remove from keys
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	// remove from values
	delete(o.values, key)
}

func (o *OrderedMap) Keys() []string {
	return o.keys
}

// SortKeys Sort the map keys using your sort func
func (o *OrderedMap) SortKeys(sortFunc func(keys []string)) {
	sortFunc(o.keys)
}

// Sort Sort the map using your sort func
func (o *OrderedMap) Sort(lessFunc func(a *Pair, b *Pair) bool) {
	pairs := make([]*Pair, len(o.keys))
	for i, key := range o.keys {
		pairs[i] = &Pair{key, o.values[key]}
	}

	sort.Sort(ByPair{pairs, lessFunc})

	for i, pair := range pairs {
		o.keys[i] = pair.key
	}
}

func (o *OrderedMap) UnmarshalJSON(b []byte) error {
	if o.values == nil {
		o.values = map[string]interface{}{}
	}
	err := json.Unmarshal(b, &o.values)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	if _, err = dec.Token(); err != nil { // skip '{'
		return err
	}
	o.keys = make([]string, 0, len(o.values))
	return decodeOrderedMap(dec, o)
}

func decodeOrderedMap(dec *json.Decoder, o *OrderedMap) error {
	hasKey := make(map[string]bool, len(o.values))
	for {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		if delim, ok := token.(json.Delim); ok && delim == '}' {
			return nil
		}
		key := token.(string)
		if hasKey[key] {
			// duplicate key
			for j, k := range o.keys {
				if k == key {
					copy(o.keys[j:], o.keys[j+1:])
					break
				}
			}
			o.keys[len(o.keys)-1] = key
		} else {
			hasKey[key] = true
			o.keys = append(o.keys, key)
		}

		token, err = dec.Token()
		if err != nil {
			return err
		}
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				if values, ok := o.values[key].(map[string]interface{}); ok {
					newMap := OrderedMap{
						keys:       make([]string, 0, len(values)),
						values:     values,
						escapeHTML: o.escapeHTML,
					}
					if err = decodeOrderedMap(dec, &newMap); err != nil {
						return err
					}
					o.values[key] = newMap
				} else if oldMap, ok := o.values[key].(OrderedMap); ok {
					newMap := OrderedMap{
						keys:       make([]string, 0, len(oldMap.values)),
						values:     oldMap.values,
						escapeHTML: o.escapeHTML,
					}
					if err = decodeOrderedMap(dec, &newMap); err != nil {
						return err
					}
					o.values[key] = newMap
				} else if err = decodeOrderedMap(dec, &OrderedMap{}); err != nil {
					return err
				}
			case '[':
				if values, ok := o.values[key].([]interface{}); ok {
					if err = decodeSlice(dec, values, o.escapeHTML); err != nil {
						return err
					}
				} else if err = decodeSlice(dec, []interface{}{}, o.escapeHTML); err != nil {
					return err
				}
			}
		}
	}
}

func decodeSlice(dec *json.Decoder, s []interface{}, escapeHTML bool) error {
	for index := 0; ; index++ {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case '{':
				if index < len(s) {
					if values, ok := s[index].(map[string]interface{}); ok {
						newMap := OrderedMap{
							keys:       make([]string, 0, len(values)),
							values:     values,
							escapeHTML: escapeHTML,
						}
						if err = decodeOrderedMap(dec, &newMap); err != nil {
							return err
						}
						s[index] = newMap
					} else if oldMap, ok := s[index].(OrderedMap); ok {
						newMap := OrderedMap{
							keys:       make([]string, 0, len(oldMap.values)),
							values:     oldMap.values,
							escapeHTML: escapeHTML,
						}
						if err = decodeOrderedMap(dec, &newMap); err != nil {
							return err
						}
						s[index] = newMap
					} else if err = decodeOrderedMap(dec, &OrderedMap{}); err != nil {
						return err
					}
				} else if err = decodeOrderedMap(dec, &OrderedMap{}); err != nil {
					return err
				}
			case '[':
				if index < len(s) {
					if values, ok := s[index].([]interface{}); ok {
						if err = decodeSlice(dec, values, escapeHTML); err != nil {
							return err
						}
					} else if err = decodeSlice(dec, []interface{}{}, escapeHTML); err != nil {
						return err
					}
				} else if err = decodeSlice(dec, []interface{}{}, escapeHTML); err != nil {
					return err
				}
			case ']':
				return nil
			}
		}
	}
}

func (o OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(o.escapeHTML)
	for i, k := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// add key
		if err := encoder.Encode(k); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		// add value
		if err := encoder.Encode(o.values[k]); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (o OrderedMap) MarshalYAML() (interface{}, error) {
	res := yaml.MapSlice{}

	err := ConvertToMapSlice(&o, &res)
	return res, err
}

func ConvertToMapSlice(source *OrderedMap, dest *yaml.MapSlice) error {
	for _, k := range source.Keys() {
		v, _ := source.Get(k)
		if v == nil {
			*dest = append(*dest, yaml.MapItem{Key: k, Value: nil})
			continue
		}

		vType := reflect.TypeOf(v)
		switch vType.Kind() {
		case reflect.Ptr:
			if reflect.TypeOf(v).Elem().Name() == "OrderedMap" {
				newMS := yaml.MapSlice{}
				err := ConvertToMapSlice(v.(*OrderedMap), &newMS)
				if err != nil {
					return err
				}
				*dest = append(*dest, yaml.MapItem{Key: k, Value: newMS})
			} else {
				// TODO:??
			}
		case reflect.Slice:
			valAsSlice := v.([]interface{})
			var resultSlice []interface{}
			for _, sliceItem := range valAsSlice {
				siType := reflect.TypeOf(sliceItem)
				switch siType.Kind() {
				case reflect.Ptr:
					if siType.Elem().Name() == "OrderedMap" {
						newMS := yaml.MapSlice{}
						err := ConvertToMapSlice(sliceItem.(*OrderedMap), &newMS)
						if err != nil {
							return err
						}
						resultSlice = append(resultSlice, newMS)
					} else {
						// TODO: ?
						resultSlice = append(resultSlice, sliceItem)
					}
				default:
					resultSlice = append(resultSlice, sliceItem)
				}
			}
			*dest = append(*dest, yaml.MapItem{Key: k, Value: resultSlice})
		default:
			*dest = append(*dest, yaml.MapItem{Key: k, Value: v})
		}

	}
	return nil
}
