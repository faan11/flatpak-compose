# Flatpak Compose

<img align="right"  src="/images/logo.png" alt="image" width="30%" style="float:left;" height="auto">
Flatpak Compose is a tool for managing Flatpak configurations using YAML files. It allows you to define and apply changes to Flatpak repositories and applications easily.

This project works in most of the cases but it has not been thoroughly tested. Therefore, issues may arise. Your feedback and bug reports are highly appreciated as they help us improve the project.

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

### Getting started
You can start from scratch or export the system state and generate the system flatpak-compose.yaml
```bash
flatpak-compose export-state system > flatpak-compose.yaml
```
Now you can change your flatpak-compose.yaml as you deside.
After changing the file, you can proceed by seeing which command will be applied ( plan ) or directly apply the changes ( using the apply command ).

### Examples
```yaml
# flatpak repositories
envs:
- core:
    min-free-space-size: 500MB
    mode: bare-user-only
    repo_version: "1"
    xa.pinned: runtime/org.freedesktop.Sdk/x86_64/23.08;runtime/org.gtk.Gtk3theme.Matcha-dark-sea/x86_64/3.22
  remotes:
    flathub:
      GPGKey: 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
      gpg-verify: "true"
      gpg-verify-summary: "true"
      url: https://dl.flathub.org/repo/
      xa.comment: Central repository of Flatpak applications
      xa.description: Central repository of Flatpak applications
      xa.homepage: https://flathub.org/
      xa.icon: https://dl.flathub.org/repo/logo.svg
      xa.title: Flathub


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

- **Current State**: Can be the **System State** or the **System Compose** state.
- **Desired State**: The desired state (or **Compose State**) is the yaml file. Contains applications/repos specified in the YAML file. 

Flatpak compose brings the system from the current state to the desired state.

- **System State**: Includes all applications/repos installed in the system.
- **System-Compose State**: Applications/repos common between the compose and system states (right join).

#### Help
Show usage information.
```bash
flatpak-compose help
```

## File Structure

- `internal/model/`: Contains the state definition
- `internal/state/`: Contains logic for getting the current and next states, as well as diffing them.
- `internal/utility/`: Contains functions used by the state module to read permissions and the environment from the system. 
- `internal/view/`: Handles generating commands and executing them.

## How It Works

The application reads a YAML file describing Flatpak configurations and applies the specified changes to the system.

## Assets

The logo image is taken by Flaticon.com.

## Contributing

Contributions are welcome! Feel free to open issues or pull requests for enhancements, bug fixes, or new features.

## License
This project is licensed under the [MIT] - see the [LICENSE](LICENSE) file for details.


