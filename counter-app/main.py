from flask import Flask, render_template_string, redirect, url_for, flash
import os
import pathlib

app = Flask(__name__)
app.secret_key = os.getenv("SECRET_KEY", "a-secure-development-secret")

APP_COLOR = os.getenv("COLOR", "black")
DATA_DIR = pathlib.Path(os.getenv("COUNTER_DIR", "/app/data"))
DATA_DIR.mkdir(parents=True, exist_ok=True)
COUNTER_FILE = DATA_DIR / "counter.txt"


def load_counter_from_file() -> int:
    """Reads the counter value from the file, defaulting to 0 on error."""
    try:
        return int(COUNTER_FILE.read_text().strip())
    except (FileNotFoundError, ValueError):
        # If the file doesn't exist or contains invalid data, start at 0.
        return 0

def save_counter_to_file(value: int) -> None:
    """Saves the given integer value to the counter file."""
    COUNTER_FILE.write_text(str(value), encoding="utf-8")


counter = load_counter_from_file()


HTML = """
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flask Counter</title>
    <style>
        body { font-family: sans-serif; max-width: 600px; margin: 2rem auto; background-color: #f4f4f9; color: #333; }
        h1 { color: {{ color }}; }
        button {
            padding: 0.5rem 1rem;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 1rem;
            margin-right: 0.5rem;
            color: white;
        }
        .btn-increment { background-color: #007bff; }
        .btn-save { background-color: #28a745; }
        .btn-load { background-color: #ffc107; color: #333; }
        .flash-msg { color: #155724; background-color: #d4edda; padding: 1rem; border-radius: 5px; margin-top: 1rem; }
        form { display: inline-block; }
    </style>
</head>
<body>
    <h1>Counter: {{ counter }}</h1>
    <form action="{{ url_for('increment') }}" method="post">
      <button class="btn-increment">Increment</button>
    </form>
    <form action="{{ url_for('save') }}" method="post">
      <button class="btn-save">Save to file</button>
    </form>
    <form action="{{ url_for('load') }}" method="post">
      <button class="btn-load">Load from file</button>
    </form>
    {% with messages = get_flashed_messages() %}
      {% if messages %}
        <p class="flash-msg">{{ messages[0] }}</p>
      {% endif %}
    {% endwith %}
</body>
</html>
"""

@app.route("/")
def index():
    """Renders the main page."""
    return render_template_string(HTML, counter=counter, color=APP_COLOR)

@app.route("/increment", methods=["POST"])
def increment():
    """Increments the in-memory counter."""
    global counter
    counter += 1
    return redirect(url_for("index"))

@app.route("/save", methods=["POST"])
def save():
    """Saves the current in-memory counter value to the file."""
    save_counter_to_file(counter)
    flash(f"Saved value: {counter}")
    return redirect(url_for("index"))

@app.route("/load", methods=["POST"])
def load():
    """Loads the counter value from the file into memory."""
    global counter
    counter = load_counter_from_file()
    flash(f"Loaded value: {counter}")
    return redirect(url_for("index"))


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
