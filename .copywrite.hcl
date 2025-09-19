# Terraform Provider for Peekaping Monitors
# This provider interfaces with the Peekaping API (https://github.com/0xfurai/peekaping)
# which is licensed under MIT License.
schema_version = 1

project {
  copyright_holder = "tafaust"
  license          = "MIT"
  copyright_year   = 2025

  header_ignore = [
    # internal catalog metadata (prose)
    "META.d/**/*.yaml",

    # examples used within documentation (prose)
    "examples/**",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",

    # Upstream project files (they have their own MIT license)
    "peekaping-upstream/**",
  ]
}
