module github.com/mheers/raygun2x

go 1.24.1

require github.com/rbretecher/go-postman-collection v0.9.0

require (
	github.com/aws/smithy-go v1.22.3
	github.com/getkin/kin-openapi v0.129.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.6.1
	gopkg.in/yaml.v3 v3.0.1
	raygun v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oasdiff/yaml v0.0.0-20241210131133-6b86fb107d80 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20241210130736-a94c01f36349 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace raygun => github.com/paclabsnet/raygun v0.1.4
