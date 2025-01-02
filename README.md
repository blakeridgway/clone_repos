# Clone Repos Python App

This Python app uses the GitHub CLI (`gh`) to clone all repositories from a specified GitHub user account, excluding repositories that are forked or archived. Repositories are cloned into a directory of your choice.

## Features
- Filters out forked and archived repositories.
- Clones repositories into separate subdirectories.
- Easy to use with simple prompts.

## Prerequisites
1. **Python 3.6+**
   Ensure Python is installed on your system. Check by running:
   ```bash
   python3 --version
   ```

2. **GitHub CLI (`gh`)**
   Install the GitHub CLI if not already installed. Follow the [GitHub CLI installation guide](https://cli.github.com/).

3. **Git**
   Ensure Git is installed. Check by running:
   ```bash
   git --version
   ```

## Installation
1. Clone this repository or download the script.
2. Ensure the script is executable:
   ```bash
   chmod +x clone_repos.py
   ```

## Usage
Run the script with Python:
```bash
python3 clone_repos.py
```

### Steps:
1. Enter your GitHub username when prompted.
2. Specify the directory where repositories should be cloned.
3. The app will fetch all non-forked, non-archived repositories and clone them into subdirectories within the specified folder.

## Example
```bash
$ python3 clone_repos.py
Enter your GitHub username: exampleuser
Enter the directory where you want to clone the repos: /path/to/clone
exampleuser/repo1 has been cloned into /path/to/clone/repo1
exampleuser/repo2 has been cloned into /path/to/clone/repo2
All repos have been cloned into /path/to/clone
```

## Dependencies
- `gh` (GitHub CLI)
- `git`

## Error Handling
- If the `gh` CLI is not installed or configured, the script will notify you.
- Errors during repository cloning will be logged to the console.

## Contributing
Contributions are welcome! Please submit a pull request or file an issue to suggest improvements or report bugs.

## License
This project is licensed under the GPLv3 License. See the LICENSE file for details.