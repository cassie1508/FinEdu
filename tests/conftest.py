import os
import socket
import subprocess
import time
from pathlib import Path

import pytest


REPO_ROOT = Path(__file__).resolve().parents[1]
BACKEND_DIR = REPO_ROOT / "backend"


def _pick_free_port() -> int:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.bind(("127.0.0.1", 0))
        return int(sock.getsockname()[1])


def _wait_for_server(base_url: str, timeout_seconds: float = 30.0) -> None:
    import urllib.request

    deadline = time.time() + timeout_seconds
    health_url = f"{base_url}/health"

    while time.time() < deadline:
        try:
            with urllib.request.urlopen(health_url, timeout=2) as resp:
                if resp.status == 200:
                    return
        except Exception:
            time.sleep(0.25)

    raise RuntimeError(f"Backend did not become healthy within {timeout_seconds}s")


@pytest.fixture(scope="session")
def api_base_url() -> str:
    port = _pick_free_port()
    base_url = f"http://127.0.0.1:{port}"

    env = os.environ.copy()
    env["PORT"] = str(port)
    env["USE_MOCK_DATA"] = "true"
    env["DATABASE_URL"] = ""
    # Keep these empty for deterministic tests that do not call external Finnhub.
    env["FINNHUB_API_KEY"] = ""
    env["MARKET_DATA_API_KEY"] = ""

    process = subprocess.Popen(
        ["go", "run", "./cmd/server"],
        cwd=str(BACKEND_DIR),
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
    )

    try:
        _wait_for_server(base_url)
        yield base_url
    finally:
        process.terminate()
        try:
            process.wait(timeout=10)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait(timeout=5)
