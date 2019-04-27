package todotxt

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Item represents todotxt task
type Item struct {
	Complete       bool
	Priority       *Priority
	CompletionDate *time.Time
	CreationDate   *time.Time
	Description    string
	Tags           []Tag
}

// String returns text representation of Item
func (i *Item) String() string {
	result := ""
	if i.Complete {
		result = "x "
	}
	if i.Priority != nil {
		result += "(" + i.Priority.String() + ") "
	}

	if i.CompletionDate != nil {
		result += i.CompletionDate.Format("2006-01-02") + " "
		if i.CreationDate != nil {
			result += i.CreationDate.Format("2006-01-02") + " "
		} else {
			result += time.Now().Format("2006-01-02") + " "
		}
	} else if i.CreationDate != nil {
		result += i.CreationDate.Format("2006-01-02") + " "
	}
	result += i.Description + " "
	for _, t := range i.Tags {
		switch t.Key {
		case TagContext:
			result += "@" + t.Value + " "
		case TagProject:
			result += "+" + t.Value + " "
		default:
			result += t.Key + ":" + t.Value + " "
		}
	}
	return strings.Trim(result, " \n")
}

// Parse multiline todotxt string
func Parse(todo string) ([]Item, error) {
	lines := strings.Split(todo, "\n")
	items := make([]Item, 0, len(lines))
	for ln, line := range lines {
		i, err := ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("error at line %d: %v", ln, err)
		}
		items = append(items, i)
	}
	return items, nil
}

// ParseLine parses single todotxt line
func ParseLine(line string) (Item, error) {
	i := Item{}
	tokens := strings.Split(line, " ")
	state := 0
	for _, t := range tokens {
		if state == 0 && t == "x" {
			state = 1
			i.Complete = true
			continue
		}
		if state <= 1 && len(t) == 3 && t[0] == '(' && t[2] == ')' {
			p, err := PriorityFromLetter(string(t[1]))
			if err != nil {
				return i, err
			}
			i.Priority = &p
			state = 2
			continue
		}
		if state <= 2 {
			ti, err := time.Parse("2006-01-02", t)
			if err == nil {
				i.CreationDate = &ti
				state = 3
				continue
			}
			state = 4
		}
		if state <= 3 {
			ti, err := time.Parse("2006-01-02", t)
			if err == nil {
				i.CompletionDate = i.CreationDate
				i.CreationDate = &ti
				state = 4
				continue
			}
			state = 4
		}
		if t[0] == '+' {
			i.Tags = append(i.Tags, Tag{
				Key:   TagProject,
				Value: string(t[1:]),
			})
			continue
		}
		if t[0] == '@' {
			i.Tags = append(i.Tags, Tag{
				Key:   TagContext,
				Value: string(t[1:]),
			})
			continue
		}
		kv := strings.Split(t, ":")
		if len(kv) == 2 {
			i.Tags = append(i.Tags, Tag{
				Key:   kv[0],
				Value: kv[1],
			})
			continue
		}
		if i.Description == "" {
			i.Description = t
		} else {
			i.Description += " " + t
		}
	}
	return i, nil
}

type Priority int

// String returns letter by priority (0=A, 1=B, 2=C ...)
func (p Priority) String() string {
	return string([]byte{byte(p + 65)})
}

// PriorityFromLetter returns numeric priority from letter priority (A=0, B=1, C=2 ...)
func PriorityFromLetter(letter string) (Priority, error) {
	if len(letter) != 1 {
		return 0, errors.New("incorrect priority length")
	}
	code := []byte(letter)[0]
	if code < 65 || code > 90 {
		return 0, errors.New("priority must be between A and Z")
	}
	return Priority(code - 65), nil
}

// TagContext constant is key for context tag
const TagContext = "@context"

// TagProject constant is key for project tag
const TagProject = "+project"

// Tag represents builtin and custom tags
type Tag struct {
	Key   string
	Value string
}
