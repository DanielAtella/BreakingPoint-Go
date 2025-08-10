package models

import (
    "fmt"
    "strconv"
)

// DataModelProxy represents a dynamic proxy for BPS data model objects
type DataModelProxy struct {
    wrapper   BPSWrapper
    name      string
    path      string
    modelPath string
    cache     map[string]interface{}
}

// BPSWrapper interface defines the methods needed from BPS client
type BPSWrapper interface {
    Get(path string, responseDepth *int, params map[string]string) (interface{}, error)
    Post(path string, data interface{}) (interface{}, error)
    Put(path string, value interface{}) error
    Patch(path string, value interface{}) error
    Delete(path string) (interface{}, error)
    Export(path, filepath string, params map[string]interface{}) error
    Import(path, filename string, params map[string]interface{}) (interface{}, error)
}

// NewDataModelProxy creates a new data model proxy
func NewDataModelProxy(wrapper BPSWrapper, name, path string) *DataModelProxy {
    return &DataModelProxy{
        wrapper:   wrapper,
        name:      name,
        path:      path,
        modelPath: path,
        cache:     make(map[string]interface{}),
    }
}

// NewDataModelProxyWithModel creates a new data model proxy with explicit model path
func NewDataModelProxyWithModel(wrapper BPSWrapper, name, path, modelPath string) *DataModelProxy {
    return &DataModelProxy{
        wrapper:   wrapper,
        name:      name,
        path:      path,
        modelPath: modelPath,
        cache:     make(map[string]interface{}),
    }
}

// Get retrieves data from the endpoint
func (p *DataModelProxy) Get(responseDepth *int, params map[string]string) (interface{}, error) {
    fullPath := p.fullPath()
    return p.wrapper.Get(fullPath, responseDepth, params)
}

// Set updates data at the endpoint using PATCH
func (p *DataModelProxy) Set(value interface{}) error {
    fullPath := p.fullPath()
    return p.wrapper.Patch(fullPath, value)
}

// Put updates data at the endpoint using PUT
func (p *DataModelProxy) Put(value interface{}) error {
    fullPath := p.fullPath()
    return p.wrapper.Put(fullPath, value)
}

// Delete removes data at the endpoint
func (p *DataModelProxy) Delete() (interface{}, error) {
    fullPath := p.fullPath()
    return p.wrapper.Delete(fullPath)
}

// GetItem creates a proxy for accessing array items
func (p *DataModelProxy) GetItem(index interface{}) *DataModelProxy {
    var indexStr string
    switch v := index.(type) {
    case int:
        indexStr = strconv.Itoa(v)
    case string:
        indexStr = v
    default:
        indexStr = fmt.Sprintf("%v", v)
    }

    return NewDataModelProxyWithModel(
        p.wrapper,
        indexStr,
        p.fullPath(),
        p.dataModelPath(),
    )
}

// GetField creates a proxy for accessing nested fields
func (p *DataModelProxy) GetField(fieldName string) *DataModelProxy {
    return NewDataModelProxyWithModel(
        p.wrapper,
        fieldName,
        p.fullPath(),
        p.dataModelPath(),
    )
}

// fullPath returns the complete API path
func (p *DataModelProxy) fullPath() string {
    if p.path == "" {
        return "/" + p.name
    }
    return p.path + "/" + p.name
}

// dataModelPath returns the data model path
func (p *DataModelProxy) dataModelPath() string {
    if p.modelPath == "" {
        return p.path
    }
    return p.modelPath + "/" + p.name
}

// URL returns the full API URL
func (p *DataModelProxy) URL() string {
    // This would need the host from the wrapper, simplified for now
    return fmt.Sprintf("/bps/api/v2/core%s", p.fullPath())
}

// String provides a string representation
func (p *DataModelProxy) String() string {
    return fmt.Sprintf("DataModelProxy{name: %s, path: %s}", p.name, p.fullPath())
}

// CachedGet retrieves and caches field values
func (p *DataModelProxy) CachedGet(field string) (interface{}, error) {
    if value, exists := p.cache[field]; exists {
        return value, nil
    }

    result, err := p.wrapper.Get(p.dataModelPath()+"/"+field, nil, nil)
    if err != nil {
        return nil, err
    }

    p.cache[field] = result
    return result, nil
}

