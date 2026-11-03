# 🌿 Leafy

**Leafy** is your cozy terminal companion for transferring media from external drives. It turns the chore of manual file transfers into a breezy, interactive experience right in your shell.

## 💡 Motivation

I started vlogging my motorcycle journey recently, and found files transferring from my action camera and external mic is a *very* tedious task. So **Leafy** is created to automate the process. Now I can quickly and easily transfer all the captured media to my editing PC by simply running **Leafy** in the terminal! 🙌

## 🚀 Quickstart

> [!NOTE]
> **Leafy** currently only supports **Linux**.

A binary will be released on GitHub later. Alternatively, you can install it using Go.

```bash
# Option 1: Install with Go
go install github.com/nightails/leafy@latest

# Option 2: Download the binary
# 1. Download: Download the latest binary from the Releases page.
# 2. Permissions: Ensure the binary is executable:
chmod +x leafy
# 3. Run: Execute the binary:
./leafy
```

## 📖 Usage

**Leafy** is designed for a simple, hands-off workflow:

1.  **Automatic Scanning & Mounting**: Upon launch, the app automatically scans for connected USB devices and mounts any unmounted partitions.
2.  **Media Detection**: It searches for media files (audio and video formats) on the mounted devices.
3.  **Media Selection**: Detected media files are displayed in a list. 
    - Use arrow keys or `j`/`k` to navigate.
    - Press `space` to select or deselect files.
    - Press `enter` to confirm your selection and proceed.
4.  **Destination Selection**: Enter the path where you want to copy the files (supports `~` for home directory).
5.  **Transferring**: The app copies the selected files to the destination with a real-time progress bar.
6.  **Exiting**: Once the transfer completes or fails, the app automatically unmounts the devices that were mounted by the app and exits. You can also press `ctrl+c` to exit before copying starts.

## 🤝 Contributing

To set up Leafy for local development, ensure you have the following dependencies installed:

-   **Go** (version 1.25 or later)
-   **lsblk**: The application relies on the `lsblk` command-line utility for device detection.
-   **udisksctl**: The application uses `udisksctl` for mounting and unmounting devices.

### Development Steps

```bash
# 1. Clone the repository:
git clone https://github.com/nightails/leafy.git
cd leafy

# 2. Install dependencies:
go mod download

# 3. Run locally:
go run main.go

# 4. Build:
go build -o leafy main.go
```

## 📜 License

This project is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for more details.
