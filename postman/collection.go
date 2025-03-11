package postman

import (
	"bytes"
	"fmt"
	"os"

	"raygun/types"

	postman "github.com/rbretecher/go-postman-collection"
)

func CreateCollection(suites []types.TestSuite) (*postman.Collection, error) {
	c := postman.CreateCollection("Raygun Import", "")

	for _, suite := range suites {
		group := c.AddItemGroup(suite.Name)
		group.Description = suite.Description

		for _, test := range suite.Tests {
			testEvents := []*postman.Event{}

			for _, expect := range test.Expects {
				if expect.Target == "" {
					continue
				}
				if expect.ExpectationType != "substring" {
					continue
				}
				testEvents = append(testEvents, &postman.Event{
					Listen: postman.Test,
					Script: &postman.Script{
						Exec: []string{
							`pm.test("response has substring", function () {`,
							fmt.Sprintf("    pm.expect(pm.response.text()).to.include(`%s`);", expect.Target),
							`});`,
						},
						Type: postman.Javascript,
					},
				})
			}

			testEvents = append(testEvents, &postman.Event{
				Listen: postman.Test,
				Script: &postman.Script{
					Exec: []string{
						`pm.test("Status code is 200", function () {`,
						`    pm.response.to.have.status(200);`,
						`});`,
					},
					Type: postman.Javascript,
				},
			})

			item := &postman.Items{
				Name: test.Name,
				Request: &postman.Request{
					Description: test.Description,
					URL: &postman.URL{
						Raw:  fmt.Sprintf("{{baseUrl}}%s", test.DecisionPath),
						Host: []string{"{{baseUrl}}"},
						Path: []string{test.DecisionPath},
					},
					Body: &postman.Body{
						Raw:  test.Input.Value,
						Mode: "raw",
						Options: &postman.BodyOptions{
							Raw: postman.BodyOptionsRaw{
								Language: postman.JSON,
							},
						},
					},
					Method: postman.Post,
				},
				Events: testEvents,
			}
			group.AddItem(item)
		}
	}

	return c, nil
}

func WriteToFile(c *postman.Collection, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	defer file.Close()

	if err := c.Write(file, postman.V210); err != nil {
		return fmt.Errorf("error writing collection to file: %v", err)
	}

	return nil
}

func MarshalCollection(c *postman.Collection) ([]byte, error) {
	buffer := &bytes.Buffer{}
	if err := c.Write(buffer, postman.V210); err != nil {
		return nil, fmt.Errorf("error writing collection to buffer: %v", err)
	}

	return buffer.Bytes(), nil
}
