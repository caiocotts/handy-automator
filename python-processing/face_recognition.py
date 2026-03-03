import cv2
from deepface import DeepFace
import time
from scipy.spatial.distance import cosine
import mediapipe as mp
import math
import face_embeddings as fe


 # Constants

face_classifier = cv2.CascadeClassifier(
    cv2.data.haarcascades + "haarcascade_frontalface_default.xml"
) 
last_check = 0
CHECK_INTERVAL = 1.0      # Seconds between recognition checks
THRESHOLD = 0.6           # Cosine similarity threshold
AUTH_LOCK_TIME = 7.0      # Lock duration after success

authenticated = False
last_auth_time = 0     


known_embeddings = {}  # name -> embedding
tracked_faces = []

#TODO remove when embeddings are on db and replace with this:
#known_embeddings = face_embeddings.get_faces()

known_embeddings = fe.get_faces_old()



def detecting_bounding_box(frame):
    """
    Detects faces within the given video frame, draws bounding boxes, 
    and performs periodic facial authentication.

    Args:
        - A cv2 frame

    Returns: 
        - list[tuple[tuple[int, int, int, int], str]]
        - A list containing bounding box coordinates (x, y, w, h) paired 
          with the detected or recognized name.
    """
    global last_check, authenticated, last_auth_time, tracked_faces

    if authenticated and time.time() - last_auth_time < AUTH_LOCK_TIME:
            cv2.putText(frame, "AUTHORIZED",
                        (20, 40), cv2.FONT_HERSHEY_SIMPLEX,
                        1, (0, 255, 0), 2)
    else:
        authenticated = False
    
    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
    faces = face_classifier.detectMultiScale(gray, 1.05, 5, minSize=(100, 100))

    current_time = time.time()
    current_tracked = []
    
    for (x, y, w, h) in faces:
        cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 3)
        detected_name = "Unknown"

        cx, cy = x + w/2, y + h/2
        for t_box, t_name in tracked_faces:
            tx, ty, tw, th = t_box
            tcx, tcy = tx + tw/2, ty + th/2
            # If the center of the face hasn't moved much, assume it's the same person
            if math.hypot(cx - tcx, cy - tcy) < 50: 
                detected_name = t_name
                break

        # Run DeepFace recognition if the interval has passed
        if current_time - last_check > CHECK_INTERVAL:
            face_crop = frame[y:y+h, x:x+w]
            match = authorize(face_crop)
            if match:
                detected_name = match
                
        # Store the result for this frame
        current_tracked.append(((x, y, w, h), detected_name))

    # Reset the timer if we just checked
    if current_time - last_check > CHECK_INTERVAL:
        last_check = current_time
        
    # Update our global tracker
    tracked_faces = current_tracked 
    
    return tracked_faces # Returns a list of ((x, y, w, h), "Name")


def authorize(face_img):
    """
    Handles authorization logic. Prints results to CLI
    Args:
        - face_img (numpy.ndarray): Cropped image of a detected face.

    Returns:
        - str | None
        - The name of the authenticated individual if a match is found,
          otherwise None.
    """
    match = check_faces(face_img)
    if match:
        print("Authorized:", match)
        return match # Return the matched name
    else:
        print("Unknown face")
        return None


def check_faces(face_img):
    """
    Generates a facial embedding from the provided face image and
    compares it against stored reference embeddings to determine identity.

    Args:
        - face_img (numpy.ndarray): Cropped image of a detected face.

    Returns:
        - str | None
        - The name of the matched individual if authentication succeeds,
          otherwise None.
    """
    global authenticated, last_auth_time

    try:
        live_embedding = DeepFace.represent(
            img_path=face_img,
            model_name="ArcFace",
            detector_backend="opencv",
            enforce_detection=True
        )[0]["embedding"]
    except Exception:
        return None

    for name, ref_embedding in known_embeddings.items():
        distance = cosine(ref_embedding, live_embedding)

        if distance < THRESHOLD:
            authenticated = True
            last_auth_time = time.time()
            return name

    return None
