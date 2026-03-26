import unittest
from unittest.mock import patch
import config

class TestConfig(unittest.TestCase):
    @patch('config.get_available_cameras')
    def test_get_available_cameras_more_than_one_camera(self, mock_available_camera):
        mock_available_camera.return_value = {0: "webcame", 1: "roboshot"}
        self.assertIsNotNone(config.get_ptz_camera())

    @patch('config.get_available_cameras')
    def test_get_available_cameras_no_cameras(self, mock_available_camera):
        mock_available_camera.return_value = {}
        self.assertIsNone(config.get_ptz_camera())

    @patch('config.get_available_cameras')
    def test_get_ptz_camera_roboshot(self, mock_available_camera):
        mock_available_camera.return_value = {0: "webcam", 1: "Roboshot Camera"}
        self.assertEqual(config.get_ptz_camera(), 1)
if __name__ == '__main__':
    unittest.main()