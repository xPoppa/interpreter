package evaluator

import (
	"fmt"

	"github.com/xPoppa/interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to len not supported got INTEGER")
			}
		}},

	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
			default:
				return newError(fmt.Sprintf("argument to `first` must be ARRAY, got %s", arg.Type()))
			}
			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
			default:
				return newError(fmt.Sprintf("argument to `last` must be ARRAY, got %s", arg.Type()))
			}
			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if len(arg.Elements) > 0 {
					newArray := make([]object.Object, length-1, length-1)
					copy(newArray, arg.Elements[1:])
					return &object.Array{Elements: newArray}
				}
			default:
				return newError(fmt.Sprintf("argument to `rest` must be ARRAY, got %s", arg.Type()))
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments got=%d, want=2", len(args))
			}

			arg1, ok := args[0].(*object.Array)
			if !ok {
				return newError(fmt.Sprintf("argument to `push` must be ARRAY, got %s", args[0].Type()))
			}
			length := len(arg1.Elements)
			newElements := make([]object.Object, length, length+1)
			copy(newElements, arg1.Elements)
			newElements = append(newElements, args[1])
			return &object.Array{Elements: newElements}
		},
	},
}
