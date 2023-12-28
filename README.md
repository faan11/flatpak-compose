# Flatpak Compose

Flatpak Compose is a tool for managing Flatpak configurations using YAML files. It allows you to define and apply changes to Flatpak repositories and applications easily.

## Features

- **Manage Repositories**: Add or remove Flatpak repositories.
- **Manage Applications**: Install, uninstall, and update Flatpak applications.
- **Permission Overrides**: Set or replace permissions (overrides) for installed applications.

## Usage

### Build 

1. Clone the repository.
2. Build the application using `go build`.
3. Optionally, set the generated binary in your system PATH.

### Examples
```yaml
# flatpak repositories
repos:
- name: flathub
  uri: https://dl.flathub.org/repo/
  type: system
- name: flathub-beta
  uri: https://dl.flathub.org/beta-repo/
  type: system

# application list
applications:
# keepass app
- name: org.keepassxc.KeePassXC
  repo: flathub
  branch: stable
  overrides: []
  type: system
# firefox app
- name: org.mozilla.firefox
  repo: flathub
  branch: stable
  overrides:
  - --nofilesystem=host
  - --nosocket=x11
  - --socket=fallback-x11
  - --allow=bluetooth
  type: system 

```

### Commands

#### Apply Changes
Apply changes specified in a YAML file.
```bash
flatpak-compose apply [-f file.yaml]
```

#### Plan Changes (Print Only)
Print the commands without applying changes.
```bash
flatpak-compose plan [-f file.yaml]
```

#### Export the system current state
Print the system current state in a YAML file.
```bash
flatpak-compose export
```

#### Help
Show usage information.
```bash
flatpak-compose help
```

## File Structure

- `internal/model/`: Contains logic for getting the current and next states, as well as diffing them.
- `internal/view/`: Handles generating commands and executing them.
- `flatpak-compose.yaml`: Default YAML file for configuration.

## How It Works

The application reads a YAML file describing Flatpak configurations and applies the specified changes to the system.

## Contributing

Contributions are welcome! Feel free to open issues or pull requests for enhancements, bug fixes, or new features.

## License

This project is licensed under the [MIT] - see the [LICENSE.md](LICENSE.md) file for details.


