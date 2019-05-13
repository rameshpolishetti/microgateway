package microgateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	// imports the flogo script language
	_ "github.com/project-flogo/core/data/expression/script"
	// imports the basic functions
	_ "github.com/project-flogo/contrib/function"
	// imports the microgateway specific functions
	_ "github.com/project-flogo/microgateway/internal/function"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
	logger "github.com/project-flogo/core/support/log"
	"github.com/project-flogo/microgateway/api"
	"github.com/project-flogo/microgateway/internal/core"
	"github.com/project-flogo/microgateway/internal/schema"
)

// Action is the microgateway action
type Action struct {
	id           string
	settings     Settings
	microgateway *core.Microgateway
	logger       logger.Logger
}

// Manager loads the microgateway definition resource
type Manager struct {
}

func init() {
	action.Register(&Action{}, &Factory{})
	resource.RegisterLoader("microgateway", &Manager{})
}

var actionMetadata = action.ToMetadata(&Settings{}, &Input{}, &Output{})
var resourceMap = make(map[string]*api.Microgateway)

// LoadResource loads the microgateway definition
func (m *Manager) LoadResource(config *resource.Config) (*resource.Resource, error) {
	data := config.Data

	err := schema.Validate(data)
	if err != nil {
		return nil, fmt.Errorf("error validating schema: %s", err.Error())
	}

	var definition *api.Microgateway
	err = json.Unmarshal(data, &definition)
	if err != nil {
		return nil, fmt.Errorf("error marshalling microgateway definition resource with id '%s', %s", config.ID, err.Error())
	}
	return resource.New("microgateway", definition), nil
}

// Factory is a microgateway factory
type Factory struct {
	*resource.Manager
}

type initContext struct {
	settings      map[string]interface{}
	mapperFactory mapper.Factory
	logger        logger.Logger
}

func newInitContext(name string, settings map[string]interface{}, log logger.Logger) *initContext {
	return &initContext{
		settings:      settings,
		mapperFactory: mapper.NewFactory(resolve.GetBasicResolver()),
		logger:        logger.ChildLogger(log, name),
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return i.mapperFactory
}

func (i *initContext) Logger() logger.Logger {
	return i.logger
}

// Initialize initializes the factory
func (f *Factory) Initialize(ctx action.InitContext) error {
	f.Manager = ctx.ResourceManager()
	return nil
}

func (f *Factory) getActionData(act Action) (*api.Microgateway, error) {
	var actionData *api.Microgateway
	if uri := act.settings.URI; uri != "" {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, err
		}
		if resData := api.GetResource(uri); resData != nil {
			actionData = resData
		} else if resData := resourceMap[uri]; resData != nil {
			actionData = resData
		} else if url.Scheme == "http" {
			//get resource from http
			res, err := http.Get(uri)
			if err != nil {
				return nil, fmt.Errorf("Error in accessing specified HTTP url")
			}
			resData, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("Error receving HTTP resource data")
			}
			var definition *api.Microgateway
			err = json.Unmarshal(resData, &definition)
			if err != nil {
				return nil, fmt.Errorf("error marshalling microgateway definition resource")
			}
			resourceMap[uri] = definition
			actionData = definition
		} else if url.Scheme == "file" {
			//get resource from local file system
			resData, err := ioutil.ReadFile(filepath.FromSlash(uri[7:]))
			if err != nil {
				return nil, fmt.Errorf("File reading error")
			}

			err = schema.Validate(resData)
			if err != nil {
				return nil, fmt.Errorf("error validating schema: %s", err.Error())
			}
			var definition *api.Microgateway
			err = json.Unmarshal(resData, &definition)
			if err != nil {
				return nil, fmt.Errorf("error marshalling microgateway definition resource")
			}
			resourceMap[uri] = definition
			actionData = definition
		} else if url.Scheme == "pattern" {
			definition, err := Load(uri[10:])
			if err != nil {
				return nil, err
			}
			actionData = definition
		} else {
			// Load action data from resources
			resData := f.Manager.GetResource(uri)
			if resData == nil {
				return nil, fmt.Errorf("failed to load microgateway URI data: '%s'", act.id)
			}
			actionData = resData.Object().(*api.Microgateway)
		}

		return actionData, nil
	}

	return nil, errors.New("no definition found for microgateway")
}

// GoCode is the microgateway go code representation
type GoCode struct {
	Imports []string
	Action  string
}

// AddImport adds an import to the GoCode
func (g *GoCode) AddImport(ref string) {
	for _, imp := range g.Imports {
		if imp == ref {
			return
		}
	}
}

