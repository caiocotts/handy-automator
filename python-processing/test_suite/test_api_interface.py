import numpy as np
from unittest.mock import MagicMock, patch
import api_interface as api

def test_ping_return_200():
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 200
        mock_get.return_value.json.return_value = {'message': 'pong'}
        assert api.ping() == True

def test_ping_return_500():
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 500
        assert api.ping() == False

#________________auth_user_api_call()___________________________________

def test_auth_user_api_call_return_200(): #BROKEN
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 200
        mock_post.return_value.json.return_value = {'authToken': 'test_token'}
        assert api.auth_user_api_call([1.0, 2.0, 3.0], 'test_user') == 'test_token'

def test_auth_user_api_call_return_404():
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 404
        assert api.auth_user_api_call([0.1, 0.2], 'test_user') == None
#_________________workflow_api_call()_________________________________
def test_workflow_api_call_return_gesture_id_0():
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 401
        assert api.workflow_api_call('test_user', 'test_token') == None
