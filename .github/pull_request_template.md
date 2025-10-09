name: 🚀 Pull Request
description: Submit a pull request to improve the Terraform Provider for Peekaping
title: "[PR] "
labels: ["needs-review"]
body:
  - type: markdown
    attributes:
      value: |
        Thanks for contributing to the Terraform Provider for Peekaping! 🚀

        Please fill out the sections below to help us understand and review your changes.

  - type: dropdown
    id: pr-type
    attributes:
      label: 🎯 Pull Request Type
      description: What type of change is this?
      options:
        - Bug fix
        - New feature
        - Documentation update
        - Performance improvement
        - Refactoring
        - Test improvement
        - Dependency update
        - Other
    validations:
      required: true

  - type: dropdown
    id: breaking-change
    attributes:
      label: ⚠️ Breaking Change
      description: Does this PR introduce a breaking change?
      options:
        - No breaking change
        - Yes, breaking change (explain below)
    validations:
      required: true

  - type: textarea
    id: description
    attributes:
      label: 📝 Description
      description: A clear and concise description of what this PR does
      placeholder: "This PR adds/fixes/improves..."
    validations:
      required: true

  - type: textarea
    id: motivation
    attributes:
      label: 🎯 Motivation
      description: Why is this change needed? What problem does it solve?
      placeholder: |
        This change is needed because:
        - Current behavior has issue X
        - Users need feature Y
        - Performance can be improved by Z
    validations:
      required: true

  - type: textarea
    id: changes
    attributes:
      label: 🔄 Changes Made
      description: List the specific changes made in this PR
      placeholder: |
        Changes made:
        - [ ] Added new resource type X
        - [ ] Fixed bug in Y
        - [ ] Updated documentation for Z
        - [ ] Added tests for A
        - [ ] Improved error handling for B
    validations:
      required: true

  - type: textarea
    id: testing
    attributes:
      label: 🧪 Testing
      description: How has this been tested?
      placeholder: |
        Testing performed:
        - [ ] Unit tests added/updated
        - [ ] Integration tests added/updated
        - [ ] Manual testing performed
        - [ ] Tested with different Terraform versions
        - [ ] Tested with different Peekaping versions
    validations:
      required: true

  - type: textarea
    id: test-results
    attributes:
      label: ✅ Test Results
      description: Provide test results and any relevant output
      placeholder: |
        Test results:
        ```bash
        $ go test ./...
        PASS: TestResourceMonitor (0.05s)
        PASS: TestResourceNotification (0.03s)
        ...
        ```
    validations:
      required: true

  - type: textarea
    id: documentation
    attributes:
      label: 📚 Documentation
      description: What documentation has been updated?
      placeholder: |
        Documentation updates:
        - [ ] README.md updated
        - [ ] Resource documentation updated
        - [ ] Examples updated
        - [ ] API documentation updated
        - [ ] No documentation changes needed
    validations:
      required: true

  - type: textarea
    id: breaking-changes-detail
    attributes:
      label: ⚠️ Breaking Changes Detail
      description: If this is a breaking change, explain what breaks and how to migrate
      placeholder: |
        Breaking changes:
        - Resource attribute X has been renamed to Y
        - Provider configuration Z is now required
        - Migration guide: [link or instructions]
    validations:
      required: false

  - type: checkboxes
    id: checklist
    attributes:
      label: ✅ Checklist
      description: Please check all that apply
      options:
        - label: I have read the Contributing Guidelines
          required: true
        - label: My code follows the project's style guidelines
          required: true
        - label: I have performed a self-review of my code
          required: true
        - label: I have commented my code, particularly in hard-to-understand areas
          required: false
        - label: I have made corresponding changes to the documentation
          required: false
        - label: My changes generate no new warnings
          required: true
        - label: I have added tests that prove my fix is effective or that my feature works
          required: true
        - label: New and existing unit tests pass locally with my changes
          required: true
        - label: Any dependent changes have been merged and published
          required: false

  - type: textarea
    id: related-issues
    attributes:
      label: 🔗 Related Issues
      description: Link any related issues (use "Fixes #123", "Closes #456", etc.)
      placeholder: |
        Related issues:
        - Fixes #123
        - Closes #456
        - Related to #789

  - type: textarea
    id: screenshots
    attributes:
      label: 📸 Screenshots
      description: If applicable, add screenshots to help explain your changes
      placeholder: "Screenshots or GIFs showing the changes in action..."

  - type: textarea
    id: additional-notes
    attributes:
      label: 📋 Additional Notes
      description: Any additional information that reviewers should know
      placeholder: |
        Additional notes:
        - This PR depends on upstream changes in Peekaping
        - Performance impact: minimal
        - Security considerations: none
        - Other important information...

  - type: markdown
    attributes:
      value: |
        ## 📋 Review Guidelines

        ### For Reviewers
        - ✅ Check that the PR addresses the issue/feature request
        - ✅ Verify that tests are comprehensive and pass
        - ✅ Ensure code follows project conventions
        - ✅ Confirm documentation is updated appropriately
        - ✅ Test the changes locally if significant

        ### For Contributors
        - ✅ Respond to review feedback promptly
        - ✅ Make requested changes clearly
        - ✅ Update tests if functionality changes
        - ✅ Keep the PR focused and atomic

        **Thank you for contributing!** 🙏
