# Flatpak Compose

<img align="right"  src="/images/logo.png" alt="image" width="30%" style="float:left;" height="auto">
Flatpak Compose is a tool for managing Flatpak configurations using YAML files. It allows you to define and apply changes to Flatpak repositories and applications easily.


## Features

- **Manage Repositories**: Add or remove Flatpak repositories.
- **Manage Applications**: Install, uninstall, and update Flatpak applications.
- **Permission Overrides**: Set or replace permissions (overrides) for installed applications.

## Installation

### Mac

```bash
# Download the latest release for Mac
curl -LO https://github.com/faan11/flatpak-compose/releases/latest/download/flatpak-compose-macos-amd64.zip

# Unzip the downloaded file
unzip flatpak-compose-macos-amd64.zip

# Make the binary executable
chmod +x flatpak-compose-macos-amd64

# Move the binary to a directory in your PATH (optional)
sudo mv flatpak-compose-macos-amd64 /usr/local/bin/flatpak-compose
```

### Windows

1. Open your web browser and go to the [Releases](https://github.com/faan11/flatpak-compose/releases) page of your repository.
2. Download the `flatpak-compose-windows-amd64.zip` file from the latest release.
3. Extract the downloaded ZIP file.
4. You'll find the `flatpak-compose-windows-amd64.exe` binary inside the extracted folder.

### Linux

```bash
# Download the latest release for Linux
curl -LO https://github.com/faan11/flatpak-compose/releases/latest/download/flatpak-compose-linux-amd64.zip

# Unpack the downloaded file
unzip flatpak-compose-linux-amd64.zip

# Make the binary executable
chmod +x flatpak-compose-linux-amd64

# Move the binary to a directory in your PATH (optional)
sudo mv flatpak-compose-linux-amd64 /usr/local/bin/flatpak-compose
```

These commands will download the latest release binary for each platform, extract the contents, make the binary executable, and optionally move it to a directory in your PATH for easier access. Adjust the downloaded file name and paths as needed.


## Build 

1. Clone the repository.
2. Build the application using `go build`.
3. Optionally, set the generated binary in your system PATH.

## Usage

### How to start
You can start from scratch or export the system state and generate the system flatpak-compose.yaml
```bash
flatpak-compose export-state system > flatpak-compose.yaml
```
Now you can change your flatpak-compose.yaml as you deside.
After changing the file, you can proceed by seeing which command will be applied ( plan ) or directly apply the changes ( using the apply command ).

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
  overrides_user: []
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
  overrides_user: []
  type: system 

```

### Commands

#### Apply Changes
Apply changes specified in a YAML file.
```bash
flatpak-compose apply [-f file.yaml] [-current-state=system-compose/system]
```
*Default file:* flatpak-compose.yaml / flatpak-compose.yml

#### Plan Changes (Print Only)
Print the commands without applying changes.
```bash
flatpak-compose plan [-f file.yaml] [-current-state=system-compose/system]
```
*Default file:* flatpak-compose.yaml / flatpak-compose.yml

#### Export System State
Print the current system state in a YAML file.
```bash
flatpak-compose export-state system
flatpak-compose export-state system-compose
```
The export-state will add a new field "all" for each application. This field holds all the permissions (default and static permissions).

#### Help
Show usage information.
```bash
flatpak-compose help
```

## States Explanation

- **Current State**: The currently configured state in the YAML file.
- **System State**: Includes all applications/repos installed in the system.
- **Compose State**: Contains applications/repos specified in the YAML file.
- **System-Compose State**: Applications/repos common between the compose and system states (right join).


#### Help
Show usage information.
```bash
flatpak-compose help
```

## File Structure

- `internal/model/`: Contains logic for getting the current and next states, as well as diffing them.
- `internal/view/`: Handles generating commands and executing them.

## How It Works

The application reads a YAML file describing Flatpak configurations and applies the specified changes to the system.

## Assets

The logo image is taken by Flaticon.com.

## Contributing

Contributions are welcome! Feel free to open issues or pull requests for enhancements, bug fixes, or new features.

## License
This project is licensed under the [MIT] - see the [LICENSE](LICENSE) file for details.


