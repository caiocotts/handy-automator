import cv2
from deepface import DeepFace
import time
from scipy.spatial.distance import cosine
import os

os.environ["TF_CPP_MIN_LOG_LEVEL"] = "2" # Silence warnings




# Constants 
face_classifier = cv2.CascadeClassifier(
    cv2.data.haarcascades + "haarcascade_frontalface_default.xml"
)

cam = cv2.VideoCapture(0)
DB_PATH = os.path.join(os.path.abspath(os.getcwd()), "database")
last_check = 0
CHECK_INTERVAL = 1.0      # Seconds between recognition checks
THRESHOLD = 0.6           # Cosine similarity threshold
AUTH_LOCK_TIME = 7.0      # Lock duration after success

authenticated = False
AUTH_DIR = os.path.join(DB_PATH, "auth_users")

known_embeddings = {}  # name -> embedding

#TODO move to seperate file to run as a script
for file in os.listdir(AUTH_DIR):
    if file.lower().endswith((".jpg", ".png", ".jpeg")):
        path = os.path.join(AUTH_DIR, file)
        name = os.path.splitext(file)[0]

        embedding = DeepFace.represent(
            img_path=path,
            model_name="ArcFace",
            detector_backend="retinaface",
            enforce_detection=True
        )[0]["embedding"]

        known_embeddings[name] = embedding
        print(f"Loaded embedding for {name}")


def detecting_bounding_box(frame):
    global last_check, authenticated

    # Skip detection during cooldown
    if authenticated and time.time() - last_auth_time < AUTH_LOCK_TIME:
        cv2.putText(frame, "AUTHORIZED",
                    (20, 40), cv2.FONT_HERSHEY_SIMPLEX,
                    1, (0, 255, 0), 2)
        return []

    gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
    faces = face_classifier.detectMultiScale(
        gray, 1.05, 5, minSize=(100, 100)
    )

    current_time = time.time()

    for (x, y, w, h) in faces:
        cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 3)

        if current_time - last_check > CHECK_INTERVAL:
            face_crop = frame[y:y+h, x:x+w]
            authorize(face_crop)
            last_check = current_time

    return faces



def authorize(face_img):
    match = check_faces(face_img)
    if match:
        print("Authorized:", match)
    else:
        print("Unknown face")


def check_faces(face_img):
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



#STEPS
'''
1. Pre-load all the authorized images into vectors
2. Scan for faces
2.5. Zoom in if needed and check different faces
3. Find if any face matches an auth face
4. Crop body and follow?
5.? Check for gestures within cropped frame
'''