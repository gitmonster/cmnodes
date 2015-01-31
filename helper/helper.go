package helper

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"labix.org/v2/mgo"
)

func Inspect(args ...interface{}) {
	spew.Dump(args)
}

type MongoSessionProvider struct {
	Session *mgo.Session
}

func (p MongoSessionProvider) GetMgoSession(collName string) (*mgo.Session, *mgo.Collection) {
	sess := p.Session.Copy()
	coll := sess.DB("cmnodes").C(collName)
	return sess, coll
}

////////////////////////////////////////////////////////////////////////////////
func GetTypeName(node interface{}) string {
	t := reflect.TypeOf(node)
	return strings.ToLower(t.Name())
}

////////////////////////////////////////////////////////////////////////////////////
func ConvertString(val string, tp reflect.Type) (interface{}, error) {
	switch tp.Kind() {
	case reflect.Int:
		if v, err := strconv.ParseInt(val, 10, 0); err != nil {
			return nil, err
		} else {
			return int(v), nil
		}
	case reflect.Int32:
		if v, err := strconv.ParseInt(val, 10, 32); err != nil {
			return nil, err
		} else {
			return int32(v), nil
		}
	case reflect.Int64:
		if v, err := strconv.ParseInt(val, 10, 64); err != nil {
			return nil, err
		} else {
			return int64(v), nil
		}
	case reflect.Float64:
		if v, err := strconv.ParseFloat(val, 64); err != nil {
			return nil, err
		} else {
			return float64(v), nil
		}
	case reflect.Struct:
		if tm, err := time.Parse("01/02/06 15:04:05", val); err != nil {
			return nil, err
		} else {
			return tm, nil
		}
	case reflect.String:
		return val, nil
	default:
		return nil, fmt.Errorf("ConvertString: no conversion to %v", tp.Kind())
	}
}

////////////////////////////////////////////////////////////////////////////////////
func GetStructFieldTypeFromTag(entity interface{}, searchTag string) (reflect.Type, error) {
	elem := reflect.TypeOf(entity)
	numFields := elem.NumField()
	for i := 0; i < numFields; i++ {
		field := elem.Field(i)
		tag := field.Tag.Get("json")
		if tag == searchTag {
			return field.Type, nil
		}
	}
	return nil, fmt.Errorf("GetStructFieldTypeFromTag: tag %s not found", searchTag)
}

//////////////////////////////////////////////////////////////////////////////
func AbsPath(name string) (string, error) {
	if path.IsAbs(name) {
		return name, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, name), err
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////////////////////////
func ParseRFC339Time(value string) time.Time {
	if dt, err := time.Parse("20060102T150405", value); err != nil {
		return time.Time{}
	} else {
		return dt
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////////////////////////
func PromptYesNo(text string, args ...interface{}) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y|n]:", fmt.Sprintf(text, args...))

	inp, _ := reader.ReadString('\n')
	if inp[0:1] == "y" || inp[0:1] == "Y" {
		return true
	} else if inp[0:1] == "n" || inp[0:1] == "N" {
		return false
	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func maxInt(val1, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func maxInt32(val1, val2 int32) int32 {
	if val1 > val2 {
		return val1
	}
	return val2
}
