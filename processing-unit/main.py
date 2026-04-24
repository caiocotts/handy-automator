import os
os.environ['TF_ENABLE_ONEDNN_OPTS'] = '0'
import cv2
import face_recognition as fr
import gesture_recognition as gr
import api_interface as api
import pose as pr
import mediapipe as mp
import time
import config
from UserRegistery import UserRegistry


pr_model_path = './database/models/pose_landmarker_full.task'
 
pr_options = mp.tasks.vision.PoseLandmarkerOptions(
    base_options=gr.BaseOptions(model_asset_path=pr_model_path),
    running_mode=gr.VisionRunningMode.LIVE_STREAM,
    result_callback=pr.update_pose_result
)

user_registry = UserRegistry()

def main():
    cam = cv2.VideoCapture(config.get_ptz_camera())
    # Get the dimensions of the frame
    frame_w = int(cam.get(cv2.CAP_PROP_FRAME_WIDTH))
    frame_h = int(cam.get(cv2.CAP_PROP_FRAME_HEIGHT))

    # We need both recognizers running
    with gr.create_recognizer(gr.options) as gesture_rec, \
         pr.create_pose_landmarker(pr_options) as pose_rec:

        while True:
            ret, frame = cam.read()
            if not ret: break

            # Keep the runtime user registry in sync with loaded embeddings.
            user_registry.upsert_from_embeddings(fr.known_embeddings)

            # 1. Run face detection first to check for authorized users
            face_data = fr.detecting_bounding_box(frame)

            # Only run gesture/pose detection if an authenticated user is present
            if fr.is_authenticated_user_present():
                timestamp = int(time.time() * 1000)
                rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
                mp_image = mp.Image(image_format=mp.ImageFormat.SRGB, data=rgb_frame)

                user_id = face_data[0][1]
                user = user_registry.get_user(user_id) if user_id != "Unknown" else None
                if user_id != "Unknown" and user is not None and user.auth_token is None:
                    now = time.time()
                    
                    # Only attempt authentication if the cooldown period has passed for specific user
                    if user.can_attempt_auth(now=now):
                        user.mark_auth_attempt(now=now)
                        embedding = user.to_embedding_list()
                        
                        auth_token = api.auth_user_api_call(embedding, user_id)
                        if auth_token is not None:
                            print(f'auth token generated')
                            user.authenticate(auth_token, now=now)
                            # Store the token
                            try:
                                fr.auth_tokens[user_id] = auth_token
                            except Exception as e:
                                print(f"Error storing auth token: {e}")
            
                

                # 2. Run gesture and pose detections
                gesture_rec.recognize_async(mp_image, timestamp)
                pose_rec.detect_async(mp_image, timestamp)

                # 3. Check for Pose Landmarks
                if pr.latest_pose_result:
                    # Coordinates
                    nose_x, nose_y = pr.get_nose_coords(pr.latest_pose_result, frame_w, frame_h)
                    rw_x, rw_y = pr.get_wrist_coords(pr.latest_pose_result, frame_w, frame_h, "right")
                    lw_x, lw_y = pr.get_wrist_coords(pr.latest_pose_result, frame_w, frame_h, "left")

                    # 4. Check Gesture vs. Wrists
                    if gr.latest_result and gr.latest_result.hand_landmarks:
                        for idx, hand_landmarks in enumerate(gr.latest_result.hand_landmarks):
                            gx, gy = int(hand_landmarks[0].x * frame_w), int(hand_landmarks[0].y * frame_h)

                            if pr.get_dist(gx, gy, rw_x, rw_y) < 50 or pr.get_dist(gx, gy, lw_x, lw_y) < 50:

                                # 5. Check Nose vs. Face Boxes WITH Names
                                for (box, name) in face_data:
                                    fx, fy, fw, fh = box
                                    if fx < nose_x < fx + fw and fy < nose_y < fy + fh:

                                        # SUCCESS: Format the final display string
                                        display_name = name.upper() if name != "Unknown" else "UNKNOWN"
                                        label_text = f"{display_name} + hand"
                                        gesture_id = gr.get_hardcode_index(gr.latest_result.gestures[0][0].category_name)
                                        
                                        active_user = user_registry.get_user(name)
                                        if gesture_id != 0 and active_user is not None and active_user.auth_token is not None:
                                            current_time = time.time()
                                            
                                            # Only call the API if the cooldown has expired
                                            if active_user.can_trigger_gesture(gesture_id, now=current_time):
                                                api.workflow_api_call(gesture_id, active_user.auth_token)
                                                active_user.mark_gesture_triggered(gesture_id, now=current_time)
                                                
                                        cv2.putText(frame, label_text, (fx, fy-10),
                                                    cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 255), 2)

                                        # Draw the gesture info
                                        frame = gr.draw_landmarks_on_image(frame, gr.latest_result)

            cv2.imshow("Integrated System", frame)
            if cv2.waitKey(1) & 0xFF == ord("q"): break

    cam.release()
    cv2.destroyAllWindows()

if __name__ == "__main__":
    main()