// Package oconf provides functionality to load configuration files in CJSON format
// (JSON with C-style comments) with support for file inclusion and public/private annotations.
//
// Features:
//   - C-style comments (// and /* */) in JSON files
//   - File inclusion via #include directive
//   - Public/private configuration values via #public annotation
//   - Nested and flat configuration formats
//   - JSON string output
//
// Example configuration file (config.cjson):
//
//	{
//	    // Include base configuration
//	    "#include": ["base.cjson"],
//
//	    // Public configuration visible to external systems
//	    "server#public": {
//	        "host#public": "example.com",
//	        "port#public": 8080
//	    },
//
//	    // Private configuration
//	    "database": {
//	        "password": "secret"
//	    }
//	}
//
// Basic usage:
//
//	config, err := oconf.GetConfig("config.cjson")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	serverPort := config["server"].(map[string]interface{})["port"].(float64)
package oconf

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// getAllConfig is an internal helper that loads a configuration file with all preprocessing
func getAllConfig(configFilePath string) (map[string]interface{}, error) {
	// Create a new context to track loaded files
	ctx := newLoadContext()

	// Mark the entry file as loaded
	absPath, err := filepath.Abs(configFilePath)
	if err != nil {
		return nil, &FileError{
			Path: configFilePath,
			Op:   "resolve path",
			Err:  err,
		}
	}
	ctx.markFileLoaded(absPath)

	// Read and strip comments
	data, err := getFileDataWithoutComments(absPath)
	if err != nil {
		return nil, err
	}

	// Parse configuration with includes
	dirPrefix := filepath.Dir(absPath)
	configMap, err := getConfigMap(ctx, dirPrefix, data)
	if err != nil {
		return nil, err
	}

	return configMap, nil
}

// GetConfig loads a configuration file and returns it as a nested map.
// All #public annotations are removed but the values are retained.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A map containing the configuration with nested structure
//   - An error if the file cannot be read or parsed
//
// Example:
//
//	config, err := oconf.GetConfig("config.cjson")
//	if err != nil {
//	    return err
//	}
//	dbHost := config["database"].(map[string]interface{})["host"].(string)
func GetConfig(configFilePath string) (map[string]interface{}, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	removePublicAnnotations(configMap)
	return configMap, nil
}

// GetFlatConfig loads a configuration file and returns it as a flat map
// where nested keys are joined with dots (e.g., "database.host").
// All #public annotations are removed but the values are retained.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A flat map where nested paths are joined with dots
//   - An error if the file cannot be read or parsed
//
// Example:
//
//	config, err := oconf.GetFlatConfig("config.cjson")
//	if err != nil {
//	    return err
//	}
//	dbHost := config["database.host"].(string)
//	serverPort := config["server.port"].(float64)
func GetFlatConfig(configFilePath string) (map[string]interface{}, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	removePublicAnnotations(configMap)

	flatMap := make(map[string]interface{})
	if err := createFlatMap("", flatMap, configMap); err != nil {
		return nil, err
	}
	return flatMap, nil
}

// GetConfigAsJSONString loads a configuration file and returns it as a JSON string.
// All #public annotations are removed but the values are retained.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A JSON string representation of the configuration
//   - An error if the file cannot be read, parsed, or marshaled
//
// Example:
//
//	jsonStr, err := oconf.GetConfigAsJSONString("config.cjson")
//	if err != nil {
//	    return err
//	}
//	fmt.Println(jsonStr)
func GetConfigAsJSONString(configFilePath string) (string, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return "", err
	}
	removePublicAnnotations(configMap)

	jsonBytes, err := json.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// GetPublicConfig loads a configuration file and returns only values marked with #public.
// The #public annotations are removed from the keys in the result.
//
// This is useful for exposing configuration to external systems or clients
// while keeping sensitive values (database passwords, API keys, etc.) private.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A map containing only the public configuration values
//   - An error if the file cannot be read or parsed
//
// Example:
//
// Given a config file:
//
//	{
//	    "server#public": {
//	        "host#public": "example.com"
//	    },
//	    "database": {
//	        "password": "secret"
//	    }
//	}
//
// Usage:
//
//	publicConfig, err := oconf.GetPublicConfig("config.cjson")
//	// Result: {"server": {"host": "example.com"}}
//	// Note: "database.password" is excluded because it's not marked #public
func GetPublicConfig(configFilePath string) (map[string]interface{}, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	saveOnlyWithPublicAnnotations(configMap)
	removePublicAnnotations(configMap)
	return configMap, nil
}

// GetFlatPublicConfig loads a configuration file and returns only values marked with #public
// as a flat map where nested keys are joined with dots.
// The #public annotations are removed from the keys in the result.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A flat map containing only public configuration values
//   - An error if the file cannot be read or parsed
//
// Example:
//
//	publicConfig, err := oconf.GetFlatPublicConfig("config.cjson")
//	if err != nil {
//	    return err
//	}
//	serverHost := publicConfig["server.host"].(string)
func GetFlatPublicConfig(configFilePath string) (map[string]interface{}, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	saveOnlyWithPublicAnnotations(configMap)
	removePublicAnnotations(configMap)

	flatMap := make(map[string]interface{})
	if err := createFlatMap("", flatMap, configMap); err != nil {
		return nil, err
	}
	return flatMap, nil
}

// GetPublicConfigAsJSONString loads a configuration file and returns only values
// marked with #public as a JSON string. The #public annotations are removed from the keys.
//
// Parameters:
//   - configFilePath: Path to the CJSON configuration file
//
// Returns:
//   - A JSON string containing only public configuration values
//   - An error if the file cannot be read, parsed, or marshaled
//
// Example:
//
//	jsonStr, err := oconf.GetPublicConfigAsJSONString("config.cjson")
//	if err != nil {
//	    return err
//	}
//	// Send jsonStr to external API or client
func GetPublicConfigAsJSONString(configFilePath string) (string, error) {
	configMap, err := getAllConfig(configFilePath)
	if err != nil {
		return "", err
	}
	saveOnlyWithPublicAnnotations(configMap)
	removePublicAnnotations(configMap)

	jsonBytes, err := json.Marshal(configMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public config to JSON: %w", err)
	}
	return string(jsonBytes), nil
}
