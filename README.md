# Connect Access Control

This repository manages the access control configuration for the Connect platform. It defines roles, permission groups, and their mappings to provide a flexible and secure authorization model across all Connect applications.

## Overview

The Connect access control system is based on the following key concepts:

1. Roles: High-level job functions defined in PLUMS
2. Scopes: Bounded contexts or applications within Connect
3. Permission Groups: Collections of related permissions within a scope
4. Permission Group Filters: Contextual filters for applying permission groups
5. Role Mappings: Associations between roles and permission groups

## Key Files and Directories

- `iam/config/roles.yaml`: Defines the global roles used across Connect
- `iam/config/schema/`: Contains JSON schema definitions for various configuration files
- `iam/scopes/`: Holds scope-specific configurations, including permission groups and role mappings
- `iam/client/`: Defines client-specific configurations and dependent scopes 
- `iam/client/users/`: Contains user-specific permission mappings (for development/testing purposes)

## Usage

1. To add a new scope:

   - Create a new directory under `iam/scopes/`
   - Add `scope.yaml`, `permission-groups.yaml`, and role mapping files

2. To modify permission groups:

   - Edit the relevant `permission-groups.yaml` file in the scope directory, **only removal and re-addition of permission groups is allowed**

3. To update role mappings:

   - Modify or add role mapping files in the scope's `role-mapping/` directory

4. To add or modify global roles:
   - Update `iam/config/roles.yaml`

## Governance

All changes to this repository must go through a review process:

1. Create a pull request with your proposed changes
2. The changes will be reviewed by the [Connect Authorization Governance Group](https://www.notion.so/volvocars/Connect-Authorization-Governance-Group-a0c294fe25394213889eea20a9ed6aac)
3. If approved, the changes will be merged and deployed

For more details on the governance process, please refer to the [Governance Documentation](https://www.notion.so/volvocars/ADR-415-Connect-Permission-Model-2addd29ad8254958aaf4fe3669c69a4c?pvs=4#f96ccd25d7d4465ba730fe70947b8d29).

## Development vs. Production

- In the production environment, users are assigned roles, and the UI only displays and manages roles.
- In the development environment, direct assignment of permission groups to users is allowed for testing and integration purposes.

## Contributing

Please read our [Contributing Guidelines](GUIDELINES.md) before submitting any changes.

## Support

For questions or issues, please contact the Connect Platform Framework team.

## TODOs
- [x] Access control service :: Use User CDSID instead of Plums ID.
- [ ] Access control service :: Validate data integrity.
- [ ] Access control service :: Re-structure/Cleanup.
- [ ] Schema validator       :: Validate client uniqueness
- [ ] Schema validator       :: Validate role uniqueness
- [ ] Schema validator       :: Validate permission groups in role mapping
- [ ] Schema validator       :: Validate users
- [ ] Documentation          :: Document how we define the roles mapping with other scopes
- [ ] HTTP Server            :: Deal user permissions with multiple scopes
