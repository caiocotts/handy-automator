import cv2
# Module Imports
import face_recognition as fr
import gesture_recognition as gr
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

def main():
    cam = cv2.VideoCapture(0) # Change based on default camera

    with gr.create_recognizer(options) as recognizer:

        while True:
            ret, frame = cam.read()
            if ret is False:
                break

            faces = fr.detecting_bounding_box(frame) # TODO Implement usage for tracking

            rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            mp_image = mp.Image(image_format=mp.ImageFormat.SRGB, data=rgb_frame)

            frame_timestamp_ms = int(time.time() * 1000)
            if fr.is_authenticated():
                recognizer.recognize_async(mp_image, frame_timestamp_ms)

                if gr.latest_result:
                    frame = gr.draw_landmarks_on_image(frame, gr.latest_result)

            cv2.imshow("Bounding Box Frame", frame)

            if cv2.waitKey(1) & 0xFF == ord("q"):
                break
    cam.release()
    cv2.destroyAllWindows()
if __name__ == "__main__":
    main()
