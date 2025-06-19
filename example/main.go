package main

import (
	"log"
	"strconv"
	"time"

	"github.com/mususu247/tree"
)

func Sample1(srcFile string) error {
	var root tree.Tree
	var indent string

	// json file: \t=tab
	indent = "\t"

	err := root.Create(srcFile, indent)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	root.Close()
	return nil
}

func Sample2(srcFile string) error {
	var root tree.Tree
	var names []string

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	names = []string{}
	err = root.AddNew("work", nil, names)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	err = root.Save()
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	root.Close()
	return nil
}

func Sample3(srcFile string) error {
	var root tree.Tree
	var names []string

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// AddNew childs(work_?)
	names = []string{"work"}
	for i := 0; i < 4; i++ {
		name := "work_" + strconv.Itoa(i)
		err = root.AddNew(name, i, names)
		if err != nil {
			log.Printf("(Error) %v", err)
			return err
		}
	}

	// get child name list
	list, err := root.List(names)
	for i := range list {
		names1 := append(names, list[i])
		val, err := root.GetValue(names1)
		if err != nil {
			log.Printf("(Error) %v", err)
		}
		log.Printf("[%v] name:%v value:%v", i, list[i], val)
	}

	// get child 'name=value;name=value;name=value;'
	child, err := root.GetString(";", names)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	log.Println(child)

	err = root.Save()
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	root.Close()
	return nil
}

func Sample4(srcFile string) error {
	var root tree.Tree
	var names []string

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// .AddNew()
	names = []string{"work", "work_0"}
	for i := 0; i < 3; i++ {
		name := "val_" + strconv.Itoa(i)
		err = root.AddNew(name, i*i, names)
		if err != nil {
			log.Printf("(Error) %v", err)
			return err
		}

		names1 := append(names, name)
		for j := 0; j < 3; j++ {
			name := "float_" + strconv.Itoa(j)
			f64 := float64(j) * 3.14
			err = root.AddNew(name, f64, names1)
			if err != nil {
				log.Printf("(Error) %v", err)
				return err
			}
		}
	}

	//.Copy()
	err = root.Copy([]string{"work", "work_0"}, []string{"work", "work_1"})
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// .GetValueInt() & setValue()
	names = []string{"work", "work_1"}
	list, err := root.List(names)
	names = append(names, "xxx")
	for i := range list {
		names[2] = list[i]
		v, err := root.GetValueInt(names)
		if err != nil {
			log.Printf("(Error) %v", err)
		}

		v = 1 + (v * 10)

		err = root.SetValue(v, names)
		if err != nil {
			log.Printf("(Error) %v", err)
		}
	}

	// get child 'name=value;name=value;name=value;'
	child, err := root.GetString(";", []string{"work", "work_0"})
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	log.Println(child)

	// get child 'name=value;name=value;name=value;'
	child, err = root.GetString(";", []string{"work", "work_1"})
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	log.Println(child)

	err = root.Save()
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	root.Close()
	return nil
}

func Sample5(srcFile string, dstFile string) error {
	var root tree.Tree
	var names []string

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	//.Move()
	err = root.Move([]string{"work", "work_0"}, []string{"work", "work_2"})
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// .GetValueInt() & setValue()
	names = []string{"work", "work_2"}
	list, err := root.List(names)
	names = append(names, "xxx")
	for i := range list {
		names[2] = list[i]
		v, err := root.GetValueInt(names)
		if err != nil {
			log.Printf("(Error) %v", err)
		}

		v = 2 + (v * 20)

		err = root.SetValue(v, names)
		if err != nil {
			log.Printf("(Error) %v", err)
		}
	}

	// get child 'name=value;name=value;name=value;'
	child, err := root.GetString(";", []string{"work", "work_2"})
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}
	log.Println(child)

	err = root.SaveAs(dstFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	root.Close()
	return nil
}

func Sample6(srcFile string, dstFile string) error {
	var root tree.Tree
	var names []string

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// create map[string]any
	m := make(map[string]any)
	m["name"] = "field"
	m["str"] = "str"
	m["int"] = int64(7)
	m["float"] = float64(3.14)
	m["bool"] = true
	m["date"] = time.Now()

	// write map[string]any
	names = []string{"work", "work_2"}
	err = root.SetMaping(m, names)
	if err != nil {
		log.Printf("(Error) %v", err)
	}

	// .SaveAs(fileName)
	err = root.SaveAs(dstFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// clear map[string]any
	m = nil

	// read map[string]any
	m, err = root.GetMaping("field", names)
	if err != nil {
		log.Printf("(Error) %v", err)
	}
	log.Printf("%v", m)

	// .GetString()
	names = append(names, "field")
	childs, err := root.GetString(" ", names)
	if err != nil {
		log.Printf("(Error) %v", err)
	}
	log.Printf("%v", childs)

	root.Close()
	return nil
}

func Sample7(srcFile string, dstFile string) error {
	var root tree.Tree
	var names []string
	var ms []map[string]any

	err := root.Open(srcFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// create []map[string]any
	for i := 0; i < 4; i++ {
		m := make(map[string]any)
		m["name"] = "field_" + strconv.Itoa(i)
		m["str"] = "str"
		m["int"] = i
		m["float"] = float64(3.14)
		m["bool"] = true
		m["date"] = time.Now()
		ms = append(ms, m)
	}

	// Addnew("worl","work_3","table")
	names = []string{"work", "work_3"}
	root.AddNew("table", nil, names)

	// .setMaoings()
	names = []string{"work", "work_3", "table"}
	root.SetMapings(ms, names)

	err = root.SaveAs(dstFile)
	if err != nil {
		log.Printf("(Error) %v", err)
		return err
	}

	// clrear []map[string]any
	ms = nil

	// .getMapings()
	ms, err = root.GetMapings(names)
	if err != nil {
		log.Printf("(Error) %v", err)
	}

	for i := range ms {
		m := ms[i]
		log.Printf("[%v] %v", i, m)
	}

	// .GetString()
	list, err := root.List(names)
	names = append(names, "xxx")
	for i := range list {
		names[3] = list[i]

		childs, err := root.GetString(", ", names)
		if err != nil {
			log.Printf("(Error) %v", err)
		}
		log.Printf("[%v] %v", list[i], childs)
	}

	root.Close()
	return nil
}

func main() {
	fileNmae1 := "./config.json"
	fileNmae2 := "./config2.json"
	fileNmae3 := "./config3.json"

	// .Create(): json file & .Close()
	Sample1(fileNmae1)

	// .AddNew(): tree("work"), .Save() json file
	Sample2(fileNmae1)

	// .AddNew(): tree.childs("work_?"), .List(), .GetValue(), .Save(): json file
	Sample3(fileNmae1)

	// .AddNew(): tree.childs("val_?"), .Copy(), .GetValue?() .SetValue() .SaveAs
	Sample4(fileNmae1)

	// .Move, .GetValue?(), .SetValue(), SaveAs()
	Sample5(fileNmae1, fileNmae2)

	// .SetMaping() .GetMaping() .GetString()
	Sample6(fileNmae1, fileNmae3)

	// .SetMapings() .GetMapings() .GetString()
	Sample7(fileNmae1, fileNmae3)
}