// GenerateGoCode generates go code from an action
func (f *Factory) GenerateGoCode(settingsName string, config *action.Config) (code GoCode, err error) {
	act := Action{
		id: config.Id,
	}
	if act.id == "" {
		act.id = config.Ref
	}

	err = metadata.MapToStruct(config.Settings, &act.settings, true)
	if err != nil {
		return code, err
	}

	actionData, err := f.getActionData(act)
	if err != nil {
		return code, err
	}

	code.AddImport(`microapi "github.com/project-flogo/microgateway/api"`)
	code.Action += fmt.Sprintf("var %s map[string]interface{}\n", settingsName)
	code.Action += "{\n"
	code.Action += fmt.Sprintf("gateway := microapi.New(\"%s\")\n", actionData.Name)
	services := make(map[string]string)
	for i, service := range actionData.Services {
		code.Action += fmt.Sprintf("service%d := gateway.NewService(\"%s\", &rest.Activity{})\n", i, service.Name)
		services[service.Name] = fmt.Sprintf("services%d", i)
		code.AddImport(service.Ref)
		if service.Description != "" {
			code.Action += fmt.Sprintf("service%d.SetDescription(\"%s\")\n", i, service.Description)
		}
		for key, value := range service.Settings {
			code.Action += fmt.Sprintf("service%d.AddSetting(\"%s\", %#v)\n", i, key, value)
		}
		code.Action += fmt.Sprintf("_ = service%d\n", i)
	}
	for i, step := range actionData.Steps {
		code.Action += fmt.Sprintf("step%d := gateway.NewStep(%s)\n", i, services[step.Service])
		if step.Condition != "" {
			code.Action += fmt.Sprintf("step%d.SetIf(\"%s\")\n", i, step.Condition)
		}
		for key, value := range step.Input {
			code.Action += fmt.Sprintf("step%d.AddInput(\"%s\", %#v)\n", i, key, value)
		}
		if step.HaltCondition != "" {
			code.Action += fmt.Sprintf("step%d.SetHalt(\"%s\")\n", i, step.HaltCondition)
		}
		code.Action += fmt.Sprintf("_ = step%d\n", i)
	}
	for i, response := range actionData.Responses {
		code.Action += fmt.Sprintf("response%d := gateway.NewResponse(%t)\n", i, response.Error)
		if response.Condition != "" {
			code.Action += fmt.Sprintf("response%d.SetIf(\"%s\")\n", i, response.Condition)
		}
		if response.Output.Code != nil {
			code.Action += fmt.Sprintf("response%d.SetCode(%d)\n", i, response.Output.Code)
		}
		if response.Output.Data != nil {
			code.Action += fmt.Sprintf("response%d.SetData(%#v)\n", i, response.Output.Data)
		}
		code.Action += fmt.Sprintf("_ = response%d\n", i)
	}
	code.Action += fmt.Sprintf("var err error\n")
	code.Action += fmt.Sprintf("%s, err = gateway.AddResource(app)\n", settingsName)
	code.Action += fmt.Sprintf("if err != nil {\n")
	code.Action += fmt.Sprintf("panic(err)\n")
	code.Action += fmt.Sprintf("}\n")
	code.Action += "}\n"

	return code, nil
}

