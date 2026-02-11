import cv2
# Module Imports
import face_recognition as fr 


def main():
    cam = cv2.VideoCapture(0) # Change based on default camera

    while True:
        ret, frame = cam.read()
        if ret is False:
            break
        
        faces = fr.detecting_bounding_box(frame) # TODO Implement usage for tracking
        cv2.imshow("Bounding Box Frame", frame)

        if cv2.waitKey(1) & 0xFF == ord("q"):
            break
    cam.release()
    cv2.destroyAllWindows()
if __name__ == "__main__":
    main()
