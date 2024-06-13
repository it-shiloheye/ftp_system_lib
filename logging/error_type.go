package logging

import (
	"encoding/json"
	"fmt"

	"time"
)

type LogErr = *LogItem

type Loc string

func Locf(str string, v ...any) Loc {
	return Loc(fmt.Sprintf(str, v...))
}

type LogLevel int

const (
	LogLevelWrite LogLevel = iota - 1
	/*
	 Only when I would be "tracing" the code and trying
	 to find one part of a function specifically.
	*/
	LogLevelTrace

	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug01
	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug02
	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug03
	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug04
	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug05
	/*
	   Information that is diagnostically helpful to people more than just developers (IT, sysadmins, etc.).
	*/
	LogLevelDebug06
	/*
		Generally useful information to log (service start/stop, configuration assumptions, etc).
		Info I want to always have available but usually don't care about under normal circumstances.
		This is my out-of-the-box config level.
	*/
	LogLevelInfo01
	/*
		Generally useful information to log (service start/stop, configuration assumptions, etc).
		Info I want to always have available but usually don't care about under normal circumstances.
		This is my out-of-the-box config level.
	*/
	LogLevelInfo02
	/*
		Generally useful information to log (service start/stop, configuration assumptions, etc).
		Info I want to always have available but usually don't care about under normal circumstances.
		This is my out-of-the-box config level.
	*/
	LogLevelInfo03
	/*
		Anything that can potentially cause application oddities,
		but for which I am automatically recovering.
		(Such as switching from a primary to backup server,
		retrying an operation, missing secondary data, etc.)
	*/
	LogLevelWarn
	/*
		Any error which is fatal to the operation, but not the service or application
		(can't open a required file, missing data, etc.).
		These errors will force user (administrator, or direct user) intervention.
		These are usually reserved (in my apps) for incorrect connection strings, missing services, etc.
	*/
	LogLevelError01
	/*
		Any error which is fatal to the operation, but not the service or application
		(can't open a required file, missing data, etc.).
		These errors will force user (administrator, or direct user) intervention.
		These are usually reserved (in my apps) for incorrect connection strings, missing services, etc.
	*/
	LogLevelError02
	/*
		Any error that is forcing a shutdown of the service or application to prevent data loss (or further data loss).
		I reserve these only for the most heinous errors and situations
		where there is guaranteed to have been data corruption or loss.
	*/
	LogLevelFatal
)

func (ll LogLevel) String() string {
	switch ll {
	case LogLevelTrace:
		return "[TRACE]"
	case LogLevelDebug01:
		return "[DEBUG-01]"
	case LogLevelDebug02:
		return "[DEBUG-02]"
	case LogLevelDebug03:
		return "[DEBUG-03]"
	case LogLevelDebug04:
		return "[DEBUG-04]"
	case LogLevelDebug05:
		return "[DEBUG-05]"
	case LogLevelDebug06:
		return "[DEBUG-06]"
	case LogLevelInfo01:
		return "[INFO-01]"
	case LogLevelInfo02:
		return "[INFO-02]"
	case LogLevelInfo03:
		return "[INFO-03]"
	case LogLevelWarn:
		return "[WARN]"
	case LogLevelError01:
		return "[ERROR-01]"
	case LogLevelError02:
		return "[ERROR-02]"
	case LogLevelFatal:
		return "[FATAL]"

	}
	return "[LOG]"
}

type LogItem struct {
	Location  string         `json:"location"`
	Time      time.Time      `json:"time"`
	After     string         `json:"after"`
	Body      map[string]any `json:"body"`
	Message   string         `json:"message"`
	Level     LogLevel       `json:"log_level"`
	CallStack []error        `json:"call_stack"`
}

func (li *LogItem) Error() string {
	if li.Level < LogLevelWarn {
		return ""
	}
	return li.to_string()
}

func (li *LogItem) String() string {
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

	stp_3 := fmt.Sprint("%s ", li.Level.String(), stp_2)
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

	li.CallStack = append(li.CallStack, err...)
	return li
}

func NewLogItem(loc string, log_level LogLevel) (lt *LogItem) {
	lt = &LogItem{
		Level:    log_level,
		Location: loc,
		Time:     time.Now(),
		Body:     map[string]any{},
	}
	return
}

func (lt *LogItem) SetMessage(v ...any) *LogItem {
	lt.Message = fmt.Sprint(v...)

	return lt
}

func (lt *LogItem) SetMessagef(str string, v ...any) *LogItem {
	lt.Message = fmt.Sprintf(str, v...)
	return lt
}

type PrintOptions = int

const (
	PO_LINE PrintOptions = iota + 1
	PO_JSON
	PO_PLAIN
)

func (lt *LogItem) Print(print_option ...PrintOptions) string {
	if lt.Level == LogLevelWrite {
		return lt.Message
	}

	if len(print_option) < 1 {
		return lt.to_string()
	}
	po := print_option[0]

	var msg string
	switch po {
	case PO_PLAIN:
		msg = fmt.Sprintf("%s: %s", lt.Level.String(), fmt.Sprint(lt.Time))

		if len(lt.Message) > 0 {
			msg += "\nMessage: " + lt.Message
		}
		return msg
	case PO_LINE:
		msg = fmt.Sprintf("%s: %s", lt.Level.String(), fmt.Sprint(lt.Time))

		if len(lt.Location) > 0 {
			msg += "\nLoc: " + lt.Location
		}
		if len(lt.After) > 0 && lt.Level > LogLevelInfo02 {
			msg += "\nAfter: " + lt.After
		}
		if len(lt.Message) > 0 {
			msg += "\nMessage: " + lt.Message
		}

		if lt.Level >= LogLevelError01 {
			for _, line := range lt.CallStack {
				msg += "\n"

				tmp, ok := line.(LogErr)
				if !ok {
					msg += fmt.Sprint(line)
					continue
				}
				msg += tmp.Print(PO_LINE)
			}

		}

		return msg

	case PO_JSON:
		d, _ := json.Marshal(lt)
		return string(d)

	}

	return lt.String()
}
