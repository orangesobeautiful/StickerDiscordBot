package hs

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
)

type typeHandle struct {
	idx      int
	isPtr    bool
	dealFunc func(reflect.Value, bool, []string) error
}

// Engine is the http server engine
type Engine struct {
	router *httprouter.Router
}

// New create a new Engine
func New() *Engine {
	return &Engine{
		router: httprouter.New(),
	}
}

func handleFuncCheck(handle any) (
	handleValue reflect.Value, handleInNum int, reqType reflect.Type, paramIdxMap, queryIdxMap map[string]typeHandle, respErrIdx int, respDataIdx int) {
	handleValue = reflect.ValueOf(handle)
	handleType := handleValue.Type()

	// handle func format check

	handleInNum = handleType.NumIn()
	handleOutNum := handleType.NumOut()
	if handleInNum < 1 || handleInNum > 2 {
		panic("handle function is not valid")
	}
	if handleOutNum > 2 {
		panic("handle function is not valid")
	}
	if handleType.In(0) != reflect.TypeOf(&Context{}) {
		panic("handle function's 1th input parameter(req) type need to be *Context")
	}

	if handleInNum == 2 {
		reqType = handleType.In(1)
		if reqType.Kind() != reflect.Pointer {
			panic("handle function's 2th input parameter(req) need to be pointer")
		}
	}

	respDataIdx, respErrIdx = -1, -1
	switch handleOutNum {
	case 1:
		respDataIdx = 0
		if handleType.Out(0) == reflect.TypeOf(&ErrResp{}) {
			respDataIdx = -1
			respErrIdx = 0
		}
	case 2:
		if handleType.Out(0).Kind() != reflect.Pointer {
			panic("handle function's data output parameter type need to be pointer")
		}
		if handleType.Out(1) != reflect.TypeOf(&ErrResp{}) {
			panic("handle function's last output parameter type need to be *ErrResp")
		}

		respDataIdx = 0
		respErrIdx = 1
	}

	queryIdxMap = make(map[string]typeHandle)
	paramIdxMap = make(map[string]typeHandle)
	if handleInNum == 2 {
		for i := 0; i < reqType.Elem().NumField(); i++ {
			field := reqType.Elem().Field(i)
			queryTag := field.Tag.Get("query")
			paramTag := field.Tag.Get("param")
			if queryTag == "-" {
				queryTag = ""
			}
			if paramTag == "-" {
				paramTag = ""
			}
			if queryTag != "" && paramTag != "" {
				panic("query tag and param tag can not be used together")
			}

			if queryTag != "" {
				queryIdxMap[queryTag] = structFieldHandle(i, field)
			} else if paramTag != "" {
				paramIdxMap[paramTag] = structFieldHandle(i, field)
			}
		}
	}

	return handleValue, handleInNum, reqType, paramIdxMap, queryIdxMap, respErrIdx, respDataIdx
}

func structFieldHandle(index int, field reflect.StructField) typeHandle {
	var res typeHandle
	res.idx = index

	if field.Type.Kind() == reflect.Ptr {
		res.isPtr = true
		field.Type = field.Type.Elem()
	}

	switch field.Type.Kind() {
	case reflect.String:
		res.dealFunc = setString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res.dealFunc = setInt64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res.dealFunc = setUint64
	case reflect.Float32, reflect.Float64:
		res.dealFunc = setFloat64
	case reflect.Bool:
		res.dealFunc = setBool
	case reflect.Complex64, reflect.Complex128:
		res.dealFunc = setComplex128
	}

	return res
}

// Handle register a handle function to the router
func (e *Engine) Handle(method, path string, handle any) {
	handleValue, handleInNum, reqType, paramIdxMap, queryIdxMap, respErrIdx, respDataIdx := handleFuncCheck(handle)

	e.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := newContext(r, w)
		defer putContext(ctx)

		in := make([]reflect.Value, 0, handleInNum)
		in = append(in, reflect.ValueOf(ctx))

		// decode request

		if handleInNum == 2 {
			req := reflect.New(reqType.Elem()).Interface()

			var err error
			dec := json.NewDecoder(r.Body)
			if err = dec.Decode(req); err != nil {
				// request decode failed

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				errResp := ErrResp{
					Message: "json body of request decode failed",
					Detail:  []string{err.Error()},
				}
				_, _ = w.Write(errResp.ToJSONBytes())
				return
			}
			pLen := len(p)
			for i := 0; i < pLen; i++ {
				handleInfo, exist := paramIdxMap[p[i].Key]
				if exist {
					err = handleInfo.dealFunc(
						reflect.ValueOf(req).Elem().Field(handleInfo.idx),
						handleInfo.isPtr,
						[]string{p[i].Value},
					)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
						errResp := ErrResp{
							Message: "param of request decode failed",
							Detail:  []string{err.Error()},
						}
						_, _ = w.Write(errResp.ToJSONBytes())
						return
					}
				}
			}
			querys := r.URL.Query()
			for k, v := range querys {
				handleInfo, exist := queryIdxMap[k]
				if exist {
					err = handleInfo.dealFunc(
						reflect.ValueOf(req).Elem().Field(handleInfo.idx),
						handleInfo.isPtr,
						v)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
						errResp := ErrResp{
							Message: "param of request decode failed",
							Detail:  []string{err.Error()},
						}
						_, _ = w.Write(errResp.ToJSONBytes())
						return
					}
				}
			}

			in = append(in, reflect.ValueOf(req))
		}

		out := handleValue.Call(in)

		// error handle

		var errResp *ErrResp
		if respErrIdx >= 0 {
			errResp = out[respErrIdx].Interface().(*ErrResp)
		}
		if errResp != nil {
			ctx.setJSONHeader()
			w.WriteHeader(errResp.Status)
			_, _ = w.Write(errResp.ToJSONBytes())
			return
		}

		if respDataIdx >= 0 {
			_ = ctx.writeJSON(out[respDataIdx].Interface())
		}
	})
}

// GET register a GET handle function to the router
func (e *Engine) GET(path string, handle any) {
	e.Handle(http.MethodGet, path, handle)
}

// POST register a POST handle function to the router
func (e *Engine) POST(path string, handle any) {
	e.Handle(http.MethodPost, path, handle)
}

// ServeHTTP implements http.Handler interface
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.ServeHTTP(w, req)
}
