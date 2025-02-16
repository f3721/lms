// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

var logger Logger = &DefaultLogger{}

// SetLogger sets the current logger.
func SetLogger(l Logger) {
	logger = l
}

// GetLogger returns the current logger.
func GetLogger() Logger {
	return logger
}

// LogModel logs the model information.
func LogModel(model [][]string) {
	logger.LogModel(model)
}

// LogEnforce logs the enforcer information.
func LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	logger.LogEnforce(matcher, request, result, explains)
}

// LogRole log info related to role.
func LogRole(roles []string) {
	logger.LogRole(roles)
}

// LogPolicy logs the policy information.
func LogPolicy(policy map[string][][]string) {
	logger.LogPolicy(policy)
}
