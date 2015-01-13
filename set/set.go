// A simple set data structure.
package set

import "sync"
import "fmt"
import "strings"

type any interface{}

type Set struct {
    m map[any]struct{}
    sync.RWMutex
}


// Determine if a variable could be a member of Set.
func IsLegal(v any) bool {
    switch v.(type) {
        // byte is alias of uint8, rune is alias of uint32,
        // so byte and rune are not in case clause
        case bool, error, string, complex64, complex128, float32, float64,
             int, int8, int16, int32, int64,
             uint, uint8, uint16, uint32, uint64, uintptr:
                return true
    }

    return false
}

func New() *Set {
    s := &Set{}
    s.m = make(map[any]struct{})
    return s
}


// Add item(s) to set.
// If there's an item that is not a legal type, return false
// If success, return true
func (s *Set) Add(items ...any) bool {
    s.Lock()
    defer s.Unlock()

    // If any item of items is not legal, return false
    for _, i := range(items) {
        if IsLegal(i) == false {
            return false
        }
    }

    for _, i := range(items) {
        // struct{}{} is a struct{} instance
        s.m[i] = struct{}{}
    }

    return true
}


// Add item(s) to set. if there's an item that is not a legal type, panic
func (s *Set) MustAdd(items ...any) {
    if s.Add(items...) == false {
        value := fmt.Sprintf("%v", items)

        // remove the outside square brackets
        if value[0] == '[' && value[len(value)-1] == ']' {
            value = value[1:len(value)-1]
        }
        panic(fmt.Sprintf("Value is not legal for adding to Set: %v\n", value))
    }
}

func (s *Set) Remove(item any) {
    s.Lock()
    s.Unlock()
    delete(s.m, item)
}

func (s *Set) Has(item any) bool {
    s.RLock()
    defer s.RUnlock()
    _, ok := s.m[item]
    return ok
}

func (s *Set) Len() int {
    return len(s.m)
}

func (s *Set) Clear() {
    s.Lock()
    defer s.Unlock()
    s.m = make(map[any]struct{})
}

func (s *Set) IsEmpty() bool {
    if len(s.m) == 0 {
    return true
    }
    return false
}

func (s *Set) List() []any {
    var l []any
    for i := range(s.m) {
        l = append(l, i)
    }
    return l
}

func (s *Set) String() string {
    items := make([]string, 0, s.Len())

    for i := range(s.m) {
        items = append(items, fmt.Sprintf("%v", i))
    }

    return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

