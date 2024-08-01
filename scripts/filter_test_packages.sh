#!/bin/bash

# List of packages to exclude (adjust as needed)
exclude_packages=(
  "github.com/beka-birhanu/finance-go/api"
  "github.com/beka-birhanu/finance-go/docs"
  "github.com/beka-birhanu/finance-go/scripts"
  "github.com/beka-birhanu/finance-go/api/user/dto"
  "github.com/beka-birhanu/finance-go/api/expense/dto"
  "github.com/beka-birhanu/finance-go/application/authentication/common"
  "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
  "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
  "github.com/beka-birhanu/finance-go/application/common/interface/hash"
  "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
  "github.com/beka-birhanu/finance-go/application/common/interface/repository"
  "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
  "github.com/beka-birhanu/finance-go/infrastructure/db"
  "github.com/beka-birhanu/finance-go/infrastructure/repository/user"
)

# Find all packages with .go files
all_packages=$(go list ./api/... ./application/... ./infrastructure/...)

# Filter out the excluded packages
for exclude in "${exclude_packages[@]}"; do
  all_packages=$(echo "$all_packages" | grep -v "^${exclude}$")
done

# Convert to a space-separated list
test_packages_list=$(echo "$all_packages" | tr '\n' ' ')

# Print the list of packages
echo $test_packages_list

