# BLENDS

## Getting Started

All instructions are done in Bash for Linux/macOS and Powershell (>=2.0) for Windows.

1. Install [Python 3.](https://python.org)
2. Make sure you've installed `venv`.
    ```bash
    $ python3 -m venv
    usage: venv [-h] [--system-site-packages] [--symlinks | --copies] [--clear]
                [--upgrade] [--without-pip] [--prompt PROMPT]
                ENV_DIR [ENV_DIR ...]
    venv: error: the following arguments are required: ENV_DIR
    ```
3. Some Linux disto doesn't install `venv` automatically. For Debian/Ubuntu, you can install `venv` by following:
    ```bash
    $ sudo apt install python3-venv
    ```
4. Initialize new virtual environment.  This will create a directory "venv" in the current directory.
    ```bash
    $ python3 -m venv venv
    ```
5. Activate the virtual environment. This will start a virtual environment.
    * For Linux/macOS
        ```bash
        $ source venv/bin/activate
        ```
    * For Windows
        ```powershell
        PS > .\venv\Scripts\activate
        ```
6. Install Requirements. This will install the following 3 modules: jupyter 1.0.0, requests 2.0.0, and Flask 1.0.2.
    ```bash
    (venv) $ pip install -r requirements.txt
    ```
7. Start a Jupyter notebook.  You will be directed to a new Jupyter notebook page on your default browser.
    ```bash
    (venv) $ jupyter notebook
    ```

