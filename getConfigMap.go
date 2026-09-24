package oconf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// loadContext tracks the chain of files currently being loaded to prevent circular includes.
// Only files on the active include chain are tracked, so the same file may still be
// included from several independent places (e.g. two sibling nested objects).
type loadContext struct {
	loadingFiles map[string]bool
}

// newLoadContext creates a new context for tracking loaded files
func newLoadContext() *loadContext {
	return &loadContext{
		loadingFiles: make(map[string]bool),
	}
}

// enterFile marks a file as being loaded on the current include chain
func (ctx *loadContext) enterFile(filePath string) {
	ctx.loadingFiles[filePath] = true
}

// leaveFile removes a file from the current include chain once it has been loaded
func (ctx *loadContext) leaveFile(filePath string) {
	delete(ctx.loadingFiles, filePath)
}

// isFileLoading checks if a file is already on the current include chain
func (ctx *loadContext) isFileLoading(filePath string) bool {
	return ctx.loadingFiles[filePath]
}

// getFromIncludes processes the #include directive and loads referenced files
// It returns a merged configuration from all included files
func getFromIncludes(ctx *loadContext, dirPrefix string, includedFiles []interface{}) (map[string]interface{}, error) {
	mergedConfig := make(map[string]interface{})

	for i, includeEntry := range includedFiles {
		// Safely assert the type
		includedFilePath, ok := includeEntry.(string)
		if !ok {
			return nil, &TypeError{
				Key:      fmt.Sprintf("%s[%d]", IncludeKey, i),
				Expected: "string",
				Actual:   fmt.Sprintf("%T", includeEntry),
			}
		}

		// Resolve the file path
		resolvedPath, resolvedDir, err := resolveIncludePath(dirPrefix, includedFilePath)
		if err != nil {
			return nil, err
		}

		// Check for circular includes
		if ctx.isFileLoading(resolvedPath) {
			continue // Skip files already on the include chain
		}

		// Load and parse the included file
		includeConfig, err := loadConfigFile(ctx, resolvedDir, resolvedPath)
		if err != nil {
			return nil, &FileError{
				Path: resolvedPath,
				Op:   "include",
				Err:  err,
			}
		}

		// Merge the included configuration
		// Fixed: Changed len(mergedConfig) > 1 to > 0
		if len(mergedConfig) > 0 {
			mergedConfig = getMergedMap(mergedConfig, includeConfig)
		} else {
			mergedConfig = includeConfig
		}
	}

	return mergedConfig, nil
}

// resolveIncludePath resolves an include path relative to the current directory
// Returns the absolute file path and the directory containing it
func resolveIncludePath(currentDir, includePath string) (absPath string, dir string, err error) {
	if filepath.IsAbs(includePath) {
		// Absolute path - use as-is
		return includePath, filepath.Dir(includePath), nil
	}

	// Relative path - resolve relative to current directory
	absPath = filepath.Join(currentDir, includePath)
	absPath = filepath.Clean(absPath)
	dir = filepath.Dir(absPath)

	return absPath, dir, nil
}

// getFileDataWithoutComments reads a file and removes all comments
func getFileDataWithoutComments(configFilePath string) ([]byte, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, &FileError{
			Path: configFilePath,
			Op:   "read",
			Err:  err,
		}
	}

	noCommentsData := removeComments(data)
	return noCommentsData, nil
}

// processIncludes resolves #include directives at the current map level and recursively
// in all nested objects. Nested objects are resolved first, so their own includes count
// as part of the current map. Included files' keys are then merged in underneath; the
// current map's keys override them. Returns the updated map with all includes resolved.
func processIncludes(ctx *loadContext, dirPrefix string, m map[string]interface{}) (map[string]interface{}, error) {
	for key, value := range m {
		nestedMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		processed, err := processIncludes(ctx, dirPrefix, nestedMap)
		if err != nil {
			return nil, fmt.Errorf("at key '%s': %w", key, err)
		}
		m[key] = processed
	}

	if includes, ok := m[IncludeKey]; ok {
		includedFiles, ok := includes.([]interface{})
		if !ok {
			return nil, &TypeError{
				Key:      IncludeKey,
				Expected: "array",
				Actual:   fmt.Sprintf("%T", includes),
			}
		}

		includesConfig, err := getFromIncludes(ctx, dirPrefix, includedFiles)
		if err != nil {
			return nil, err
		}

		m = getMergedMap(includesConfig, m)
		delete(m, IncludeKey)
	}

	return m, nil
}

// getConfigMap parses JSON data and processes any #include directives
// It returns a merged configuration map
func getConfigMap(ctx *loadContext, dirPrefix string, data []byte) (map[string]interface{}, error) {
	var parsedData interface{}

	if err := json.Unmarshal(data, &parsedData); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	configMap, ok := parsedData.(map[string]interface{})
	if !ok {
		return nil, &TypeError{
			Expected: "object (map)",
			Actual:   fmt.Sprintf("%T", parsedData),
		}
	}

	return processIncludes(ctx, dirPrefix, configMap)
}

// loadConfigFile loads and parses a configuration file
func loadConfigFile(ctx *loadContext, dirPrefix, filePath string) (map[string]interface{}, error) {
	// Track this file on the include chain to prevent circular includes
	ctx.enterFile(filePath)
	defer ctx.leaveFile(filePath)

	// Read and strip comments
	data, err := getFileDataWithoutComments(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the configuration
	return getConfigMap(ctx, dirPrefix, data)
}
