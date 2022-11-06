package lambda

import (
	"time"

	"github.com/Shopify/go-lua"
	"github.com/sirupsen/logrus"
)

func createBaseState(se *StorageEntry) *lua.State {
	l := lua.NewState()
	lua.BaseOpen(l)

	_ = lua.NewMetaTable(l, "neosMetaTable")
	lua.SetFunctions(l, []lua.RegistryFunction{
		{
			Name: "update",
			Function: func(l *lua.State) int {
				varName := lua.CheckString(l, 1)
				varVal, exists := se.persistence[varName]
				if !exists {
					return 1
				}

				select {
				case se.writeChan <- MessageTpl{
					VarName: varName,
					VarVal:  varVal,
				}:
					logrus.WithFields(logrus.Fields{
						"namespace": se.Namespace,
						"var_name":  varName,
						"var_value": varVal,
					}).Info("send var to ws")
				case <-time.After(time.Second):
					logrus.WithFields(logrus.Fields{
						"namespace": se.Namespace,
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
					"namespace": se.Namespace,
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
					"namespace": se.Namespace,
					"var_name":  varName,
				}).Info("read persistent var")

				l.PushLightUserData(se.persistence[varName])
				return 1
			},
		},
	}, 0)
	l.SetGlobal("neos")
	lua.SetMetaTableNamed(l, "neosMetaTable")

	return l
}
