import cv2
from cv2_enumerate_cameras import enumerate_cameras
import platform


def get_available_cameras()-> dict[int, str]:
    """
    This function will search for all camera inputs and save their index and 
    their name
    Returns: 
        - Dict {int Index : str Camera Name}
    """
    os = platform.system()

    available_cameras = {}
    if os == "Windows":
        for camera_info in enumerate_cameras(cv2.CAP_MSMF):
            if camera_info is not None:
                available_cameras[camera_info.index] = camera_info.name
        return available_cameras
    
    elif os == "Linux":
        for camera_info in enumerate_cameras(cv2.CAP_V4L2):
            if camera_info is not None:
                available_cameras[camera_info.index] = camera_info.name
        return available_cameras
    
    elif os == "Darwin":
        for camera_info in enumerate_cameras(cv2.CAP_AVFOUNDATION):
            if camera_info is not None:
                available_cameras[camera_info.index] = camera_info.name
        return available_cameras

def get_ptz_camera()-> int|None:
    """
    This function will check whether the VADDIO PTZ Camera is detected.

    Returns:
        - Int Index of the camera
        - None if not found
    """
    available = get_available_cameras()
    for keys in available:
        if "vaddio" in available[keys].lower() or "roboshot" in available[keys].lower():
            return keys
    return 0


print(get_available_cameras())