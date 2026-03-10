import cv2
# Module Imports
import face_recognition as fr
import gesture_recognition as gr
import pose as pr
import mediapipe as mp
import time

model_path = './gesture_recognizer.task'
BaseOptions = mp.tasks.BaseOptions
GestureRecognizer = mp.tasks.vision.GestureRecognizer
GestureRecognizerOptions = mp.tasks.vision.GestureRecognizerOptions
GestureRecognizerResult = mp.tasks.vision.GestureRecognizerResult
VisionRunningMode = mp.tasks.vision.RunningMode
options = GestureRecognizerOptions( # TODO configure for gesture recognizer
    base_options=BaseOptions(model_asset_path=model_path), # finds which model to use (hand_landmarker.task)
    running_mode=VisionRunningMode.LIVE_STREAM, # whether an image, video, or live-stream (can only use detect_async() if live stream)
    result_callback=gr.update_result, # where the result is sent
    min_hand_detection_confidence=0.7,
    num_hands=2
)
pr_model_path = './pose_landmarker_full.task' 
pr_options = mp.tasks.vision.PoseLandmarkerOptions(
    base_options=BaseOptions(model_asset_path=pr_model_path),
    running_mode=VisionRunningMode.LIVE_STREAM,
    result_callback=pr.update_pose_result
)

def main():
    cam = cv2.VideoCapture(0)
    # Get the dimensions of the frame
    frame_w = int(cam.get(cv2.CAP_PROP_FRAME_WIDTH))
    frame_h = int(cam.get(cv2.CAP_PROP_FRAME_HEIGHT))

    # We need both recognizers running
    with gr.create_recognizer(options) as gesture_rec, \
         pr.create_pose_landmarker(pr_options) as pose_rec:

        while True:
            ret, frame = cam.read()
            if not ret: break

            # 1. Run face detection first to check for authorized users
            face_data = fr.detecting_bounding_box(frame)

            # Only run gesture/pose detection if an authenticated user is present
            if fr.is_authenticated_user_present():
                timestamp = int(time.time() * 1000)
                rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
                mp_image = mp.Image(image_format=mp.ImageFormat.SRGB, data=rgb_frame)

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