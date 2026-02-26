import mediapipe as mp
import cv2

# Global variable to store the latest pose landmarks
latest_pose_result = None

def update_pose_result(result: mp.tasks.vision.PoseLandmarkerResult, output_image: mp.Image, timestamp_ms: int):
    global latest_pose_result
    if result.pose_landmarks:
        latest_pose_result = result.pose_landmarks
    else:
        latest_pose_result = None

def create_pose_landmarker(options):
    return mp.tasks.vision.PoseLandmarker.create_from_options(options)

def get_nose_coords(pose_landmarks, frame_w, frame_h):
    """Returns the (x, y) pixel coordinates of the nose."""
    nose = pose_landmarks[0][mp.tasks.vision.PoseLandmark.NOSE]
    return int(nose.x * frame_w), int(nose.y * frame_h)

def get_wrist_coords(pose_landmarks, frame_w, frame_h, side="left"):
    """Returns pixel coordinates for the requested wrist."""
    idx = mp.tasks.vision.PoseLandmark.LEFT_WRIST if side == "left" else mp.tasks.vision.PoseLandmark.RIGHT_WRIST
    wrist = pose_landmarks[0][idx]
    return int(wrist.x * frame_w), int(wrist.y * frame_h)