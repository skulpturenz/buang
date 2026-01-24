package workers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"skulpture/buang/components/o11y"
	enumsdiagnosticlogtype "skulpture/buang/enums/diagnostic_log_type"
	enumsdurableexecutors "skulpture/buang/enums/durable_executors"

	"github.com/DataDog/gostackparse"
)

var ErrWorkflowPanic = errors.New("workflow panic")

func RecoverWorkflowPanic(durableExecutor enumsdurableexecutors.DurableExecutor, name string, err error) {
	if r := recover(); r != nil {
		switch x := r.(type) {
		case error:
			err = x
		default:
			err = ErrWorkflowPanic
		}

		slog.ErrorContext(context.Background(), fmt.Sprintf("[%v] %v panic: %v", durableExecutor, name, r))

		stack := debug.Stack()
		goroutines, _ := gostackparse.Parse(bytes.NewReader(stack))

		p := o11y.CreateDiagnosticLogParams{
			Type: enumsdiagnosticlogtype.Panic,
			Log: map[string]any{
				"source": durableExecutor.String(),
				"panic":  fmt.Sprintf("%v", r),
				"stack":  goroutines,
			},
		}
		p.Exec(context.Background())
	}
}
