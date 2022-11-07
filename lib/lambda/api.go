package lambda

import (
	"time"

	"github.com/Shopify/go-lua"
	"github.com/sirupsen/logrus"
)

func createBaseState(se *StorageEntry) *lua.State {
	l := lua.NewState()
	lua.Require(l, "math", lua.MathOpen, true)
	lua.Require(l, "string", lua.StringOpen, true)
	lua.Require(l, "table", lua.TableOpen, true)

	_ = lua.NewMetaTable(l, "neosMetaTable")
	lua.SetFunctions(l, []lua.RegistryFunction{
		{
			Name: "send",
			Function: func(l *lua.State) int {
				varName := lua.CheckString(l, 1)
				varVal := l.ToValue(2)

				select {
				case se.writeChan <- MessageTpl{
					VarName: varName,
					VarVal:  varVal,
				}:
					logrus.WithFields(logrus.Fields{
						"namespace": se.session,
						"var_name":  varName,
						"var_value": varVal,
					}).Info("send var to ws")
				case <-time.After(time.Second):
					logrus.WithFields(logrus.Fields{
						"namespace": se.session,
					}).Info("send timeout")
				}

				return 0
			},
		},
		{
			Name: "persist",
			Function: func(l *lua.State) int {
				varName := lua.CheckString(l, 1)
				varVal := l.ToValue(2)

				se.persistence[varName] = varVal
				logrus.WithFields(logrus.Fields{
					"namespace": se.session,
					"var_name":  varName,
					"var_value": varVal,
				}).Info("updated persistent var")
				return 0
			},
		},
		{
			Name: "load",
			Function: func(l *lua.State) int {
				varName := lua.CheckString(l, 1)

				logrus.WithFields(logrus.Fields{
					"namespace": se.session,
					"var_name":  varName,
				}).Info("read persistent var")

				l.PushLightUserData(se.persistence[varName])
				return 1
			},
		},
		{
			Name: "tonumber",
			Function: func(l *lua.State) int {
				source, _ := l.ToNumber(1)
				l.PushLightUserData(source)
				return 1
			},
		},
		{
			Name: "tostring",
			Function: func(state *lua.State) int {
				lua.CheckAny(l, 1)
				lua.ToStringMeta(l, 1)
				return 1
			},
		},
	}, 0)
	l.SetGlobal("neos")
	lua.SetMetaTableNamed(l, "neosMetaTable")

	return l
}
