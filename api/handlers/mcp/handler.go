package mcp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"skulpture/buang/app"
	constantsenvs "skulpture/buang/constants/envs"
	"skulpture/buang/handlers/mcp/tools"
	"strings"
)

type jsonRpcRequest struct {
	Jsonrpc string         `json:"jsonrpc"`
	Id      any            `json:"id,omitempty"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

type jsonRpcResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      any           `json:"id,omitempty"`
	Result  any           `json:"result,omitempty"`
	Error   *jsonRpcError `json:"error,omitempty"`
}

type jsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolCallResult struct {
	Content []toolContent `json:"content"`
	IsError bool          `json:"isError"`
}

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolDescription struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func Handler(s app.ApplicationServices) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req jsonRpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(jsonRpcResponse{
				Jsonrpc: "2.0",
				Error:   &jsonRpcError{Code: -32700, Message: "parse error"},
			})
			return
		}
		defer r.Body.Close()

		// notifications have no id and require no response
		if req.Id == nil && strings.HasPrefix(req.Method, "notifications/") {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		var result any
		var rpcErr *jsonRpcError

		switch req.Method {
		case "initialize":
			result = map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]any{
					"tools": map[string]any{},
				},
				"serverInfo": map[string]any{
					"name":    "buang",
					"version": constantsenvs.BUANG_VERSION,
				},
			}

		case "tools/list":
			list := make([]toolDescription, 0, len(tools.AllTools))
			for _, t := range tools.AllTools {
				list = append(list, toolDescription{
					Name:        t.Name,
					Description: t.Description,
					InputSchema: t.InputSchema,
				})
			}
			result = map[string]any{"tools": list}

		case "tools/call":
			name, _ := req.Params["name"].(string)
			arguments, _ := req.Params["arguments"].(map[string]any)
			if arguments == nil {
				arguments = map[string]any{}
			}

			var matched *tools.Tool
			for i := range tools.AllTools {
				if tools.AllTools[i].Name == name {
					matched = &tools.AllTools[i]
					break
				}
			}

			if matched == nil {
				result = toolCallResult{
					Content: []toolContent{{Type: "text", Text: "unknown tool: " + name}},
					IsError: true,
				}
				break
			}

			text, err := matched.Handler(r.Context(), s, arguments)
			if err != nil {
				slog.ErrorContext(r.Context(), "mcp tool call", "tool", name, "err", err.Error())
				result = toolCallResult{
					Content: []toolContent{{Type: "text", Text: err.Error()}},
					IsError: true,
				}
				break
			}

			result = toolCallResult{
				Content: []toolContent{{Type: "text", Text: text}},
				IsError: false,
			}

		default:
			rpcErr = &jsonRpcError{Code: -32601, Message: "method not found: " + req.Method}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jsonRpcResponse{
			Jsonrpc: "2.0",
			Id:      req.Id,
			Result:  result,
			Error:   rpcErr,
		})
	}
}
