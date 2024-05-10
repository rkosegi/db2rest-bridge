// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

// NameList List of names, such as backends or entities
type NameList = []string

// UntypedDto Unstructured content, dictionary of string-to-any values.
type UntypedDto map[string]interface{}

// Backend defines model for backend.
type Backend = string

// Entity defines model for entity.
type Entity = string

// ListItemsParams defines parameters for ListItems.
type ListItemsParams struct {
	// PageOffset Page offset
	PageOffset *int `form:"page-offset,omitempty" json:"page-offset,omitempty"`

	// PageSize Page size
	PageSize *int `form:"page-size,omitempty" json:"page-size,omitempty"`

	// Order List of order instructions in form of `key=direction`.
	// Key represents entity field (column) and direction is one of `ASC` or `DESC`,
	// for example `name=ASC` or `id=DESC`.
	// Direction can be omitted, in such case `ASC` is assumed.
	Order *[]string `form:"order[],omitempty" json:"order[],omitempty"`

	// Filter Filter is JSON-encoded FilterExpression.
	// Currently supported types are `simple`, `not` and `junction`.
	// Examples:
	//
	// - `{"simple": { "name": "id", "op": "=", "val" : 1}}`
	//
	//    is equivalent to SQL `id=1`
	//
	// - `{"not": { "simple": { "name": "id", "op": ">", "val" : 100}}}`
	//
	//    is equivalent to SQL `NOT (id>100)`
	//
	// - `{"junction": {"op": "AND", "sub" : [{"simple": { "name": "age", "op": ">", "val" : 35}}, {"simple": { "name": "salary", "op": ">", "val" : 5000}}]}}`
	//
	//    is equivalent to SQL `(age>35) AND (salary > 5000)`
	Filter *string `form:"filter,omitempty" json:"filter,omitempty"`
}

// CreateItemJSONRequestBody defines body for CreateItem for application/json ContentType.
type CreateItemJSONRequestBody = UntypedDto

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all configured backends
	// (GET /backends)
	ListBackends(w http.ResponseWriter, r *http.Request)
	// List all known entities within backend
	// (GET /{backend}/entities)
	ListEntities(w http.ResponseWriter, r *http.Request, backend Backend)
	// List entity items
	// (GET /{backend}/{entity})
	ListItems(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity, params ListItemsParams)
	// Create new entity item
	// (POST /{backend}/{entity})
	CreateItem(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity)
	// Delete entity item by ID
	// (DELETE /{backend}/{entity}/{id})
	DeleteItemById(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity, id string)
	// Get entity item by ID
	// (GET /{backend}/{entity}/{id})
	GetItemById(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity, id string)
	// Check for existence of entity item by ID
	// (HEAD /{backend}/{entity}/{id})
	ExistsItemById(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity, id string)
	// Update entity item in-place by ID
	// (PUT /{backend}/{entity}/{id})
	UpdateItemById(w http.ResponseWriter, r *http.Request, backend Backend, entity Entity, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ListBackends operation middleware
func (siw *ServerInterfaceWrapper) ListBackends(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListBackends(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ListEntities operation middleware
func (siw *ServerInterfaceWrapper) ListEntities(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListEntities(w, r, backend)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ListItems operation middleware
func (siw *ServerInterfaceWrapper) ListItems(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params ListItemsParams

	// ------------- Optional query parameter "page-offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "page-offset", r.URL.Query(), &params.PageOffset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page-offset", Err: err})
		return
	}

	// ------------- Optional query parameter "page-size" -------------

	err = runtime.BindQueryParameter("form", true, false, "page-size", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page-size", Err: err})
		return
	}

	// ------------- Optional query parameter "order[]" -------------

	err = runtime.BindQueryParameter("form", true, false, "order[]", r.URL.Query(), &params.Order)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "order[]", Err: err})
		return
	}

	// ------------- Optional query parameter "filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter", r.URL.Query(), &params.Filter)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListItems(w, r, backend, entity, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateItem operation middleware
func (siw *ServerInterfaceWrapper) CreateItem(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateItem(w, r, backend, entity)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// DeleteItemById operation middleware
func (siw *ServerInterfaceWrapper) DeleteItemById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteItemById(w, r, backend, entity, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetItemById operation middleware
func (siw *ServerInterfaceWrapper) GetItemById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetItemById(w, r, backend, entity, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ExistsItemById operation middleware
func (siw *ServerInterfaceWrapper) ExistsItemById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ExistsItemById(w, r, backend, entity, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// UpdateItemById operation middleware
func (siw *ServerInterfaceWrapper) UpdateItemById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "backend" -------------
	var backend Backend

	err = runtime.BindStyledParameterWithOptions("simple", "backend", mux.Vars(r)["backend"], &backend, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "backend", Err: err})
		return
	}

	// ------------- Path parameter "entity" -------------
	var entity Entity

	err = runtime.BindStyledParameterWithOptions("simple", "entity", mux.Vars(r)["entity"], &entity, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "entity", Err: err})
		return
	}

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", mux.Vars(r)["id"], &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateItemById(w, r, backend, entity, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/backends", wrapper.ListBackends).Methods("GET")

	r.HandleFunc(options.BaseURL+"/{backend}/entities", wrapper.ListEntities).Methods("GET")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}", wrapper.ListItems).Methods("GET")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}", wrapper.CreateItem).Methods("POST")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}/{id}", wrapper.DeleteItemById).Methods("DELETE")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}/{id}", wrapper.GetItemById).Methods("GET")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}/{id}", wrapper.ExistsItemById).Methods("HEAD")

	r.HandleFunc(options.BaseURL+"/{backend}/{entity}/{id}", wrapper.UpdateItemById).Methods("PUT")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
