import numpy as np
from unittest.mock import MagicMock, patch
import unittest
import face_recognition as fr

class TestFaceRecognition(unittest.TestCase):

    def setUp(self):
        # Reset global state to prevent test interference
        fr.tracked_faces = []
        fr.auth_times = {}

    @patch('face_recognition.face_classifier')
    def test_detect_no_faces(self, mock_face_classifier):
        """Verifies that the function returns an empty list when no faces are detected."""
        mock_face_classifier.detectMultiScale.return_value = []
        frame = np.zeros((200, 200, 3), dtype=np.uint8)
        result = fr.detecting_bounding_box(frame)
        
        self.assertEqual(len(result), 0, "Result should be empty when no faces are present.")

    @patch('face_recognition.face_classifier')
    def test_detect_one_face(self, mock_face_classifier):
        """Verifies that the function correctly identifies and labels a single face."""
        mock_face_classifier.detectMultiScale.return_value = [(50, 50, 100, 100)]
        frame = np.zeros((400, 400, 3), dtype=np.uint8)
        result = fr.detecting_bounding_box(frame)
        
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0][0], (50, 50, 100, 100))
        self.assertEqual(result[0][1], "Unknown", "A new face should initially be 'Unknown'.")

    @patch('face_recognition.face_classifier')
    def test_detect_multiple_faces(self, mock_face_classifier):
        """Verifies that the function identifies multiple faces in the same frame."""
        mock_face_classifier.detectMultiScale.return_value = [(10, 10, 50, 50), (120, 120, 60, 60)]
        frame = np.zeros((400, 400, 3), dtype=np.uint8)
        result = fr.detecting_bounding_box(frame)
        
        self.assertEqual(len(result), 2)
        self.assertEqual(result[0][0], (10, 10, 50, 50))
        self.assertEqual(result[1][0], (120, 120, 60, 60))

    @patch('face_recognition.DeepFace.represent')
    def test_authorize_success(self, mock_represent):
        # Mock the DeepFace.represent function to return a valid embedding
        mock_represent.return_value = [{"embedding": [0.1] * 128}]

        # Mock the check_faces function to return a name
        with patch('face_recognition.check_faces', return_value="test_user"):
            name = fr.authorize(np.zeros((100, 100, 3), dtype=np.uint8))
            self.assertEqual(name, "test_user")

if __name__ == '__main__':
    unittest.main()

    