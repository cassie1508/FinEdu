import json
import uuid
import urllib.error
import urllib.request


def _request_json(method: str, url: str, payload=None):
    body = None
    headers = {}
    if payload is not None:
        body = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = urllib.request.Request(url=url, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            raw = resp.read().decode("utf-8")
            return resp.status, json.loads(raw) if raw else {}
    except urllib.error.HTTPError as exc:
        raw = exc.read().decode("utf-8")
        parsed = json.loads(raw) if raw else {}
        return exc.code, parsed


def test_health(api_base_url):
    status, payload = _request_json("GET", f"{api_base_url}/health")
    assert status == 200
    assert payload.get("status") == "ok"




def test_learning_flashcard_crud_and_review(api_base_url):
    status, list_payload = _request_json("GET", f"{api_base_url}/api/v1/learning/flashcards")
    assert status == 200
    assert isinstance(list_payload.get("data"), list)

    status, single_payload = _request_json("GET", f"{api_base_url}/api/v1/learning/flashcards/fc-1")
    assert status == 200
    assert single_payload["data"]["id"] == "fc-1"

    status, not_found_payload = _request_json("GET", f"{api_base_url}/api/v1/learning/flashcards/not-found")
    assert status == 404
    assert "error" in not_found_payload

    unique_title = f"Test Flashcard {uuid.uuid4()}"
    create_body = {
        "title": unique_title,
        "category": "Unit Test",
        "whyItMatters": "Verifies API behavior",
        "definition": "An API contract validation card",
        "example": "Sample example",
        "commonMisconception": ["This is production data"],
    }
    status, created_payload = _request_json("POST", f"{api_base_url}/api/v1/learning/flashcards", payload=create_body)
    assert status == 201
    created_id = created_payload["data"]["id"]

    update_body = {
        "title": f"{unique_title} Updated",
        "category": "Unit Test",
        "whyItMatters": "Still verifies API behavior",
        "definition": "Updated definition",
        "example": "Updated example",
        "commonMisconception": ["Still production data"],
    }
    status, updated_payload = _request_json(
        "PUT",
        f"{api_base_url}/api/v1/learning/flashcards/{created_id}",
        payload=update_body,
    )
    assert status == 200
    assert updated_payload["data"]["title"] == update_body["title"]

    status, reviewed_payload = _request_json(
        "POST",
        f"{api_base_url}/api/v1/learning/flashcards/{created_id}/review",
    )
    assert status == 200
    assert reviewed_payload["data"]["reviewCount"] == 1

    status, deleted_payload = _request_json("DELETE", f"{api_base_url}/api/v1/learning/flashcards/{created_id}")
    assert status == 200
    assert deleted_payload.get("message") == "flashcard deleted"

    status, deleted_missing_payload = _request_json("GET", f"{api_base_url}/api/v1/learning/flashcards/{created_id}")
    assert status == 404
    assert "error" in deleted_missing_payload


def test_learning_flashcard_validation_and_not_found_paths(api_base_url):
    bad_create = {"title": "Missing required fields"}
    status, payload = _request_json("POST", f"{api_base_url}/api/v1/learning/flashcards", payload=bad_create)
    assert status == 400
    assert payload.get("error") == "invalid request body"

    update_body = {
        "title": "x",
        "category": "x",
        "whyItMatters": "x",
        "definition": "x",
        "example": "x",
        "commonMisconception": [],
    }
    status, payload = _request_json(
        "PUT",
        f"{api_base_url}/api/v1/learning/flashcards/not-found",
        payload=update_body,
    )
    assert status == 404

    status, payload = _request_json("DELETE", f"{api_base_url}/api/v1/learning/flashcards/not-found")
    assert status == 404

    status, payload = _request_json("POST", f"{api_base_url}/api/v1/learning/flashcards/not-found/review")
    assert status == 404


def test_learning_center_endpoints_without_finnhub_key(api_base_url):
    status, payload = _request_json("GET", f"{api_base_url}/api/v1/learning_center/resources")
    assert status == 500
    assert "error" in payload

    status, payload = _request_json("GET", f"{api_base_url}/api/v1/learning_center/abc")
    assert status == 400
    assert payload.get("error") == "invalid resources_id"

    status, payload = _request_json("GET", f"{api_base_url}/api/v1/learning_center/123")
    assert status == 500
    assert "error" in payload

