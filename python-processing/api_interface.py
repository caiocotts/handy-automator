import requests
import json
from deepface import DeepFace
import face_recognition as fr

def get_public_ipv4():
    try:
        # Use an external API that returns only the IP address
        response = requests.get('https://api.ipify.org/?format=json')
        if response.status_code == 200:
            return response.json()['ip']
        else:
            return f"Error: {response.status_code}"
    except requests.exceptions.RequestException as e:
        return f"An error occurred: {e}"


# Module that sends tempEmb (with a toList() )

def send_embeddings(uid):
    # convert embeddings into a list
    embedding = fr.live_embedding
    
    if embedding is None:
        return

    payload = {
        "user_id": uid,
        "embedding": embedding
    }

    serialized_embeddings = json.dumps(payload)

    print("Sending NEW embedding...")
    print(serialized_embeddings[:100], "...")

    fr.live_embedding = None

    # requests.post('', json=embedding)
    # response = requests.post('https://api.ipify.org', json=serialized_embeddings)
    # requests.post(url, data={key: value}, json={key: value}, args)

def ping(ip_address=get_public_ipv4(), port="3000"):
    r = requests.get("http://" + ip_address + ":" + port + "/api/ping")
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

