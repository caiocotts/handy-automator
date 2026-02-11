# 1️⃣ Install uv (if not installed)
Windows (PowerShell)
irm https://astral.sh/uv/install.ps1 | iex

macOS / Linux
curl -LsSf https://astral.sh/uv/install.sh | sh


or you can use `pip` if you have that installed
```
pip install uv
```

Restart your terminal and verify:

```
uv --version
```

# 2️⃣ Create the Virtual Environment

From the project root (where `pyproject.toml` is located):

```
uv venv
```

# 3️⃣ Install All Dependencies

All dependencies are already defined in `pyproject.toml`:

```
uv sync
```

# 4️⃣ Activate the Environment
Windows (PowerShell)

```
.venv\Scripts\activate
```

macOS / Linux source 
```
.venv/bin/activate
```

# 5️⃣ Run the Application
```
python main.py
```


Press `Q` to exit.