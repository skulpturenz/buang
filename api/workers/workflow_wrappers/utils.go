package workflowwrappers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"

	"github.com/DataDog/gostackparse"
)

var ErrWorkflowWrapperPanic = errors.New("workflow wrapper panic")

func RecoverWorkflowPanic(name string, err error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case error:
			err = x
		default:
			err = ErrWorkflowWrapperPanic
		}

		slog.ErrorContext(context.Background(), fmt.Sprintf("[workflow wrapper] %v panic: %v", name, r))

		stack := debug.Stack()
		goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

		p := o11y.CreateDiagnosticLogParams{
			Type: enumsdiagnosticlogtype.Panic,
			Log: map[string]any{
				"source": "workflow_wrapper",
				"panic":  fmt.Sprintf("%v", r),
				"stack":  goroutines,
			},
		}
		p.Exec(context.Background())
	}
}
