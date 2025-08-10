package operations

type ClientWrapper interface {
	Get(path string, responseDepth *int, params map[string]string) (interface{}, error)
	Post(path string, data interface{}) (interface{}, error)
	Put(path string, value interface{}) error
	Patch(path string, value interface{}) error
	Delete(path string) (interface{}, error)
	Export(path, filepath string, params map[string]interface{}) error
	Import(path, filename string, params map[string]interface{}) (interface{}, error)

	EnableProfiling(enabled bool)
	PrintVersions()
	PrintProfilingData()
}