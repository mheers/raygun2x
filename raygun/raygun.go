package raygun

import (
	"fmt"

	"raygun/config"
	"raygun/finder"
	"raygun/log"
	"raygun/parser"
	"raygun/types"
)

func GetTestSuites(raygunFolder string) ([]types.TestSuite, error) {
	var entities = make([]string, 0)

	entities = append(entities, raygunFolder)

	/*
	 *  Find the raygun files amidst the files and directories specified on the command line
	 */
	finder := finder.NewFinder(".raygun")

	suite_files, err := finder.FindTargets(entities)

	if err != nil {
		return nil, fmt.Errorf("error finding test suites: %v", err)
	}

	if len(suite_files) == 0 {
		return nil, fmt.Errorf("no .raygun files found in specified location(s)")
	}

	/*
	 *  Parse the .raygun files that we found in the previous step
	 */

	log.Verbose("Parsing Raygun files: %v", suite_files)

	parser := parser.NewRaygunParser(config.SkipOnParseError)

	test_suite_list, err := parser.Parse(suite_files)
	if err != nil {
		return nil, fmt.Errorf("unable to parse test files: %v", err)
	}

	return test_suite_list, nil
}
