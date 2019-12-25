package cmd

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/lighttiger2505/liary/internal"
	"github.com/urfave/cli"
)

var ConfigCommand = cli.Command{
	Name:    "config",
	Aliases: []string{"c"},
	Usage:   "modify config",
	Action:  ConfigAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "list config",
		},
		cli.StringFlag{
			Name:  "get",
			Usage: "get config value",
		},
	},
	// --get
	// --get-all
	// --set
	// --un-set
	// --list
}

func ConfigAction(c *cli.Context) error {
	cfg, err := internal.GetConfig()
	if err != nil {
		return err
	}

	if c.Bool("list") {
		err := EachField(cfg, func(depth int, name string, value interface{}, kind reflect.Kind) {
			// indent := strings.Repeat("  ", depth)
			// if kind == reflect.Struct {
			// 	format := "%s- FieldName: %s, Type: %v\n"
			// 	fmt.Printf(format, indent, name, kind)
			// } else {
			// 	format := "%s- FieldName: %s, Value: %v, Type: %v\n"
			// 	fmt.Printf(format, indent, name, value, kind)
			// }

			switch kind {
			case reflect.Struct:
				fmt.Printf("%s=%v\n", name, value)
			case reflect.Map:
				fmt.Printf("%s=%v\n", name, value)
			case reflect.Slice:
				fmt.Printf("%s=%v\n", name, value)
			default:
				fmt.Printf("%s=%v\n", name, value)
			}
		})
		if err != nil {
			fmt.Println(err)
		}
		return nil
	}

	if c.String("get") != "" {
		switch c.String("get") {
		case "diarydir":
			fmt.Println(cfg.DiaryDir)
		case "editor":
			fmt.Println(cfg.Editor)
		case "grepcmd":
			fmt.Println(cfg.GrepCmd)
		default:
			return errors.New("key does not contain a section")
		}
		return nil
	}

	internal.OpenEditor(cfg.Editor, cfg.Path())
	return nil
}

func toReflectValue(i interface{}) (reflect.Value, error) {
	v, ok := i.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(i)
	}
	switch v.Kind() {
	case reflect.Ptr:
		// ポインタ以外になるまで reflect.Indirect する
		return toReflectValue(reflect.Indirect(v))
	case reflect.Struct:
		return v, nil
	default:
		return v, errors.New("Not a struct")
	}
}

type EachFunc func(int, string, interface{}, reflect.Kind)

func EachField(s interface{}, fn EachFunc) error {
	depth := 0
	return eachField(depth, s, fn)
}

func eachField(depth int, s interface{}, fn EachFunc) error {
	v, err := toReflectValue(s)
	if err != nil {
		return fmt.Errorf("toReflectValue: %s", err)
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		k := f.Type.Kind()
		vv := v.FieldByName(f.Name)
		fn(depth, f.Name, vv, k)
		if k == reflect.Struct {
			eachField(depth+1, vv, fn)
		}
	}
	return nil
}
