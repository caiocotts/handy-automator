import requests
import json

def ping():
    r = requests.get("http://localhost:3000/api/ping")
    print(r.text)
    return r.status_code == 200

def api_call(gesture_id:str, auth_token):
    r = requests.get(f"http://localhost:3000/api/workflows{gesture_id}")
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