// New creates a new microgateway
func (f *Factory) New(config *action.Config) (action.Action, error) {
	log := logger.ChildLogger(logger.RootLogger(), "microgateway")
	act := Action{
		id:     config.Id,
		logger: log,
	}
	if act.id == "" {
		act.id = config.Ref
	}

	err := metadata.MapToStruct(config.Settings, &act.settings, true)
	if err != nil {
		return nil, err
	}

	actionData, err := f.getActionData(act)
	if err != nil {
		return nil, err
	}

	envFlags := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envFlags[pair[0]] = pair[1]
	}
	executionContext := map[string]interface{}{
		"async": act.settings.Async,
		"env":   envFlags,
		"conf":  config.Settings,
	}
	scope := data.NewSimpleScope(executionContext, nil)

	expressionFactory := expression.NewFactory(resolve.GetBasicResolver())
	getExpression := func(name string, value interface{}) (*core.Expr, error) {
		if stringValue, ok := value.(string); ok && len(stringValue) > 0 && stringValue[0] == '=' {
			expr, err := expressionFactory.NewExpr(stringValue[1:])
			if err != nil {
				return nil, err
			}
			return core.NewExpr(name, stringValue, expr), nil
		}
		return core.NewExpr(name, fmt.Sprintf("%v", value), expression.NewLiteralExpr(value)), nil
	}

	services := make(map[string]*core.Service, len(actionData.Services))
	for i := range actionData.Services {
		name := actionData.Services[i].Name
		if _, ok := services[name]; ok {
			return nil, fmt.Errorf("duplicate service name: %s", name)
		}

		values, index := make([]*core.Expr, len(actionData.Services[i].Settings)), 0
		for key, value := range actionData.Services[i].Settings {
			values[index], err = getExpression(key, value)
			if err != nil {
				log.Infof("expression parsing error: %s", value)
				return nil, err
			}
			index++
		}

		settingsMap, err := core.TranslateMappings(scope, values, log)
		if err != nil {
			return nil, err
		}
		settings := make([]core.Setting, len(settingsMap))
		index = 0
		for key, value := range settingsMap {
			settings[index] = core.Setting{
				Name:  key,
				Value: value,
			}
			index++
		}

		if ref := actionData.Services[i].Ref; ref != "" {
			if factory := activity.GetFactory(ref); factory != nil {
				actvt, err := factory(newInitContext(name, settingsMap, log))
				if err != nil {
					return nil, err
				}
				services[name] = &core.Service{
					Name:     name,
					Settings: settings,
					Activity: actvt,
				}
				continue
			}
			actvt := activity.Get(ref)
			if actvt == nil {
				return nil, fmt.Errorf("can't find activity %s", ref)
			}
			services[name] = &core.Service{
				Name:     name,
				Settings: settings,
				Activity: actvt,
			}
		} else if handler := actionData.Services[i].Handler; handler != nil {
			services[name] = &core.Service{
				Name:     name,
				Settings: settings,
				Activity: &core.Adapter{Handler: handler},
			}
		} else {
			return nil, fmt.Errorf("no ref or handler for service: %s", name)
		}
	}

	steps, responses := actionData.Steps, actionData.Responses
	microgateway := core.Microgateway{
		Name:          actionData.Name,
		Async:         act.settings.Async,
		Steps:         make([]core.Step, len(steps)),
		Responses:     make([]core.Response, len(responses)),
		Configuration: config.Settings,
	}
	for j := range steps {
		if condition := steps[j].Condition; condition != "" {
			expr, err := expressionFactory.NewExpr(condition)
			if err != nil {
				log.Infof("condition parsing error: %s", condition)
				return nil, err
			}
			microgateway.Steps[j].Condition = core.NewExpr("condition", condition, expr)
		}

		service := services[steps[j].Service]
		if service == nil {
			return nil, fmt.Errorf("service not found: %s", steps[j].Service)
		}
		microgateway.Steps[j].Service = service

		input := steps[j].Input
		inputExpression, index := make([]*core.Expr, len(input)), 0
		for key, value := range input {
			inputExpression[index], err = getExpression(key, value)
			if err != nil {
				return nil, err
			}
			index++
		}
		microgateway.Steps[j].Input = inputExpression

		if condition := steps[j].HaltCondition; condition != "" {
			expr, err := expressionFactory.NewExpr(condition)
			if err != nil {
				log.Infof("halt condition parsing error: %s", condition)
				return nil, err
			}
			microgateway.Steps[j].HaltCondition = core.NewExpr("halt", condition, expr)
		}
	}

	for j := range responses {
		if condition := responses[j].Condition; condition != "" {
			expr, err := expressionFactory.NewExpr(condition)
			if err != nil {
				log.Infof("condition parsing error: %s", condition)
				return nil, err
			}
			microgateway.Responses[j].Condition = core.NewExpr("condition", condition, expr)
		}

		microgateway.Responses[j].Error = responses[j].Error

		microgateway.Responses[j].Output.Code, err = getExpression("code", responses[j].Output.Code)
		if err != nil {
			return nil, err
		}

		data := responses[j].Output.Data
		if hashMap, ok := data.(map[string]interface{}); ok {
			dataExpressions, index := make([]*core.Expr, len(hashMap)), 0
			for key, value := range hashMap {
				dataExpressions[index], err = getExpression(key, value)
				if err != nil {
					return nil, err
				}
				index++
			}
			microgateway.Responses[j].Output.Datum = dataExpressions
		} else {
			microgateway.Responses[j].Output.Data, err = getExpression("data", data)
			if err != nil {
				return nil, err
			}
		}
	}

	act.microgateway = &microgateway

	return &act, nil
}

// Metadata returns the metadata for the microgateway
func (a *Action) Metadata() *action.Metadata {
	return actionMetadata
}

// IOMetadata returns the iometadata for the microgateway
func (a *Action) IOMetadata() *metadata.IOMetadata {
	return actionMetadata.IOMetadata
}

// Run executes the microgateway
func (a *Action) Run(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	code, mData, err := core.Execute(a.id, input, a.microgateway, a.IOMetadata(), a.logger)
	output := make(map[string]interface{}, 8)
	output["code"] = code
	output["data"] = mData

	return output, err
}
