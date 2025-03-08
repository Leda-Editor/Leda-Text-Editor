# Leda - A Text Editor Built with Fyne  

Leda is a modern and lightweight text editor built using the [Fyne](https://fyne.io/) framework. It offers a simple, intuitive interface for editing text files while providing powerful functionality under the hood.  

---

## Build Showcase

![LedaEditor](https://github.com/Leda-Editor/Leda-Text-Editor/blob/main/assets/team3_LedaEditor.png)

## Features  

- Basic text editing functions (cut, copy, paste, undo, redo)  
- Cross-platform support (Linux, macOS, Windows)  
- Clean and minimalistic design  
- Autosave functionality  
- Dark/Light mode  
- Markdown parsing  
- Open, edit, and save files  
- Custom UI presets/layouts
- JSON configuration

---

## Project Purpose And Architecture Overview  

### Purpose

Leda was developed to fill a gap in the text editor landscape, offering a streamlined and efficient
experience for users seeking a balance between simplicity and functionality.

We recognized that there is a stark contrast between the complexity of a text editor like Vim and
the feature-rich bloat of IDE’s such as VSCode, so we aimed to create a lightweight yet powerful
alternative. Leda provides essential features such as customization, terminal integration, and
syntax highlighting, helping users focus on their content - be it code, configuration files, or plain
text - with no unnecessary distractions.

Our focused approach makes Leda ideal for people who prioritize speed and simplicity within
their development workflow.

---

Leda is built using **Go** and the **Fyne** UI framework, designed to be a cross-platform and lightweight text editor. The architecture is modular, ensuring flexibility and maintainability.  

### **Core Technologies**  

- **Go (Golang):** A fast and efficient language, chosen for its performance and simplicity in building cross-platform applications.  
- **Fyne:** A native GUI framework for Go, used to create a consistent UI experience across Linux, macOS, and Windows.  

### **How We Use Fyne in Leda**  

- **UI Rendering:**  
  - We use `fyne.Container` to organize the layout of the text editor, ensuring a responsive and clean interface.  
  - The theme engine in Fyne allows Leda to support **dark and light modes**, improving user experience.  

- **Text Editing with Fyne Widgets:**  
  - The editor is built around `widget.Entry`, a Fyne text input field that provides native text editing functionality.  
  - We extend `widget.Entry` to handle features like **syntax highlighting** and **Markdown preview**.  

- **File Management and Autosave:**  
  - We use `dialog.ShowFileOpen` and `dialog.ShowFileSave` from Fyne’s built-in dialog system to handle file opening and saving.  
  - Autosave functionality is implemented using Go’s **goroutines**, ensuring background saving without blocking the UI.  

- **Cross-Platform Compatibility:**  
  - Fyne abstracts OS-level differences, allowing Leda to run seamlessly on different operating systems with a single codebase.  

This architecture allows Leda to be **lightweight, responsive, and easy to extend**, making it a great choice for users who need a simple yet powerful text editor.  

---

## Features  

### **1. Markdown Preview with Real-Time Rendering**  
Leda includes a **real-time Markdown preview pane** that updates instantly as you type. This allows users to see formatted Markdown output without needing to manually refresh or switch tabs.  

- Uses **Markdown parsing libraries** to ensure fast and accurate rendering.  
- Supports **common Markdown syntax**, including headers, lists, links, and code blocks.  
- Live preview is displayed on the **right side** of the editor in split-screen mode.  

### **2. Split-Screen Layout (Editor + Markdown Preview)**  
The interface is divided into two panes:  

- **Left pane:** The main text editor where users type and edit text.  
- **Right pane:** The live Markdown preview.  

This setup provides a seamless **WYSIWYG (What You See Is What You Get) experience**, making it ideal for Markdown-based documentation, note-taking, and blogging.  

### **3. Menu Bar with File, View, and Help Options**  
Leda’s menu bar provides quick access to **essential features**, including:  

- **File Menu:**  
  - `Open File`: Load a text or Markdown file into the editor.  
  - `Save File`: Save the current content to a file.  
  - `Clear Text`: Clear all text from the editor.  

- **View Menu:**  
  - `Zoom In/Out`: Adjust the font size dynamically.  
  - `Toggle Dark Mode`: Switch between light and dark themes for better readability.  

- **Help Menu:**  
  - `About`: Displays information about Leda.  
  - `Shortcuts`: Lists keyboard shortcuts for increased efficiency.  

### **4. Fully Customizable UI (JSON-Based or In-App Settings)**  
Leda’s UI can be customized in two ways:  

- **JSON Configuration:**  
  - Users can modify a `config.json` file to change themes, font sizes, colors, and layout preferences.  

- **In-App Theme Menu:**  
  - A graphical settings menu allows users to adjust UI elements **without manually editing JSON files**.  
  - Changes are applied instantly and saved for future sessions.  

This makes Leda **highly adaptable** to different workflows and user preferences.  

### **5. Responsive UI That Adjusts to Window Resizing**  
Leda’s UI automatically resizes and adapts when the window is resized:  

- Uses **Fyne’s responsive layout system** to dynamically adjust component sizes.  
- Ensures that the editor, preview pane, and menus remain usable at all screen sizes.  
- Prevents overlapping or misaligned UI elements.  

### **6. Live Statistics (Character & Line Count, Like Vim)**  
Leda provides **real-time text statistics**, displayed at the bottom of the editor:  

- **Character count:** Shows the total number of characters in the document.  
- **Line count:** Displays the number of lines in the editor, similar to Vim.  
- Updates dynamically as users type, giving immediate feedback on text length.  

### **7. Integrated Live Terminal (Rendered at the Bottom of the Editor)**  
Leda includes a **built-in terminal** at the bottom of the editor, allowing users to:  

- Run shell commands directly within the editor.  
- Execute scripts without switching to an external terminal.  
- See command outputs in real-time, making it useful for developers and power users.  

The terminal is fully integrated, meaning users can work on code or documentation while running commands **without leaving the editor**.

## Challenges & Lessons Learned  

### **1. Performance Issues with Large Files**  
#### **Challenge:**  
When loading large text files (over 100,000 lines), Leda experienced severe lag and unresponsiveness. The issue stemmed from Fyne’s `widget.Entry`, which wasn’t optimized for handling massive amounts of text in real-time.  

#### **Solution:**  
- Implemented **lazy rendering** by only loading visible portions of the text instead of rendering the entire document at once.  
- Optimized the **Markdown preview** by only updating sections that changed instead of re-parsing the entire document on every keystroke.  

### **2. UI Scaling Problems on High-Resolution Displays**  
#### **Challenge:**  
On high-DPI screens (4K and Retina displays), some UI elements appeared too small, making the text editor difficult to use. Fyne’s default scaling did not always apply correctly on all operating systems.  

#### **Solution:**  
- Added a **manual scaling factor** in the settings, allowing users to adjust UI scaling as needed.  
- Implemented **automatic DPI detection**, adjusting font sizes and layout proportions accordingly.  

### **3. Markdown Preview Syncing with Editor Cursor**  
#### **Challenge:**  
Users wanted the Markdown preview to **follow the cursor position** in the editor. Initially, the preview always scrolled to the top when the content updated, causing frustration.  

#### **Solution:**  
- Implemented a **scroll position tracker** that synchronizes the editor and preview pane.  
- Ensured that when a user moves the cursor in the editor, the preview updates to the corresponding section.  

### **4. Terminal Integration and Process Handling**  
#### **Challenge:**  
The built-in terminal initially struggled to handle **long-running processes** and failed to capture certain system outputs, especially on Windows.  

#### **Solution:**  
- Switched from using `os/exec.Command` to a more robust **pseudo-terminal (pty) implementation**, allowing better process management.  
- Added an **output buffer** to handle continuous command outputs without freezing the UI.  

### **5. Cross-Platform Compatibility (Linux, macOS, Windows)**  
#### **Challenge:**  
Leda’s file dialog and keyboard shortcuts behaved **differently on each OS**, leading to inconsistent user experiences.  

#### **Solution:**  
- Used **Fyne’s built-in abstraction layers** to standardize file handling across platforms.  
- Added OS-specific keybindings (e.g., `Cmd` for macOS, `Ctrl` for Windows/Linux).

---

## Team Contributions  

| Team Member       | Contributions                 |
|-------------------|-------------------------------|
| **Aleksy Siek**   | **Key voice in team discussions** - actively shaping project direction and feature priorities. <br> **Leda pull request reviews** - providing feedback and ensuring high code quality. <br> **Directed feature implementation** - deciding what to build and how to approach key functionalities. <br> **Set up and maintained the CI pipeline** - ensuring stable builds across all platforms. <br> **Refactored and optimized code** - improving performance, maintainability, and fixing critical bugs. <br> **Enhanced documentation** - contributing to the **About section** and **README** for better project clarity. |
| **Mate Saary**    | _ |
| **Emon Monsur**   | _ |
| **Oisin Portley** | _ |
| **Eoghan Murch**  | _ |


---


## Installation & Setup  

### Prerequisites  
Before you start developing with Leda, ensure you have the following installed:  

- **Go** (Version 1.18 or higher)  
- **Fyne Framework** (install using `go get`)  

### Clone the Repository  

```sh
git clone git@github.com:Leda-Editor/Leda-Text-Editor.git
cd Leda-Text-Editor
```

### Install Dependencies  

```sh
go get .
```

### Build & Run  

To build the project:  

```sh
go build .
```

To run the project without creating an executable:  

```sh
go run .
```

## License  

This project is licensed under the **Apache 2.0 License**. See the [LICENSE](LICENSE) file for details.  

---

## Additional Resources  

- [Fyne Documentation](https://developer.fyne.io/)  
- [Go Programming Language](https://golang.org/)  

---
