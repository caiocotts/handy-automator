
import mediapipe as mp
import cv2 as cv
import numpy as np
import time
from mediapipe.tasks import python
from mediapipe.tasks.python import vision

# --- Configuration ---
model_path = './gesture_recognizer.task'
BaseOptions = mp.tasks.BaseOptions
GestureRecognizer = mp.tasks.vision.GestureRecognizer
GestureRecognizerOptions = mp.tasks.vision.GestureRecognizerOptions
GestureRecognizerResult = mp.tasks.vision.GestureRecognizerResult
VisionRunningMode = mp.tasks.vision.RunningMode

# --- Constants ---
HAND_CONNECTIONS = [
    (0, 1), (1, 2), (2, 3), (3, 4),
    (0, 5), (5, 6), (6, 7), (7, 8),
    (5, 9), (9, 10), (10, 11), (11, 12),
    (9, 13), (13, 14), (14, 15), (15, 16),
    (13, 17), (17, 18), (18, 19), (19, 20),
    (0, 17)
]

# Global variable to store the latest result
latest_result = None

# --- Custom Drawing Function (Pure OpenCV) ---
# No more dependency on mp.solutions.drawing_utils or Protobufs!
def draw_landmarks_on_image(rgb_image, detection_result):
    hand_landmarks_list = detection_result.hand_landmarks
    gestures_list = detection_result.gestures  # This is the new list we need
    annotated_image = np.copy(rgb_image)
    height, width, _ = annotated_image.shape

    # Zip through landmarks and gestures together
    # enumerate helps us position the text if there are multiple hands
    for idx, hand_landmarks in enumerate(hand_landmarks_list):

        # --- 1. Draw Landmarks & Connections ---
        pixel_coords = []
        for landmark in hand_landmarks:
            cx, cy = int(landmark.x * width), int(landmark.y * height)
            pixel_coords.append((cx, cy))
            cv.circle(annotated_image, (cx, cy), 5, (0, 255, 0), -1)

        for connection in HAND_CONNECTIONS:
            cv.line(annotated_image, pixel_coords[connection[0]],
                    pixel_coords[connection[1]], (255, 0, 0), 2)

        # --- 2. Draw Gesture Text ---
        if gestures_list and idx < len(gestures_list):
            # Get the top gesture for this hand
            top_gesture = gestures_list[idx][0]
            gesture_name = top_gesture.category_name
            confidence = round(top_gesture.score, 2)

            # Position text slightly above the wrist (landmark 0)
            text_x, text_y = pixel_coords[0]
            display_text = f"{gesture_name} ({confidence})"

            # Draw a simple background rectangle for readability
            cv.rectangle(annotated_image, (text_x, text_y - 30),
                         (text_x + 200, text_y), (0, 0, 0), -1)

            # Write the gesture name
            cv.putText(annotated_image, display_text, (text_x, text_y - 10),
                       cv.FONT_HERSHEY_SIMPLEX, 0.8, (255, 255, 255), 2, cv.LINE_AA)

    return annotated_image

# --- Callback ---

def update_result(result: GestureRecognizerResult, output_image: mp.Image, timestamp_ms: int):
    global latest_result
    latest_result = result

def create_recognizer(options):
    return GestureRecognizer.create_from_options(options)
