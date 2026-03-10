import cv2
from pygrabber.dshow_graph import FilterGraph

def get_available_cameras()-> dict[int, str]:
    """
    This function will search for all camera inputs and save their index and 
    their name
    Returns: 
        - Dict {int Index : str Camera Name}
    """

    devices = FilterGraph().get_input_devices()

    available_cameras = {}

    for device_index, device_name in enumerate(devices):
        available_cameras[device_index] = device_name

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
        if available[keys] == "VADDIO":
            return keys
    # Change to None when development is done
    return 0


print(get_available_cameras())