package tree

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type any = interface{}

type Twig struct {
	Name   string `json:"name"`
	Kind   string `json:"kind"`
	Value  any    `json:"value"`
	Childs []Twig `json:"childs"`
}

type Tree struct {
	Base     *Twig
	fileName string
	Indent   string
}

const treeLayout = "2006-01-02T15:04:05.0000000-07:00"

func fullPath(fileName string) (string, error) {
	dir, fn := filepath.Split(fileName)
	dir = filepath.ToSlash(dir)
	if (dir == "./") || (dir == "") {
		exeFull, _ := os.Executable()
		exePath, _ := filepath.Split(exeFull)
		fileName = filepath.Join(exePath, fn)
	}

	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	fileName = filepath.ToSlash(fileName)

	return fileName, nil
}

func (tw *Twig) set(name string, value any) error {
	tw.Name = name
	tw.setKind(value)

	return nil
}

func (tw *Twig) setKind(value any) error {
	switch x := value.(type) {
	case string:
		tw.Value = x
		tw.Kind = "string"
	case int:
		tw.Value = int64(x)
		tw.Kind = "integer"
	case int8:
		tw.Value = int64(x)
		tw.Kind = "integer"
	case int16:
		tw.Value = int64(x)
		tw.Kind = "integer"
	case int32:
		tw.Value = int64(x)
		tw.Kind = "integer"
	case int64:
		tw.Value = x
		tw.Kind = "integer"
	case float32:
		tw.Value = float64(x)
		tw.Kind = "float"
	case float64:
		tw.Value = x
		tw.Kind = "float"
	case time.Time:
		tw.Value = x.Format(treeLayout)
		tw.Kind = "datetime"
	case bool:
		tw.Value = x
		tw.Kind = "bool"
	case nil:
		tw.Value = nil
		tw.Kind = ""
	default:
		tw.Value = nil
		tw.Kind = ""
		return fmt.Errorf("not found type:%v", x)
	}

	return nil
}

func (tw *Twig) get(kind ...string) any {
	var result any
	var k string
	//var k string

	if len(kind) != 0 {
		k = kind[0]
	} else {
		k = tw.Kind
	}

	k = strings.ToLower(k)
	switch k {
	case "string":

	case "integer":

	case "float":

	case "bool":

	case "datetime":

	default:
		return nil
	}

	result = nil
	switch k {
	case "string":
		result = ""

		switch x := tw.Value.(type) {
		case string:
			result = x
		case int64:
			result = strconv.FormatInt(int64(x), 10)
		case float64:
			result = strconv.FormatFloat(x, 'f', -1, 64)
		case bool:
			result = strconv.FormatBool(x)
		default:
			return result
		}
	case "integer":
		result = int64(0)

		switch x := tw.Value.(type) {
		case string:
			result, _ = strconv.ParseInt(x, 10, 64)
		case int64:
			result = x
		case float64:
			result = int64(x)
		case bool:
			if x {
				result = int64(1)
			}
		default:
			return result
		}
	case "float":
		result = float64(0)

		switch x := tw.Value.(type) {
		case string:
			result, _ = strconv.ParseFloat(x, 64)
		case int64:
			result = float64(x)
		case float64:
			result = x
		case bool:
			if x {
				result = float64(1)
			}
		default:
			return result
		}
	case "bool":
		result = false

		switch x := tw.Value.(type) {
		case string:
			s := strings.Trim(x, " ")
			s = strings.ToLower(s)

			if (s[:1] == "t") || (s[:1] == "y") {
				result = true
			}
		case int64:
			if x != 0 {
				result = true
			}
		case float64:
			if x != 0 {
				result = true
			}
		case bool:
			result = x
		default:
			return result
		}
	case "datetime":
		var tt time.Time
		result = tt

		switch x := tw.Value.(type) {
		case string:
			tt, err := time.Parse(treeLayout, x)
			if err != nil {
				return nil
			}
			result = tt
		case time.Time:
			result = tw.Value.(time.Time)
		default:
			return result
		}
	default:
		result = nil
	}
	return result
}

