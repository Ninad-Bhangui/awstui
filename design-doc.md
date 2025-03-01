# AWS Terminal UI (TUI) Application Design

## Overview
The AWS TUI application is a terminal-based interface (written in Go) for browsing AWS resources in a style similar to **k9s** (the popular Kubernetes CLI TUI). It aims to improve productivity by allowing quick navigation and inspection of AWS resources without leaving the terminal. The initial version will support browsing **EC2 instances**, **ECR repositories**, **Lambda functions**, and **Secrets Manager secrets**. 

The interface will feature Vim-inspired key bindings (along with arrow key support) for navigation and actions. Users can quickly switch between resource views, inspect resource details in JSON form, and copy those details to their clipboard for use in other tools or scripts.

---

## Goals and Features

### Core Features
- **Multi-Service Resource Browser**: View lists of resources for EC2, ECR, Lambda, and Secrets Manager.
- **K9s-Like UI/UX**:
  - Terminal UI with a main table/list view of resources.
  - **Vim-like navigation**: `j/k` or arrow keys to move up/down, `h/l` or left/right for navigation.
  - **Command hints**: A list of available key bindings is always visible (e.g., in the top-right corner).
  - **Help screen**: Press `?` to display available key bindings and commands.
- **Resource Detail View**:
  - Fetch detailed information (in JSON format) via the AWS SDK.
  - Display the JSON details in a scrollable text window or modal.
  - Offer a key binding to copy JSON output to the system clipboard.
- **Profile Selection and Authentication**:
  - Show available AWS profiles on startup.
  - If `AWS_PROFILE` is set, use that profile by default.
  - Display active profile (and region) in the UI.
- **Clipboard Support**: Copy JSON details to the clipboard.
  - **Mac & Linux Priority**: Use `pbcopy` (macOS) and `xclip/xsel` (Linux).
  - Windows support is optional.
- **Snappy Performance**:
  - Use Go concurrency where required to improve UI responsiveness.
  - Load data efficiently without unnecessary background processing.
- **Fixed Set of AWS Services** (Customization Not a Priority):
  - EC2, ECR, Lambda, and Secrets Manager.
  - No configuration options for adding/removing services initially.

---

## User Interface Design

### Layout
- **Header Bar**: Displays the active AWS profile, region, and selected resource view.
- **Main Table View**: List of resources, with sortable columns.
- **Footer/Status Bar**: Displays action messages, errors, or prompts.
- **Key Binding Hint Panel**: Located in the top-right, listing common key bindings.
- **Detailed View Modal**:
  - Opens when describing a resource (`d` key).
  - Displays JSON details in a scrollable format.
  - Allows copying JSON to clipboard (`c` key).
  - Close with `Esc` or `q`.
- **Profile Selection Menu** (Shown on Startup if Needed):
  - Lists available AWS profiles for selection.

### Navigation and Key Bindings
#### **Global Navigation**
- `j/k` (or ↓/↑): Move up/down in resource list.
- `h/l` (or ←/→): Potentially switch between resource views.
- `g/G`: Jump to top/bottom.
- `/`: Start typing to filter (future enhancement).
- `?`: Show help screen.
- `q`: Quit.

#### **Resource Actions**
- `Enter/o`: Select/Open (contextual, might just trigger describe).
- `d`: Describe resource (show JSON details).
- `c`: Copy JSON details to clipboard.
- `r`: Refresh resource list.
- `:` (Command Mode, future feature).

#### **Switching Resource Views**
- `1`: EC2 Instances
- `2`: ECR Repositories
- `3`: Lambda Functions
- `4`: Secrets Manager

---

## AWS Integration and Backend Design

### AWS SDK
- Use **AWS SDK for Go v2**.
- Authentication via AWS credentials file or `AWS_PROFILE` env var.

### Listing Resources (APIs Used)
- **EC2**: `DescribeInstances`
- **ECR**: `DescribeRepositories`
- **Lambda**: `ListFunctions`
- **Secrets Manager**: `ListSecrets`

### Fetching Resource Details
- **EC2**: `DescribeInstances`
- **ECR**: `DescribeRepositories`
- **Lambda**: `GetFunction`
- **Secrets Manager**: `DescribeSecret`

---

## Technology Choices

### TUI Library: `tview` vs `bubbletea`
- **Decision: Use `tview`** for its prebuilt widgets and ease of use.
- `bubbletea` could be an option for future versions if more customization is needed.

### Clipboard Support Implementation
- **MacOS**: `pbcopy`
- **Linux**: `xclip` or `xsel`
- **Windows (Optional)**: Windows clipboard API (via `atotto/clipboard`).

### Application Structure
```
aws-tui/
├── main.go           # Entry point
├── ui/               # UI components (layout, navigation, event handlers)
├── aws/              # AWS SDK logic (list, describe, etc.)
├── clipboard/        # Clipboard utility
└── config/           # AWS profile/region selection
```

---

## Conclusion
This design outlines a lightweight AWS resource browser for the terminal. It emphasizes **speed, usability, and a k9s-like experience** while focusing on four key AWS services. The initial implementation will be read-only with essential navigation and clipboard functionality, leaving room for future enhancements.
