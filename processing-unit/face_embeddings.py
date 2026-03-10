import os
from deepface import DeepFace

DB_PATH = os.path.join(os.path.abspath(os.getcwd()), "database") 
AUTH_DIR = os.path.join(DB_PATH, "auth_users")


def get_faces_old() -> dict:
    '''
    Old version remove when DB implemented
    '''
    known_embeddings = {}
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
    
    if not known_embeddings:
        print(f"There is no images in the database")
    return known_embeddings

def get_faces():
    '''
    Querey the DB and pull all the face embeddings {Name: VECTOR}
    '''
    return saved_embeddings

def add_face()-> bool: 
    '''
    When the user enters the setup mode to add a new user to their list of authroized users, this should be run. 
    It should search the frame for a new face and once detected it should pull an embedding and save it to the db with 
    the user inputted name.
    '''
