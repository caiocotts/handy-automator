import mediapipe as mp
import numpy as np

# Global variable to store the latest pose landmarks
latest_pose_result = None

def update_pose_result(result: mp.tasks.vision.PoseLandmarkerResult, output_image: mp.Image, timestamp_ms: int):
    """
    Updates the global pose result with the latest detected landmarks.

    Args:
        - result (mp.tasks.vision.PoseLandmarkerResult):
          The pose detection result returned by the MediaPipe model.
        - output_image (mp.Image):
          The processed image associated with the detection result.
        - timestamp_ms (int):
          The timestamp (in milliseconds) corresponding to the frame.
    """
    global latest_pose_result
    if result.pose_landmarks:
        latest_pose_result = result.pose_landmarks
    else:
        latest_pose_result = None

def create_pose_landmarker(options):
    return mp.tasks.vision.PoseLandmarker.create_from_options(options)

def get_nose_coords(pose_landmarks, frame_w, frame_h):
    """
    Retrieves the pixel coordinates of the nose from pose landmarks.

    Args:
        - pose_landmarks (list): List containing detected pose landmarks
          from MediaPipe PoseLandmarker.
        - frame_w (int): Width of the video frame in pixels.
        - frame_h (int): Height of the video frame in pixels.

    Returns:
        - tuple[int, int]:
          The (x, y) pixel coordinates of the nose.
    """
    nose = pose_landmarks[0][mp.tasks.vision.PoseLandmark.NOSE]
    return int(nose.x * frame_w), int(nose.y * frame_h)

def get_wrist_coords(pose_landmarks, frame_w, frame_h, side="left"):
    """
    Returns pixel coordinates for the requested wrist.

    Args:
        - pose_landmarks (list): List containing detected pose landmarks
          from MediaPipe PoseLandmarker.
        - frame_w (int): Width of the video frame in pixels.
        - frame_h (int): Height of the video frame in pixels.
        - side (str, optional): Specifies which wrist to retrieve.
        - Accepts "left" or "right". Defaults to "left".

    Returns:
        - tuple[int, int]:
          The (x, y) pixel coordinates of the requested wrist.
    """
    idx = mp.tasks.vision.PoseLandmark.LEFT_WRIST if side == "left" else mp.tasks.vision.PoseLandmark.RIGHT_WRIST
    wrist = pose_landmarks[0][idx]
    return int(wrist.x * frame_w), int(wrist.y * frame_h)

def get_dist(p1_x, p1_y, p2_x, p2_y):
    """
    Calculates the Euclidean distance between two points in 2D space.

    Args:
        - p1_x (float): X-coordinate of the first point.
        - p1_y (float): Y-coordinate of the first point.
        - p2_x (float): X-coordinate of the second point.
        - p2_y (float): Y-coordinate of the second point.

    Returns:
        - float
          The Euclidean distance between the two points.
    """
    return np.sqrt((p1_x - p2_x)**2 + (p1_y - p2_y)**2)