import requests
import json

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

