import os
import subprocess
import sys

def run_all_tests():
    """
    Finds all test files in the test_suite directory and runs them one by one.
    """
    # Get the directory where this script is located
    test_suite_dir = os.path.dirname(os.path.abspath(__file__))

    # List all files in the directory
    files = os.listdir(test_suite_dir)

    # Filter for python files that start with 'test_' and are not this file
    test_files = [f for f in files if f.startswith("test_") and f.endswith(".py") and f != os.path.basename(__file__)]

    print(f"Found {len(test_files)} test files to run.")

    for test_file in test_files:
        print(f"--- Running {test_file} ---")
        subprocess.run([sys.executable, "-m", "pytest", os.path.join(test_suite_dir, test_file)])
        print(f"--- Finished {test_file} ---\n")

if __name__ == "__main__":
    run_all_tests()
