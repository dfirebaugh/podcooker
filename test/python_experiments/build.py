import platform
import subprocess

# Determine the current operating system
current_os = platform.system().lower()

# Define the PyInstaller options for each operating system
if current_os == "windows":
    options = ["--onefile", "audio_processing.py", "--name", "audio_processing.exe"]
elif current_os == "darwin":
    options = ["--onefile", "--windowed", "audio_processing.py", "--name", "audio_processing"]
elif current_os == "linux":
    options = ["--onefile", "--console", "audio_processing.py", "--name", "audio_processing"]
else:
    print(f"Unsupported operating system: {current_os}")
    exit()

# Call PyInstaller with the appropriate options
subprocess.run(["pyinstaller"] + options)