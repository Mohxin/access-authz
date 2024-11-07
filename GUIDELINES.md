# Contributing Guidelines for Connect Access Control

Thank you for your interest in contributing to the Connect Access Control repository. This document outlines the process for contributing changes and the standards we expect contributors to follow.

## Table of Contents

1. [Before You Start](#before-you-start)
2. [Making Changes](#making-changes)
3. [Submitting Changes](#submitting-changes)
4. [Review Process](#review-process)
5. [Coding Standards](#coding-standards)
6. [Documentation](#documentation)
7. [Testing](#testing)
8. [Questions and Support](#questions-and-support)

## Before You Start

- Ensure you have read and understood the README.md file in the repository root.
- Familiarize yourself with the current structure and naming conventions used in the repository.
- Check existing issues and pull requests to see if someone else has already addressed the change you want to make.

## Making Changes

1. Create a new branch from the `main` branch for your changes.
2. Make your changes in your branch, following the repository structure and coding standards.
3. Commit your changes with clear, descriptive commit messages.
4. Push your changes to your branch on GitHub.

## Submitting Changes

1. Create a pull request from your branch to the main repository's `main` branch.
2. In your pull request description:
   - Clearly describe the purpose of your changes
   - Link to any relevant issues
   - Provide context on how your changes affect the access control system

## Review Process

1. All changes must be reviewed by the Connect Authorization Governance Group.
2. Reviewers may ask for clarifications or request changes.
3. Address any feedback or comments from reviewers.
4. Once approved, your changes will be merged by a repository maintainer.

## Coding Standards

- Use YAML for configuration files.
- Follow consistent indentation (2 spaces) in YAML files.
- Use clear, descriptive names for roles, permission groups, and scopes.
- Ensure all configuration files pass schema validation.

## Documentation

- Update the [README.md](README.md) file if your changes affect the overall structure or usage of the repository.
- Provide clear comments in configuration files where necessary.
- If adding a new scope or significant feature, consider adding documentation to explain its purpose and usage.

## Testing

- Test your changes thoroughly in a non-production environment before submitting.
- Ensure your changes don't unintentionally affect existing roles or permissions.
- If possible, provide test cases or scenarios that demonstrate the correctness of your changes.

## Questions and Support

If you have any questions about contributing or need support, please:

1. Check the existing documentation in the repository.
2. Review past issues and pull requests for similar topics.
3. If you still need help, contact the Connect Shared Platform Services team through [Slack](https://volvocars.enterprise.slack.com/archives/C062BJJ79MH).

Thank you for contributing to the Connect Access Control system. Your efforts help make our platform more secure and efficient for all users.