import os
import numpy as np
from deepface import DeepFace

DB_PATH = os.path.join(os.path.abspath(os.getcwd()), "database") 
AUTH_DIR = os.path.join(DB_PATH, "auth_users")
EMBEDDINGS_DIR = os.path.join(DB_PATH, "embeddings")

if not os.path.exists(EMBEDDINGS_DIR):
    os.makedirs(EMBEDDINGS_DIR)

def get_faces_old() -> dict:
    '''
    Loads embeddings from individual files. If files don't exist, creates them from images.
    '''
    known_embeddings = {}
    if not os.path.exists(AUTH_DIR):
        return known_embeddings

    for file in os.listdir(AUTH_DIR):
        if file.lower().endswith((".jpg", ".png", ".jpeg")):
            name = os.path.splitext(file)[0]
            embedding_path = os.path.join(EMBEDDINGS_DIR, f"{name}.npy")

            if os.path.exists(embedding_path):
                known_embeddings[name] = np.load(embedding_path)
                print(f"Loaded existing embedding for {name}")
            else:
                path = os.path.join(AUTH_DIR, file)
                embedding = DeepFace.represent(
                    img_path=path,
                    model_name="ArcFace",
                    detector_backend="retinaface",
                    enforce_detection=True
                )[0]["embedding"]
                
                np.save(embedding_path, embedding)
                known_embeddings[name] = embedding
                print(f"Created and loaded embedding for {name}")
    
    if not known_embeddings:
        print(f"There are no images/embeddings in the database")
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
