import requests


# Need to adjust the header names according to what the api specs say

def ping():
    try:
        r = requests.get("http://localhost:3000/api/ping")
        print(r.json())
        return r.status_code == 200
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return False


def auth_user_api_call(embeddings: list[float], username: str) -> str | None:
    try:
        payload = {"userId": username, "embedding": embeddings}
        r = requests.post("http://localhost:3000/api/auth/login/face", json=payload, headers=f"")
        if r.status_code == 200:
            token_json = r.json()
            return token_json["accessToken"]
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None


def workflow_api_call(gesture_id: int, auth_token: str):
    try:
        payload = {"gestureId": gesture_id}
        header = {"Authorization": f"Bearer {auth_token}"}
        r = requests.post("http://localhost:3000/api/workflow/trigger", json=payload, headers=header)
        if r.status_code != 200:
            print(f"Workflow API call failed with status code {r.status_code}: {r.text}")
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")
        return None


""" HIGH PRIO
    - Get all workflows
    - Get workflow by id
"""
""" LOW PRIO
    - Panning
    - Tilting
    - Zooming
"""
""" OPTIONAL
    - Adjust camera settings (ISO, exposure, etc.) or just auto
"""
