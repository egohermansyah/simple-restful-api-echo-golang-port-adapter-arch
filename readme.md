# GoLang API use Hexagonal Architecture

## Logger
1. Log all error, except in controller.

## Migration
1. Create manually migration code in repository.

## Seeder
1. Seeder data used for init data staging and production which has don't have feature insert and delete, if need seeder in test, don't add in seeder, create manually in your test code using faker.
2. Don't add credential or sensitive value in seeder. If not null field, add with empty value example : "".

## Linter
1. Run linter <code>go vet ./...</code>

## Test & Coverage
1. Run all test <code>go test ./controller/... ./service/... ./repository/... -v -coverpkg=./controller/... ./service/... ./repository/... -coverprofile=coverage.out</code>
2. See coverage detail in html <code>go tool cover -html=coverage.out</code>
3. Run single test <code>go test -failfast dir/file_name.go</code>
4. Run function test <code>go test -failfast dir/package -run "^(YourTestFunction)$"</code>
