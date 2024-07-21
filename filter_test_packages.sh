#!/bin/bash

# List of packages to exclude (adjust as needed)
exclude_packages=(
  "github.com/beka-birhanu/finance-go/api/users/dto"
  "github.com/beka-birhanu/finance-go/application/authentication/common"
  "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
  "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
  "github.com/beka-birhanu/finance-go/application/common/interfaces/hash"
  "github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
  "github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
  "github.com/beka-birhanu/finance-go/infrastructure/db"
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

