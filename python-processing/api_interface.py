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


def auth_user_api_call(embeddings: list[float], username : str) -> str | None:
    try:
        payload = {"userId":username,"embedding":embeddings}
        r = requests.post("http://localhost:3000/api/auth/login/face", json=payload)
        if r.status_code == 200:
            token_json = r.json()
            return token_json["accessToken"]
    except requests.exceptions.RequestException as e: 
        print(f"Error: {e}")
        return None
    
    

def workflow_api_call(gesture_id:str, auth_token):
    try:
        payload = {"gestureId": gesture_id, "authToken": auth_token} #Still need to check header names
        r = requests.post("http://localhost:3000/api/workflows",json=payload)
        if r.status_code == 200:
            return r.json() #Check what the output is
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