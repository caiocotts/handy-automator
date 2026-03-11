
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
latest_result = None # This will be the one to be accessed when you want the classified gesture
"""
This variable will serve as the container for the classified gesture

    -> gestures: List[List[category_module.Category]] 
        -> category_module.Category: [str: category_name, float: score, int: index]
    -> handedness: List[List[category_module.Category]] (similar)
        -> category_module.Category: [str: category_name, float: score, int: index]
    -> hand_landmarks: List[List[landmark_module.NormalizedLandmark]]
        -> landmark_module.NormalizedLandmark: [float: x, float: y, float: z]
    -> hand_world_landmarks: List[List[landmark_module.Landmark]]
        -> landmark_module.Landmark: [float: x, float: y, float: z]

"""

def update_result(result: GestureRecognizerResult, output_image: mp.Image, timestamp_ms: int):
    """
    Receives the result of the classification when 
    gesture_rec.recognize_async() is performed and updates the global
    latest_result.

    Args: 
        - a GestureRecognizerResult object 
            -> gestures: Recognized hand gestures of detected hands where its index is always -1
            -> handedness: Classification of handedness
            -> hand_world_landmarks: Detected hand landmarks in normalized image coordinates
        - an Image object
        - a timestamp of when the frame is taken

    Returns: None

    """
    global latest_result
    latest_result = result


options = GestureRecognizerOptions( 
    base_options=BaseOptions(model_asset_path=model_path), # finds which model to use (hand_landmarker.task)
    running_mode=VisionRunningMode.LIVE_STREAM, # whether an image, video, or live-stream (can only use detect_async() if live stream)
    result_callback=update_result, # where the result is sent
    min_hand_detection_confidence=0.7,
    num_hands=2
)

# --- Custom Drawing Function ---
def draw_landmarks_on_image(rgb_image, detection_result):
    """
    Draws the landmarks on the frames of the detected hand along with the type of 
    gesture the hands are and their handedness.

    Args: 
        - a cv2 frame
        - a GestureRecognizerResult object 
    
    Returns: 
        - an annotated frame
    """
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

def create_recognizer(options):
    """
    Loading of the model from the gesture_recognizer.task file that contains the 
    set of instructions needed to perform a gesture classification

    Args: 
        - a GestureRecognizerOptions object
            -> provides customizations to the gesture recognizer task
    
    Returns: 
        - a GestureRecognizer object created from the options
    """
    return GestureRecognizer.create_from_options(options)