func (root *Tree) Create(fileName string, indent string) error {
	var rt Twig
	var op Twig
	var id Twig
	rt.set("root", nil)
	op.set("options", nil)
	id.set("indent", indent)

	op.Childs = append(op.Childs, id)
	rt.Childs = append(rt.Childs, op)
	root.Base = &rt

	fileName, err := fullPath(fileName)
	if err != nil {
		return err
	}

	_, err = os.Stat(fileName)
	if err == nil {
		err = os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	buf, _ := json.MarshalIndent(root.Base, "", indent)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	f.Write(buf)
	f.Close()

	root.fileName = fileName
	root.Indent = indent
	return nil
}

func (root *Tree) Open(fileName string) error {
	fileName, err := fullPath(fileName)
	if err != nil {
		return err
	}

	_, err = os.Stat(fileName)
	if err != nil {
		return err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	var jsonData []byte
	jsonData, err = io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &root.Base)
	if err != nil {
		return err
	}

	root.fileName = fileName

	root.Indent, err = root.GetValueStr([]string{"options", "indent"})
	if err != nil {
		return err
	}

	return nil
}

func (root *Tree) Close() {
	root.Base = nil
	root.fileName = ""
	root.Indent = ""
}

func (root *Tree) Reload() error {
	fileName := root.fileName
	root.Base = nil

	err := root.Open(fileName)
	if err != nil {
		return err
	}
	return nil
}

func (root *Tree) SaveAs(fileName string) error {
	dir, fn := filepath.Split(fileName)
	if dir == "./" {
		exeFull, _ := os.Executable()
		exePath, _ := filepath.Split(exeFull)
		fileName = filepath.Join(exePath, fn)
	}

	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	_, err = os.Stat(fileName)
	if err == nil {
		err = os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	var buf []byte
	if len(root.Indent) == 0 {
		buf, _ = json.Marshal(root.Base)
	} else {
		buf, _ = json.MarshalIndent(root.Base, "", root.Indent)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	f.Write(buf)
	f.Close()

	return nil
}

func (root *Tree) Save() error {
	_, err := os.Stat(root.fileName)
	if err == nil {
		err = os.Remove(root.fileName)
		if err != nil {
			return err
		}
	}

	var buf []byte
	if len(root.Indent) == 0 {
		buf, _ = json.Marshal(root.Base)
	} else {
		buf, _ = json.MarshalIndent(root.Base, "", root.Indent)
	}

	f, err := os.Create(root.fileName)
	if err != nil {
		return err
	}

	f.Write(buf)
	f.Close()

	return nil
}

func (root *Tree) Find(names []string) (*Twig, error) {
	var count int = 0

	if root.Base == nil {
		return nil, fmt.Errorf("root is NULL: %v", names)
	}

	r := root.Base
	for i := range names {
		for j := range r.Childs {
			//log.Printf(".name: %v", names[i])

			if r.Childs[j].Name == names[i] {
				r = &r.Childs[j]
				count++
				break
			}
		}
	}

	if count == len(names) {
		return r, nil
	}

	return nil, fmt.Errorf("not found: %v", names[count])
}

func (root *Tree) findPlus(names []string) (*Twig, *Twig, error) {
	var count int

	if root.Base == nil {
		return nil, nil, fmt.Errorf("root is NULL: %v", names)
	}

	r := root.Base
	v := root.Base
	for i := range names {
		for j := range r.Childs {
			//log.Printf(".name: %v", names[i])

			if r.Childs[j].Name == names[i] {
				v = r
				r = &r.Childs[j]
				count++
				break
			}
		}
	}

	if count == len(names) {
		return r, v, nil
	}

	return nil, nil, fmt.Errorf("not found: %v", names[count])
}

func (root *Tree) List(names []string) ([]string, error) {
	var list []string

	tw, err := root.Find(names)
	if err != nil {
		return list, err
	}

	for i := range tw.Childs {
		list = append(list, tw.Childs[i].Name)
	}

	return list, nil
}

func (root *Tree) AddNew(name string, value any, dst []string) error {
	if name == "" {
		return fmt.Errorf(".Name blank.%v", name)
	}

	fw, err := root.Find(dst)
	if err != nil {
		return err
	}

	for i := range fw.Childs {
		if fw.Childs[i].Name == name {
			return fmt.Errorf("duplicate name. %v", name)
		}
	}

	var tw Twig
	err = tw.set(name, value)
	if err != nil {
		return err
	}

	fw.Childs = append(fw.Childs, tw)
	return nil
}

func (root *Tree) Delete(dst []string) error {
	var _delete func(t *Twig) error

	_delete = func(t *Twig) error {
		for i := range t.Childs {
			_delete(&t.Childs[i])
		}
		t.Childs = nil
		return nil
	}

	td, tr, err := root.findPlus(dst)
	if err != nil {
		return err
	}

	for i := range td.Childs {
		_delete(&td.Childs[i])
	}
	td.Childs = nil

	var tw Twig
	tw.Childs = tr.Childs
	tr.Childs = nil
	for i := range tw.Childs {
		if tw.Childs[i].Name != td.Name {
			tr.Childs = append(tr.Childs, tw.Childs[i])
		}
	}

	return nil
}

func (root *Tree) Copy(src []string, dst []string) error {
	var _copy func(src *Twig, dst *Twig) error

	_copy = func(src *Twig, dst *Twig) error {
		for i := range src.Childs {
			child := src.Childs[i]

			var w Twig
			w.set(child.Name, child.Value)
			err := _copy(&child, &w)
			if err != nil {
				return err
			}
			dst.Childs = append(dst.Childs, w)
		}
		return nil
	}

	fs, err := root.Find(src)
	if err != nil {
		return err
	}

	fd, err := root.Find(dst)
	if err != nil {
		return err
	}

	var tw Twig
	tw.set(fs.Name, fs.Value)
	_copy(fs, &tw)
	for i := range tw.Childs {
		fd.Childs = append(fd.Childs, tw.Childs[i])
	}

	return nil
}

func (root *Tree) Move(src []string, dst []string) error {
	ts, err := root.Find(src)
	if err != nil {
		return err
	}

	td, err := root.Find(dst)
	if err != nil {
		return err
	}
	list, err := root.List(dst)
	if err != nil {
		return err
	}

	for i := range ts.Childs {
		sw := true
		for j := range list {
			if list[j] == ts.Childs[i].Name {
				sw = false
				break
			}
		}

		if sw {
			td.Childs = append(td.Childs, ts.Childs[i])
		}
	}
	ts.Childs = nil
	return nil
}

func (root *Tree) GetValue(src []string) (any, error) {
	const emsg = "do not match kind: %v"
	var result any

	tw, err := root.Find(src)
	if err != err {
		return nil, err
	}

	if tw.Value == nil {
		return nil, fmt.Errorf("value is null: %v", src)
	}

	if tw.Kind == "" {
		return nil, fmt.Errorf("kind is null: %v", src)
	}

	result = tw.get()

	return result, nil
}

func (root *Tree) GetValueStr(src []string) (string, error) {
	const emsg = "do not match kind: %v"
	var result string

	tw, err := root.Find(src)
	if err != nil {
		return result, err
	}

	if tw.Value == nil {
		return result, fmt.Errorf("value is null: %v", src)
	}

	result = tw.get("string").(string)

	return result, nil
}

func (root *Tree) GetValueInt(src []string) (int64, error) {
	const emsg = "do not match kind: %v"
	var result int64

	tw, err := root.Find(src)
	if err != nil {
		return result, err
	}

	if tw.Value == nil {
		return result, fmt.Errorf("value is null: %v", src)
	}

	result = tw.get("integer").(int64)

	return result, nil
}

func (root *Tree) GetValueFloat(src []string) (float64, error) {
	const emsg = "do not match kind: %v"
	var result float64

	tw, err := root.Find(src)
	if err != nil {
		return result, err
	}

	if tw.Value == nil {
		return result, fmt.Errorf("value is null: %v", src)
	}

	result = tw.get("float").(float64)

	return result, nil
}

func (root *Tree) GetValueBool(src []string) (bool, error) {
	const emsg = "do not match kind: %v"
	var result bool

	tw, err := root.Find(src)
	if err != nil {
		return result, err
	}

	if tw.Value == nil {
		return result, fmt.Errorf("value is null: %v", src)
	}

	result = tw.get("bool").(bool)

	return result, nil
}

func (root *Tree) SetValue(value any, src []string) error {
	tw, err := root.Find(src)
	if err != nil {
		return err
	}

	err = tw.set(tw.Name, value)
	if err != nil {
		return err
	}

	return nil
}

func (root *Tree) GetString(sep string, src []string) (string, error) {
	var result string
	var v string

	if sep == "" {
		sep = " "
	}

	tw, err := root.Find(src)
	if err != nil {
		return result, err
	}

	for i := range tw.Childs {
		ctw := tw.Childs[i]
		v = ctw.get("string").(string)

		result = result + ctw.Name + "=" + v

		if i < len(tw.Childs)-1 {
			result = result + sep
		}
	}
	return result, nil
}

func (root *Tree) SetMaping(m map[string]any, src []string) error {
	var name string

	if v, ok := m["name"]; ok {
		name = v.(string)
	} else {
		return fmt.Errorf("format is incorrect. %v", m)
	}

	dst := append(src, name)
	_, err := root.Find(dst)
	if err != nil {
		err := root.AddNew(name, nil, src)
		if err != nil {
			return err
		}
	}

	for k, _ := range m {
		if k == "name" {
			continue
		}

		chd := append(dst, k)
		_, err := root.Find(chd)
		if err != nil {
			err = root.AddNew(k, m[k], dst)
			if err != nil {
				return err
			}
		} else {
			err = root.SetValue(m[k], chd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (root *Tree) SetMapings(ms []map[string]any, src []string) error {
	for i := range ms {
		m := ms[i]
		err := root.SetMaping(m, src)
		if err != nil {
			return err
		}
	}
	return nil
}

func (root *Tree) GetMaping(name string, src []string) (map[string]any, error) {
	dst := append(src, name)
	tw, err := root.Find(dst)
	if err != nil {
		return nil, err
	}

	m := make(map[string]any)
	m["name"] = tw.Name

	for i := range tw.Childs {
		ctw := tw.Childs[i]
		n := ctw.Name
		v := ctw.get()

		m[n] = v
	}

	return m, nil
}

func (root *Tree) GetMapings(src []string) ([]map[string]any, error) {
	var result []map[string]any

	tw, err := root.Find(src)
	if err != nil {
		return nil, err
	}

	for i := range tw.Childs {
		ctw := tw.Childs[i]
		m := make(map[string]any)
		m, err = root.GetMaping(ctw.Name, src)
		if err != nil {
			return result, err
		}
		result = append(result, m)
	}

	return result, nil
}
