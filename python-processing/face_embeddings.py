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
    
    # Check EMBEDDINGS_DIR first to load pre-calculated embeddings
    if os.path.exists(EMBEDDINGS_DIR):
        for file in os.listdir(EMBEDDINGS_DIR):
            if file.lower().endswith(".npy"):
                name = os.path.splitext(file)[0]
                embedding_path = os.path.join(EMBEDDINGS_DIR, file)
                try:
                    known_embeddings[name] = np.load(embedding_path)
                    print(f"Loaded existing embedding for {name}")
                except Exception as e:
                    print(f"Error loading embedding for {name}: {e}")

    # Then check AUTH_DIR for any new images that don't have embeddings yet
    if os.path.exists(AUTH_DIR):
        for file in os.listdir(AUTH_DIR):
            if file.lower().endswith((".jpg", ".png", ".jpeg")):
                name = os.path.splitext(file)[0]
                
                # Skip if we already loaded an embedding for this user
                if name in known_embeddings:
                    continue

                path = os.path.join(AUTH_DIR, file)
                embedding_path = os.path.join(EMBEDDINGS_DIR, f"{name}.npy")
                
                try:
                    embedding = DeepFace.represent(
                        img_path=path,
                        model_name="ArcFace",
                        detector_backend="retinaface",
                        enforce_detection=True
                    )[0]["embedding"]
                    
                    embedding_np = np.array(embedding)
                    np.save(embedding_path, embedding_np)
                    known_embeddings[name] = embedding_np    
                    print(f"Created and loaded embedding for {name}")
                except Exception as e:
                    print(f"Error creating embedding for {name}: {e}")
    
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
