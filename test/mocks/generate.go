package mocks

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name RefTypeRepository --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name RefTypeBroker --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name PropertyRepository --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name PropertyBroker --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name RecordRepository --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name RecordBroker --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name ValueRepository --output "."

//go:generate go run github.com/vektra/mockery/v2@latest --dir ../../internal/domain --name ValueBroker --output "."
