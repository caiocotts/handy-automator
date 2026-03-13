import requests
import json
# Need to adjust the header names according to what the api specs say

def ping():
    r = requests.get("http://localhost:3000/api/ping")
    print(r.json())
    return r.status_code == 200    

def auth_user_api_call(embeddings: list[float], username : str) -> str | None:
    r = requests.post("http:localhost:3000/api/login/face", data={"Face Embeddings": embeddings, "Username": username})
    if r.status_code == 200:
        return r.text()
    else:
        return None
    

def workflow_api_call(gesture_id:str, auth_token):
    r = requests.post("http://localhost:3000/api/workflows",data={"Gesture Id": gesture_id, "Auth Token": auth_token})
    print(r.text)
    return r.status_code == 200
    
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

# auth_user_api_call(r"C:\Users\chris\OneDrive\Desktop\handy-automator\python-processing\database\embeddings\PXL_20260312_014832345.MP.npy", "Chris")