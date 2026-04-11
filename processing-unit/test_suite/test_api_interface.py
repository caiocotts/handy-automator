import numpy as np
from unittest.mock import MagicMock, patch
import api_interface as api


def test_ping_return_200():
    """Verifies that ping() returns True when the API responds with status code 200."""
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 200
        mock_get.return_value.json.return_value = {'message': 'pong'}
        assert api.ping() == True


def test_ping_return_500():
    """Verifies that ping() returns False when the API responds with status code 500."""
    with patch('requests.get') as mock_get:
        mock_get.return_value.status_code = 500
        assert api.ping() == False

def test_auth_user_api_call_return_200():
    """Verifies that auth_user_api_call() returns the access token when the API responds with status code 200."""
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 200
        mock_post.return_value.json.return_value = {'accessToken': 'test_token'}
        assert api.auth_user_api_call([1.0, 2.0, 3.0], 'test_user') == 'test_token'

def test_auth_user_api_call_return_404():
    """Verifies that auth_user_api_call() returns None when the API responds with status code 404."""
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 404
        assert api.auth_user_api_call([0.1, 0.2], 'test_user') == None


def test_workflow_api_call_return_gesture_id_0():
    """Verifies that workflow_api_call() returns None when the gesture_id is 0."""
    with patch('requests.post') as mock_post:
        mock_post.return_value.status_code = 401
        assert api.workflow_api_call(0, 'test_token') == None