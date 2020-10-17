package commands

import (
	"bytes"
	"errors"
	"flag"
	"reflect"
	"strconv"
	"strings"
)

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type Runner interface {
	Run(params Params) error
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type Params struct {
	Str string

	Info interface{}

	Args []string

	Bundle map[string]interface{}
}

type ParamSupport struct {
	FiledName string

	Name string

	Usage string

	BaseStruct reflect.Type

	Kind reflect.Kind

	DefaultValue interface{}
}

func NewParamSupport(baseStruct reflect.Type, filedName string, name string, value string, usage string) ParamSupport {
	filed, _ := baseStruct.FieldByName(filedName)
	return ParamSupport{
		FiledName:    filedName,
		Name:         name,
		Usage:        usage,
		Kind:         filed.Type.Kind(),
		BaseStruct:   baseStruct,
		DefaultValue: value,
	}
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type ModelProvider interface {
	GetParamSupport(filedName string) (ParamSupport, bool)

	GetParamsModel() interface{}
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type ModelProviderBase struct {
	ModelProvider
	modelType reflect.Type
	supports  map[string]ParamSupport
}

func NewDefaultModelProvider(modelType reflect.Type) (ModelProvider, error) {
	var supportMap map[string]ParamSupport = map[string]ParamSupport{}

	//参数模型的参数设置
	for i := 0; i < modelType.NumField(); i++ {

		filed := modelType.Field(i)
		tag := filed.Tag

		//设置默认值
		var defaultValue interface{}
		var err error
		if v, exist := tag.Lookup("value"); exist {
			switch filed.Type.Kind() {
			case reflect.String:
				defaultValue = v
			case reflect.Bool:
				defaultValue, err = strconv.ParseBool(v)
			case reflect.Int:
				defaultValue, err = strconv.Atoi(v)
			case reflect.Int8:
				defaultValue, err = strconv.ParseInt(v, 10, 8)
				defaultValue = int8(defaultValue.(int64))
			case reflect.Int16:
				defaultValue, err = strconv.ParseInt(v, 10, 16)
				defaultValue = int16(defaultValue.(int64))
			case reflect.Int32:
				defaultValue, err = strconv.ParseInt(v, 10, 32)
				defaultValue = int32(defaultValue.(int64))
			case reflect.Int64:
				defaultValue, err = strconv.ParseInt(v, 10, 64)
				defaultValue = defaultValue.(int64)
			case reflect.Uint:
				defaultValue, err = strconv.ParseUint(v, 10, 32)
				defaultValue = uint(defaultValue.(uint64))
			case reflect.Uint8:
				defaultValue, err = strconv.ParseUint(v, 10, 8)
				defaultValue = uint8(defaultValue.(uint64))
			case reflect.Uint16:
				defaultValue, err = strconv.ParseUint(v, 10, 16)
				defaultValue = uint16(defaultValue.(uint64))
			case reflect.Uint32:
				defaultValue, err = strconv.ParseUint(v, 10, 32)
				defaultValue = uint32(defaultValue.(uint64))
			case reflect.Uint64:
				defaultValue, err = strconv.ParseUint(v, 10, 64)
				defaultValue = defaultValue.(uint64)
			case reflect.Float32:
				defaultValue, err = strconv.ParseFloat(v, 32)
				defaultValue = float32(defaultValue.(float64))
			case reflect.Float64:
				defaultValue, err = strconv.ParseFloat(v, 64)
				defaultValue = defaultValue.(float64)
				break
			}

			if err != nil {
				return nil, err
			}
		} else {
			switch filed.Type.Kind() {
			case reflect.String:
				defaultValue = ""
			case reflect.Bool:
				defaultValue = false
			case reflect.Int:
				defaultValue = int(0)
			case reflect.Int8:
				defaultValue = int8(0)
			case reflect.Int16:
				defaultValue = int16(0)
			case reflect.Int32:
				defaultValue = int32(0)
			case reflect.Int64:
				defaultValue = int64(0)
			case reflect.Uint:
				defaultValue = uint(0)
			case reflect.Uint8:
				defaultValue = uint8(0)
			case reflect.Uint16:
				defaultValue = uint16(0)
			case reflect.Uint32:
				defaultValue = uint32(0)
			case reflect.Uint64:
				defaultValue = uint64(0)
			case reflect.Float32:
				defaultValue = float32(0)
			case reflect.Float64:
				defaultValue = float64(0)
				break
			}
		}

		//获取使用方法
		var usage string
		if v, exist := tag.Lookup("usage"); exist {
			usage = v
		} else {
			usage = filed.Name + " for " + modelType.Name()
		}

		//获取参数名
		var supportName string
		if v, exist := tag.Lookup("name"); exist {
			supportName = v
		} else {
			//将驼峰命名转换为 小写+横杠 的模式
			//如 : NewBee --> new-bee
			strBytes := []byte(filed.Name)
			buffer := new(bytes.Buffer)
			for i := 0; i < len(strBytes); i++ {
				nowChar := strBytes[i]
				if 'A' < nowChar && nowChar < 'Z' {
					if i > 0 {
						buffer.WriteByte('-') //添加横杠
					}
					buffer.WriteByte(nowChar + 32) //转换
				} else {
					buffer.WriteByte(nowChar)
				}
			}
			supportName = buffer.String()
		}

		//生成ParamSupport
		support := ParamSupport{
			FiledName:    filed.Name,
			Name:         supportName,
			Usage:        usage,
			BaseStruct:   modelType,
			Kind:         filed.Type.Kind(),
			DefaultValue: defaultValue,
		}

		supportMap[filed.Name] = support

	}

	return &ModelProviderBase{
		modelType: modelType,
		supports:  supportMap,
	}, nil
}

func NewModelProvider(modelType reflect.Type, paramSupports []ParamSupport) ModelProvider {
	var supports map[string]ParamSupport
	for i := 0; i < len(paramSupports); i++ {
		supports[paramSupports[i].FiledName] = paramSupports[i]
	}
	return &ModelProviderBase{
		modelType: modelType,
		supports:  supports,
	}
}

func (provider *ModelProviderBase) GetParamsModel() interface{} {
	return reflect.New(provider.modelType).Interface()
}

func (provider *ModelProviderBase) GetParamSupport(filedName string) (ParamSupport, bool) {
	if support, exist := provider.supports[filedName]; exist {
		return support, true
	} else {
		return ParamSupport{}, false
	}
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type ErrorHandling flag.ErrorHandling

const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2) or for -h/-help Exit(0).
	PanicOnError                         // Call panic with a descriptive error.
)

type CommandBase struct {
	Command

	runner func(params Params) error

	flags *flag.FlagSet

	modelProvider ModelProvider
}

func NewCommand(commandName string, errorHandling ErrorHandling, runner func(params Params) error, provider ModelProvider) Command {
	flags := flag.NewFlagSet(commandName, flag.ErrorHandling(errorHandling))
	return &CommandBase{
		runner:        runner,
		flags:         flags,
		modelProvider: provider}
}

func CreateDefaultCommand(command string, modelType reflect.Type, runner func(params Params) error) Command {
	provider, err := NewDefaultModelProvider(modelType)
	if err != nil {
		return nil
	}
	return NewCommand(command, ContinueOnError, runner, provider)
}

func (cmd *CommandBase) Execute(str string, bundle map[string]interface{}) error {

	//构建模型
	model := cmd.modelProvider.GetParamsModel()

	//复制flag解析器
	flagCopy := flag.NewFlagSet(cmd.flags.Name(), cmd.flags.ErrorHandling())

	//获取参数
	modelType := reflect.TypeOf(model).Elem()
	valueMap := map[string]interface{}{}

	//参数模型的参数设置
	for i := 0; i < modelType.NumField(); i++ {
		filed := modelType.Field(i)
		support, exist := cmd.modelProvider.GetParamSupport(filed.Name)

		if !exist {
			continue
		} else {

			//判断该参数所属的结构体是否正确
			if support.BaseStruct != modelType {
				return errors.New(modelType.Name() + ":类型实例获取不正确！")
			}

			//检查内部参数类型
			if support.Kind == filed.Type.Kind() {
				//类型转换
				switch support.Kind {
				case reflect.String:
					tmpStr := flagCopy.String(support.Name, support.DefaultValue.(string), support.Usage)
					valueMap[filed.Name] = tmpStr
					break
				case reflect.Bool:
					tmpBool := flagCopy.Bool(support.Name, support.DefaultValue.(bool), support.Usage)
					valueMap[filed.Name] = tmpBool
					break
				case reflect.Int:
					tmpInt := flagCopy.Int(support.Name, support.DefaultValue.(int), support.Usage)
					valueMap[filed.Name] = tmpInt
					break
				case reflect.Int64:
					tmpInt64 := flagCopy.Int64(support.Name, support.DefaultValue.(int64), support.Usage)
					valueMap[filed.Name] = tmpInt64
					break
				case reflect.Uint:
					tmpUint := flagCopy.Uint(support.Name, support.DefaultValue.(uint), support.Usage)
					valueMap[filed.Name] = tmpUint
					break
				case reflect.Uint64:
					tmpUint64 := flagCopy.Uint64(support.Name, support.DefaultValue.(uint64), support.Usage)
					valueMap[filed.Name] = tmpUint64
					break
				case reflect.Float64:
					tmpFloat64 := flagCopy.Float64(support.Name, support.DefaultValue.(float64), support.Usage)
					valueMap[filed.Name] = tmpFloat64
					break
				default:
					return errors.New(filed.Name + ":不能解析的类型[" + filed.Type.Name() + "]")
				}

			} else {
				return errors.New(filed.Name + ":参数类型有误！")
			}

		}
	}

	//对命令进行解析
	err := flagCopy.Parse(strings.Fields(str)[1:])
	if err != nil {
		return err
	}

	//赋值
	modelValue := reflect.ValueOf(model)
	modelValue = modelValue.Elem() //添加这一个才能赋值

	for i := 0; i < modelType.NumField(); i++ {
		filed := modelType.Field(i)
		value := modelValue.FieldByName(filed.Name)

		//类型转换与赋值
		switch filed.Type.Kind() {
		case reflect.String:
			value.SetString(*(valueMap[filed.Name].(*string)))
			break
		case reflect.Bool:
			value.SetBool(*(valueMap[filed.Name].(*bool)))
			break
		case reflect.Int:
			value.SetInt(int64(*(valueMap[filed.Name].(*int))))
			break
		case reflect.Int64:
			value.SetInt(*(valueMap[filed.Name].(*int64)))
			break
		case reflect.Uint:
			value.SetUint(uint64(*(valueMap[filed.Name].(*uint))))
			break
		case reflect.Uint64:
			value.SetUint(*(valueMap[filed.Name].(*uint64)))
			break
		case reflect.Float64:
			value.SetFloat(*(valueMap[filed.Name].(*float64)))
			break
		default:
			return errors.New(filed.Name + ":不能解析的类型[" + filed.Type.Name() + "]")
		}

	}

	//设置参数
	params := Params{
		Str:    str,
		Info:   model,
		Args:   flagCopy.Args(),
		Bundle: bundle,
	}

	//将参数传入runner
	return cmd.runner(params)
}
