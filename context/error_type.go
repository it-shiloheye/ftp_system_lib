package ftp_context

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type LogErr = *LogItem

type LogItem struct {
	Location  string         `json:"location"`
	Time      time.Time      `json:"time"`
	After     string         `json:"after"`
	Body      map[string]any `json:"body"`
	Message   string         `json:"message"`
	Err       bool           `json:"is_error"`
	CallStack []error        `json:"call_stack"`
}

func (li *LogItem) Error() string {
	if !li.Err {
		return ""
	}
	return li.to_string()
}

func (li *LogItem) to_string() string {
	stp_1 := func() string {
		b, err := json.MarshalIndent(li, "\t", " ")
		if err != nil {
			panic(fmt.Sprint("LogItem.Error json.MarshalIndent ", err.Error()))
		}
		return string(b)
	}()

	stp_2 := fmt.Sprintf("%s:\n%s:\n%s", li.Time.Format(time.RFC822Z), li.Location, stp_1)

	stp_3 := func() string {
		if li.Err {
			return fmt.Sprint("[ERR] ", stp_2)
		}
		return fmt.Sprint("[LOG] ", stp_2)
	}()

	return stp_3
}

func (li *LogItem) Set(key string, value any) *LogItem {
	if li.Body == nil {
		li.Body = make(map[string]any)
	}
	li.Body[key] = value

	return li
}

func (li *LogItem) Setf(key string, str string, value ...any) *LogItem {
	if li.Body == nil {
		li.Body = make(map[string]any)
	}
	li.Body[key] = fmt.Sprintf(str, value...)

	return li
}

func (li *LogItem) SetAfter(after string) *LogItem {
	li.After = after

	return li
}

func (li *LogItem) SetAfterf(after string, v ...any) *LogItem {
	li.After = fmt.Sprintf(after, v...)

	return li
}
func (li *LogItem) Get(key string) (it any, ok bool) {
	if li.Body == nil {
		li.Body = make(map[string]any)
		return
	}
	it, ok = li.Body[key]
	return
}

func (li *LogItem) AppendParentError(err ...error) *LogItem {
	li.Err = true
	li.CallStack = append(li.CallStack, err...)
	return li
}

func NewLogItem(loc string, err bool) (lt *LogItem) {
	lt = &LogItem{
		Err:      err,
		Location: loc,
		Time:     time.Now(),
		Body:     map[string]any{},
	}
	return
}

func (lt *LogItem) SetMessage(v ...any) *LogItem {
	lt.Message = fmt.Sprint(v...)
	lt.Message = strings.ReplaceAll(lt.Message, "\\\\", "\\")
	lt.Message = strings.ReplaceAll(lt.Message, "\\n", "\n")
	lt.Message = strings.ReplaceAll(lt.Message, "\\t", "\t")
	return lt
}

func (lt *LogItem) SetMessagef(str string, v ...any) *LogItem {
	lt.Message = fmt.Sprintf(str, v...)
	return lt
}
