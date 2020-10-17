package commands

import (
	"errors"
	"flag"
	"reflect"
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

func NewDefaultModelProvider(modelType reflect.Type) ModelProvider {
	var supportMap map[string]ParamSupport

	//参数模型的参数设置
	for i := 0; i < modelType.NumField(); i++ {
		filed := modelType.Field(i)
		support := ParamSupport{
			FiledName:  filed.Name,
			Name:       filed.Name,
			Usage:      "",
			BaseStruct: modelType,
			Kind:       filed.Type.Kind(),
		}
		supportMap[filed.Name] = support
	}

	return &ModelProviderBase{
		modelType: modelType,
		supports:  supportMap,
	}
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

type CommandBase struct {
	Command

	runner Runner

	flags flag.FlagSet

	modelProvider ModelProvider
}

func NewCommand(runner Runner, flags flag.FlagSet, provider ModelProvider) Command {
	return &CommandBase{runner: runner, flags: flags, modelProvider: provider}
}

func (cmd *CommandBase) Execute(str string) error {

	//构建模型
	model := cmd.modelProvider.GetParamsModel()

	//复制flag解析器
	flagCopy := flag.NewFlagSet(cmd.flags.Name(), cmd.flags.ErrorHandling())
	StructCopy(flagCopy, cmd.flags)

	//获取参数
	modelType := reflect.TypeOf(model)
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
	err := flagCopy.Parse(strings.Fields(str))
	if err != nil {
		return err
	}

	//赋值
	modelValue := reflect.ValueOf(&model)
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
	params := Params{Str: str, Info: model}

	//将参数传入runner
	return cmd.runner.Run(params)
}
